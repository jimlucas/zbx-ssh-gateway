# Development

The project is organized into small packages: `internal/api` owns the HTTP surface; `internal/config` parses and validates YAML; `internal/operations` validates parameters and builds allowlisted commands; `internal/gateway` owns queueing/execution; and `internal/sshclient` owns SSH authentication and command execution.

GitHub Actions runs `go test ./...`, `go vet ./...`, and builds both commands on every push and pull request.

Security-sensitive changes should add negative tests. In particular, tests should prove that arbitrary JSON fields, unknown parameters, invalid enum/regex values, and future command-injection vectors are rejected.
