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
		BlockSize:    8,
		MaxDisparity: 64,
	}
)

// NewServer creates a new web-ui server
func NewServer(
	logger *logger.Logger,
	params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
) (http.Handler, error) {
	mux := http.NewServeMux()
	err := AddRoutes(mux, logger, params, leftStream, rightStream, outputStream)
	if err != nil {
		return nil, err
	}
	slogLogHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"request",
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
	leftStreamManager := camera.NewStreamManager(leftCamera, &logger)
	rightStreamManager := camera.NewStreamManager(rightCamera, &logger)
	outputCamera := camera.NewOutputCamera(
		&logger,
		&defaultParams,
		leftStreamManager,
		rightStreamManager,
	)
	outputStreamManager := camera.NewStreamManager(outputCamera, &logger)
	handler, err := NewServer(
		&logger,
		&defaultParams,
		leftStreamManager,
		rightStreamManager,
		outputStreamManager,
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
		return gracefulShutdown(innerCtx, shutdownTimeout, &wg, httpServer)
	case <-ctx.Done(): // Parent context cancelled
		slog.Info("parent context cancelled, shutting down...")
		return gracefulShutdown(ctx, shutdownTimeout, &wg, httpServer)
	}
}

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
