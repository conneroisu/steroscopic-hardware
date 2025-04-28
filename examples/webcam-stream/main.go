// Package main demonstrates how to stream video from a webcam using WebSockets and FFmpeg.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		return true // Allow any origin for this example
	},
}

// ClientManager handles connected websocket clients
type ClientManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
	ctx        context.Context
}

// NewClientManager creates a new manager for websocket clients.
func NewClientManager(ctx context.Context) *ClientManager {
	return &ClientManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		ctx:        ctx,
	}
}

// Start the client manager in a goroutine
func (manager *ClientManager) start() {
	for {
		select {
		case <-manager.ctx.Done():
			// Context cancelled, exit the loop
			log.Println("Client manager shutting down...")

			// Close all client connections
			manager.mutex.Lock()
			for conn := range manager.clients {
				err := conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				)
				if err != nil {
					slog.Error(
						"Error sending close message:",
						slog.String("err", err.Error()),
					)
				}
				conn.Close()
				delete(manager.clients, conn)
			}
			manager.mutex.Unlock()
			return

		case conn := <-manager.register:
			manager.mutex.Lock()
			manager.clients[conn] = true
			manager.mutex.Unlock()
			log.Println("New client connected. Total clients:", len(manager.clients))

		case conn := <-manager.unregister:
			manager.mutex.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				conn.Close()
			}
			manager.mutex.Unlock()
			log.Println("Client disconnected. Total clients:", len(manager.clients))

		case message := <-manager.broadcast:
			manager.mutex.Lock()
			for conn := range manager.clients {
				err := conn.WriteMessage(websocket.BinaryMessage, message)
				if err != nil {
					conn.Close()
					delete(manager.clients, conn)
				}
			}
			manager.mutex.Unlock()
		}
	}
}

// Handle new WebSocket connections
func handleWebSocket(manager *ClientManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to websocket:", err)
			return
		}

		// Register new client
		manager.register <- conn

		// Keep connection open until client disconnects or context is cancelled
		for {
			select {
			case <-manager.ctx.Done():
				// Send close message and close connection
				err := conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(
						websocket.CloseNormalClosure,
						"Server shutting down",
					),
				)
				if err != nil {
					slog.Error(
						"Error sending close message:",
						slog.String("err", err.Error()),
					)
				}
				conn.Close()
				return
			default:
				// Read message (just to detect disconnects)
				_, _, err := conn.ReadMessage()
				if err != nil {
					manager.unregister <- conn
					return
				}
				// Small sleep to prevent tight loop
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// Start streaming from webcam using ffmpeg
func startWebcamStream(ctx context.Context, manager *ClientManager) {
	log.Println("Starting webcam capture...")

	// Use platform detection to get the right command for the current OS
	cmd := GetFFmpegWebcamCommand(ctx)

	// For debugging: print the command
	log.Println("Executing command:", cmd.String())

	// Get pipe to read ffmpeg output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating stdout pipe:", err)
	}

	// Create a stderr pipe and log the output
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("Error creating stderr pipe:", err)
	}
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stderr.Read(buffer)
			if err != nil {
				break
			}
			if n > 0 {
				log.Printf("FFmpeg: %s", buffer[:n])
			}
		}
	}()

	// Use context for command execution
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}

	// Start the ffmpeg process
	if err := cmd.Start(); err != nil {
		log.Fatal("Error starting ffmpeg:", err)
	}

	// Ensure cleanup on context cancellation
	go func() {
		<-ctx.Done()
		log.Println("Stopping webcam capture...")
		err := cmd.Cancel()
		if err != nil {
			slog.Error("Error stopping ffmpeg:", slog.String("err", err.Error()))
		}
	}()

	// Buffer for reading JPEG delimiters
	buffer := make([]byte, 4*1024*1024) // 4MB buffer, adjust as needed
	position := 0

	// Read from ffmpeg output
	for {
		select {
		case <-ctx.Done():
			log.Println("Webcam capture loop terminated by context")
			return
		default:
			// Read a chunk from stdout
			n, err := stdout.Read(buffer[position:])
			if err != nil {
				log.Println("Error reading from ffmpeg:", err)
				return
			}

			position += n

			// Search for JPEG markers (SOI: 0xFF 0xD8, EOI: 0xFF 0xD9)
			start := 0
			for i := range position - 1 {
				// Find JPEG Start Of Image marker
				if buffer[i] == 0xFF && buffer[i+1] == 0xD8 {
					start = i
				}

				// Find JPEG End Of Image marker
				if buffer[i] == 0xFF && buffer[i+1] == 0xD9 && i > start {
					// We have a complete JPEG image
					frame := make([]byte, i-start+2)
					copy(frame, buffer[start:i+2])

					// Broadcast the frame to all clients
					select {
					case manager.broadcast <- frame:
						// Frame sent successfully
					case <-ctx.Done():
						return
					default:
						// Channel would block, skip this frame
						// This happens if there are too many frames and not enough clients
					}

					// Move remaining data to the beginning of buffer
					copy(buffer, buffer[i+2:position])
					position = position - (i + 2)

					// Control frame rate to not overwhelm clients
					select {
					case <-ctx.Done():
						return
					case <-time.After(time.Millisecond * 3): // ~30 FPS
						// Continue after delay
					}
					break
				}
			}

			// If buffer is getting full without finding a complete JPEG,
			// reset it to avoid overflows (this is a safety measure)
			if position > 3*1024*1024 {
				position = 0
				log.Println("Warning: Buffer overflow, resetting")
			}
		}
	}
}

// GetFFmpegWebcamCommand returns the appropriate ffmpeg command for the current platform
func GetFFmpegWebcamCommand(ctx context.Context) *exec.Cmd {
	platform := runtime.GOOS

	switch platform {
	case "linux":
		return getLinuxCommand(ctx)
	case "darwin":
		return getMacOSCommand(ctx)
	case "windows":
		return getWindowsCommand(ctx)
	default:
		// Fallback to Linux command
		fmt.Printf("Warning: Unknown platform %s, using Linux command as fallback\n", platform)
		return getLinuxCommand(ctx)
	}
}

// Linux webcam command using v4l2
func getLinuxCommand(ctx context.Context) *exec.Cmd {
	// Try to find the first video device
	device := "/dev/video0"

	// Check if the default device exists
	if _, err := exec.Command("ls", device).Output(); err != nil {
		// Try to find an alternative video device
		out, err := exec.Command("ls", "/dev").Output()
		if err == nil {
			devices := strings.SplitSeq(string(out), "\n")
			for d := range devices {
				if strings.HasPrefix(d, "video") {
					device = "/dev/" + d
					break
				}
			}
		}
	}

	fmt.Println("Using Linux webcam device:", device)

	return exec.CommandContext(
		ctx,
		"ffmpeg",
		"-f",
		"v4l2",
		"-framerate",
		"15",
		"-video_size",
		"640x480",
		"-i",
		device,
		"-f",
		"image2pipe",
		"-pix_fmt",
		"yuv420p",
		"-vcodec",
		"mjpeg",
		"-q:v",
		"5",
		"-",
	)
}

// macOS webcam command using avfoundation
func getMacOSCommand(ctx context.Context) *exec.Cmd {
	// On macOS, "0" typically refers to the default camera
	// "0:none" means video device 0 with no audio
	device := "0:none"

	fmt.Println("Using macOS webcam device:", device)

	return exec.CommandContext(
		ctx,
		"ffmpeg",
		"-f",
		"avfoundation",
		"-framerate",
		"15",
		"-video_size",
		"640x480",
		"-i",
		device,
		"-f",
		"image2pipe",
		"-pix_fmt",
		"yuv420p",
		"-vcodec",
		"mjpeg",
		"-q:v",
		"5",
		"-",
	)
}

// Windows webcam command using DirectShow
func getWindowsCommand(ctx context.Context) *exec.Cmd {
	// Try to detect available cameras
	device := "video=Integrated Camera" // Default device name

	// Try to list available devices
	out, err := exec.Command(
		"ffmpeg",
		"-list_devices",
		"true",
		"-f",
		"dshow",
		"-i",
		"dummy",
	).CombinedOutput()
	if err == nil {
		// Parse the output to find camera names
		output := string(out)
		lines := strings.Split(output, "\n")

		for i, line := range lines {
			if strings.Contains(line, "DirectShow video devices") && i+1 < len(lines) {
				// Next lines should contain camera names
				for j := i + 1; j < len(lines) && j < i+10; j++ {
					if strings.Contains(lines[j], "Alternative name") {
						continue
					}
					if strings.Contains(lines[j], "\"") {
						// Extract the camera name between quotes
						parts := strings.Split(lines[j], "\"")
						if len(parts) >= 3 {
							device = "video=" + parts[1]
							break
						}
					}
				}
				break
			}
		}
	}

	fmt.Println("Using Windows webcam device:", device)

	return exec.CommandContext(
		ctx,
		"ffmpeg",
		"-f",
		"dshow",
		"-framerate",
		"15",
		"-video_size",
		"640x480",
		"-i",
		device,
		"-f",
		"image2pipe",
		"-pix_fmt",
		"yuv420p",
		"-vcodec",
		"mjpeg",
		"-q:v",
		"5",
		"-",
	)
}

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Handle shutdown signals
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v\n", sig)
		cancel() // Cancel the context

		// Give goroutines a chance to clean up
		time.Sleep(500 * time.Millisecond)

		// Exit if cleanup takes too long
		select {
		case sig = <-sigChan:
			log.Printf("Received second signal: %v, forcing exit\n", sig)
			os.Exit(1)
		case <-time.After(3 * time.Second):
			log.Println("Forcing exit after timeout")
			os.Exit(1)
		}
	}()

	// Create a client manager with the context
	manager := NewClientManager(ctx)

	// Start the client manager
	go manager.start()

	// Start webcam capture in a separate goroutine
	go startWebcamStream(ctx, manager)

	// Create a server
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	// WebSocket endpoint
	mux.HandleFunc("/ws/video", handleWebSocket(manager))

	// Create a server with graceful shutdown
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// Start the server
	go func() {
		fmt.Printf("Server starting at http://localhost:%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Server is shutting down...")

	// Gracefully shut down the HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}
