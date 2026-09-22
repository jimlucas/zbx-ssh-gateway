# Installation and Upgrade

## Debian 12 installation

Create a dedicated account and protected configuration directories:

```sh
sudo useradd --system --home /var/lib/zbx-ssh-gateway --create-home --shell /usr/sbin/nologin zbx-ssh-gateway
sudo install -d -o root -g zbx-ssh-gateway -m 0750 /etc/zbx-ssh-gateway /etc/zbx-ssh-gateway/secrets /etc/zbx-ssh-gateway/tls
```

Build and install:

```sh
go build -o zbx-ssh-gateway ./cmd/zbx-ssh-gateway
go build -o zbx-ssh-gatewayctl ./cmd/zbx-ssh-gatewayctl
sudo install -o root -g root -m 0755 zbx-ssh-gateway /usr/sbin/zbx-ssh-gateway
sudo install -o root -g root -m 0755 zbx-ssh-gatewayctl /usr/bin/zbx-ssh-gatewayctl
sudo install -o root -g zbx-ssh-gateway -m 0640 configs/gateway.example.yaml /etc/zbx-ssh-gateway/gateway.local.yaml
sudo cp packaging/systemd/zbx-ssh-gateway.service /etc/systemd/system/
```

`gateway.local.yaml` is the site's persistent configuration. It is not a version-controlled deployment file and must not be overwritten during upgrades. Put custom operations under its top-level `operations:` section.

Create the bearer-token file and TLS certificate/key, populate `known_hosts`, and protect configuration permissions. Restrict TCP/9443 at the firewall to the primary Zabbix Server.

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now zbx-ssh-gateway
sudo systemctl status zbx-ssh-gateway
```

## Upgrade

1. Back up `/etc/zbx-ssh-gateway`, including `gateway.local.yaml`.
2. Review release notes and compare `configs/gateway.example.yaml` with the installed `gateway.local.yaml` for newly introduced settings.
3. Do **not** copy the example YAML over `gateway.local.yaml` during an upgrade.
4. Build or obtain the new binary and run the test suite for source builds.
5. Stop the service and replace the binary while preserving local configuration.
6. Install an updated systemd unit if the release changes it, then run `systemctl daemon-reload`.
7. Start the service and verify `/api/v1/health`.
8. Confirm a test Zabbix HTTP Agent item before broad polling resumes.

Rollback by restoring the prior binary and configuration backup and restarting the service.
