# haul

A fast CLI tool for syncing environment configs across multiple remote servers via SSH.

---

## Installation

```bash
go install github.com/yourusername/haul@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/haul/releases).

---

## Usage

Define your servers and config files in a `haul.yaml` file, then run:

```bash
haul sync --config haul.yaml
```

**Example `haul.yaml`:**

```yaml
servers:
  - host: web1.example.com
    user: deploy
  - host: web2.example.com
    user: deploy

files:
  - src: .env.production
    dest: /var/www/app/.env
```

**Common commands:**

```bash
haul sync               # Sync configs to all servers
haul sync --dry-run     # Preview changes without applying
haul check              # Verify SSH connectivity to all servers
haul --help             # Show help
```

---

## Requirements

- Go 1.21+
- SSH access to target servers

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)