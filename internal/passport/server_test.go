package passport

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddrLocal(t *testing.T) {
	srv := NewTestServer()
	assert.Equal(t, "localhost:3001", srv.addr())
}

func TestAddrNonLocal(t *testing.T) {
	srv := NewServer(
		NewUserService(CreateMockDataSet()),
		NewPassportService(CreateMockPassportDataSet()),
		slog.Default(),
		ServerOptions{
			Env:  "PRD",
			Port: "8080",
		},
	)
	assert.Equal(t, ":8080", srv.addr())
}

func TestNewServerWithRateLimiter(t *testing.T) {
	srv := NewServer(
		NewUserService(CreateMockDataSet()),
		NewPassportService(CreateMockPassportDataSet()),
		slog.Default(),
		ServerOptions{
			Env:       "LOCAL",
			Port:      "3001",
			RateLimit: 10,
			RateBurst: 20,
		},
	)
	assert.NotNil(t, srv.rateLimiter)
}

func TestNewServerWithoutRateLimiter(t *testing.T) {
	srv := NewTestServer()
	assert.Nil(t, srv.rateLimiter)
}

func TestMiddlewareWithCORSAndRateLimiter(t *testing.T) {
	srv := NewServer(
		NewUserService(CreateMockDataSet()),
		NewPassportService(CreateMockPassportDataSet()),
		slog.Default(),
		ServerOptions{
			Env:         "LOCAL",
			Port:        "3001",
			CORSOrigins: "http://localhost:3000",
			RateLimit:   100,
			RateBurst:   100,
		},
	)
	handler := srv.middleware(srv.routes())

	r := httptest.NewRequest(http.MethodGet, "/ready", nil)
	r.RemoteAddr = "127.0.0.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
}
