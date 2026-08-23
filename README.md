# windsea

WindSea is a backend operations platform for offshore wind farms. It coordinates turbine campaigns, inspections, maintenance work orders, parts reservations, telemetry alerts, contractor handoffs, audit events and retryable background dispatch.

## Run

```text
GOTOOLCHAIN=local go test ./... -count=1
GOTOOLCHAIN=local go test -race ./... -count=1
GOTOOLCHAIN=local go vet ./...
GOTOOLCHAIN=local go build ./...
DATABASE_URL=file:windsea.db?_pragma=foreign_keys(1)
go run ./cmd/server
```

The default database is SQLite so local development and tests use a real relational database without an external service. Migrations are applied on startup and are safe to re-run. The HTTP API exposes `/healthz`, `/readyz`, campaign, work-order, reservation, alert, telemetry and handoff flows.

