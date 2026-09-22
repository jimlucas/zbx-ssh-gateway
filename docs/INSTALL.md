# Installation and Upgrade

This guide assumes Debian 12 and an administrative login with `sudo`. Build and test as your normal administrative user. Use `sudo` only for system installation and service administration. The `zbx-ssh-gateway` account is a non-interactive runtime account; **never switch to it for installation or upgrades.**

## Quick installation

### 1. Install prerequisites

**User:** normal administrative user  
**Directory:** home directory

```sh
cd ~
sudo apt update
sudo apt install -y git golang-go make ca-certificates openssl
git --version
go version
make --version
```

The required Go version is defined in `go.mod`. If Debian's packaged Go is older, install a compatible Go release before continuing.

### 2. Obtain the source

Repository: https://github.com/jimlucas/zbx-ssh-gateway

**User:** normal administrative user

```sh
sudo mkdir -p /usr/local/src
sudo chown "$(id -u):$(id -g)" /usr/local/src
cd /usr/local/src
git clone https://github.com/jimlucas/zbx-ssh-gateway.git
cd /usr/local/src/zbx-ssh-gateway
```

The expected source directory is `/usr/local/src/zbx-ssh-gateway`.

Releases, when published, are available at https://github.com/jimlucas/zbx-ssh-gateway/releases. Until a release explicitly supplies supported pre-built binaries, use this source installation procedure.

### 3. Test and build

**User:** normal administrative user  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
go mod download
make check
make build
```

Do not install if `make check` fails.

### 4. Install

**User:** root privileges via `sudo`  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
sudo make install
```

The install target:

- creates the `zbx-ssh-gateway` system user/group if absent;
- creates protected configuration directories;
- builds and installs both binaries;
- installs the systemd unit;
- creates `gateway.local.yaml` and `operations.local.yaml` from their examples **only if those local files do not already exist**;
- runs `systemctl daemon-reload`.

It deliberately does **not** start or enable the daemon. Complete security configuration first.

Installed layout:

```text
/usr/sbin/zbx-ssh-gateway
/usr/bin/zbx-ssh-gatewayctl

/etc/zbx-ssh-gateway/
    gateway.local.yaml
    operations.local.yaml
    known_hosts
    secrets/
        api-token
    tls/
        server.crt
        server.key

/etc/systemd/system/zbx-ssh-gateway.service
```

The source/build tree remains in `/usr/local/src/zbx-ssh-gateway`.

## Configure the installation

### Main configuration

```sh
sudo editor /etc/zbx-ssh-gateway/gateway.local.yaml
```

This contains listener, TLS, API, SSH credentials/password list, resource limits, and the path to the operations file.

### Operations

```sh
sudo editor /etc/zbx-ssh-gateway/operations.local.yaml
```

All custom Zabbix polling operations belong here.

Both `*.local.yaml` files are persistent site configuration. Future `sudo make install` runs preserve them.

### API token

```sh
openssl rand -hex 32 | sudo tee /etc/zbx-ssh-gateway/secrets/api-token >/dev/null
sudo chown root:zbx-ssh-gateway /etc/zbx-ssh-gateway/secrets/api-token
sudo chmod 0640 /etc/zbx-ssh-gateway/secrets/api-token
```

Never commit this token.

### SSH known_hosts

Strict host-key verification uses:

```text
/etc/zbx-ssh-gateway/known_hosts
```

Populate it with host keys verified through a trusted source, then:

```sh
sudo chown root:zbx-ssh-gateway /etc/zbx-ssh-gateway/known_hosts
sudo chmod 0640 /etc/zbx-ssh-gateway/known_hosts
```

Do not disable host-key verification to simplify testing.

### HTTPS certificate

The example configuration expects:

```text
/etc/zbx-ssh-gateway/tls/server.crt
/etc/zbx-ssh-gateway/tls/server.key
```

#### Obtaining a Let's Encrypt certificate with Certbot

The gateway does not need Apache or nginx. For a host with a public DNS name, Certbot's standalone authenticator is a simple way to obtain a certificate. The DNS name must resolve publicly to this server, and inbound TCP/80 must be reachable from the Internet while the ACME HTTP-01 challenge is performed.

Install Certbot using the installation method recommended for your Debian system. Then, before the gateway is started, request a certificate, replacing `gateway.example.com` and the email address with your own values:

```sh
sudo certbot certonly --standalone \
  --domain gateway.example.com \
  --email admin@example.com \
  --agree-tos \
  --no-eff-email
```

The resulting certificate is normally maintained beneath:

```text
/etc/letsencrypt/live/gateway.example.com/
    fullchain.pem
    privkey.pem
```

Do not copy the private key into the source tree. The gateway service runs as `zbx-ssh-gateway`, so certificate deployment must preserve appropriate permissions while allowing the service to read its certificate and private key. One straightforward deployment is to copy the current Certbot-managed files into the gateway's protected TLS directory:

```sh
sudo install -o root -g zbx-ssh-gateway -m 0640 \
  /etc/letsencrypt/live/gateway.example.com/fullchain.pem \
  /etc/zbx-ssh-gateway/tls/server.crt

sudo install -o root -g zbx-ssh-gateway -m 0640 \
  /etc/letsencrypt/live/gateway.example.com/privkey.pem \
  /etc/zbx-ssh-gateway/tls/server.key
```

Certbot renews its managed certificate, not these copied files. Therefore a production deployment using this method must use a Certbot deploy hook to refresh the gateway copies and restart or reload the gateway after successful renewal. Do not rely on the initial copies indefinitely.

Test Certbot's renewal configuration with:

```sh
sudo certbot renew --dry-run
```

If TCP/80 cannot be exposed to the Internet, use a supported Certbot DNS authenticator instead of `--standalone`. DNS validation is also the appropriate approach for wildcard certificates.

For current installation choices, DNS plugins, renewal hooks, and expanded Certbot instructions, see:

- https://certbot.eff.org/instructions
- https://eff-certbot.readthedocs.io/en/stable/using.html

#### Using an existing certificate

If you obtain the certificate by another method, install the certificate and private key at the paths configured in `gateway.local.yaml`, then protect them:

```sh
sudo chown root:zbx-ssh-gateway /etc/zbx-ssh-gateway/tls/server.crt /etc/zbx-ssh-gateway/tls/server.key
sudo chmod 0640 /etc/zbx-ssh-gateway/tls/server.crt /etc/zbx-ssh-gateway/tls/server.key
```

## Start the service

Only after the local configuration, API token, TLS files, and `known_hosts` are ready:

```sh
sudo systemctl enable --now zbx-ssh-gateway
sudo systemctl status zbx-ssh-gateway
```

Logs:

```sh
sudo journalctl -u zbx-ssh-gateway -n 100 --no-pager
sudo journalctl -u zbx-ssh-gateway -f
```

Verify that the daemon runs as the unprivileged service account:

```sh
ps -o user,group,pid,cmd -C zbx-ssh-gateway
```

You should never need to `su` or log in as `zbx-ssh-gateway`.

## Network requirements

Restrict the HTTPS listener (TCP/9443 in the example) to the primary Zabbix Server and explicitly authorized administrative/test hosts. The gateway requires outbound TCP/22 access to the network equipment it polls.

## Make targets

From `/usr/local/src/zbx-ssh-gateway`, the supported targets are:

```text
make build             Build both binaries into ./bin
make test              Run Go tests
make vet               Run go vet
make check             Run tests and vet

sudo make install      Safe complete system installation/update
sudo make install-user Create the service account/group if needed
sudo make install-dirs Create protected configuration directories
sudo make install-binaries
                       Build and install binaries
sudo make install-config
                       Create missing local configs; never overwrite existing ones
sudo make install-systemd
                       Install the systemd unit and daemon-reload

sudo make config-diff  Compare examples with installed local configuration
sudo make uninstall    Remove binaries/unit but preserve local configuration
make clean             Remove ./bin
```

### Important uninstall behavior

`sudo make uninstall` intentionally preserves:

```text
/etc/zbx-ssh-gateway/
/var/lib/zbx-ssh-gateway/
```

This prevents an uninstall/reinstall cycle from deleting credentials, operations, TLS material, host keys, or other local state.

## Upgrade procedure

### 1. Use your normal administrative account

Do not switch to `zbx-ssh-gateway`.

### 2. Back up local configuration

```sh
sudo rm -rf /etc/zbx-ssh-gateway.backup
sudo cp -a /etc/zbx-ssh-gateway /etc/zbx-ssh-gateway.backup
```

### 3. Update the checkout

```sh
cd /usr/local/src/zbx-ssh-gateway
git status
git pull --ff-only
```

If `git status` reports modifications, review them before pulling. Site-specific configuration should not normally be stored in the checkout.

### 4. Review configuration changes

```sh
cd /usr/local/src/zbx-ssh-gateway
sudo make config-diff
```

This compares the current project examples against the installed local configuration. Differences are informational; `config-diff` does not modify anything.

Manually add newly required settings to the local files. Never replace existing local files with the examples during an upgrade.

### 5. Test

```sh
cd /usr/local/src/zbx-ssh-gateway
go mod download
make check
```

### 6. Install the update

```sh
cd /usr/local/src/zbx-ssh-gateway
sudo systemctl stop zbx-ssh-gateway
sudo make install
sudo systemctl start zbx-ssh-gateway
```

`make install` replaces the binaries and project-managed systemd unit while preserving both local YAML files.

### 7. Verify

```sh
sudo systemctl status zbx-ssh-gateway
sudo journalctl -u zbx-ssh-gateway -n 100 --no-pager
```

Confirm a test Zabbix HTTP Agent item before returning the gateway to broad polling.

## Manual installation / troubleshooting

The Makefile performs standard Unix installation operations. If troubleshooting requires doing them manually, the equivalent high-level sequence is:

1. Build both Go commands from the repository root.
2. Create the `zbx-ssh-gateway` system user/group.
3. Create `/etc/zbx-ssh-gateway/{secrets,tls}`.
4. Install the daemon to `/usr/sbin/zbx-ssh-gateway`.
5. Install the control utility to `/usr/bin/zbx-ssh-gatewayctl`.
6. On a new installation only, copy the example YAML files to their corresponding `*.local.yaml` paths.
7. Install the systemd unit under `/etc/systemd/system`.
8. Run `systemctl daemon-reload`.
9. Configure credentials, known hosts, TLS, and network policy before starting the service.

The Makefile itself is the definitive reference for the exact file modes and ownership used by automated installation.

## Rollback

Stop the service, restore the previous binaries if necessary, restore `/etc/zbx-ssh-gateway` from the backup, run `systemctl daemon-reload` if the unit changed, and restart the service.
