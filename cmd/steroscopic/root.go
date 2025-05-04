package steroscopic

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

const (
	defaultHost       = "0.0.0.0"
	defaultPort       = "8080"
	shutdownTimeout   = 10 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 999 * time.Second
	idleTimeout       = 60 * time.Second
	readHeaderTimeout = 5 * time.Second
)

var (
	defaultParams = despair.Parameters{
		BlockSize:    16,
		MaxDisparity: 32,
	}
)

// NewServer creates a new web-ui server
func NewServer(
	logger *logger.Logger,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
	leftCam, rightCam camera.Camer,
) (http.Handler, error) {
	mux := http.NewServeMux()
	err := AddRoutes(mux, logger, params, leftStream, rightStream, outputStream, leftCam, rightCam)
	if err != nil {
		return nil, err
	}
	slogLogHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"reqwuest",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
			)

			mux.ServeHTTP(w, r)
		})
	var handler http.Handler = slogLogHandler
	return handler, nil
}

// Run is the entry point for the application.
func Run(
	ctx context.Context,
	onStart func(),
) error {
	var wg sync.WaitGroup

	start := time.Now()

	innerCtx, cancel := signal.NotifyContext(
		context.Background(), // Fresh Context
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer cancel()

	logger := logger.NewLogger()

	leftCamera := camera.NewStaticCamera("./testdata/L_00001.png")
	rightCamera := camera.NewStaticCamera("./testdata/R_00001.png")
	leftStreamManager := camera.NewStreamManager(leftCamera)
	rightStreamManager := camera.NewStreamManager(rightCamera)
	outputCamera := camera.NewOutputCamera(
		&logger,
		&defaultParams,
		leftStreamManager,
		rightStreamManager,
	)
	outputStreamManager := camera.NewStreamManager(outputCamera)
	handler, err := NewServer(
		&logger,
		&defaultParams,
		leftStreamManager,
		rightStreamManager,
		outputStreamManager,
		leftCamera,
		rightCamera,
	)
	if err != nil {
		return err
	}

	// Configure server with timeouts
	httpServer := &http.Server{
		Addr: net.JoinHostPort(
			defaultHost,
			defaultPort,
		),
		Handler:           handler,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	serverErrors := make(chan error, 1)

	// Start server
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info(
			"server starting",
			slog.String("address", httpServer.Addr),
			slog.String("setup-time", time.Since(start).String()),
		)
		onStart()
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for either server error or shutdown signal
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-innerCtx.Done():
		// Signal received, initiate graceful shutdown
		slog.Info("shutdown signal received, shutting down server...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), // Fresh Context for shutdown
			shutdownTimeout,
		)
		defer cancel()

		// Attempt graceful shutdown
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("error during server shutdown",
				slog.String("error", err.Error()),
				slog.Duration("timeout", shutdownTimeout),
			)
		}

		// Wait for all goroutines to finish
		slog.Info("waiting for server shutdown to complete")
		wg.Wait()
		slog.Info("server shutdown completed")
		return nil
	case <-ctx.Done():
		// Parent context cancelled
		slog.Info("parent context cancelled, shutting down...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(), // Use a fresh context for shutdown
			shutdownTimeout,
		)
		defer cancel()

		// Attempt graceful shutdown
		err := httpServer.Shutdown(shutdownCtx)
		if err != nil {
			slog.Error("error during server shutdown",
				slog.String("error", err.Error()),
				slog.Duration("timeout", shutdownTimeout),
			)
		}

		// Wait for all goroutines to finish
		wg.Wait()
		return nil
	}
}
