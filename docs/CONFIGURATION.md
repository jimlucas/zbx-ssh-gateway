# Configuration

The gateway reads YAML from `/etc/zbx-ssh-gateway/gateway.yaml` by default.

## SSH credentials

`ssh.username` is the configured SSH username. `ssh.password_list.name` is an administrative name and `ssh.password_list.passwords` is ordered. **Every new request begins with the first password. The gateway does not cache, promote, or remember which password succeeded.**

Only an SSH authentication failure advances to the next password. Connection, negotiation, and host-key failures stop the attempt.

## Operations

Zabbix supplies an operation name, never command text. Example:

```yaml
operations:
  radio.stats:
    command: "radio stats ${interface}"
    parameters:
      interface:
        type: enum
        values: ["wlan0", "wlan1"]
        required: true
```

Request:

```json
{"target":"192.0.2.10","operation":"radio.stats","parameters":{"interface":"wlan0"}}
```

Unknown JSON fields and unknown operation parameters are rejected. Validator types are `enum`, `regex`, and bounded `string`; prefer `enum` whenever possible.

## Zabbix HTTP Agent

POST JSON to `https://GATEWAY:9443/api/v1/execute` with `Content-Type: application/json` and `Authorization: Bearer <token>`. Store the API token as a protected Zabbix secret macro where available. SSH passwords must never be stored in Zabbix.

For scalar checks, JSONPath preprocessing can extract `$.value`. For commands returning multiple metrics, use a master HTTP Agent item and dependent items.

## TLS, host keys, and limits

Use TLS in production and permit the gateway port only from the primary Zabbix Server. Strict SSH host-key verification uses `known_hosts`; changed keys are errors. `max_workers`, `max_queue`, `max_request_bytes`, and `max_output_bytes` bound resource consumption.
