package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	// DefaultStartSeq is the default start marker for image data.
	DefaultStartSeq = []byte{0xff, 0xd8}
	// DefaultEndSeq is the default end marker for image data.
	DefaultEndSeq = []byte{0xff, 0xd9}
	// Default image dimensions
	DefaultImageWidth  = 320
	DefaultImageHeight = 240
)

func main() {
	// Parse command line flags
	leftPort := flag.String("left-port", "/dev/ttyV0", "Left camera virtual port")
	rightPort := flag.String("right-port", "/dev/ttyV2", "Right camera virtual port")
	imagePath := flag.String("image", "", "Path to image file to stream (PNG format)")
	width := flag.Int("width", DefaultImageWidth, "Image width in pixels")
	height := flag.Int("height", DefaultImageHeight, "Image height in pixels")
	interval := flag.Duration("interval", 500*time.Millisecond, "Interval between image sends")
	flag.Parse()

	if *imagePath == "" {
		log.Fatal("Image path is required. Use -image to specify a PNG file.")
	}

	// Load the image
	imgData, err := loadImageData(*imagePath, *width, *height)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	log.Printf("Loaded image data: %d bytes (%dx%d pixels)", len(imgData), *width, *height)
	if len(imgData) != *width**height {
		log.Fatalf("Image data size mismatch: expected %d bytes, got %d bytes", *width**height, len(imgData))
	}

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Create socat processes to set up virtual serial ports
	leftPair := []string{*leftPort, *leftPort + "-sim"}
	rightPair := []string{*rightPort, *rightPort + "-sim"}

	cmd1, err := startSocat(leftPair[0], leftPair[1])
	if err != nil {
		log.Fatalf("Failed to start socat for left camera: %v", err)
	}
	defer func() {
		if err := cmd1.Process.Kill(); err != nil {
			log.Printf("Failed to kill socat process 1: %v", err)
		}
	}()

	cmd2, err := startSocat(rightPair[0], rightPair[1])
	if err != nil {
		log.Fatalf("Failed to start socat for right camera: %v", err)
	}
	defer func() {
		if err := cmd2.Process.Kill(); err != nil {
			log.Printf("Failed to kill socat process 2: %v", err)
		}
	}()

	// Wait for the virtual ports to be ready
	time.Sleep(2 * time.Second)

	// Start camera simulators
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		simulateCameraOnPort(ctx, leftPair[1], imgData, *interval, *width, *height)
	}()

	go func() {
		defer wg.Done()
		simulateCameraOnPort(ctx, rightPair[1], imgData, *interval, *width, *height)
	}()

	log.Printf("Camera simulators started on ports %s and %s", leftPair[1], rightPair[1])
	log.Printf("The application can connect to ports %s and %s", leftPair[0], rightPair[0])
	log.Println("Press Ctrl+C to stop")

	// Wait for all simulators to complete (should only happen when context is cancelled)
	wg.Wait()
	log.Println("Camera simulators stopped")
}

// startSocat starts a socat process to create a pair of virtual serial ports
func startSocat(port1, port2 string) (*exec.Cmd, error) {
	// Remove any existing ports with these names
	os.Remove(port1)
	os.Remove(port2)

	// Start socat to create the virtual serial ports
	cmd := exec.Command("socat", "-d",
		fmt.Sprintf("pty,raw,echo=0,link=%s,mode=0666", port1),
		fmt.Sprintf("pty,raw,echo=0,link=%s,mode=0666", port2))

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start socat: %w", err)
	}

	return cmd, nil
}

// loadImageData loads an image file and converts it to grayscale raw pixel data
func loadImageData(imagePath string, targetWidth, targetHeight int) ([]byte, error) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening image file: %w", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding PNG image: %w", err)
	}

	// Create a grayscale image of the target dimensions
	grayImg := image.NewGray(image.Rect(0, 0, targetWidth, targetHeight))

	// Scale and convert to grayscale
	bounds := img.Bounds()
	srcWidth, srcHeight := bounds.Dx(), bounds.Dy()

	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			// Simple scaling by interpolation
			srcX := x * srcWidth / targetWidth
			srcY := y * srcHeight / targetHeight

			// Convert to grayscale
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(srcX, srcY)))
		}
	}

	// The raw pixel data is directly accessible in the Pix field
	return grayImg.Pix, nil
}

// simulateCameraOnPort simulates a camera on the given port
func simulateCameraOnPort(ctx context.Context, portPath string, imgData []byte, interval time.Duration, width, height int) {
	// Open the virtual port
	port, err := os.OpenFile(portPath, os.O_RDWR, 0666)
	if err != nil {
		log.Printf("Failed to open virtual port %s: %v", portPath, err)
		return
	}
	defer port.Close()

	startSeq := DefaultStartSeq
	endSeq := DefaultEndSeq
	ack := []byte{0x01} // ACK byte

	log.Printf("Camera simulator started on %s", portPath)

	// Buffer for reading commands
	buf := make([]byte, 256)

	// Main simulation loop
	for {
		select {
		case <-ctx.Done():
			log.Printf("Camera simulator on %s stopping", portPath)
			return
		default:
			// Set a read deadline to avoid blocking indefinitely
			port.SetDeadline(time.Now().Add(500 * time.Millisecond))

			// Check if there are any commands to read
			n, err := port.Read(buf)
			if err == nil && n >= 2 {
				// Check for start sequence
				if buf[0] == startSeq[0] && buf[1] == startSeq[1] {
					log.Printf("Start sequence received on %s, sending image", portPath)

					// Send ACK
					if _, err := port.Write(ack); err != nil {
						log.Printf("Error sending ACK on %s: %v", portPath, err)
						continue
					}

					// Send image data in chunks
					chunkSize := 1024
					for i := 0; i < len(imgData); i += chunkSize {
						end := i + chunkSize
						if end > len(imgData) {
							end = len(imgData)
						}

						if _, err := port.Write(imgData[i:end]); err != nil {
							log.Printf("Error writing image data to %s: %v", portPath, err)
							break
						}

						// Brief pause to prevent overwhelming the receiver
						time.Sleep(10 * time.Millisecond)
					}
					log.Printf("Image data sent on %s (%d bytes)", portPath, len(imgData))
				}

				// Handle end sequence
				if buf[0] == endSeq[0] && buf[1] == endSeq[1] {
					log.Printf("End sequence received on %s", portPath)
				}
			}

			// Small sleep to avoid busy waiting
			time.Sleep(50 * time.Millisecond)
		}
	}
}
