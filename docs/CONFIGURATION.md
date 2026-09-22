# Configuration

The gateway uses two site-local YAML files:

- `/etc/zbx-ssh-gateway/gateway.local.yaml` — listener, TLS, API authentication, SSH credentials, limits, and the operations-file location.
- `/etc/zbx-ssh-gateway/operations.local.yaml` — the allowlisted operations Zabbix is permitted to request.

Files matching `*.local.yaml` are intentionally ignored by Git so deployments and upgrades do not overwrite site-specific configuration.

The version-controlled templates are `configs/gateway.example.yaml` and `configs/operations.example.yaml`.

## Linking the files

`gateway.local.yaml` contains:

```yaml
operations_file: "/etc/zbx-ssh-gateway/operations.local.yaml"
```

A relative path is also supported and is resolved relative to the directory containing `gateway.local.yaml`.

## SSH credentials

`ssh.username` is the configured SSH username. `ssh.password_list.name` is an administrative name and `ssh.password_list.passwords` is ordered. **Every new request begins with the first password. The gateway does not cache, promote, or remember which password succeeded.**

Only an SSH authentication failure advances to the next password. Connection, negotiation, and host-key failures stop the attempt.

## Custom operations

Put all site-specific operations in `/etc/zbx-ssh-gateway/operations.local.yaml`:

```yaml
operations:
  radio.stats:
    command: "radio stats ${interface}"
    parameters:
      interface:
        type: enum
        values: ["wlan0", "wlan1"]
        required: true
    timeout_seconds: 10
    max_output_bytes: 65536
```

Zabbix supplies the operation name, never command text. Unknown JSON fields and unknown operation parameters are rejected. Validator types are `enum`, `regex`, and bounded `string`; prefer `enum` whenever possible.

## Zabbix HTTP Agent

POST JSON to `https://GATEWAY:9443/api/v1/execute` with `Content-Type: application/json` and `Authorization: Bearer <token>`. SSH passwords must never be stored in Zabbix.

For scalar checks, JSONPath preprocessing can extract `$.value`. For commands returning multiple metrics, use a master HTTP Agent item and dependent items.

## TLS, host keys, and limits

Use TLS in production and permit the gateway port only from the primary Zabbix Server. Strict SSH host-key verification uses `known_hosts`; changed keys are errors. `max_workers`, `max_queue`, `max_request_bytes`, and `max_output_bytes` bound resource consumption.
