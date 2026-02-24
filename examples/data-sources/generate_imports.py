#!/usr/bin/env python3
"""
PowerScale Terraform Import Generator

This script reads the terraform output from getall.tf and generates:
1. Resource blocks (.tf files) for all importable resources
2. Import commands (import.sh) to import existing resources into Terraform state

Usage:
    terraform output -json > output.json
    python3 generate_imports.py output.json

    # Or directly:
    terraform output -json | python3 generate_imports.py -

Output:
    imports/
    ├── import.sh              # All import commands
    ├── accesszones.tf         # Access zone resources
    ├── quotas.tf              # Quota resources
    ├── nfs_exports.tf         # NFS export resources
    └── ...                    # Other resource files
"""

import json
import sys
import os
import re
from pathlib import Path


# Mapping: datasource output name -> (resource_type, id_field, name_field, import_format_func)
RESOURCE_MAPPINGS = {
    "accesszones": {
        "resource_type": "powerscale_accesszone",
        "data_key": "access_zones_details",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name"],
        "skip_fields": ["id", "zone_id", "system"],
    },
    "groupnets": {
        "resource_type": "powerscale_groupnet",
        "data_key": "groupnets_details",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name"],
        "skip_fields": ["id"],
    },
    "subnets": {
        "resource_type": "powerscale_subnet",
        "data_key": "subnets",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item['groupnet']}.{item['name']}",
        "required_fields": ["name", "groupnet", "addr_family", "prefixlen"],
        "skip_fields": ["id", "pools", "base_addr"],
    },
    "networkpools": {
        "resource_type": "powerscale_networkpool",
        "data_key": "network_pools_details",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item['groupnet']}.{item['subnet']}.{item['name']}",
        "required_fields": ["name", "groupnet", "subnet"],
        "skip_fields": ["id", "addr_family", "rules", "sc_suspended_nodes"],
    },
    "quotas": {
        "resource_type": "powerscale_quota",
        "data_key": "quotas",
        "id_field": "id",
        "name_field": "path",
        "import_id": lambda item: item["id"],
        "required_fields": ["path", "type"],
        "skip_fields": ["id", "linked", "ready", "usage"],
    },
    "nfs_exports": {
        "resource_type": "powerscale_nfs_export",
        "data_key": "nfs_exports",
        "id_field": "id",
        "name_field": "id",
        "import_id": lambda item: f"{item.get('zone', 'System')}:{item['id']}" if item.get('zone') else str(item['id']),
        "required_fields": ["paths"],
        "skip_fields": ["id", "conflicting_paths", "unresolved_clients"],
    },
    "nfs_aliases": {
        "resource_type": "powerscale_nfs_alias",
        "data_key": "nfs_aliases",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item.get('zone', 'System')}:{item['name']}" if item.get('zone') else item['name'],
        "required_fields": ["name", "path"],
        "skip_fields": ["id"],
    },
    "smb_shares": {
        "resource_type": "powerscale_smb_share",
        "data_key": "smb_shares",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item.get('zone', 'System')}:{item['name']}",
        "required_fields": ["name", "path"],
        "skip_fields": ["id", "zid"],
    },
    "snapshots": {
        "resource_type": "powerscale_snapshot",
        "data_key": "snapshots",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name", "path"],
        "skip_fields": ["id", "created", "expires", "has_locks", "pct_filesystem", "pct_reserve", "size", "state", "target_id", "target_name"],
    },
    "snapshot_schedules": {
        "resource_type": "powerscale_snapshot_schedule",
        "data_key": "snapshot_schedules",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name", "path", "schedule"],
        "skip_fields": ["id", "next_run", "next_snapshot"],
    },
    "users": {
        "resource_type": "powerscale_user",
        "data_key": "users",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item.get('query_zone', 'System')}:{item['name']}" if item.get('query_zone') else item['name'],
        "required_fields": ["name"],
        "skip_fields": ["id", "dn", "dns_domain", "expired", "generated_gid", "generated_uid", "generated_upn",
                       "gid", "locked", "max_password_age", "password_expired", "password_expiry", "password_last_set",
                       "primary_group_sid", "provider_name", "sam_account_name", "type", "upn", "user_can_change_password"],
    },
    "user_groups": {
        "resource_type": "powerscale_user_group",
        "data_key": "user_groups",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item.get('query_zone', 'System')}:{item['name']}" if item.get('query_zone') else item['name'],
        "required_fields": ["name"],
        "skip_fields": ["id", "dn", "dns_domain", "generated_gid", "provider_name", "sam_account_name", "type"],
    },
    "roles": {
        "resource_type": "powerscale_role",
        "data_key": "roles",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name"],
        "skip_fields": ["id"],
    },
    "s3_buckets": {
        "resource_type": "powerscale_s3_bucket",
        "data_key": "s3_buckets",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: f"{item.get('zone', 'System')}:{item['name']}" if item.get('zone') else item['name'],
        "required_fields": ["name", "path"],
        "skip_fields": ["id", "zid"],
    },
    "synciq_policies": {
        "resource_type": "powerscale_synciq_policy",
        "data_key": "policies",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["id"],
        "required_fields": ["name", "action", "source_root_path", "target_host", "target_path"],
        "skip_fields": ["id", "last_job_state", "last_started", "last_success", "next_run", "linked_service_policies"],
    },
    "filepool_policies": {
        "resource_type": "powerscale_filepool_policy",
        "data_key": "filepool_policies",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name"],
        "skip_fields": ["id", "state", "state_details"],
    },
    "adsproviders": {
        "resource_type": "powerscale_adsprovider",
        "data_key": "ads_providers",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name"],
        "skip_fields": ["id", "site", "status", "forest", "dc_name", "recommended_spns", "spns"],
    },
    "ldap_providers": {
        "resource_type": "powerscale_ldap_provider",
        "data_key": "ldap_providers",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name", "server_uris", "base_dn"],
        "skip_fields": ["id", "status"],
    },
    "ntpservers": {
        "resource_type": "powerscale_ntpserver",
        "data_key": "ntp_servers",
        "id_field": "id",
        "name_field": "name",
        "import_id": lambda item: item["name"],
        "required_fields": ["name"],
        "skip_fields": ["id"],
    },
}


def sanitize_name(name):
    """Convert a name to a valid Terraform resource name."""
    # Replace invalid characters with underscores
    sanitized = re.sub(r'[^a-zA-Z0-9_]', '_', str(name))
    # Ensure it starts with a letter or underscore
    if sanitized and sanitized[0].isdigit():
        sanitized = '_' + sanitized
    # Avoid empty names
    if not sanitized:
        sanitized = '_unnamed'
    return sanitized.lower()


def format_value(value, indent=4):
    """Format a value for Terraform HCL."""
    indent_str = ' ' * indent

    if value is None:
        return 'null'
    elif isinstance(value, bool):
        return 'true' if value else 'false'
    elif isinstance(value, (int, float)):
        return str(value)
    elif isinstance(value, str):
        # Escape special characters
        escaped = value.replace('\\', '\\\\').replace('"', '\\"').replace('\n', '\\n')
        return f'"{escaped}"'
    elif isinstance(value, list):
        if not value:
            return '[]'
        # Check if it's a list of simple values or objects
        if all(isinstance(v, (str, int, float, bool)) or v is None for v in value):
            items = [format_value(v) for v in value]
            return '[' + ', '.join(items) + ']'
        else:
            # List of objects
            lines = ['[']
            for item in value:
                lines.append(indent_str + '  {')
                for k, v in item.items():
                    lines.append(f'{indent_str}    {k} = {format_value(v, indent + 4)}')
                lines.append(indent_str + '  },')
            lines.append(indent_str + ']')
            return '\n'.join(lines)
    elif isinstance(value, dict):
        if not value:
            return '{}'
        lines = ['{']
        for k, v in value.items():
            lines.append(f'{indent_str}  {k} = {format_value(v, indent + 2)}')
        lines.append(indent_str + '}')
        return '\n'.join(lines)
    else:
        return f'"{value}"'


def generate_resource_block(resource_type, name, item, required_fields, skip_fields):
    """Generate a Terraform resource block."""
    lines = [f'resource "{resource_type}" "{sanitize_name(name)}" {{']

    # Add required fields first
    for field in required_fields:
        if field in item:
            lines.append(f'  {field} = {format_value(item[field])}')

    # Add optional fields as comments
    lines.append('')
    lines.append('  # Optional fields (uncomment and modify as needed):')
    for key, value in sorted(item.items()):
        if key not in required_fields and key not in skip_fields and value is not None:
            formatted = format_value(value)
            # Multi-line values need special handling in comments
            if '\n' in formatted:
                lines.append(f'  # {key} = ...')
            else:
                lines.append(f'  # {key} = {formatted}')

    lines.append('}')
    lines.append('')
    return '\n'.join(lines)


def process_datasource(name, data, mapping):
    """Process a single datasource and generate resources + imports."""
    resources = []
    imports = []

    # Get the actual data from the value
    value = data.get('value', {})
    data_key = mapping['data_key']
    items = value.get(data_key, [])

    if not items:
        return [], []

    resource_type = mapping['resource_type']
    name_field = mapping['name_field']

    for item in items:
        if not isinstance(item, dict):
            continue

        item_name = item.get(name_field, 'unknown')

        # Generate resource block
        resource_block = generate_resource_block(
            resource_type,
            item_name,
            item,
            mapping['required_fields'],
            mapping['skip_fields']
        )
        resources.append(resource_block)

        # Generate import command
        try:
            import_id = mapping['import_id'](item)
            import_cmd = f'terraform import {resource_type}.{sanitize_name(item_name)} "{import_id}"'
            imports.append(import_cmd)
        except (KeyError, TypeError) as e:
            imports.append(f'# ERROR: Could not generate import for {item_name}: {e}')

    return resources, imports


def main():
    if len(sys.argv) < 2:
        print(__doc__)
        print("Error: Please provide the JSON file path or use '-' for stdin")
        sys.exit(1)

    input_path = sys.argv[1]

    # Read input
    if input_path == '-':
        data = json.load(sys.stdin)
    else:
        with open(input_path, 'r') as f:
            data = json.load(f)

    # Create output directory
    output_dir = Path('imports')
    output_dir.mkdir(exist_ok=True)

    all_imports = ['#!/bin/bash', '', '# PowerScale Terraform Import Commands',
                   '# Generated by generate_imports.py', '', 'set -e', '']

    # Process each datasource
    for ds_name, mapping in RESOURCE_MAPPINGS.items():
        if ds_name not in data:
            print(f"Skipping {ds_name}: not found in output")
            continue

        resources, imports = process_datasource(ds_name, data[ds_name], mapping)

        if resources:
            # Write resource file
            resource_file = output_dir / f'{ds_name}.tf'
            with open(resource_file, 'w') as f:
                f.write(f'# {mapping["resource_type"]} resources\n')
                f.write(f'# Generated from PowerScale cluster data\n\n')
                f.write('\n'.join(resources))
            print(f"Generated {resource_file} ({len(resources)} resources)")

            # Add imports to the batch file
            all_imports.append(f'# === {ds_name} ===')
            all_imports.extend(imports)
            all_imports.append('')

    # Write import script
    import_file = output_dir / 'import.sh'
    with open(import_file, 'w') as f:
        f.write('\n'.join(all_imports))
    os.chmod(import_file, 0o755)
    print(f"\nGenerated {import_file}")

    print(f"\n✅ Done! Files generated in '{output_dir}/' directory")
    print("\nNext steps:")
    print("  1. Review and adjust the generated .tf files")
    print("  2. Copy provider.tf and variables.tf to the imports/ directory")
    print("  3. Run: cd imports && terraform init")
    print("  4. Run: ./import.sh")
    print("  5. Run: terraform plan (should show no changes)")


if __name__ == '__main__':
    main()
