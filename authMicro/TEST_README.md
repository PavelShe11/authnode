# Testing Guide

## Running Tests

### Unit Tests (Fast)
```bash
go test ./...
```

### Integration Tests (Requires Docker)
```bash
go test -tags=integration ./...
```

### E2E Tests (Requires Docker)
```bash
go test -tags=e2e ./...
```

### All Tests Together
```bash
go test -tags=integration,e2e ./...
```

### Specific Test Suite
```bash
# Service tests
go test -v ./internal/service/

# Repository tests
go test -tags=integration -v ./internal/infrastructure/adapter/repository/

# Handler tests
go test -tags=e2e -v ./internal/api/rest/handler/
```

### With Coverage
```bash
go test -tags=integration,e2e -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Test Organization

- **Unit Tests:** `internal/service/*_test.go` - Business logic with mocks
- **Integration Tests:** `internal/infrastructure/adapter/repository/*_test.go` - Real PostgreSQL via testcontainers
- **E2E Tests:** `internal/api/rest/handler/*_test.go` - Full stack HTTP tests with real DB

## Test Count

| Layer | Tests | Description |
|-------|-------|-------------|
| Unit (service) | 5 | Register (3) + ConfirmRegistration (2) |
| Integration (repository) | 6 | Save, FindByEmail, Delete, CleanExpired |
| E2E (handler) | 8 | HTTP endpoints with full stack |
| **Total** | **19** | |

## Regenerating Mocks

```bash
cd authMicro
mockery
```

Configuration is in `.mockery.yaml`.

## Test Helpers

Located in `test/helpers/`:
- `NoopLogger` - silent logger for tests
- `SetupSuccessfulValidation` - mock helper
- `SetupSessionFound` / `SetupSessionNotFound` - mock helpers

## Test Fixtures

Located in `test/fixtures/`:
- `NewValidUserData()` - valid registration data

## Troubleshooting

### Testcontainers fails
- Ensure Docker is running: `docker ps`
- Increase timeout in TestMain if needed

### Mocks not found
- Run `mockery` to regenerate
- Check `.mockery.yaml` configuration

### Build tags
- Integration tests require `-tags=integration`
- E2E tests require `-tags=e2e`