package cmd

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
	"github.com/conneroisu/steroscopic-hardware/pkg/homedir"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

const (
	// defaultHost is the IP address the server binds to (0.0.0.0 = all interfaces).
	defaultHost = "0.0.0.0"

	// defaultPort is the TCP port the server listens on.
	defaultPort = "8080"

	// shutdownTimeout is the maximum time allowed for the server to complete a graceful shutdown.
	shutdownTimeout = 10 * time.Second

	// readTimeout is the maximum duration for reading the entire request.
	readTimeout = 15 * time.Second

	// writeTimeout is the maximum duration before timing out writes of the response
	// Set very high to accommodate long-running stream connections.
	writeTimeout = 999 * time.Second

	// idleTimeout is the maximum amount of time to wait for the next request.
	idleTimeout = 60 * time.Second

	// readHeaderTimeout is the amount of time allowed to read request headers.
	readHeaderTimeout = 5 * time.Second
)

// Run is the entry point for the application that starts the HTTP server and
// manages its lifecycle.
//
// Process:
//  1. Sets up signal handling for graceful shutdown
//  2. Initializes the logger and camera system
//  3. Creates and configures the HTTP server with appropriate timeouts
//  4. Starts the server and monitors for shutdown signals
//  5. Performs graceful shutdown when terminated
func Run(ctx context.Context, onStart func()) error {
	// Use a WaitGroup to track background goroutines
	var wg sync.WaitGroup
	start := time.Now()

	// Create a context with signal handling
	innerCtx, cancel := signal.NotifyContext(
		context.Background(), // Fresh Context
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	defer cancel()

	// Initialize logger
	logger := logger.NewLogger()

	// Initialize camera system
	initCameras(ctx)
	defer func() {
		err := camera.CloseAll()
		if err != nil {
			slog.Error("Failed to close cameras", "err", err)
		}
		err = homedir.SaveFile(logger.Bytes())
		if err != nil {
			slog.Error("Failed to save log file", "err", err)
		}
	}()

	// Create HTTP server
	handler, err := NewServer(
		ctx,
		&logger,
		cancel,
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

	// Channel for server errors
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info(
			"server starting",
			"address", httpServer.Addr,
			"setup_time", time.Since(start).String(),
		)

		// Execute the onStart callback
		onStart()

		// Start the server
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-innerCtx.Done():
		slog.Info("shutdown signal received, shutting down server...")

		return gracefulShutdown(
			innerCtx,
			shutdownTimeout,
			&wg,
			httpServer,
		)

	case <-ctx.Done():
		slog.Info("parent context cancelled, shutting down...")

		return gracefulShutdown(
			ctx,
			shutdownTimeout,
			&wg,
			httpServer,
		)
	}
}

// initCameras initializes the camera system with default cameras.
func initCameras(ctx context.Context) {
	// Initialize left camera with static test image
	leftCam := camera.NewStaticCamera(ctx, "./testdata/L_00001.png", camera.LeftCameraType)
	err := camera.SetCamera(ctx, camera.LeftCameraType, leftCam)
	if err != nil {
		slog.Error("failed to initialize left camera", "error", err)

		return
	}

	// Initialize right camera with static test image
	rightCam := camera.NewStaticCamera(ctx, "./testdata/R_00001.png", camera.RightCameraType)
	err = camera.SetCamera(ctx, camera.RightCameraType, rightCam)
	if err != nil {
		slog.Error("failed to initialize right camera", "error", err)

		return
	}

	// Initialize output camera for depth mapping
	outputCam := camera.NewOutputCamera(ctx)
	err = camera.SetCamera(ctx, camera.OutputCameraType, outputCam)
	if err != nil {
		slog.Error("failed to initialize output camera", "error", err)

		return
	}

	slog.Info("camera system initialized")
}

// gracefulShutdown manages the orderly shutdown of the HTTP server.
//
// It creates a timeout context for the shutdown operation, attempts to close all active
// connections, and waits for all background goroutines to complete before returning.
func gracefulShutdown(
	ctx context.Context,
	shutdownTimeout time.Duration,
	wg *sync.WaitGroup,
	server *http.Server,
) error {
	// Create timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	// Attempt to shut down the server
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		return fmt.Errorf("error during server shutdown: %w", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}

// NewServer creates a new web-ui server with all necessary routes and handlers configured.
//
// It sets up the HTTP server with routes for camera streaming, configuration, and depth map generation.
// The server includes logging middleware that captures request information.
//
// Parameters:
//   - logger: The application logger for recording events and errors
//   - params: Stereoscopic algorithm parameters (block size, max disparity)
//   - leftStream: Stream manager for the left camera
//   - rightStream: Stream manager for the right camera
//   - outputStream: Stream manager for the generated depth map output
//   - cancel: CancelFunc to gracefully shut down the application
//
// Returns an http.Handler and any error encountered during setup.
func NewServer(
	ctx context.Context,
	logger *logger.Logger,
	cancel context.CancelFunc,
) (http.Handler, error) {
	mux := http.NewServeMux()
	err := AddRoutes(
		ctx,
		mux,
		logger,
		cancel,
	)
	if err != nil {
		return nil, err
	}
	slogLogHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// if r.Method == http.MethodGet && r.URL.Path == "/healthcheck" {
			// 	logger.Info(
			// 		"request",
			// 		slog.String("method", r.Method),
			// 		slog.String("url", r.URL.String()),
			// 		slog.String("pattern", r.Pattern),
			// 	)
			// }
			mux.ServeHTTP(w, r)
		})
	var handler http.Handler = slogLogHandler

	return handler, nil
}
