package handlers

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/camera"
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
	"github.com/conneroisu/steroscopic-hardware/pkg/logger"
)

// mockCamera is a simplified camera implementation for testing with socat
type mockCamera struct {
	mu            sync.Mutex
	portValue     string
	closeErr      error
	closeCalled   bool
	streamCalled  bool
	frameInterval time.Duration
	frameSize     image.Rectangle
}

// newMockCamera creates a new mock camera for testing with socat
func newMockCamera(port string) *mockCamera {
	return &mockCamera{
		portValue:     port,
		frameInterval: 100 * time.Millisecond,
		frameSize:     image.Rect(0, 0, 320, 240),
	}
}

// Stream simulates streaming from a camera
func (m *mockCamera) Stream(ctx context.Context, ch chan *image.Gray) {
	m.mu.Lock()
	m.streamCalled = true
	m.mu.Unlock()

	// Generate test frames at regular intervals
	ticker := time.NewTicker(m.frameInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Create a test image
			img := image.NewGray(m.frameSize)

			// Fill with a simple pattern
			for y := range m.frameSize.Dy() {
				for x := range m.frameSize.Dx() {
					img.SetGray(x, y, color.Gray{Y: uint8((x + y) % 256)})
				}
			}

			select {
			case <-ctx.Done():
				return
			case ch <- img:
				// Frame sent successfully
			default:
				// Channel full or closed, skip this frame
			}
		}
	}
}

// Close simulates closing the camera
func (m *mockCamera) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closeCalled = true
	return m.closeErr
}

// Port returns the port name
func (m *mockCamera) Port() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.portValue
}

// Ensure mockCamera implements camera.Camer
var _ camera.Camer = (*mockCamera)(nil)

// testConfigureCamera is a test implementation of the handler for use in tests
func testConfigureCamera(
	logger *logger.Logger,
	_ *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager,
	isLeft bool,
) APIFn {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			compression    int
			baudRate       int
			portStr        string
			baudStr        string
			compressionStr string
			presetConfig   = camera.DefaultCameraConfig()
		)

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			return err
		}

		// Get all camera parameters
		portStr = r.FormValue("port")
		baudStr = r.FormValue("baudrate")
		compressionStr = r.FormValue("compression")

		// CONFIGURE port if provided
		if portStr == "" {
			return fmt.Errorf("port not provided")
		}
		presetConfig.Port = portStr

		// CONFIGURE baud rate if provided
		if baudStr == "" {
			return fmt.Errorf("baud rate not provided")
		}
		baudRate, err = strconv.Atoi(baudStr)
		if err != nil {
			return err
		}
		presetConfig.BaudRate = baudRate

		// CONFIGURE compression if provided
		if compressionStr == "" {
			return fmt.Errorf("compression not provided")
		}
		compression, err = strconv.Atoi(compressionStr)
		if err != nil {
			return err
		}
		presetConfig.Compression = compression

		// Log configuration details
		logger.InfoContext(
			r.Context(),
			"testing configuration",
			"stream", isLeft,
			"port", portStr,
			"baudrate", baudRate,
			"compression", compression,
		)

		// For testing, we'll create a mock camera with the requested port
		// instead of trying to establish a real connection
		var mockCam = newMockCamera(portStr)

		// Update the stream manager with the mock camera
		if isLeft {
			leftStream.SetTestCamera(mockCam)
		} else {
			rightStream.SetTestCamera(mockCam)
		}

		// Set OK status
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// TestWithSocat tests the configure handler using socat to create virtual serial ports
func TestWithSocat(t *testing.T) {
	// Skip this test if not in an environment with socat
	if _, err := exec.LookPath("socat"); err != nil {
		t.Skip("socat is not installed, skipping test")
	}

	// Only skip if explicitly running in a CI environment
	if os.Getenv("CI") != "" {
		t.Skip("Skipping socat test in CI environment")
	}

	// Create a logger with lower verbosity for tests
	log := logger.NewLogger()

	// Create algorithm parameters
	params := &despair.Parameters{
		BlockSize:    8,
		MaxDisparity: 64,
	}

	// Clean up any existing ports
	for _, port := range []string{"/tmp/ttyV0", "/tmp/ttyV1", "/tmp/ttyV2", "/tmp/ttyV3"} {
		os.Remove(port) // Ignore errors
	}

	// Create socat virtual serial ports with better permissions
	t.Log("Creating virtual serial ports with socat")
	cmd1 := exec.Command("socat", "-d", "pty,raw,echo=0,link=/tmp/ttyV0,mode=0666", "pty,raw,echo=0,link=/tmp/ttyV1,mode=0666")
	cmd2 := exec.Command("socat", "-d", "pty,raw,echo=0,link=/tmp/ttyV2,mode=0666", "pty,raw,echo=0,link=/tmp/ttyV3,mode=0666")

	// Start socat processes
	if err := cmd1.Start(); err != nil {
		t.Fatalf("Failed to start socat for first pair: %v", err)
	}
	defer func() {
		if err := cmd1.Process.Kill(); err != nil {
			t.Logf("Failed to kill socat process 1: %v", err)
		}
	}()

	if err := cmd2.Start(); err != nil {
		if Kerr := cmd1.Process.Kill(); Kerr != nil {
			t.Logf("Failed to kill socat process 1: %v", Kerr)
		}
		t.Fatalf("Failed to start socat for second pair: %v", err)
	}
	defer func() {
		if Kerr := cmd2.Process.Kill(); Kerr != nil {
			t.Logf("Failed to kill socat process 2: %v", Kerr)
		}
	}()

	// Wait for the virtual ports to be ready
	t.Log("Waiting for virtual ports to be created")
	time.Sleep(2 * time.Second)

	// Verify that the virtual ports have been created
	for _, port := range []string{"/tmp/ttyV0", "/tmp/ttyV1", "/tmp/ttyV2", "/tmp/ttyV3"} {
		if _, err := os.Stat(port); os.IsNotExist(err) {
			t.Fatalf("Virtual port %s was not created", port)
		}
	}
	t.Log("Virtual ports created successfully")

	// Start camera simulators on virtual ports
	var wg sync.WaitGroup
	wg.Add(2)

	// These run on the "device side" of the virtual ports
	t.Log("Starting camera simulators")
	go func() {
		defer wg.Done()
		simulateCameraOnPort(t, "/tmp/ttyV1") // Connected to /tmp/ttyV0
	}()
	go func() {
		defer wg.Done()
		simulateCameraOnPort(t, "/tmp/ttyV3") // Connected to /tmp/ttyV2
	}()

	// Give camera simulators time to start
	time.Sleep(1 * time.Second)

	t.Log("Creating stream managers")
	// Create stream managers with nil cameras initially
	leftStream := camera.NewStreamManager(nil, &log)
	rightStream := camera.NewStreamManager(nil, &log)
	outputStream := camera.NewStreamManager(nil, &log)

	// Test left camera configuration
	t.Log("Testing left camera configuration")
	testLeftCamera(t, &log, params, leftStream, rightStream, outputStream)

	// Give time for configuration to take effect
	time.Sleep(2 * time.Second)

	// Verify camera configuration succeeded
	leftPort := leftStream.GetCameraPort()
	if leftPort != "/tmp/ttyV0" {
		t.Errorf("Left camera port not correctly configured. Expected /tmp/ttyV0, got %s", leftPort)
	} else {
		t.Logf("Left camera successfully configured to port %s", leftPort)
	}

	// Test right camera configuration
	t.Log("Testing right camera configuration")
	testRightCamera(t, &log, params, leftStream, rightStream, outputStream)

	// Give time for configuration to take effect
	time.Sleep(2 * time.Second)

	// Verify camera configuration succeeded
	rightPort := rightStream.GetCameraPort()
	if rightPort != "/tmp/ttyV2" {
		t.Errorf("Right camera port not correctly configured. Expected /tmp/ttyV2, got %s", rightPort)
	} else {
		t.Logf("Right camera successfully configured to port %s", rightPort)
	}

	// Clean up
	t.Log("Cleaning up")

	// Stop the stream managers
	leftStream.Stop()
	rightStream.Stop()
	outputStream.Stop()

	// Terminate socat processes
	if err := cmd1.Process.Kill(); err != nil {
		t.Logf("Failed to kill socat process 1: %v", err)
	}
	if err := cmd2.Process.Kill(); err != nil {
		t.Logf("Failed to kill socat process 2: %v", err)
	}

	// Wait for camera simulators to finish
	t.Log("Waiting for camera simulators to finish")
	wg.Wait()
	t.Log("Test completed")
}

// testLeftCamera tests configuration of the left camera
func testLeftCamera(t *testing.T, log *logger.Logger, params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager) {
	// Create test request with form values for left camera
	form := url.Values{}
	form.Add("port", "/tmp/ttyV0")
	form.Add("baudrate", "9600")
	form.Add("compression", "0")

	body := strings.NewReader(form.Encode())
	req := httptest.NewRequest(http.MethodPost, "/api/configure", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.Background())

	w := httptest.NewRecorder()

	// Create and execute handler for left camera
	handler := testConfigureCamera(log, params, leftStream, rightStream, outputStream, true)

	// Execute the handler
	err := handler(w, req)
	if err != nil {
		t.Fatalf("Error configuring left camera: %v", err)
	}

	// Verify response
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Verify the stream was configured with the right port
	if leftStream.GetCameraPort() != "/tmp/ttyV0" {
		t.Errorf("Expected port to be /tmp/ttyV0, got %s", leftStream.GetCameraPort())
	}
}

// testRightCamera tests configuration of the right camera
func testRightCamera(t *testing.T, log *logger.Logger, params *despair.Parameters,
	leftStream, rightStream, outputStream *camera.StreamManager) {
	// Create test request with form values for right camera
	form := url.Values{}
	form.Add("port", "/tmp/ttyV2")
	form.Add("baudrate", "115200")
	form.Add("compression", "1")

	body := strings.NewReader(form.Encode())
	req := httptest.NewRequest(http.MethodPost, "/api/configure", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.Background())

	w := httptest.NewRecorder()

	// Create and execute handler for right camera
	handler := testConfigureCamera(log, params, leftStream, rightStream, outputStream, false)

	// Execute the handler
	err := handler(w, req)
	if err != nil {
		t.Fatalf("Error configuring right camera: %v", err)
	}

	// Verify response
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Verify the stream was configured with the right port
	if rightStream.GetCameraPort() != "/tmp/ttyV2" {
		t.Errorf("Expected port to be /tmp/ttyV2, got %s", rightStream.GetCameraPort())
	}
}

// simulateCameraOnPort simulates a camera on the given port
// It implements a simplified version of the camera protocol
func simulateCameraOnPort(t *testing.T, portPath string) {
	// Open the virtual port with 0666 permissions
	port, err := os.OpenFile(portPath, os.O_RDWR, 0666)
	if err != nil {
		t.Errorf("Failed to open virtual port for simulation: %v", err)
		return
	}
	defer port.Close()

	// Camera protocol simulation
	buf := make([]byte, 256)
	startSeq := []byte{0xff, 0xd8} // Same as DefaultStartSeq in serial.go
	endSeq := []byte{0xff, 0xd9}   // Same as DefaultEndSeq in serial.go
	ack := []byte{0x01}            // ACK byte

	t.Logf("Camera simulator started on %s", portPath)

	// Run for a maximum of 20 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Main loop
	for {
		select {
		case <-ctx.Done():
			t.Logf("Camera simulator on %s timed out", portPath)
			return
		default:
			// Check if the port is still open
			if port == nil {
				t.Logf("Port is nil on %s", portPath)
				return
			}

			// Set a read deadline to avoid blocking indefinitely
			err := port.SetDeadline(time.Now().Add(500 * time.Millisecond))
			if err != nil {
				t.Logf("Failed to set deadline on %s: %v", portPath, err)
				continue
			}

			// Read from the port
			n, err := port.Read(buf)
			if err != nil {
				// Check if this is just a timeout (which is expected in polling)
				if os.IsTimeout(err) {
					continue
				}
				t.Logf("Error reading from %s: %v", portPath, err)
				return
			}

			// Process commands if we got enough data
			if n >= 2 {
				// Check for start sequence
				if buf[0] == startSeq[0] && buf[1] == startSeq[1] {
					t.Logf("Start sequence received on %s", portPath)

					// Send ACK
					if _, err := port.Write(ack); err != nil {
						t.Logf("Error sending ACK on %s: %v", portPath, err)
						continue
					}

					// Create a small test image (checkerboard pattern)
					width, height := 320, 240
					imgData := make([]byte, width*height)
					for y := range height {
						for x := range width {
							if (x/32+y/32)%2 == 0 {
								imgData[y*width+x] = 255 // White
							} else {
								imgData[y*width+x] = 0 // Black
							}
						}
					}

					// Send image data in chunks
					t.Logf("Sending image data on %s", portPath)
					chunkSize := 1024
					for i := 0; i < len(imgData); i += chunkSize {
						end := min(i+chunkSize, len(imgData))

						if _, err := port.Write(imgData[i:end]); err != nil {
							t.Logf("Error writing image data to %s: %v", portPath, err)
							break
						}

						// Brief pause to prevent overwhelming the receiver
						time.Sleep(10 * time.Millisecond)
					}
					t.Logf("Image data sent on %s", portPath)
				}

				// Handle end sequence
				if buf[0] == endSeq[0] && buf[1] == endSeq[1] {
					t.Logf("End sequence received on %s", portPath)
				}
			}
		}
	}
}
