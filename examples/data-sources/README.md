# PowerScale Data Sources - Discovery & Import

This directory contains tools to **discover** all resources on a PowerScale cluster and **generate Terraform import configurations**.

## Overview

| File | Purpose |
|------|---------|
| `getall.tf` | Retrieves all data from a PowerScale cluster |
| `generate_imports.py` | Generates `.tf` resource files and import commands |

## Complete Workflow

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  1. Discovery   │ --> │  2. Generation   │ --> │   3. Import     │
│   getall.tf     │     │ generate_imports │     │  terraform      │
│   terraform     │     │     .py          │     │   import        │
└─────────────────┘     └──────────────────┘     └─────────────────┘
```

---

## Step 1: Discovery - Retrieve Cluster Data

### 1.1 Create a working directory

```bash
mkdir ~/powerscale-import
cd ~/powerscale-import
cp /path/to/examples/data-sources/getall.tf .
cp /path/to/examples/data-sources/generate_imports.py .
```

### 1.2 Create `terraform.tfvars`

```hcl
powerscale_endpoint = "https://192.168.1.100:8080"
powerscale_username = "root"
powerscale_password = "your_password"
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

### 1.3 Initialize and retrieve data

```bash
terraform init
terraform apply
```

### 1.4 Export to JSON

```bash
terraform output -json > cluster_data.json
```

---

## Step 2: Generate Import Configuration

### 2.1 Run the generator

```bash
python3 generate_imports.py cluster_data.json
```

This creates an `imports/` directory with:

```
imports/
├── import.sh              # Shell script with all import commands
├── accesszones.tf         # Access zone resources
├── quotas.tf              # Quota resources
├── nfs_exports.tf         # NFS export resources
├── smb_shares.tf          # SMB share resources
├── users.tf               # User resources
├── user_groups.tf         # User group resources
└── ...                    # Other resource files
```

### 2.2 Review generated files

Each `.tf` file contains:
- Required fields populated with actual values
- Optional fields commented out for reference

Example `quotas.tf`:
```hcl
resource "powerscale_quota" "ifs_data" {
  path = "/ifs/data"
  type = "directory"

  # Optional fields (uncomment and modify as needed):
  # include_snapshots = false
  # thresholds = { ... }
}
```

---

## Step 3: Import Resources into Terraform

### 3.1 Setup the import directory

```bash
cd imports

# Copy provider configuration
cat > provider.tf << 'EOF'
terraform {
  required_providers {
    powerscale = {
      source = "registry.terraform.io/dell/powerscale"
    }
  }
}

provider "powerscale" {
  endpoint = var.powerscale_endpoint
  username = var.powerscale_username
  password = var.powerscale_password
  insecure = var.powerscale_insecure
}

variable "powerscale_endpoint" { type = string }
variable "powerscale_username" { type = string }
variable "powerscale_password" { type = string; sensitive = true }
variable "powerscale_insecure" { type = bool; default = true }
EOF

# Copy your tfvars
cp ../terraform.tfvars .

# Initialize
terraform init
```

### 3.2 Run the imports

```bash
# Import all resources
./import.sh

# Or import selectively
terraform import powerscale_quota.ifs_data "AABpAQEAAAA..."
terraform import powerscale_nfs_export.export_1 "1"
```

### 3.3 Verify

```bash
# Should show no changes if import was successful
terraform plan
```

If `terraform plan` shows differences:
1. Review the generated `.tf` file
2. Adjust attributes to match the actual state
3. Re-run `terraform plan` until clean

---

## Supported Resources

The generator supports these resource types:

| Datasource | Resource | Import ID Format |
|------------|----------|------------------|
| `accesszones` | `powerscale_accesszone` | `name` |
| `groupnets` | `powerscale_groupnet` | `name` |
| `subnets` | `powerscale_subnet` | `groupnet.subnet` |
| `networkpools` | `powerscale_networkpool` | `groupnet.subnet.pool` |
| `quotas` | `powerscale_quota` | `quota_id` |
| `nfs_exports` | `powerscale_nfs_export` | `[zone:]id` |
| `nfs_aliases` | `powerscale_nfs_alias` | `[zone:]name` |
| `smb_shares` | `powerscale_smb_share` | `zone:name` |
| `snapshots` | `powerscale_snapshot` | `name` |
| `snapshot_schedules` | `powerscale_snapshot_schedule` | `name` |
| `users` | `powerscale_user` | `[zone:]username` |
| `user_groups` | `powerscale_user_group` | `[zone:]groupname` |
| `roles` | `powerscale_role` | `name` |
| `s3_buckets` | `powerscale_s3_bucket` | `[zone:]name` |
| `synciq_policies` | `powerscale_synciq_policy` | `policy_id` |
| `filepool_policies` | `powerscale_filepool_policy` | `name` |
| `adsproviders` | `powerscale_adsprovider` | `name` |
| `ldap_providers` | `powerscale_ldap_provider` | `name` |
| `ntpservers` | `powerscale_ntpserver` | `name` |

---

## Troubleshooting

### Some datasources fail during discovery

Comment out failing datasources in `getall.tf` (e.g., if SyncIQ is not licensed).

### Import fails with "resource already exists"

The resource is already in your state. Either:
- Remove it: `terraform state rm <resource>`
- Or skip importing it

### Plan shows differences after import

This is normal. The generated `.tf` may not have all attributes. Review and adjust:
1. Check `terraform state show <resource>` for actual values
2. Update the `.tf` file to match
3. Re-run `terraform plan`

### Python script errors

Ensure Python 3.6+ is installed:
```bash
python3 --version
```

---

## Example: Full Import Session

```bash
# 1. Setup
mkdir ~/powerscale-import && cd ~/powerscale-import
cp /path/to/examples/data-sources/{getall.tf,generate_imports.py} .

# 2. Configure
cat > terraform.tfvars << 'EOF'
powerscale_endpoint = "https://10.0.0.1:8080"
powerscale_username = "root"
powerscale_password = "Password123!"
powerscale_insecure = true
EOF

# 3. Discover
terraform init
terraform apply -auto-approve
terraform output -json > cluster_data.json

# 4. Generate
python3 generate_imports.py cluster_data.json

# 5. Import
cd imports
cp ../terraform.tfvars .
cat > provider.tf << 'EOF'
terraform {
  required_providers {
    powerscale = { source = "registry.terraform.io/dell/powerscale" }
  }
}
provider "powerscale" {
  endpoint = var.powerscale_endpoint
  username = var.powerscale_username
  password = var.powerscale_password
  insecure = var.powerscale_insecure
}
variable "powerscale_endpoint" { type = string }
variable "powerscale_username" { type = string }
variable "powerscale_password" { type = string; sensitive = true }
variable "powerscale_insecure" { type = bool; default = true }
EOF

terraform init
./import.sh
terraform plan
```
