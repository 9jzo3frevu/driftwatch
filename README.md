# driftwatch

> CLI tool that detects configuration drift between deployed services and their declared infrastructure-as-code state.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your IaC definition and a target environment to scan for drift:

```bash
# Check drift against a Terraform state file
driftwatch scan --config ./infra/terraform.tfstate --env production

# Output results as JSON
driftwatch scan --config ./infra/terraform.tfstate --env staging --output json

# Watch for drift continuously (every 5 minutes)
driftwatch watch --config ./infra/terraform.tfstate --env production --interval 5m
```

Example output:

```
[DRIFT DETECTED] service: api-gateway
  expected: instance_type = t3.medium
  actual:   instance_type = t3.large

[OK] service: auth-service
[OK] service: database-primary

2 services checked. 1 drift(s) found.
```

---

## Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Path to IaC state or config file | required |
| `--env` | Target environment to inspect | required |
| `--output` | Output format: `text`, `json` | `text` |
| `--interval` | Poll interval for `watch` mode | `10m` |

---

## License

MIT © 2024 yourusername