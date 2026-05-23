# driftwatch

A daemon that detects infrastructure config drift between Terraform state and live cloud resources.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Point driftwatch at your Terraform state file and let it run as a background daemon:

```bash
driftwatch --state ./terraform.tfstate --interval 5m --cloud aws
```

It will periodically compare your Terraform state against live cloud resources and report any drift via logs or a configured webhook.

### Example config (`driftwatch.yaml`)

```yaml
state: ./terraform.tfstate
interval: 5m
cloud: aws
region: us-east-1
notify:
  webhook: https://hooks.example.com/alerts
```

Run with a config file:

```bash
driftwatch --config driftwatch.yaml
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--state` | `terraform.tfstate` | Path to Terraform state file |
| `--interval` | `10m` | How often to check for drift |
| `--cloud` | `aws` | Cloud provider (`aws`, `gcp`, `azure`) |
| `--config` | `` | Path to config file |

---

## License

MIT © yourorg