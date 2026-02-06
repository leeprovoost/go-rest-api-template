package passport

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
)

// Server holds application dependencies and provides HTTP handlers.
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

// ServerOptions configures the server.
type ServerOptions struct {
	Version     string
	Env         string
	Port        string
	CORSOrigins string
	RateLimit   float64 // requests per second; 0 disables rate limiting
	RateBurst   int     // burst size for rate limiter
}

// NewServer creates a new Server with the given dependencies.
func NewServer(
	userStore models.UserStorage,
	passportStore models.PassportStorage,
	logger *slog.Logger,
	opts ServerOptions,
) *Server {
	var rl *rateLimiter
	if opts.RateLimit > 0 {
		rl = newRateLimiter(opts.RateLimit, opts.RateBurst)
	}
	return &Server{
		userStore:     userStore,
		passportStore: passportStore,
		logger:        logger,
		version:       opts.Version,
		env:           opts.Env,
		port:          opts.Port,
		corsOrigins:   opts.CORSOrigins,
		rateLimiter:   rl,
	}
}

// NewTestServer creates a Server configured for testing.
func NewTestServer() *Server {
	return NewServer(
		NewUserService(CreateMockDataSet()),
		NewPassportService(CreateMockPassportDataSet()),
		slog.Default(),
		ServerOptions{
			Version: "0.0.0",
			Env:     "LOCAL",
			Port:    "3001",
		},
	)
}

// Run starts the HTTP server and blocks until shutdown.
func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.addr(),
		Handler:      s.middleware(s.routes()),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("starting server",
			"version", s.version,
			"env", s.env,
			"addr", srv.Addr,
		)
		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-quit:
		s.logger.Info("shutting down server", "signal", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func (s *Server) addr() string {
	if s.env == "LOCAL" {
		return "localhost:" + s.port
	}
	return ":" + s.port
}

// middleware chains all middleware in order of execution.
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

// responseWriter wraps http.ResponseWriter to capture the status code for logging.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		s.logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", time.Since(start),
			"request_id", w.Header().Get("X-Request-ID"),
		)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

func clacksOverhead(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
		next.ServeHTTP(w, r)
	})
}
