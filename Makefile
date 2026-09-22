SHELL := /bin/sh

BINDIR := bin
PREFIX ?= /usr
SBINDIR ?= $(PREFIX)/sbin
USERBINDIR ?= $(PREFIX)/bin
SYSCONFDIR ?= /etc/zbx-ssh-gateway
SYSTEMD_DIR ?= /etc/systemd/system
SERVICE_USER ?= zbx-ssh-gateway
SERVICE_GROUP ?= zbx-ssh-gateway

.PHONY: all build test vet check install install-user install-dirs install-binaries install-config install-systemd config-diff uninstall clean

all: build

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/zbx-ssh-gateway ./cmd/zbx-ssh-gateway
	go build -o $(BINDIR)/zbx-ssh-gatewayctl ./cmd/zbx-ssh-gatewayctl

test:
	go test ./...

vet:
	go vet ./...

check: test vet

install: install-user install-dirs install-binaries install-config install-systemd
	@echo "Installation complete."
	@echo "Local configuration was preserved if it already existed."
	@echo "Review $(SYSCONFDIR), configure secrets/TLS/known_hosts, then enable and start the service."

install-user:
	@if ! getent group "$(SERVICE_GROUP)" >/dev/null 2>&1; then groupadd --system "$(SERVICE_GROUP)"; fi
	@if ! id -u "$(SERVICE_USER)" >/dev/null 2>&1; then useradd --system --gid "$(SERVICE_GROUP)" --home /var/lib/zbx-ssh-gateway --create-home --shell /usr/sbin/nologin "$(SERVICE_USER)"; fi

install-dirs:
	install -d -o root -g "$(SERVICE_GROUP)" -m 0750 "$(SYSCONFDIR)"
	install -d -o root -g "$(SERVICE_GROUP)" -m 0750 "$(SYSCONFDIR)/secrets"
	install -d -o root -g "$(SERVICE_GROUP)" -m 0750 "$(SYSCONFDIR)/tls"

install-binaries: build
	install -o root -g root -m 0755 "$(BINDIR)/zbx-ssh-gateway" "$(SBINDIR)/zbx-ssh-gateway"
	install -o root -g root -m 0755 "$(BINDIR)/zbx-ssh-gatewayctl" "$(USERBINDIR)/zbx-ssh-gatewayctl"

install-config: install-dirs
	@if [ ! -e "$(SYSCONFDIR)/gateway.local.yaml" ]; then 		install -o root -g "$(SERVICE_GROUP)" -m 0640 configs/gateway.example.yaml "$(SYSCONFDIR)/gateway.local.yaml"; 		echo "Created $(SYSCONFDIR)/gateway.local.yaml"; 	else 		echo "Preserving existing $(SYSCONFDIR)/gateway.local.yaml"; 	fi
	@if [ ! -e "$(SYSCONFDIR)/operations.local.yaml" ]; then 		install -o root -g "$(SERVICE_GROUP)" -m 0640 configs/operations.example.yaml "$(SYSCONFDIR)/operations.local.yaml"; 		echo "Created $(SYSCONFDIR)/operations.local.yaml"; 	else 		echo "Preserving existing $(SYSCONFDIR)/operations.local.yaml"; 	fi

install-systemd:
	install -o root -g root -m 0644 packaging/systemd/zbx-ssh-gateway.service "$(SYSTEMD_DIR)/zbx-ssh-gateway.service"
	systemctl daemon-reload

config-diff:
	@echo "=== gateway.example.yaml vs gateway.local.yaml ==="
	@if [ -f "$(SYSCONFDIR)/gateway.local.yaml" ]; then diff -u configs/gateway.example.yaml "$(SYSCONFDIR)/gateway.local.yaml" || true; else echo "$(SYSCONFDIR)/gateway.local.yaml does not exist"; fi
	@echo
	@echo "=== operations.example.yaml vs operations.local.yaml ==="
	@if [ -f "$(SYSCONFDIR)/operations.local.yaml" ]; then diff -u configs/operations.example.yaml "$(SYSCONFDIR)/operations.local.yaml" || true; else echo "$(SYSCONFDIR)/operations.local.yaml does not exist"; fi

uninstall:
	@echo "Stopping/disabling service if present..."
	-systemctl disable --now zbx-ssh-gateway.service
	rm -f "$(SYSTEMD_DIR)/zbx-ssh-gateway.service"
	rm -f "$(SBINDIR)/zbx-ssh-gateway"
	rm -f "$(USERBINDIR)/zbx-ssh-gatewayctl"
	systemctl daemon-reload
	@echo "Preserved $(SYSCONFDIR) and /var/lib/zbx-ssh-gateway."

clean:
	rm -rf "$(BINDIR)"
