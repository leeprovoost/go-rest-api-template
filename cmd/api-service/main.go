package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	passport "github.com/leeprovoost/go-rest-api-template/internal/passport"
	vparse "github.com/leeprovoost/go-rest-api-template/pkg/version"
)

func main() {
	env := strings.ToUpper(os.Getenv("ENV"))
	port := os.Getenv("PORT")
	versionPath := os.Getenv("VERSION")
	corsOrigins := os.Getenv("CORS_ORIGINS")
	rateLimit, _ := strconv.ParseFloat(os.Getenv("RATE_LIMIT"), 64)
	rateBurst, _ := strconv.Atoi(os.Getenv("RATE_BURST"))

	// Configure structured logging
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
	slog.SetDefault(logger)

	// Read version
	version, err := vparse.ParseVersionFile(versionPath)
	if err != nil {
		logger.Error("can't find VERSION file",
			"env", env,
			"path", versionPath,
			"error", err,
		)
		os.Exit(1)
	}
	logger.Info("loaded VERSION file", "env", env, "version", version)

	// Initialise data storage
	userStore := passport.NewUserService(passport.CreateMockDataSet())
	passportStore := passport.NewPassportService(passport.CreateMockPassportDataSet())

	// Create and run server
	srv := passport.NewServer(userStore, passportStore, logger, passport.ServerOptions{
		Version:     version,
		Env:         env,
		Port:        port,
		CORSOrigins: corsOrigins,
		RateLimit:   rateLimit,
		RateBurst:   rateBurst,
	})
	if err := srv.Run(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
