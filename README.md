# go-rest-api-template

A template for building REST web services in Go using the standard library, [x/time](https://pkg.go.dev/golang.org/x/time/rate) for rate limiting, and [testify](https://github.com/stretchr/testify) for testing. Requires Go 1.23+.

## Introduction

### Why?

After writing many REST APIs with Java Dropwizard, Node.js/Express and Go, I wanted to distil my lessons learned into a reusable template for writing REST APIs in Go.

It's mainly for myself. I don't want to keep reinventing the wheel and just want to get the foundation of my REST API "ready to go" so I can focus on the business logic and integration with other systems and data stores.

This is not a framework. It uses Go's excellent standard library for almost everything:

* `net/http` with the enhanced `ServeMux` (Go 1.22+) for routing with method and path parameter support
* `log/slog` for structured logging (Go 1.21+)
* `encoding/json` for JSON serialisation
* Standard middleware pattern using `http.Handler` wrapping
* [stretchr/testify](https://github.com/stretchr/testify) for readable test assertions
* [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate) for per-IP rate limiting

### Previous versions

Previous versions of this template used external packages such as gorilla/mux, negroni, logrus, and unrolled/render. Since Go 1.22, the standard library's `net/http.ServeMux` supports method-based routing and path parameters, making most routing libraries unnecessary. Similarly, `log/slog` (Go 1.21+) provides structured logging out of the box. This version has been modernised to rely on the standard library wherever possible.

### Knowledge of Go

If you're new to programming in Go, I would highly recommend these resources:

* [A Tour of Go](https://go.dev/tour/)
* [Effective Go](https://go.dev/doc/effective_go)
* [Go by Example](https://gobyexample.com/)
* [Learn X in Y Minutes](https://learnxinyminutes.com/docs/go/)
* [50 Shades of Go](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)

A good project to discover Go packages is [awesome-go](https://github.com/avelino/awesome-go).

### Development tools

The choice of editor is personal, but these are the most popular:

* [Visual Studio Code](https://code.visualstudio.com/) with the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
* IntelliJ [GoLand](https://www.jetbrains.com/go/)
* Neovim with [gopls](https://github.com/golang/tools/tree/master/gopls)

### How to run

Using the root Makefile:

```bash
make run
```

Or manually from the `cmd/api-service` directory:

```bash
export ENV=LOCAL
export PORT=3001
export VERSION=VERSION
go run .
```

The app will bind to `localhost:3001` in LOCAL mode. In other environments (DEV, STG, PRD) it binds to `0.0.0.0:<PORT>`, which is appropriate when running behind a load balancer or proxy.

> **Security note:** Don't bind your application to ports below 1024 on a server. Those are privileged ports that require elevated privileges. Run your app on a higher port (e.g. 8080) with a restricted user, and use a reverse proxy (nginx, AWS ALB, etc.) to handle ports 80/443.

### Build and test

```bash
# Build binary
make build

# Run tests
make test

# Run linter (requires golangci-lint)
make lint

# Build Docker image
make docker

# Clean build artefacts
make clean
```

### Docker

Build and run in a container:

```bash
docker build -t go-rest-api-template .
docker run -p 8080:8080 go-rest-api-template
```

The Dockerfile uses a multi-stage build: a `golang:1.23-alpine` builder stage compiles the binary with `CGO_ENABLED=0`, then the final image is based on `alpine:3.20` (a few MB in size).

## Project Structure

```
.
├── api/
│   └── openapi.yaml           # OpenAPI 3.1 specification
├── cmd/
│   └── api-service/
│       ├── main.go             # Application entry point: config, logging, startup
│       ├── Makefile             # Build and run tasks
│       └── VERSION              # Semantic version file
├── internal/
│   └── passport/
│       ├── models/
│       │   ├── user.go          # User struct and UserStorage interface
│       │   └── passport.go      # Passport struct and PassportStorage interface
│       ├── server.go            # Server struct, constructor, middleware, graceful shutdown
│       ├── routes.go            # Route registration (maps URLs to handlers)
│       ├── handlers.go          # HTTP handler implementations
│       ├── handlers_test.go     # Handler integration tests
│       ├── middleware.go        # Request ID, CORS, rate limiting middleware
│       ├── middleware_test.go   # Middleware unit tests
│       ├── server_test.go       # Server configuration tests
│       ├── db_user.go           # In-memory UserStorage implementation
│       ├── db_user_test.go      # User storage unit tests
│       ├── db_passport.go       # In-memory PassportStorage implementation
│       └── db_passport_test.go  # Passport storage unit tests
├── pkg/
│   ├── health/
│   │   └── check.go             # Health check response struct
│   ├── status/
│   │   └── response.go          # Error/validation response struct
│   └── version/
│       ├── parser.go            # VERSION file parser with semver validation
│       └── parser_test.go       # Version parser tests
├── .github/
│   └── workflows/
│       └── ci.yml               # GitHub Actions CI pipeline
├── .golangci.yml                # Linter configuration
├── Dockerfile                   # Multi-stage Docker build
├── Makefile                     # Root build tasks
├── .gitignore
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

**Directory conventions:**
- `cmd/` - Application entry points. Each subdirectory is a separate binary.
- `internal/` - Application code that should not be imported by other projects. The Go compiler enforces this.
- `pkg/` - Library code that could be reused by other projects.
- `api/` - API specification files (OpenAPI/Swagger).

## Architecture

### The Server struct

The central concept is the `Server` struct in `internal/passport/server.go`:

```go
type Server struct {
    userStore     models.UserStorage
    passportStore models.PassportStorage
    logger        *slog.Logger
    version       string
    env           string
    port          string
    corsOrigins   string
    rateLimiter   *rateLimiter
}
```

This holds all application dependencies and provides HTTP handlers as methods. This pattern avoids global variables and makes dependencies explicit. In Go, this is the idiomatic alternative to dependency injection frameworks - you simply pass dependencies through struct fields.

The server is configured through `ServerOptions`:

```go
type ServerOptions struct {
    Version     string
    Env         string
    Port        string
    CORSOrigins string
    RateLimit   float64 // requests per second; 0 disables rate limiting
    RateBurst   int     // burst size for rate limiter
}

srv := passport.NewServer(userStore, passportStore, logger, passport.ServerOptions{
    Version:     version,
    Env:         env,
    Port:        port,
    CORSOrigins: corsOrigins,
    RateLimit:   rateLimit,
    RateBurst:   rateBurst,
})
```

### Routing with net/http (Go 1.22+)

Since Go 1.22, the standard `net/http.ServeMux` supports method-based routing and path parameters:

```go
func (s *Server) routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /healthcheck", s.handleHealthcheck)
    mux.HandleFunc("GET /ready", s.handleReady)

    // Users
    mux.HandleFunc("GET /users", s.handleListUsers)
    mux.HandleFunc("GET /users/{id}", s.handleGetUser)
    mux.HandleFunc("POST /users", s.handleCreateUser)
    mux.HandleFunc("PUT /users/{id}", s.handleUpdateUser)
    mux.HandleFunc("DELETE /users/{id}", s.handleDeleteUser)

    // Passports
    mux.HandleFunc("GET /users/{uid}/passports", s.handleListUserPassports)
    mux.HandleFunc("GET /passports/{id}", s.handleGetPassport)
    mux.HandleFunc("POST /users/{uid}/passports", s.handleCreatePassport)
    mux.HandleFunc("PUT /passports/{id}", s.handleUpdatePassport)
    mux.HandleFunc("DELETE /passports/{id}", s.handleDeletePassport)

    return mux
}
```

Path parameters are accessed via `r.PathValue("id")` in handlers. This eliminates the need for external routers like gorilla/mux.

### Middleware

Middleware uses the standard `http.Handler` wrapping pattern. Each middleware function takes a handler and returns a new handler:

```go
func securityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        next.ServeHTTP(w, r)
    })
}
```

The middleware chain is composed in `server.go`:

```go
func (s *Server) middleware(next http.Handler) http.Handler {
    h := next
    h = clacksOverhead(h)
    h = securityHeaders(h)
    if s.corsOrigins != "" {
        h = cors(s.corsOrigins)(h)
    }
    if s.rateLimiter != nil {
        h = s.rateLimiter.middleware(h)
    }
    h = s.requestLogger(h)
    h = requestID(h)
    return h
}
```

This chains six middleware layers (in order of execution):
1. **Request ID** - reads `X-Request-ID` from the incoming request or generates a UUID using `crypto/rand`
2. **Request logging** - logs method, path, status code, duration and request ID using `slog`
3. **Rate limiting** (optional) - per-IP token bucket rate limiter using `golang.org/x/time/rate`
4. **CORS** (optional) - sets `Access-Control-Allow-*` headers and handles OPTIONS preflight requests
5. **Security headers** - sets `X-Content-Type-Options` and `X-Frame-Options`
6. **Clacks overhead** - adds `X-Clacks-Overhead: GNU Terry Pratchett` (a [Terry Pratchett tribute](http://www.gnuterrypratchett.com/))

### Input validation

Create and update handlers validate the request body and return `422 Unprocessable Entity` with field-level errors:

```go
func validateUser(u models.User) []string {
    var errs []string
    if u.FirstName == "" {
        errs = append(errs, "firstName is required")
    }
    // ...
    return errs
}
```

Response for a failed validation:

```json
{
    "status": "422",
    "message": "validation failed",
    "errors": [
        "firstName is required",
        "lastName is required",
        "dateOfBirth is required",
        "locationOfBirth is required"
    ]
}
```

### Pagination

The `GET /users` endpoint supports offset/limit pagination:

```
GET /users?offset=0&limit=10
```

Defaults: `offset=0`, `limit=25`. Values outside 1–100 reset to the default of 25. The response includes metadata:

```json
{
    "users": [...],
    "count": 10,
    "total": 42,
    "offset": 0,
    "limit": 10
}
```

### Structured logging with slog

Go 1.21 introduced `log/slog`, a structured logging package in the standard library. It outputs text in LOCAL mode and JSON in other environments:

```go
// LOCAL mode (human-readable)
var logger *slog.Logger
if env == "LOCAL" {
    logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
} else {
    logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}
```

Usage in handlers:

```go
s.logger.Error("user not found", "id", uid, "error", err)
s.logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration", elapsed)
```

### Graceful shutdown

The server handles `SIGINT` and `SIGTERM` signals for graceful shutdown, giving in-flight requests up to 30 seconds to complete:

```go
func (s *Server) Run() error {
    srv := &http.Server{
        Addr:         s.addr(),
        Handler:      s.middleware(s.routes()),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // Start server in background
    errCh := make(chan error, 1)
    go func() {
        errCh <- srv.ListenAndServe()
    }()

    // Wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    select {
    case err := <-errCh:
        return err
    case <-quit:
        // Graceful shutdown with timeout
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        return srv.Shutdown(ctx)
    }
}
```

### Environment configuration

The application is configured through environment variables:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `ENV` | Environment name (controls logging format and bind address) | - | `LOCAL`, `DEV`, `STG`, `PRD` |
| `PORT` | Server port | - | `3001` |
| `VERSION` | Path to the VERSION file | - | `VERSION` |
| `CORS_ORIGINS` | Allowed CORS origin (empty disables CORS) | - | `http://localhost:3000` |
| `RATE_LIMIT` | Requests per second per IP (0 disables) | `0` | `10` |
| `RATE_BURST` | Burst size for rate limiter | `0` | `20` |

- **LOCAL**: Text logging at DEBUG level, binds to `localhost:PORT`
- **Other**: JSON logging at INFO level, binds to `:PORT` (all interfaces)

## Data

### Data model

We use a travel Passport for our example. A passport belongs to a user and a user can have one or more passports.

```go
type User struct {
    ID              int       `json:"id"`
    FirstName       string    `json:"firstName"`
    LastName        string    `json:"lastName"`
    DateOfBirth     time.Time `json:"dateOfBirth"`
    LocationOfBirth string    `json:"locationOfBirth"`
}

type Passport struct {
    ID           string    `json:"id"`
    DateOfIssue  time.Time `json:"dateOfIssue"`
    DateOfExpiry time.Time `json:"dateOfExpiry"`
    Authority    string    `json:"authority"`
    UserID       int       `json:"userId"`
}
```

**JSON field naming:** Field names use camelCase (e.g. `firstName`) because the "JS" in JSON stands for JavaScript, where camelCase is the convention.

**Exported vs unexported:** In Go, uppercase field names are exported (public) and lowercase are unexported (private). Fields must be exported for `encoding/json` to marshal them. The `json:"..."` struct tags control the JSON field names.

**Dates:** We use `time.Time` which automatically marshals to/from RFC 3339 format (`2006-01-02T15:04:05Z07:00`). This is an unambiguous international standard (ISO 8601).

### Data access layer

Storage interfaces in `internal/passport/models/` define the contracts for data operations:

```go
type UserStorage interface {
    ListUsers(ctx context.Context) ([]User, error)
    GetUser(ctx context.Context, id int) (User, error)
    AddUser(ctx context.Context, u User) (User, error)
    UpdateUser(ctx context.Context, u User) (User, error)
    DeleteUser(ctx context.Context, id int) error
}

type PassportStorage interface {
    ListPassportsByUser(ctx context.Context, userID int) ([]Passport, error)
    GetPassport(ctx context.Context, id string) (Passport, error)
    AddPassport(ctx context.Context, p Passport) (Passport, error)
    UpdatePassport(ctx context.Context, p Passport) (Passport, error)
    DeletePassport(ctx context.Context, id string) error
}
```

All methods accept `context.Context` as their first parameter, following Go conventions. This allows propagating request cancellation and timeouts to the data layer, which becomes essential when you swap the in-memory store for a real database.

These interfaces allow swapping the implementation. Currently we use in-memory mocks (`UserService` and `PassportService`), but you could implement the same interfaces for PostgreSQL, SQLite, or any other store.

**Compile-time interface check:** To ensure an implementation satisfies its interface, we use this pattern:

```go
var _ models.UserStorage = (*UserService)(nil)
var _ models.PassportStorage = (*PassportService)(nil)
```

This causes a compile error if any interface method is missing. This is important because Go uses implicit interface satisfaction - there's no `implements` keyword like in Java.

### Mock data

The `CreateMockDataSet()` and `CreateMockPassportDataSet()` functions initialise test data:

```go
// Users
list[0] = models.User{
    ID:              0,
    FirstName:       "John",
    LastName:        "Doe",
    DateOfBirth:     dt,  // 1985-12-31T00:00:00Z
    LocationOfBirth: "London",
}

// Passports
list["012345678"] = models.Passport{
    ID:           "012345678",
    DateOfIssue:  doi, // 2020-01-15T00:00:00Z
    DateOfExpiry: doe, // 2030-01-15T00:00:00Z
    Authority:    "HMPO",
    UserID:       0,
}
```

## API

### Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/healthcheck` | `handleHealthcheck` | Health check with app name and version |
| GET | `/ready` | `handleReady` | Readiness probe (returns `{"status":"ok"}`) |
| GET | `/users` | `handleListUsers` | List all users (paginated) |
| GET | `/users/{id}` | `handleGetUser` | Get a single user |
| POST | `/users` | `handleCreateUser` | Create a new user (validates input) |
| PUT | `/users/{id}` | `handleUpdateUser` | Update an existing user (validates input) |
| DELETE | `/users/{id}` | `handleDeleteUser` | Delete a user |
| GET | `/users/{uid}/passports` | `handleListUserPassports` | List passports for a user |
| GET | `/passports/{id}` | `handleGetPassport` | Get a single passport |
| POST | `/users/{uid}/passports` | `handleCreatePassport` | Create a passport for a user (validates input) |
| PUT | `/passports/{id}` | `handleUpdatePassport` | Update a passport (validates input) |
| DELETE | `/passports/{id}` | `handleDeletePassport` | Delete a passport |

The full API is documented in [api/openapi.yaml](api/openapi.yaml) (OpenAPI 3.1).

### Handlers

Handlers are methods on the `Server` struct, giving them access to the stores, logger, and configuration:

```go
func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
    uid, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        respond(w, http.StatusBadRequest, status.Response{
            Status:  strconv.Itoa(http.StatusBadRequest),
            Message: "invalid user id",
        })
        return
    }
    user, err := s.userStore.GetUser(r.Context(), uid)
    if err != nil {
        s.logger.Error("user not found", "id", uid, "error", err)
        respond(w, http.StatusNotFound, status.Response{
            Status:  strconv.Itoa(http.StatusNotFound),
            Message: "can't find user",
        })
        return
    }
    respond(w, http.StatusOK, user)
}
```

**Error handling pattern:** Check for errors immediately and return early. Don't leak internal error details to the client - log the real error server-side and send a sanitised `status.Response` to the client.

**JSON responses:** The `respond` helper sets `Content-Type: application/json` and encodes the response:

```go
func respond(w http.ResponseWriter, code int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if data != nil {
        json.NewEncoder(w).Encode(data)
    }
}
```

### Health check and readiness

Two operational endpoints are provided:

- **`GET /healthcheck`** - returns the app name and version. Useful for monitoring tools and smoke tests in deployment pipelines.
- **`GET /ready`** - returns `{"status":"ok"}` when the service is ready to accept traffic. Use this for Kubernetes readiness probes or load balancer health checks.

`GET /healthcheck`:

```json
{
    "appName": "go-rest-api-template",
    "version": "1.0.0"
}
```

`GET /ready`:

```json
{
    "status": "ok"
}
```

### Example responses

**List users** (`GET /users`):

```json
{
    "count": 2,
    "total": 2,
    "offset": 0,
    "limit": 25,
    "users": [
        {
            "id": 0,
            "firstName": "John",
            "lastName": "Doe",
            "dateOfBirth": "1985-12-31T00:00:00Z",
            "locationOfBirth": "London"
        },
        {
            "id": 1,
            "firstName": "Jane",
            "lastName": "Doe",
            "dateOfBirth": "1992-01-01T00:00:00Z",
            "locationOfBirth": "Milton Keynes"
        }
    ]
}
```

**List user passports** (`GET /users/0/passports`):

```json
{
    "count": 1,
    "passports": [
        {
            "id": "012345678",
            "dateOfIssue": "2020-01-15T00:00:00Z",
            "dateOfExpiry": "2030-01-15T00:00:00Z",
            "authority": "HMPO",
            "userId": 0
        }
    ]
}
```

**Validation error** (`POST /users` with missing required fields):

```json
{
    "status": "422",
    "message": "validation failed",
    "errors": [
        "firstName is required",
        "lastName is required",
        "dateOfBirth is required",
        "locationOfBirth is required"
    ]
}
```

**Error response** (`GET /users/99`):

```json
{
    "status": "404",
    "message": "can't find user"
}
```

## Testing

### Running tests

```bash
# Run all tests
make test

# Or directly
go test ./... -v -cover
```

### Test structure

Tests are organised into four categories across six test files:

**Storage layer tests** (`db_user_test.go`, `db_passport_test.go`) - test CRUD operations on the in-memory stores:

```go
func TestGetUserSuccess(t *testing.T) {
    srv := NewTestServer()
    dt, _ := time.Parse(time.RFC3339, "1985-12-31T00:00:00Z")
    u, err := srv.userStore.GetUser(context.Background(), 0)
    require.NoError(t, err)
    assert.Equal(t, 0, u.ID)
    assert.Equal(t, "John", u.FirstName)
    assert.Equal(t, dt, u.DateOfBirth)
}
```

**Handler tests** (`handlers_test.go`) - test HTTP endpoints end-to-end using `httptest`:

```go
func TestHealthcheckHandler(t *testing.T) {
    handler := newTestHandler()
    r := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
    w := httptest.NewRecorder()
    handler.ServeHTTP(w, r)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "GNU Terry Pratchett", w.Header().Get("X-Clacks-Overhead"))
    assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}
```

**Middleware tests** (`middleware_test.go`) - test request ID generation, CORS headers, and rate limiting behaviour.

**Server and utility tests** (`server_test.go`, `pkg/version/parser_test.go`) - test server configuration and the VERSION file parser.

The `NewTestServer()` helper creates a fully configured server with mock data, making tests self-contained. Handler tests go through the full middleware chain (including request ID, logging, security headers) for realistic integration testing.

**Current coverage:** 72 tests, 93.5% coverage on `internal/passport`, 90.9% on `pkg/version`.

### Manual testing with curl

```bash
# Health check
curl -s http://localhost:3001/healthcheck | jq

# Readiness
curl -s http://localhost:3001/ready | jq

# List all users
curl -s http://localhost:3001/users | jq

# List users with pagination
curl -s "http://localhost:3001/users?offset=0&limit=1" | jq

# Get a specific user
curl -s http://localhost:3001/users/0 | jq

# Create a user
curl -s -X POST http://localhost:3001/users \
  -H "Content-Type: application/json" \
  -d '{"firstName":"Apple","lastName":"Jack","dateOfBirth":"1972-03-07T00:00:00Z","locationOfBirth":"Cambridge"}' | jq

# Update a user
curl -s -X PUT http://localhost:3001/users/0 \
  -H "Content-Type: application/json" \
  -d '{"id":0,"firstName":"John","lastName":"Updated","dateOfBirth":"1985-12-31T00:00:00Z","locationOfBirth":"Manchester"}' | jq

# Delete a user
curl -s -X DELETE http://localhost:3001/users/1 -w "\n%{http_code}\n"

# List passports for a user
curl -s http://localhost:3001/users/0/passports | jq

# Get a passport
curl -s http://localhost:3001/passports/012345678 | jq

# Create a passport
curl -s -X POST http://localhost:3001/users/0/passports \
  -H "Content-Type: application/json" \
  -d '{"id":"111222333","dateOfIssue":"2024-01-01T00:00:00Z","dateOfExpiry":"2034-01-01T00:00:00Z","authority":"HMPO"}' | jq

# Delete a passport
curl -s -X DELETE http://localhost:3001/passports/012345678 -w "\n%{http_code}\n"
```

## CI/CD

### GitHub Actions

The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that runs on every push and pull request to `master`:

- **test** job: runs `go test ./... -v -cover` and `go vet ./...`
- **lint** job: runs [golangci-lint](https://golangci-lint.run/) with the configuration in `.golangci.yml`

Enabled linters: `errcheck`, `govet`, `staticcheck`, `unused`, `gosimple`, `ineffassign`.

## Production deployment

### Binary

Copy the binary and VERSION file to your server:

```bash
# Build
make build

# Deploy
scp bin/api-service cmd/api-service/VERSION yourserver:/opt/go-rest-api-template/
```

Run as a systemd service:

```ini
[Unit]
Description=Go REST API Template
After=network.target

[Service]
Type=simple
User=appuser
Environment=ENV=PRD
Environment=PORT=8080
Environment=VERSION=/opt/go-rest-api-template/VERSION
ExecStart=/opt/go-rest-api-template/api-service
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

The server handles SIGTERM for graceful shutdown, so systemd's `stop` command will allow in-flight requests to complete.

### Docker

```bash
docker build -t go-rest-api-template .
docker run -p 8080:8080 \
  -e ENV=PRD \
  -e PORT=8080 \
  -e CORS_ORIGINS="https://myapp.example.com" \
  -e RATE_LIMIT=10 \
  -e RATE_BURST=20 \
  go-rest-api-template
```

## What changed from the previous version?

The entire codebase was refactored and modernised using [Claude Code](https://docs.anthropic.com/en/docs/claude-code), Anthropic's CLI tool for Claude. The refactoring covered upgrading all dependencies, rewriting every source file to follow current idiomatic Go, expanding the test suite, and rewriting this README.

### Summary of changes

**Go version and dependencies:**
- Upgraded from Go 1.15 to Go 1.23+
- Removed 6 external dependencies (gorilla/mux, urfave/negroni, unrolled/render, unrolled/secure, sirupsen/logrus, palantir/stacktrace)
- Only 2 dependencies remain: stretchr/testify (testing) and golang.org/x/time (rate limiting)
- Removed the `/vendor` directory (modern Go uses the module cache)

**Code modernisation:**
- Replaced gorilla/mux with `net/http.ServeMux` (Go 1.22+ supports `"GET /users/{id}"` routing)
- Replaced logrus with `log/slog` (structured logging in the stdlib since Go 1.21)
- Replaced unrolled/render with `encoding/json`
- Replaced palantir/stacktrace with `fmt.Errorf` and `%w` error wrapping
- Replaced `ioutil.ReadFile` with `os.ReadFile` (ioutil deprecated since Go 1.16)
- Replaced the AppEnv struct + MakeHandler closure with a Server struct and handler methods (cleaner dependency injection)
- Replaced urfave/negroni with stdlib middleware using the standard `http.Handler` wrapping pattern
- Replaced unrolled/secure with a simple custom middleware for security headers
- Added `context.Context` to all storage interface methods

**New features:**
- Full passport CRUD endpoints (previously returned 501 Not Implemented)
- `PassportStorage` interface with in-memory implementation
- Input validation returning 422 with field-level errors
- Pagination on list endpoints (offset/limit)
- Request ID middleware (generates UUID via `crypto/rand`)
- CORS middleware (configurable allowed origins)
- Per-IP rate limiting middleware (token bucket via `golang.org/x/time/rate`)
- `/ready` endpoint for Kubernetes readiness probes
- `ServerOptions` configuration pattern
- Graceful shutdown with SIGINT/SIGTERM signal handling
- HTTP server timeouts (read, write, idle) for production safety
- Request logging middleware with status code, duration, and request ID

**Infrastructure:**
- Root Makefile with run, build, test, lint, docker, clean targets
- Multi-stage Dockerfile (alpine-based, ~10MB image)
- GitHub Actions CI pipeline (test + lint on push/PR)
- golangci-lint configuration
- OpenAPI 3.1 specification

**Testing:**
- 72 tests, 93.5% coverage on `internal/passport`, 90.9% on `pkg/version`
- Handler tests cover: healthcheck, readiness, user CRUD, pagination, validation errors, passport CRUD
- Tests go through the full middleware chain for realistic integration testing
- Added `require.NoError` for clearer test failure messages

### Migration table

| Before | After | Why |
|--------|-------|-----|
| Go 1.15 | Go 1.23+ | Access to modern stdlib features |
| gorilla/mux | `net/http.ServeMux` | Go 1.22+ has method routing and path params |
| urfave/negroni | stdlib middleware | Standard `http.Handler` wrapping is simpler |
| unrolled/render | `encoding/json` | No need for external JSON rendering |
| sirupsen/logrus | `log/slog` | Structured logging in the stdlib since Go 1.21 |
| palantir/stacktrace | `fmt.Errorf` with `%w` | Standard error wrapping is idiomatic |
| unrolled/secure | Custom middleware | Simple headers don't need a library |
| Vendor directory | Go module cache | Modern Go doesn't need vendoring |
| `ioutil.ReadFile` | `os.ReadFile` | `ioutil` deprecated since Go 1.16 |
| AppEnv + MakeHandler | Server struct with methods | Cleaner dependency injection pattern |
| No graceful shutdown | Signal handling + `srv.Shutdown` | Production-ready server lifecycle |
| Passport stubs (501) | Full passport CRUD | Complete domain model implementation |
| No validation | 422 with field-level errors | Proper input validation |
| No pagination | offset/limit pagination | Scalable list endpoints |
| No CI | GitHub Actions | Automated testing and linting |
| No Dockerfile | Multi-stage Dockerfile | Container-ready deployment |

## Useful references

### Go standard library

* [net/http ServeMux routing](https://pkg.go.dev/net/http#ServeMux) - Go 1.22+ enhanced routing
* [log/slog](https://pkg.go.dev/log/slog) - Structured logging
* [encoding/json](https://pkg.go.dev/encoding/json) - JSON encoding/decoding
* [net/http/httptest](https://pkg.go.dev/net/http/httptest) - HTTP testing utilities

### General Go

* [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
* [Effective Go](https://go.dev/doc/effective_go)
* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

### HTTP and REST

* [JSON and Go](https://go.dev/blog/json)
* [Go by Example: HTTP Server](https://gobyexample.com/http-servers)

### Testing

* [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
* [testify](https://github.com/stretchr/testify) - Test assertions

## License

MIT
