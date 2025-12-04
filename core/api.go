// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clivern/zewi/api"
	"github.com/clivern/zewi/db"
	"github.com/clivern/zewi/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// SetupAPI creates and configures the HTTP server for the API
func SetupAPI(Static embed.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	if viper.GetInt("app.timeout") > 0 {
		timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second
		r.Use(chimiddleware.Timeout(timeout))
	}
	r.Use(middleware.CORS)
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.Logger)

	// Routes
	r.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	r.Get("/_health", api.HealthAction)
	r.Get("/_ready", api.ReadyAction)
	r.With(middleware.BasicAuth(
		viper.GetString("app.metrics.username"),
		viper.GetString("app.metrics.secret"),
	)).Get("/_metrics", promhttp.Handler().ServeHTTP)

	// State endpoints
	r.Get("/api/v1/state", api.GetStateAction)
	r.Put("/api/v1/state", api.UpdateStateAction)

	return r
}

// InitDatabaseAPI initializes the database connection from configuration
func InitDatabaseAPI() error {
	dbConfig := db.Config{
		Driver:          viper.GetString("app.database.driver"),
		Host:            viper.GetString("app.database.host"),
		Port:            viper.GetInt("app.database.port"),
		Username:        viper.GetString("app.database.username"),
		Password:        viper.GetString("app.database.password"),
		Database:        viper.GetString("app.database.name"),
		MaxOpenConns:    viper.GetInt("app.database.max_open_conns"),
		MaxIdleConns:    viper.GetInt("app.database.max_idle_conns"),
		ConnMaxLifetime: viper.GetInt("app.database.conn_max_lifetime"),
		DataSource:      viper.GetString("app.database.datasource"),
	}

	return db.InitDB(dbConfig)
}

// RunAPI starts the HTTP server for the API with graceful shutdown support
func RunAPI(handler http.Handler) error {
	if err := InitDatabaseAPI(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Error().Err(err).Msg("Error closing database connection")
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("app.port"))),
		Handler: handler,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info().
			Int("port", viper.GetInt("app.port")).
			Bool("tls", viper.GetBool("app.tls.status")).
			Msg("Starting HTTP server")

		var err error
		if viper.GetBool("app.tls.status") {
			err = srv.ListenAndServeTLS(
				viper.GetString("app.tls.crt_path"),
				viper.GetString("app.tls.key_path"),
			)
		} else {
			err = srv.ListenAndServe()
		}

		// Ignore ErrServerClosed as it's expected during graceful shutdown
		if err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info().
			Str("signal", sig.String()).
			Msg("Received shutdown signal")

		shutdownTimeout := 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		log.Info().
			Dur("timeout", shutdownTimeout).
			Msg("Gracefully shutting down server")

		// Shutdown with timeout to allow in-flight requests to complete
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		log.Info().Msg("Server shutdown complete")
	}

	return nil
}
