# Installation and Upgrade

This guide assumes Debian 12 and an administrative login that can use `sudo`. Commands show which user should run them and which working directory should be used.

The daemon itself runs as the unprivileged `zbx-ssh-gateway` system account. **Do not log in or build the software as that account.** Source retrieval and compilation are performed as your normal administrative user; installation and system configuration use `sudo`.

## 1. Install prerequisites

**User:** normal administrative user  
**Directory:** your home directory

```sh
cd ~
sudo apt update
sudo apt install -y git golang-go ca-certificates openssl
```

Verify the tools:

```sh
git --version
go version
```

The project's Go version requirement is defined in `go.mod`. If Debian's packaged Go version is older than that requirement, install a compatible Go release before building.

## 2. Obtain the source

The canonical repository is:

https://github.com/jimlucas/zbx-ssh-gateway

For a source installation, keep the checkout under `/usr/local/src/zbx-ssh-gateway`. This directory contains source code only; the running daemon does not execute from this directory.

**User:** normal administrative user

First create the source directory and give your administrative user ownership of the project checkout:

```sh
sudo mkdir -p /usr/local/src
sudo chown "$(id -u):$(id -g)" /usr/local/src
cd /usr/local/src
```

Clone the repository:

```sh
git clone https://github.com/jimlucas/zbx-ssh-gateway.git
cd /usr/local/src/zbx-ssh-gateway
```

Confirm that you are in the correct directory:

```sh
pwd
git status
```

The expected working directory is:

```text
/usr/local/src/zbx-ssh-gateway
```

### Installing from a release

When packaged releases are available, they can be downloaded from:

https://github.com/jimlucas/zbx-ssh-gateway/releases

A release containing pre-built Debian/Linux binaries may be installed without cloning the repository. Follow the release notes for that version. Until a release explicitly supplies supported pre-built binaries, use the source-build procedure in this document.

## 3. Run the test suite

**User:** normal administrative user  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
go mod download
go test ./...
go vet ./...
```

Do not proceed with installation if the tests or vet checks fail.

## 4. Build the binaries

**User:** normal administrative user  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
mkdir -p bin
go build -o bin/zbx-ssh-gateway ./cmd/zbx-ssh-gateway
go build -o bin/zbx-ssh-gatewayctl ./cmd/zbx-ssh-gatewayctl
```

The resulting files are:

```text
/usr/local/src/zbx-ssh-gateway/bin/zbx-ssh-gateway
/usr/local/src/zbx-ssh-gateway/bin/zbx-ssh-gatewayctl
```

## 5. Create the service account

The service account is deliberately non-interactive. You do **not** switch to this user during installation.

**User:** normal administrative user using `sudo`

```sh
sudo useradd --system \
  --home /var/lib/zbx-ssh-gateway \
  --create-home \
  --shell /usr/sbin/nologin \
  zbx-ssh-gateway
```

If the account already exists, `useradd` will report that fact; do not recreate it.

## 6. Create runtime configuration directories

**User:** normal administrative user using `sudo`

```sh
sudo install -d -o root -g zbx-ssh-gateway -m 0750 /etc/zbx-ssh-gateway
sudo install -d -o root -g zbx-ssh-gateway -m 0750 /etc/zbx-ssh-gateway/secrets
sudo install -d -o root -g zbx-ssh-gateway -m 0750 /etc/zbx-ssh-gateway/tls
```

The intended layout is:

```text
/etc/zbx-ssh-gateway/
    gateway.local.yaml
    operations.local.yaml
    known_hosts
    secrets/
        api-token
    tls/
        server.crt
        server.key
```

## 7. Install the binaries

**User:** normal administrative user using `sudo`  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
sudo install -o root -g root -m 0755 bin/zbx-ssh-gateway /usr/sbin/zbx-ssh-gateway
sudo install -o root -g root -m 0755 bin/zbx-ssh-gatewayctl /usr/bin/zbx-ssh-gatewayctl
```

The source checkout remains under `/usr/local/src`; only the built binaries are copied into the system executable directories.

## 8. Create the local configuration

**User:** normal administrative user using `sudo`  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

For a **new installation only**, copy the example files:

```sh
cd /usr/local/src/zbx-ssh-gateway
sudo install -o root -g zbx-ssh-gateway -m 0640 configs/gateway.example.yaml /etc/zbx-ssh-gateway/gateway.local.yaml
sudo install -o root -g zbx-ssh-gateway -m 0640 configs/operations.example.yaml /etc/zbx-ssh-gateway/operations.local.yaml
```

Do not run those two commands over an existing installation. The `*.local.yaml` files belong to the local system and are deliberately excluded from Git.

Edit the main configuration:

```sh
sudo editor /etc/zbx-ssh-gateway/gateway.local.yaml
```

Edit the operation allowlist:

```sh
sudo editor /etc/zbx-ssh-gateway/operations.local.yaml
```

All custom polling commands belong in `operations.local.yaml`.

## 9. Configure the API token

Generate a random bearer token:

```sh
openssl rand -hex 32 | sudo tee /etc/zbx-ssh-gateway/secrets/api-token >/dev/null
sudo chown root:zbx-ssh-gateway /etc/zbx-ssh-gateway/secrets/api-token
sudo chmod 0640 /etc/zbx-ssh-gateway/secrets/api-token
```

Do not place this token in the Git repository.

## 10. Configure SSH host keys

The gateway uses strict SSH host-key verification. Create and maintain:

```text
/etc/zbx-ssh-gateway/known_hosts
```

Populate it only with host keys you have verified through a trusted source. Do not disable host-key verification simply to make initial testing easier.

After creating it:

```sh
sudo chown root:zbx-ssh-gateway /etc/zbx-ssh-gateway/known_hosts
sudo chmod 0640 /etc/zbx-ssh-gateway/known_hosts
```

## 11. Configure HTTPS certificates

Place the server certificate and private key at the locations configured in `gateway.local.yaml`. The example configuration uses:

```text
/etc/zbx-ssh-gateway/tls/server.crt
/etc/zbx-ssh-gateway/tls/server.key
```

Protect the private key:

```sh
sudo chown root:zbx-ssh-gateway /etc/zbx-ssh-gateway/tls/server.crt /etc/zbx-ssh-gateway/tls/server.key
sudo chmod 0640 /etc/zbx-ssh-gateway/tls/server.crt /etc/zbx-ssh-gateway/tls/server.key
```

## 12. Install the systemd service

**User:** normal administrative user using `sudo`  
**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
sudo install -o root -g root -m 0644 packaging/systemd/zbx-ssh-gateway.service /etc/systemd/system/zbx-ssh-gateway.service
sudo systemctl daemon-reload
sudo systemctl enable zbx-ssh-gateway
```

Do not start the service until the local YAML files, bearer token, TLS files, and `known_hosts` are ready.

## 13. Start and verify the daemon

```sh
sudo systemctl start zbx-ssh-gateway
sudo systemctl status zbx-ssh-gateway
```

View daemon logs with:

```sh
sudo journalctl -u zbx-ssh-gateway -n 100 --no-pager
```

Follow logs interactively with:

```sh
sudo journalctl -u zbx-ssh-gateway -f
```

The service process should run as `zbx-ssh-gateway`:

```sh
ps -o user,group,pid,cmd -C zbx-ssh-gateway
```

You should never need to `su` or log in as `zbx-ssh-gateway`.

## 14. Network access

Restrict the configured HTTPS listener (9443 in the example configuration) so it is reachable only from the primary Zabbix Server and any explicitly authorized administrative/test hosts.

The gateway itself also requires outbound TCP/22 access to the network equipment it polls.

## Upgrade procedure

### 1. Become your normal administrative user

Do not perform upgrades as the `zbx-ssh-gateway` service account.

### 2. Back up local configuration

```sh
sudo cp -a /etc/zbx-ssh-gateway /etc/zbx-ssh-gateway.backup
```

### 3. Update the source checkout

**Directory:** `/usr/local/src/zbx-ssh-gateway`

```sh
cd /usr/local/src/zbx-ssh-gateway
git status
git pull --ff-only
```

If `git status` reports local modifications, review them before pulling. Site configuration should not normally exist inside this checkout.

### 4. Review configuration changes

Compare:

```text
configs/gateway.example.yaml   -> /etc/zbx-ssh-gateway/gateway.local.yaml
configs/operations.example.yaml -> /etc/zbx-ssh-gateway/operations.local.yaml
```

Manually incorporate any newly required settings. **Never copy the example files over existing local files during an upgrade.**

### 5. Test and rebuild

```sh
cd /usr/local/src/zbx-ssh-gateway
go mod download
go test ./...
go vet ./...
mkdir -p bin
go build -o bin/zbx-ssh-gateway ./cmd/zbx-ssh-gateway
go build -o bin/zbx-ssh-gatewayctl ./cmd/zbx-ssh-gatewayctl
```

### 6. Stop, replace, and restart

```sh
sudo systemctl stop zbx-ssh-gateway
sudo install -o root -g root -m 0755 bin/zbx-ssh-gateway /usr/sbin/zbx-ssh-gateway
sudo install -o root -g root -m 0755 bin/zbx-ssh-gatewayctl /usr/bin/zbx-ssh-gatewayctl
sudo install -o root -g root -m 0644 packaging/systemd/zbx-ssh-gateway.service /etc/systemd/system/zbx-ssh-gateway.service
sudo systemctl daemon-reload
sudo systemctl start zbx-ssh-gateway
```

### 7. Verify

```sh
sudo systemctl status zbx-ssh-gateway
sudo journalctl -u zbx-ssh-gateway -n 100 --no-pager
```

Confirm a test Zabbix HTTP Agent item before returning the gateway to broad polling.

## Rollback

Stop the service, restore the previous binaries if necessary, restore the configuration backup, run `systemctl daemon-reload` if the unit file changed, and restart the service.
