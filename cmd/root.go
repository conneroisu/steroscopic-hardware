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
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
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
	// writeTimeout is the maximum duration before timing out writes of the response.
	// Set very high to accommodate long-running stream connections (Like our streams).
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
//  2. Initializes the logger and camera stream managers
//  3. Creates and configures the HTTP server with appropriate timeouts
//  4. Starts the server and monitors for shutdown signals
//  5. Performs graceful shutdown when terminated
//
// Parameters:
//   - ctx: Parent context for controlling the application lifecycle
//   - onStart: Callback function executed after server initialization but
//
// before accepting connections
//
// It returns any unexpected error encountered during server startup or shutdown.
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

	camera.SetLeftCamera(ctx,
		&camera.Camera{
			Camer: camera.NewStaticCamera("./testdata/L_00001.png", camera.LeftCh()),
		})
	camera.SetRightCamera(ctx,
		&camera.Camera{
			Camer: camera.NewStaticCamera("./testdata/R_00001.png", camera.RightCh()),
		})
	camera.SetOutputCamera(ctx,
		camera.NewOutputCamera(despair.DefaultParams()))
	defer func() {
		CloseErr := camera.Left().Close()
		if CloseErr != nil {
			fmt.Println("Failed to close left camera" + CloseErr.Error())
		}
	}()
	defer func() {
		CloseErr := camera.Right().Close()
		if CloseErr != nil {
			fmt.Println("Failed to close right camera" + CloseErr.Error())
		}
	}()
	defer func() {
		CloseErr := camera.Output().Close()
		if CloseErr != nil {
			fmt.Println("Failed to close output camera" + CloseErr.Error())
		}
	}()
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

	select { // Wait for either server error or shutdown signal
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-innerCtx.Done(): // Signal received, initiate graceful shutdown
		slog.Info("shutdown signal received, shutting down server...")
		return gracefulShutdown(
			innerCtx,
			shutdownTimeout,
			&wg,
			httpServer,
		)
	case <-ctx.Done(): // Parent context cancelled
		slog.Info("parent context cancelled, shutting down...")
		return gracefulShutdown(
			ctx,
			shutdownTimeout,
			&wg,
			httpServer,
		)
	}
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
			logger.Info(
				"request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("pattern", r.Pattern),
			)

			mux.ServeHTTP(w, r)
		})
	var handler http.Handler = slogLogHandler
	return handler, nil
}

// gracefulShutdown manages the orderly shutdown of the HTTP server.
//
// It creates a timeout context for the shutdown operation, attempts to close all active
// connections, and waits for all background goroutines to complete before returning.
//
// Parameters:
//   - ctx: Context that may be canceled to abort the shutdown
//   - shutdownTimeout: Maximum duration to wait for connections to close
//   - wg: WaitGroup tracking background goroutines
//   - server: The HTTP server to shut down
//
// Returns an error if the server fails to shut down cleanly within the timeout period.
func gracefulShutdown(
	ctx context.Context,
	shutdownTimeout time.Duration,
	wg *sync.WaitGroup,
	server *http.Server,
) error {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error during server shutdown: %w", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	return nil
}
