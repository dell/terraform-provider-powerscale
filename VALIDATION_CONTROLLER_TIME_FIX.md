# Validation Guide for controller_time Fix

## Overview
This guide describes how to validate the fix for the `controller_time` inconsistency error in the `powerscale_adsprovider` resource.

## Problem
The `controller_time` attribute is a live epoch timestamp from the domain controller that changes on every API read. This caused Terraform to raise "Provider produced inconsistent result after apply" errors when any other attribute triggered an update cycle.

## Solution
Changed `controller_time` schema from `Optional: true, Computed: true` to `Computed: true` only. This makes the attribute computed-only, which means:
- Terraform will always show it as (known after apply) in plans
- The post-apply consistency check will accept any value the provider returns
- `ignore_changes` becomes unnecessary

## Changes Made

### 1. Schema Change
**File:** `powerscale/provider/ads_provider_resource.go`
- Removed `Optional: true` from `controller_time` attribute (line 107)
- Kept only `Computed: true`

### 2. Test Addition
**File:** `powerscale/provider/ads_provider_resource_test.go`
- Added `TestAccAdsProviderControllerTimeComputedOnly` test
- Validates that `controller_time` is present in state as a computed field
- Validates that updates don't fail due to `controller_time` inconsistency errors

## Validation Steps

### Option 1: Run Unit Tests
```bash
# From the provider repository root
go test -v -run TestAccAdsProviderControllerTimeComputedOnly ./powerscale/provider/
```

### Option 2: Manual Terraform Test
1. Build the provider with the fix:
```bash
make build
```

2. Create a test Terraform configuration:
```hcl
terraform {
  required_providers {
    powerscale = {
      source = "dell/powerscale"
      version = "1.8.1" # or your custom built version
    }
  }
}

provider "powerscale" {
  username = var.username
  password = var.password
  endpoint = var.endpoint
  insecure = true
}

resource "powerscale_adsprovider" "example" {
  name     = "EXAMPLE.COM"
  user     = var.ad_user
  password = var.ad_password
}
```

3. Import an existing AD provider:
```bash
terraform import 'powerscale_adsprovider.example' 'EXAMPLE.COM'
```

4. Run terraform plan:
```bash
terraform plan
```

5. Run terraform apply:
```bash
terraform apply
```

6. Verify that the apply succeeds without the inconsistency error.

7. Make a change to trigger an update (e.g., change `lookup_users`):
```hcl
resource "powerscale_adsprovider" "example" {
  name         = "EXAMPLE.COM"
  user         = var.ad_user
  password     = var.ad_password
  lookup_users = false
}
```

8. Run terraform plan and apply again:
```bash
terraform plan
terraform apply
```

9. Verify that the update succeeds without the inconsistency error.

### Option 3: Validate in terraform-powerscale-test Repository
If you have access to the terraform-powerscale-test repository, you can:
1. Update the provider reference to use your custom-built provider
2. Run the existing ADS provider acceptance tests
3. The new test `TestAccAdsProviderControllerTimeComputedOnly` will specifically validate the fix

## Expected Results
- ✅ Terraform apply succeeds without "Provider produced inconsistent result after apply" error
- ✅ Subsequent terraform plan shows no changes (after initial import)
- ✅ Updates to other attributes succeed without controller_time causing errors
- ✅ `controller_time` appears in state as a computed field
- ✅ No need for `lifecycle { ignore_changes = [controller_time] }` block

## Verification Checklist
- [ ] Schema change is present in `ads_provider_resource.go`
- [ ] Test `TestAccAdsProviderControllerTimeComputedOnly` exists
- [ ] Manual terraform apply succeeds without inconsistency error
- [ ] Updates to other attributes succeed
- [ ] `ignore_changes` for controller_time is no longer needed

## Rollback Plan
If issues are encountered, revert the commits:
```bash
git revert 37c22cf bc06b30
```
