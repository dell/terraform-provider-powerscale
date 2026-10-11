# AGENTS.md - Dell Terraform Provider for PowerScale

## Project Overview

This is the Terraform provider for Dell PowerScale (Isilon) scale-out NAS storage. It implements resources and data sources using HashiCorp's Terraform Plugin Framework, enabling infrastructure-as-code management of PowerScale clusters.

- **Language:** Go 1.25
- **Module path:** `terraform-provider-powerscale`
- **Terraform Plugin Framework:** v1.15.1
- **SDK:** `dell/powerscale-go-client` (vendored, local)
- **Registry address:** `registry.terraform.io/dell/powerscale`
- **License:** Mozilla Public License 2.0

## Architecture

The provider follows the standard Terraform Plugin Framework architecture. It runs as a gRPC server that Terraform Core communicates with to manage PowerScale resources.

### Provider Configuration

The provider authenticates to a PowerScale cluster (Platform API, typically port 8080) using endpoint, username, and password. Configuration can be supplied via HCL provider block or environment variables (`POWERSCALE_ENDPOINT`, `POWERSCALE_USERNAME`, `POWERSCALE_PASSWORD`, `POWERSCALE_INSECURE`, `POWERSCALE_TIMEOUT`).

### SDK Strategy

Uses a **vendored SDK** — `dell/powerscale-go-client` lives inside the provider repo as a local directory. The `go.mod` declares:

```go
require dell/powerscale-go-client v0.0.0
replace dell/powerscale-go-client => ./powerscale-go-client
```

SDK and provider release together. Changes to the SDK require changes in the same repo.

### Resources and Data Sources

The provider exposes approximately 48 resources and 32 data sources covering PowerScale entities such as access zones, ACLs, quotas, SMB shares, NFS exports, snapshot schedules, user groups, file systems, network settings, and more.

## Directory Structure

```
main.go                           Entry point (providerserver.Serve)
powerscale/
  provider/
    provider.go                   Provider configuration, resource/datasource registration
    *_resource.go                 Resource implementations
    *_resource_schema.go          Resource schema definitions
    *_datasource.go               Data source implementations
    *_datasource_schema.go        Data source schema definitions
    *_test.go                     Unit and acceptance tests
  helper/                         Shared helper functions
  models/                         Terraform state model structs
  constants/                      Shared constants
client/                           PowerScale client wrapper
powerscale-go-client/             Vendored PowerScale Go SDK (local dependency)
goClientZip/                      SDK distribution archives
examples/                         Example HCL configurations
docs/                             Generated documentation
templates/                        Documentation templates
tools/                            Build and generation tools
about/                            Provider metadata
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Compile the provider binary |
| `make install` | Build and install to `~/.terraform.d/plugins/` |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests (`TF_ACC=1`, requires live hardware) |
| `make check` | Run `gofmt`, `golangci-lint`, `go vet` |
| `make gosec` | Run security scan with `gosec` |
| `make cover` | Generate HTML coverage report |
| `make generate` | Run `go generate` (docs generation) |

## Testing

### Unit Tests (mockey)

- Test files follow `*_test.go` convention alongside source files in `powerscale/provider/`.
- Frameworks: `github.com/stretchr/testify` (assertions), `github.com/bytedance/mockey` (function-level mocking).
- Run with `make test`.
- No hardware required.

### Acceptance Tests (terraform-plugin-testing)

- **Requires live PowerScale hardware** with credentials set via environment variables.
- Creates real resources — clean up after failures.
- Run with `make testacc`.

### Running Tests

```bash
# Unit tests (no hardware)
make test

# Acceptance tests (requires live hardware)
export POWERSCALE_ENDPOINT="https://10.0.0.1:8080"
export POWERSCALE_USERNAME="admin"
export POWERSCALE_PASSWORD="secret"
export POWERSCALE_INSECURE="true"
make testacc
```

## Code Style and Conventions

### Code Organization Patterns

- **Resource pattern:** Each resource has up to three files: `<name>_resource.go`, `<name>_resource_schema.go`, plus helpers.
- **Models:** Terraform state structs in `powerscale/models/` using `tfsdk` struct tags.
- **Helpers:** API-to-Terraform mapping functions in `powerscale/helper/`.

### File Header

All source files must include the Dell copyright and MPL 2.0 license header.

## Common Development Tasks

### Adding a New Resource

1. Create resource, schema, and model files following existing patterns.
2. Add helper functions for API-to-Terraform mapping.
3. Register in `powerscale/provider/provider.go`.
4. Add unit and acceptance tests.
5. Create example HCL in `examples/resources/powerscale_<name>/`.
6. Run `make generate` to produce documentation.

### Updating the Vendored SDK

Edit files directly in `./powerscale-go-client/`. No `go mod tidy` needed (local replace directive). Commit SDK and provider changes together.

## CI/CD

GitHub Actions workflows in `.github/workflows/`. GoReleaser configuration in `.goreleaser.yml` builds cross-platform binaries.

## Code Ownership

All files are owned by the maintainers defined in `.github/CODEOWNERS`.
