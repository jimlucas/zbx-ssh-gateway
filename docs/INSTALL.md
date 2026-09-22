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
sudo install -o root -g zbx-ssh-gateway -m 0640 configs/operations.example.yaml /etc/zbx-ssh-gateway/operations.local.yaml
sudo cp packaging/systemd/zbx-ssh-gateway.service /etc/systemd/system/
```

Both `gateway.local.yaml` and `operations.local.yaml` are persistent site configuration and must not be overwritten during upgrades. Put custom operations only in `operations.local.yaml`.

Create the bearer-token file and TLS certificate/key, populate `known_hosts`, and protect configuration permissions. Restrict TCP/9443 at the firewall to the primary Zabbix Server.

```sh
sudo systemctl daemon-reload
sudo systemctl enable --now zbx-ssh-gateway
sudo systemctl status zbx-ssh-gateway
```

## Upgrade

1. Back up `/etc/zbx-ssh-gateway`, including both local YAML files.
2. Compare the repository's example YAML files with the installed local files for newly introduced settings or operation syntax.
3. Do **not** copy either example YAML over an existing local YAML during an upgrade.
4. Build or obtain the new binaries and run the test suite for source builds.
5. Stop the service and replace the binaries while preserving local configuration.
6. Install an updated systemd unit if the release changes it, then run `systemctl daemon-reload`.
7. Start the service and verify `/api/v1/health`.
8. Confirm a test Zabbix HTTP Agent item before broad polling resumes.

Rollback by restoring the prior binaries and configuration backup and restarting the service.
