# portwatch

Lightweight daemon that monitors port bindings and alerts on unexpected listeners via webhook or desktop notification.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file:

```bash
portwatch --config /etc/portwatch/config.yaml
```

Example `config.yaml`:

```yaml
interval: 10s
allowed_ports:
  - 22
  - 80
  - 443
alerts:
  webhook: "https://hooks.example.com/notify"
  desktop: true
```

When an unexpected port binding is detected, portwatch fires a webhook POST request and/or a desktop notification with details about the new listener, including PID, process name, and port number.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `./config.yaml` | Path to config file |
| `--interval` | `10s` | Poll interval |
| `--dry-run` | `false` | Log alerts without sending |

---

## License

MIT © 2024 yourusername