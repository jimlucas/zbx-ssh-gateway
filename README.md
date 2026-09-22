# zbx-ssh-gateway

A standalone, allowlisted SSH polling gateway for Zabbix HTTP Agent items.

Zabbix sends a named operation and validated parameters over HTTPS. The gateway owns the SSH username, ordered password list, host-key policy, allowed commands, queue limits, and execution limits. It does **not** accept arbitrary command text.

> Status: early development / MVP scaffold. Do not deploy as a production security boundary until the implementation and security tests are complete.

## Architecture

`Zabbix Server -> HTTP Agent/HTTPS -> zbx-ssh-gateway -> SSH -> network equipment`

The Zabbix Proxy and Agent 2 are not in this polling path.

## Development

```sh
go test ./...
go vet ./...
go build ./cmd/zbx-ssh-gateway ./cmd/zbx-ssh-gatewayctl
```

See `docs/INSTALL.md` and `docs/CONFIGURATION.md`.
