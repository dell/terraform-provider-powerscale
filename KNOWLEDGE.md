# KNOWLEDGE.md — terraform-provider-powerscale

<!-- yaml-metadata-start -->
scope_paths: ["./"]
capture_git_sha: "88363a0b5011c1541b3d877d1ba48771432c2b0f"
status: "current"
auto_update: false
preview_before_apply: true
scaffold_version: "1.0"
# session_state: { is_complete: true }
<!-- yaml-metadata-end -->

<!-- quick-reference-start -->
## Agent Quick Reference

| Section | Heading | Summary | never_again_count |
|---------|---------|---------|-------------------|
| Component Overview | `## Component Overview` | Dell PowerScale (Isilon) scale-out NAS storage provider | — |
| Architectural Rationale | `## Architectural Rationale` | Vendored SDK strategy; Plugin Framework architecture | — |
| Failure Modes & Gotchas | `## Failure Modes & Gotchas` | Endpoint format, SDK versioning, state secrets | 0 |
| Implicit Contracts | `## Implicit Contracts` | Env var precedence, auth validation, TLS defaults | — |
<!-- quick-reference-end -->

## Five Questions Quick Reference

### What does it do?
Terraform provider for Dell PowerScale (Isilon) scale-out NAS storage. Exposes 48 resources and 32 data sources covering access zones, ACLs, quotas, SMB shares, NFS exports, snapshot schedules, user groups, file systems, network settings, cluster configuration, S3 buckets, SyncIQ policies, and more
through HashiCorp's Terraform Plugin Framework. Communicates with
the hardware REST API via `dell/powerscale-go-client` (vendored, local).

### How do you modify it?
Create `resource_<name>.go` (or `*_resource.go`) implementing
`resource.Resource`, add model structs, register in `provider.go`,
add unit tests with mockey mocks, add acceptance tests, create
example HCL, and run `make generate` for docs.

### What breaks?
**Endpoint must use the Platform API port** — typically port 8080, not the default HTTPS port 443. Acceptance tests against live hardware create real
resources — failed test runs may leave orphaned resources. State files
contain secrets — use encrypted remote backends.

### What depends on it?
Terraform Core (gRPC go-plugin), `dell/powerscale-go-client` (vendored, local),
`hashicorp/terraform-plugin-framework` v1.15.1.

### What's undocumented?
Client wraps the vendored `powerscale-go-client`. SDK archive distributions are in `goClientZip/`.

---

## Component Overview

Terraform provider for Dell PowerScale (Isilon) scale-out NAS storage.
48 resources and 32 data sources covering access zones, ACLs, quotas, SMB shares, NFS exports, snapshot schedules, user groups, file systems, network settings, cluster configuration, S3 buckets, SyncIQ policies, and more. Resources use `*_resource.go` naming under `powerscale/provider/`. Schema definitions are in separate `*_resource_schema.go` files.

---

## Architectural Rationale

The provider follows the standard Terraform Plugin Framework architecture
— a standalone Go binary communicating with Terraform Core over gRPC.

**SDK strategy (Vendored):** Vendored SDK — lives inside the provider repo as `./powerscale-go-client/`. The `go.mod` declares `require dell/powerscale-go-client v0.0.0` with `replace => ./powerscale-go-client`. SDK and provider release together. Edit SDK files directly in `./powerscale-go-client/`. No `go mod tidy` needed (local replace directive). Commit SDK and provider changes together.

All providers in the Dell Terraform family share this architecture:
Terraform Plugin Framework interfaces, `resource.Resource` for CRUD
resources, `datasource.DataSource` for read-only queries, models with
`tfsdk` struct tags, and mockey-based unit testing.

### Evolution

TBD — requires SME input on how the architecture changed over time.

---

## Failure Modes & Gotchas

### 1. Endpoint URL format

**Endpoint must use the Platform API port** — typically port 8080, not the default HTTPS port 443.

### 2. Sensitive attributes must be marked

All credential fields must have `Sensitive: true` in the schema.
Without this, passwords appear in `terraform plan` output and state
files. This is enforced by code convention, not by the framework.

### 3. State file contains secrets

Terraform state files contain full resource representations including
credentials. Always use encrypted remote backends (S3+KMS, Terraform
Cloud) in production.

### 4. Vendored SDK updates

SDK changes require editing files in `./powerscale-go-client/` directly. No `go mod tidy` needed. Commit SDK and provider changes together. The `goClientZip/` directory holds SDK distribution archives.

### 5. Schema split pattern

Unlike other providers, PowerScale separates resource logic (`*_resource.go`) from schema definitions (`*_resource_schema.go`). Both must be updated when modifying a resource's attributes.

### Never Again

No incident-derived constraints recorded. If you know of past
incidents affecting this component, please record them during the
next Knowledge Extraction session.

### Evolution

TBD — requires SME input.

---

## Performance Characteristics

TBD — requires SME input for bottlenecks, scaling limits, tuning
parameters, benchmarks, and known performance cliffs.

### Evolution

TBD — requires SME input.

---

## Implicit Contracts

**Environment variable precedence:** env vars (`POWERSCALE_*`)
override HCL provider block values when both are set. This is
implemented in `Configure()` and is not documented as an explicit
contract.

**Authentication validation:** `Configure()` makes a dummy API call
to validate credentials before any resource operations proceed. If
this call fails, all resource operations are blocked.

**TLS verification default:** `insecure` defaults to `false` —
TLS verification is on by default. Setting `insecure = true` is
a lab-only setting and must never be used in production.

**Acceptance test gating:** tests guarded by `TF_ACC=1` — never
run without live hardware credentials. Tests create real resources
that must be cleaned up manually if the test run fails.

### Evolution

TBD — requires SME input.

---

## Threading & Synchronization

Terraform Plugin Framework handles concurrency at the provider level.
Individual resource operations are not concurrent by default.

### Evolution

TBD — requires SME input.

---

## Build System & Configuration

Standard Makefile targets shared across all Dell Terraform providers:

| Target | Purpose | Hardware Required |
|--------|---------|-------------------|
| `make build` | Compile provider binary | No |
| `make install` | Install to `~/.terraform.d/plugins/` | No |
| `make test` | Run unit tests | No |
| `make testacc` | Run acceptance tests | **Yes** |
| `make check` | Format, lint, vet | No |
| `make gosec` | Security scan | No |
| `make cover` | Generate coverage report | No |
| `make generate` | Generate documentation | No |

GoReleaser configuration: CGO_ENABLED=0, platforms (freebsd, windows,
linux, darwin), architectures (amd64, 386, arm, arm64).

### Evolution

TBD — requires SME input.

---

## Operational Knowledge

**Unit tests:** `bytedance/mockey` for runtime function patching.
No hardware required. Run with `make test`.

**Acceptance tests:** `terraform-plugin-testing` against live hardware.
Creates real resources. Run with `TF_ACC=1 make testacc`. Clean up
manually if tests fail mid-run.

### Evolution

TBD — requires SME input.

---

## General Context

### Open Issues

TBD — requires code scanning for TODO/FIXME/HACK markers.

### Glossary

| Term | Definition |
|------|------------|
| Plugin Framework | HashiCorp's Terraform Plugin Framework (`terraform-plugin-framework`) |
| mockey | `bytedance/mockey` — runtime function patching for unit tests |
| POWERSCALE | Environment variable prefix for this provider |

---

## References

- [Terraform Plugin Framework Docs](https://developer.hashicorp.com/terraform/plugin/framework)
- [Dell Terraform Registry](https://registry.terraform.io/namespaces/dell)

---

## Governance Spec Discrepancies

No discrepancies detected between code/SME knowledge and loaded
governance specs.
