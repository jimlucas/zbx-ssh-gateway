# zbx-ssh-gateway

A standalone, allowlisted SSH polling gateway for Zabbix HTTP Agent items.

Zabbix sends a named operation and validated parameters over HTTPS. The gateway owns the SSH username, ordered password list, host-key policy, allowed commands, queue limits, and execution limits. It does **not** accept arbitrary command text.

> Status: early development / MVP scaffold. Do not deploy as a production security boundary until the implementation and security tests are complete.

## Architecture

`Zabbix Server -> HTTP Agent/HTTPS -> zbx-ssh-gateway -> SSH -> network equipment`

The Zabbix Proxy and Agent 2 are not in this polling path.

## Development

```sh
make check
make build
```

## Installation

On Debian 12, obtain the source and install from the repository root:

```sh
cd /usr/local/src
git clone https://github.com/jimlucas/zbx-ssh-gateway.git
cd zbx-ssh-gateway
make check
sudo make install
```

The installer preserves existing `/etc/zbx-ssh-gateway/gateway.local.yaml` and `operations.local.yaml` files. Complete TLS, API token, SSH host-key, credential, and network configuration before starting the service.

See `docs/INSTALL.md` for the complete installation/upgrade procedure and `docs/CONFIGURATION.md` for configuration details.
