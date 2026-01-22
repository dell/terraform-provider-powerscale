# PowerScale Data Sources - Get All Configuration

This directory contains `getall.tf`, a Terraform configuration that retrieves **all** data from a PowerScale cluster using every available data source.

## Purpose

Use this configuration to:
- **Audit** your PowerScale cluster configuration
- **Export** the current state to a Terraform state file
- **Discover** all configured resources (exports, shares, quotas, users, etc.)
- **Generate** a baseline before making changes

## Quick Start

### 1. Create a working directory

```bash
mkdir ~/powerscale-audit
cd ~/powerscale-audit
cp /path/to/terraform-provider-powerscale/examples/data-sources/getall.tf .
```

### 2. Create `terraform.tfvars`

Create a `terraform.tfvars` file with your PowerScale credentials:

```hcl
powerscale_endpoint = "https://192.168.1.100:8080"
powerscale_username = "root"
powerscale_password = "your_password_here"
powerscale_insecure = true
```

| Variable | Description | Example |
|----------|-------------|---------|
| `powerscale_endpoint` | PowerScale API URL with port 8080 | `https://isilon.example.com:8080` |
| `powerscale_username` | API username (usually `root`) | `root` |
| `powerscale_password` | API password | `secretpassword` |
| `powerscale_insecure` | Skip SSL verification (`true` for self-signed certs) | `true` |
| `powerscale_timeout` | API timeout in ms (optional, default: 2000) | `5000` |

> **Security Note:** Never commit `terraform.tfvars` to version control. Add it to `.gitignore`.

### 3. Initialize Terraform

```bash
terraform init
```

### 4. Retrieve all data

```bash
# First run - creates terraform.tfstate with all cluster data
terraform apply

# Or use refresh to update an existing state
terraform refresh
```

### 5. View the results

```bash
# Show all outputs
terraform output

# Show specific output (e.g., NFS exports)
terraform output nfs_exports

# Export to JSON for processing
terraform output -json > powerscale_config.json

# Show state in human-readable format
terraform show
```

## Updating the State

To refresh the state with current cluster data:

```bash
terraform refresh
```

This updates `terraform.tfstate` without making any changes to the cluster.

## Available Data Sources

The configuration retrieves:

| Category | Data Sources |
|----------|--------------|
| **Cluster** | cluster, cluster_email |
| **Networking** | groupnets, subnets, networkpools, network_settings, network_rules |
| **Auth** | accesszones, aclsettings, adsproviders, ldap_providers, user_mapping_rules |
| **Users** | users, user_groups, roles, role_privileges |
| **Quotas** | quotas |
| **NFS** | nfs_exports, nfs_aliases, nfs_global_settings, nfs_export_settings, nfs_zone_settings |
| **SMB** | smb_shares, smb_server_settings, smb_share_settings |
| **S3** | s3_buckets |
| **Snapshots** | snapshots, snapshot_schedules, writable_snapshots |
| **SyncIQ** | synciq_global_settings, synciq_policies, synciq_rules, synciq_peer_certificates, synciq_replication_jobs, synciq_replication_reports |
| **Storage** | smartpool_settings, storagepool_tiers, filepool_policies |
| **NTP** | ntpservers, ntpsettings |

## Troubleshooting

### Some data sources fail

If certain data sources fail (e.g., SyncIQ not licensed), comment them out in `getall.tf`:

```hcl
# Comment out if SyncIQ is not configured
# data "powerscale_synciq_policy" "all" {}
# output "synciq_policies" { value = data.powerscale_synciq_policy.all }
```

### Connection timeout

Increase the timeout in your `terraform.tfvars` (default is 2000ms):

```hcl
powerscale_timeout = 5000
```

### SSL certificate errors

Set `powerscale_insecure = true` in your `terraform.tfvars` for self-signed certificates.

## Example Output

```bash
$ terraform output cluster
{
  "config" = {
    "name" = "MyCluster"
    "guid" = "xxxx-xxxx-xxxx"
    ...
  }
  "nodes" = [...]
}

$ terraform output nfs_exports
{
  "exports" = [
    {
      "id" = 1
      "paths" = ["/ifs/data"]
      "clients" = ["*"]
      ...
    }
  ]
}
```
