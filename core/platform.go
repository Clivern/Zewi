// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package core

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clivern/zewi/api"
	"github.com/clivern/zewi/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// SetupPlatform creates and configures the HTTP server for the platform
func SetupPlatform(Static embed.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	if viper.GetInt("app.timeout") > 0 {
		timeout := time.Duration(viper.GetInt("app.timeout")) * time.Second
		r.Use(chimiddleware.Timeout(timeout))
	}
	r.Use(middleware.PrometheusMiddleware)
	r.Use(middleware.Logger)

	// Routes
	r.Get("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	r.Get("/_health", api.HealthAction)
	r.With(middleware.BasicAuth(
		viper.GetString("app.metrics.username"),
		viper.GetString("app.metrics.secret"),
	)).Get("/_metrics", promhttp.Handler().ServeHTTP)

	dist, err := fs.Sub(Static, "web/dist")
	if err != nil {
		panic(fmt.Sprintf(
			"Error while accessing dist files: %s",
			err.Error(),
		))
	}

	// Serve static assets (CSS, JS, images, etc.)
	r.Handle("/assets/*", http.StripPrefix("/", http.FileServer(http.FS(dist))))

	// SPA fallback: serve index.html for all other routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := dist.Open("index.html")
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		// Read the HTML content
		htmlContent, err := io.ReadAll(indexFile)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Get API base URL from environment variable or config
		apiBaseURL := os.Getenv("API_BASE_URL")
		if apiBaseURL == "" {
			apiBaseURL = viper.GetString("app.api_base_url")
		}

		// Inject script tag to set API_BASE_URL before the main script loads
		htmlStr := string(htmlContent)
		if apiBaseURL != "" {
			scriptTag := fmt.Sprintf(`<script>window.API_BASE_URL=%q;</script>`, apiBaseURL)
			// Insert before the closing </head> tag
			headCloseIdx := -1
			for i := 0; i < len(htmlStr)-6; i++ {
				if htmlStr[i:i+7] == "</head>" {
					headCloseIdx = i
					break
				}
			}
			if headCloseIdx != -1 {
				htmlStr = htmlStr[:headCloseIdx] + scriptTag + "\n    " + htmlStr[headCloseIdx:]
			}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlStr))
	})

	return r
}

// RunPlatform starts the HTTP server for the platform with graceful shutdown support
func RunPlatform(handler http.Handler) error {
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
