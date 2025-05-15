# steroscopic-hardware

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)
<img class="badge" tag="github.com/conneroisu/steroscopic-hardware" src="https://goreportcard.com/badge/github.com/conneroisu/conneroh.com" alt="go report card badge">

---

## Overview

**steroscopic-hardware** is an open-source project for real-time stereoscopic depth mapping using Zedboards and a Go-based webserver. It streams synchronized video feeds from two Zedboards, computes a depth map in hardware, and provides a WebUI for visualization and control.

### Key Features

- Real-time stereo video streaming from dual Zedboards
- Hardware-accelerated depth map calculation
- Web-based user interface for live viewing and control
- Prebuilt binaries for easy deployment
- Nix-based reproducible development environment

---

## Project Structure

- `main.go` – Entry point for the Go webserver
- `cmd/` – Command-line and web server components
- `pkg/` – Core Go packages (camera, despair, logger, lzma, web, etc.)
- `static/` – Static assets for the WebUI (JS, CSS, icons)
- `assets/` – Images and UI previews
- `image_capture/`, `image_receive/` – C code for image acquisition/processing
- `testdata/` – Example images for testing

---

## Hardware & Software Requirements

- **Hardware:** 2x Zedboards (or compatible FPGA boards)
- **Software:**
  - [Go](https://go.dev/doc/install) (for server development)
  - [Nix](https://nixos.org/download.html) (optional, for reproducible dev environment)
  - [direnv](https://direnv.net/docs/installation.html) (for environment management)
  - Modern web browser (for WebUI)

---

## Architecture

1. **Zedboards** capture synchronized video streams and send data to the Go webserver.
2. **Go Webserver** receives, processes, and streams the feeds, computes the depth map, and serves the WebUI.
3. **WebUI** displays live video, depth map, and provides controls for users.

![WebUI Preview showing the MVP software interface](assets/WebUI_Preview.png)

---

## Download

Download the latest release [here](https://github.com/conneroisu/steroscopic-hardware/releases)

---

## Usage

Included in the repository is a prebuilt webserver binary. (See the release section)

To run it, simply download the respective binary for your platform and run it.

---

## Development

### Simple

To develop the webserver, you need to have the following installed:

- [Go](https://go.dev/doc/install)

Then, run the following commands (from the root of the repository):

```bash
# Install dependencies
go mod tidy

# Run Code Generation Step
go generate ./...

# Run the webserver
go run main.go
```

This will start the webserver on port 8080.

### Advanced

To develop using the development environment, you need to have [nix](https://nixos.org/download.html) installed.

- Best [Nix](https://docs.determinate.systems/) Installer
- [direnv](https://direnv.net/docs/installation.html)

From the root of the repository, run the following commands:

```bash
direnv allow
```

This will allow direnv to automatically load the environment variables and development dependencies.

---

## Contributing

Contributions, bug reports, and feature requests are welcome! Please open an issue or submit a pull request.

## Contact

For questions or support, open an issue or contact the maintainer via GitHub.



# cmd

```go
import "github.com/conneroisu/steroscopic-hardware/cmd"
```

Package cmd implements the application's web server and HTTP API for stereoscopic image processing.

The cmd package serves as the entry point for the application, providing:

- A web server with a UI for controlling the stereoscopic cameras
- API endpoints for camera configuration and image streaming
- Depth map generation from stereo image pairs
- Graceful shutdown handling

The main packages are:

- Server: HTTP server implementation with proper timeouts
- Routes: API endpoint definitions for camera control and streaming
- Components: Templ\-based UI components for web interface
- Handlers: HTTP handlers for API endpoints

## Index

- [func AddRoutes\(ctx context.Context, mux \*http.ServeMux, logger \*logger.Logger, cancel context.CancelFunc\) error](<#AddRoutes>)
- [func NewServer\(ctx context.Context, logger \*logger.Logger, cancel context.CancelFunc\) \(http.Handler, error\)](<#NewServer>)
- [func Run\(ctx context.Context, onStart func\(\)\) error](<#Run>)


<a name="AddRoutes"></a>
## func [AddRoutes](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/routes.go#L25-L30>)

```go
func AddRoutes(ctx context.Context, mux *http.ServeMux, logger *logger.Logger, cancel context.CancelFunc) error
```

AddRoutes configures all HTTP routes and handlers for the application.

This function registers endpoints for camera control, streaming, and UI components.

<a name="NewServer"></a>
## func [NewServer](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/root.go#L231-L235>)

```go
func NewServer(ctx context.Context, logger *logger.Logger, cancel context.CancelFunc) (http.Handler, error)
```

NewServer creates a new web\-ui server with all necessary routes and handlers configured.

It sets up the HTTP server with routes for camera streaming, configuration, and depth map generation. The server includes logging middleware that captures request information.

Parameters:

- logger: The application logger for recording events and errors
- params: Stereoscopic algorithm parameters \(block size, max disparity\)
- leftStream: Stream manager for the left camera
- rightStream: Stream manager for the right camera
- outputStream: Stream manager for the generated depth map output
- cancel: CancelFunc to gracefully shut down the application

Returns an http.Handler and any error encountered during setup.

<a name="Run"></a>
## func [Run](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/root.go#L54>)

```go
func Run(ctx context.Context, onStart func()) error
```

Run is the entry point for the application that starts the HTTP server and manages its lifecycle.

Process:

1. Sets up signal handling for graceful shutdown
2. Initializes the logger and camera system
3. Creates and configures the HTTP server with appropriate timeouts
4. Starts the server and monitors for shutdown signals
5. Performs graceful shutdown when terminated



  <type>markdown</type>
# handlers

Http handlers for the web ui can be found here.



# handlers

```go
import "github.com/conneroisu/steroscopic-hardware/cmd/handlers"
```

Package handlers contains functions for handling API requests.

This Go package \`handlers\` is part of a stereoscopic hardware system project that manages HTTP requests for a web UI controlling stereo cameras.

It handles the communication between the web interface and the physical camera hardware.

\#\# Core Components

\#\#\# API Handling Structure

- \`APIFn\` is the fundamental type \- a function signature that processes HTTP requests and returns errors
- \`Make\(\)\` converts these API functions into standard HTTP handlers, with built\-in error handling
- \`ErrorHandler\(\)\` wraps API functions to provide formatted HTML error responses \(using color\-coded success/failure messages\)

\#\#\# Key Handlers

1. \*\*Camera Configuration \(\`ConfigureCamera\`\):\*\* \- Processes form data for camera setup \(port, baud rate, compression\) \- Configures either left or right camera streams based on parameters \- Creates new output streams after successful configuration \- Includes validation and error handling for all input parameters

2. \*\*Parameters Management \(\`ParametersHandler\`\):\*\* \- Handles changes to disparity map generator parameters \- Processes form data for block size and maximum disparity values \- Uses mutex locking to ensure thread safety when updating shared parameters \- Logs parameter changes for debugging

3. \*\*Port Discovery \(\`GetPorts\`\):\*\* \- Enumerates and returns available serial ports as HTML options \- Implements retry logic \(up to 10 attempts\) if ports aren't initially found \- Formats port information for direct use in form select elements

4. \*\*Image Streaming \(\`StreamHandlerFn\`\):\*\* \- Sets up MJPEG streaming with multipart boundaries \- Manages client registration and connection lifecycle \- Implements performance optimizations: \- Buffer pooling to minimize memory allocation \- JPEG quality control and compression \- Frame rate limiting \(10 FPS\) \- Connection timeouts \(30 minutes\) \- Efficient image encoding with reusable buffers

\#\#\# UI Integration

- \`MorphableHandler\(\)\` supports HTMX integration by detecting the presence of HX\-Request headers

\#\# Technical Design Highlights

- Thread safety with mutex locks for parameter updates
- Memory efficiency through object pooling \(JPEG options\)
- Graceful error handling with formatted responses
- Efficient image streaming with buffer reuse
- Robust port detection with retry mechanisms
- Context\-aware logging throughout the system

This package serves as the interface layer between the web UI and the underlying stereoscopic hardware, providing both configuration management and real\-time image streaming capabilities.

## Index

- [func HandleLeftStream\(w http.ResponseWriter, r \*http.Request\) error](<#HandleLeftStream>)
- [func HandleOutputStream\(w http.ResponseWriter, r \*http.Request\) error](<#HandleOutputStream>)
- [func HandleRightStream\(w http.ResponseWriter, r \*http.Request\) error](<#HandleRightStream>)
- [func Make\(fn APIFn\) http.HandlerFunc](<#Make>)
- [func MorphableHandler\(wrapper func\(templ.Component\) templ.Component, morph templ.Component\) http.HandlerFunc](<#MorphableHandler>)
- [type APIFn](<#APIFn>)
  - [func ConfigureCamera\(ctx context.Context, typ camera.Type\) APIFn](<#ConfigureCamera>)
  - [func ConfigureMiddleware\(apiFn APIFn\) APIFn](<#ConfigureMiddleware>)
  - [func ErrorHandler\(fn APIFn\) APIFn](<#ErrorHandler>)
  - [func GetPorts\(logger \*logger.Logger\) APIFn](<#GetPorts>)
  - [func HandleCameraStream\(camType camera.Type\) APIFn](<#HandleCameraStream>)
  - [func ParametersHandler\(\) APIFn](<#ParametersHandler>)
  - [func UploadHandler\(appCtx context.Context, typ camera.Type\) APIFn](<#UploadHandler>)
- [type CtxKey](<#CtxKey>)


<a name="HandleLeftStream"></a>
## func [HandleLeftStream](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/stream.go#L40>)

```go
func HandleLeftStream(w http.ResponseWriter, r *http.Request) error
```

HandleLeftStream returns a handler for streaming the left camera.

<a name="HandleOutputStream"></a>
## func [HandleOutputStream](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/stream.go#L50>)

```go
func HandleOutputStream(w http.ResponseWriter, r *http.Request) error
```

HandleOutputStream returns a handler for streaming the output camera.

<a name="HandleRightStream"></a>
## func [HandleRightStream](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/stream.go#L45>)

```go
func HandleRightStream(w http.ResponseWriter, r *http.Request) error
```

HandleRightStream returns a handler for streaming the right camera.

<a name="Make"></a>
## func [Make](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/api.go#L15>)

```go
func Make(fn APIFn) http.HandlerFunc
```

Make returns a function that can be used as an http.HandlerFunc.

<a name="MorphableHandler"></a>
## func [MorphableHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/api.go#L48-L51>)

```go
func MorphableHandler(wrapper func(templ.Component) templ.Component, morph templ.Component) http.HandlerFunc
```

MorphableHandler returns a handler that checks for the presence of the hx\-trigger header and serves either the full or morphed view.

<a name="APIFn"></a>
## type [APIFn](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/api.go#L12>)

APIFn is a function that handles an API request.

```go
type APIFn func(w http.ResponseWriter, r *http.Request) error
```

<a name="ConfigureCamera"></a>
### func [ConfigureCamera](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/configure.go#L77>)

```go
func ConfigureCamera(ctx context.Context, typ camera.Type) APIFn
```

ConfigureCamera handles client requests to configure camera parameters.

<a name="ConfigureMiddleware"></a>
### func [ConfigureMiddleware](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/configure.go#L26>)

```go
func ConfigureMiddleware(apiFn APIFn) APIFn
```

ConfigureMiddleware parses camera configuration from form data.

It adds the configuration to the request context.

This middleware is required for the ConfigureCamera handler.

<a name="ErrorHandler"></a>
### func [ErrorHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/api.go#L63-L65>)

```go
func ErrorHandler(fn APIFn) APIFn
```

ErrorHandler returns a handler that returns an error response.

<a name="GetPorts"></a>
### func [GetPorts](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/ports.go#L14-L16>)

```go
func GetPorts(logger *logger.Logger) APIFn
```

GetPorts handles client requests to configure the camera.

<a name="HandleCameraStream"></a>
### func [HandleCameraStream](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/stream.go#L14>)

```go
func HandleCameraStream(camType camera.Type) APIFn
```

HandleCameraStream is a generic handler for streaming camera images.

<a name="ParametersHandler"></a>
### func [ParametersHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/params.go#L14>)

```go
func ParametersHandler() APIFn
```

ParametersHandler handles client requests to update disparity algorithm parameters.

<a name="UploadHandler"></a>
### func [UploadHandler](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/upload.go#L16>)

```go
func UploadHandler(appCtx context.Context, typ camera.Type) APIFn
```


<a name="CtxKey"></a>
## type [CtxKey](<https://github.com/conneroisu/steroscopic-hardware/blob/main/cmd/handlers/configure.go#L15>)

CtxKey is a type alias for context keys used to store camera configuration.

```go
type CtxKey string
```



  <type>markdown</type>


# despair

```go
import "github.com/conneroisu/steroscopic-hardware/pkg/despair"
```

Package despair provides a Go implementation of a stereoscopic depth mapping algorithm, designed for efficient generation of depth/disparity maps from stereo image pairs.

### Core Functionality

The package implements the Sum of Absolute Differences \(SAD\) algorithm, a common technique in stereoscopic vision that:

1. Takes left and right grayscale images from slightly different viewpoints
2. Compares blocks of pixels to find matching points between images
3. Calculates disparity \(horizontal displacement\) between matching points
4. Generates a grayscale disparity map where pixel brightness represents depth

### Data Structures

```
InputChunk: Represents a portion of the image pair to process
OutputChunk: Contains processed disparity data for a specific region
Parameters: Configuration settings for the algorithm including:
`BlockSize`: Size of pixel blocks for comparison
`MaxDisparity`: Maximum pixel displacement to check
```

### Processing Pipeline

1. \`SetupConcurrentSAD\`: Creates a pipeline with configurable worker count, returning input/output channels

2. \`RunSad\`: Convenience function that orchestrates the entire process: \- Divides images into chunks \- Distributes processing across workers \- Assembles final disparity map

3. \`AssembleDisparityMap\`: Combines processed chunks into a complete disparity map

4. \`sumAbsoluteDifferences\`: Low\-level function that calculates block matching scores

### Image Handling

The package includes efficient image handling utilities:

- PNG Loading/Saving: Optimized functions for loading and saving grayscale PNG images
- Type\-Specific Conversions: Specialized routines for different image formats \(Gray, RGBA, generic\)
- Error Handling: Both standard error\-returning functions and "Must" variants that panic on failure

### Performance Optimizations

- Concurrent Processing: Utilizes Go's concurrency with multiple worker goroutines
- Chunked Processing: Splits images into smaller regions for parallel processing
- Direct Pixel Access: Works with underlying pixel arrays rather than the higher\-level interface
- Type\-Specific Optimizations: Different code paths for different image types
- Early Termination: Breaks comparison loops when perfect matches are found
- Optimized Bounds Checking: Reduces redundant checks in inner loops
- Precomputed Lookup Tables: Uses LUTs for common conversions

Example:

```
```go
// Load stereo image pair
left := despair.MustLoadPNG("left.png")
right := despair.MustLoadPNG("right.png")

// Generate disparity map with block size 9 and max disparity 64
disparityMap := despair.RunSad(left, right, 9, 64)

// Save the result
despair.MustSavePNG("depth_map.png", disparityMap)
```
```

## Index

- [func AssembleDisparityMap\(outputChan \<\-chan OutputChunk, dimensions image.Rectangle, chunks int\) \*image.Gray](<#AssembleDisparityMap>)
- [func RunSad\(left, right \*image.Gray, blockSize, maxDisparity int\) \*image.Gray](<#RunSad>)
- [func SetDefaultParams\(params Parameters\)](<#SetDefaultParams>)
- [func SetupConcurrentSAD\(numWorkers int\) \(chan\<\- InputChunk, \<\-chan OutputChunk\)](<#SetupConcurrentSAD>)
- [func SumAbsoluteDifferences\(left, right \*image.Gray, leftX, leftY, rightX, rightY, blockSize int\) int](<#SumAbsoluteDifferences>)
- [type InputChunk](<#InputChunk>)
- [type OutputChunk](<#OutputChunk>)
- [type Parameters](<#Parameters>)
  - [func DefaultParams\(\) \*Parameters](<#DefaultParams>)


<a name="AssembleDisparityMap"></a>
## func [AssembleDisparityMap](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/sad.go#L172-L176>)

```go
func AssembleDisparityMap(outputChan <-chan OutputChunk, dimensions image.Rectangle, chunks int) *image.Gray
```

AssembleDisparityMap assembles the disparity map from output chunks.

<a name="LoadPNG"></a>
## func [LoadPNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L10>)

```go
```

LoadPNG loads a PNG image and converts it to grayscale with optimizations.

<a name="MustLoadPNG"></a>
## func [MustLoadPNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L44>)

```go
```

MustLoadPNG loads a PNG image and converts it to grayscale with optimizations and panics if an error occurs.

<a name="MustSavePNG"></a>
## func [MustSavePNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L70>)

```go
```


<a name="RunSad"></a>
## func [RunSad](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/sad.go#L119-L122>)

```go
func RunSad(left, right *image.Gray, blockSize, maxDisparity int) *image.Gray
```

RunSad is a convenience function that sets up the pipeline, feeds the images, and assembles the disparity map.

This is not used in the web UI, but is useful for testing.

<a name="SavePNG"></a>
## func [SavePNG](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/png.go#L55>)

```go
```


<a name="SetDefaultParams"></a>
## func [SetDefaultParams](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/params.go#L21>)

```go
func SetDefaultParams(params Parameters)
```

SetDefaultParams sets the default stereoscopic algorithm parameters.

<a name="SetupConcurrentSAD"></a>
## func [SetupConcurrentSAD](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/sad.go#L29-L31>)

```go
func SetupConcurrentSAD(numWorkers int) (chan<- InputChunk, <-chan OutputChunk)
```

SetupConcurrentSAD sets up a concurrent SAD processing pipeline.

It returns an input channel to feed image chunks into and an output channel to receive results from.

If the input channel is closed, the processing pipeline will stop.

<a name="SumAbsoluteDifferences"></a>
## func [SumAbsoluteDifferences](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/sad.go#L205-L208>)

```go
func SumAbsoluteDifferences(left, right *image.Gray, leftX, leftY, rightX, rightY, blockSize int) int
```

SumAbsoluteDifferences calculates SAD directly on image data.

<a name="InputChunk"></a>
## type [InputChunk](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/sad.go#L12-L15>)

InputChunk represents a portion of the image to process.

```go
type InputChunk struct {
    Left, Right *image.Gray
    Region      image.Rectangle
}
```

<a name="OutputChunk"></a>
## type [OutputChunk](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/sad.go#L18-L21>)

OutputChunk represents the processed disparity data for a region.

```go
type OutputChunk struct {
    DisparityData []uint8
    Region        image.Rectangle
}
```

<a name="Parameters"></a>
## type [Parameters](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/params.go#L34-L37>)

Parameters is a struct that holds the parameters for the stereoscopic image processing.

```go
type Parameters struct {
    BlockSize    int `json:"blockSize"`
    MaxDisparity int `json:"maxDisparity"`
}
```

<a name="DefaultParams"></a>
### func [DefaultParams](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/despair/params.go#L28>)

```go
func DefaultParams() *Parameters
```

DefaultParams returns the default stereoscopic algorithm parameters.



  <type>markdown</type>
# logger

Package logger provides a multi faceted logger that



# logger

```go
import "github.com/conneroisu/steroscopic-hardware/pkg/logger"
```

Package logger provides a multi faceted logger that can be used to log to the console and a buffer.

It's intended to be used as a default logger for the application.

Allowing for the logging of console messages both to the console and to the browser.

## Index

- [func NewLogWriter\(w io.Writer\) slog.Handler](<#NewLogWriter>)
- [type LogEntry](<#LogEntry>)
- [type Logger](<#Logger>)
  - [func NewLogger\(\) Logger](<#NewLogger>)
  - [func \(l Logger\) Bytes\(\) \[\]byte](<#Logger.Bytes>)


<a name="NewLogWriter"></a>
## func [NewLogWriter](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/logger/logger.go#L51>)

```go
func NewLogWriter(w io.Writer) slog.Handler
```

NewLogWriter returns a slog.Handler that writes to a buffer.

<a name="LogEntry"></a>
## type [LogEntry](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/logger/logger.go#L43-L48>)

LogEntry represents a structured log entry.

```go
type LogEntry struct {
    Level   slog.Level
    Time    time.Time
    Message string
    Attrs   []slog.Attr
}
```

<a name="Logger"></a>
## type [Logger](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/logger/logger.go#L15-L18>)

Logger is a slog.Logger that sends logs to a channel and also to the console.

```go
type Logger struct {
    *slog.Logger
    // contains filtered or unexported fields
}
```

<a name="NewLogger"></a>
### func [NewLogger](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/logger/logger.go#L26>)

```go
func NewLogger() Logger
```

NewLogger creates a new Logger.

<a name="Logger.Bytes"></a>
### func \(Logger\) [Bytes](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/logger/logger.go#L21>)

```go
func (l Logger) Bytes() []byte
```

Bytes returns the buffered log.



  <type>markdown</type>
# lzma

Package lzma package implements reading and writing of LZMA format compressed data.



# lzma

```go
import "github.com/conneroisu/steroscopic-hardware/pkg/lzma"
```

Package lzma package implements reading and writing of LZMA format compressed data.

Reference implementation is LZMA SDK version 4.65 originally developed by Igor Pavlov, available online at:

```
http://www.7-zip.org/sdk.html
```

Usage examples. Write compressed data to a buffer:

```
var b bytes.Buffer
w := lzma.NewWriter(&b)
w.Write([]byte("hello, world\n"))
w.Close()
```

read that data back:

```
r := lzma.NewReader(&b)
io.Copy(os.Stdout, r)
r.Close()
```

If the data is bigger than you'd like to hold into memory, use pipes. Write compressed data to an io.PipeWriter:

```
pr, pw := io.Pipe()
 go func() {
 	defer pw.Close()
	w := lzma.NewWriter(pw)
	defer w.Close()
	// the bytes.Buffer would be an io.Reader used to read uncompressed data from
	io.Copy(w, bytes.NewBuffer([]byte("hello, world\n")))
 }()
```

and read it back:

```
defer pr.Close()
r := lzma.NewReader(pr)
defer r.Close()
// the os.Stdout would be an io.Writer used to write uncompressed data to
io.Copy(os.Stdout, r)
```


\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-\-
| Offset | Size | Description |
|\-\-\-\-\-\-\-\-|\-\-\-\-\-\-|\-\-\-\-\-\-\-\-\-\-\-\-\-|
| 0 | 1 | Special LZMA properties \(lc,lp, pb in encoded form\) |
| 1 | 4 | Dictionary size \(little endian\) |
| 5 | 8 | Uncompressed size \(little endian\). Size \-1 stands for unknown size |

## Index

- [Constants](<#constants>)
- [func NewReader\(r io.Reader\) io.ReadCloser](<#NewReader>)
- [func NewWriter\(w io.Writer\) \(io.WriteCloser, error\)](<#NewWriter>)
- [func NewWriterLevel\(w io.Writer, level int\) \(io.WriteCloser, error\)](<#NewWriterLevel>)
- [func NewWriterSize\(w io.Writer, size int64\) \(io.WriteCloser, error\)](<#NewWriterSize>)
- [func NewWriterSizeLevel\(w io.Writer, size int64, level int\) \(io.WriteCloser, error\)](<#NewWriterSizeLevel>)
- [type ArgumentValueError](<#ArgumentValueError>)
  - [func \(e \*ArgumentValueError\) Error\(\) string](<#ArgumentValueError.Error>)
- [type HeaderError](<#HeaderError>)
  - [func \(e HeaderError\) Error\(\) string](<#HeaderError.Error>)
- [type NWriteError](<#NWriteError>)
  - [func \(e \*NWriteError\) Error\(\) string](<#NWriteError.Error>)
- [type Reader](<#Reader>)
- [type StreamError](<#StreamError>)
  - [func \(e \*StreamError\) Error\(\) string](<#StreamError.Error>)
- [type Writer](<#Writer>)


## Constants

<a name="BestSpeed"></a>

```go
const (
    // BestSpeed is the fastest compression level.
    BestSpeed = 1
    // BestCompression is the compression level that gives the best compression ratio.
    BestCompression = 9
    // DefaultCompression is the default compression level.
    DefaultCompression = 5
)
```

<a name="NewReader"></a>
## func [NewReader](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/decoder.go#L47>)

```go
func NewReader(r io.Reader) io.ReadCloser
```

NewReader returns a new [io.ReadCloser](<https://pkg.go.dev/io/#ReadCloser>) that can be used to read the uncompressed version of \`r\`

It is the caller's responsibility to call Close on the [io.ReadCloser](<https://pkg.go.dev/io/#ReadCloser>) when finished reading.

<a name="NewWriter"></a>
## func [NewWriter](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/encoder.go#L89>)

```go
func NewWriter(w io.Writer) (io.WriteCloser, error)
```

NewWriter creates a new Writer that compresses data to the given Writer using the default compression level.

Same as NewWriterSizeLevel\(w, \-1, DefaultCompression\).

<a name="NewWriterLevel"></a>
## func [NewWriterLevel](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/encoder.go#L61>)

```go
func NewWriterLevel(w io.Writer, level int) (io.WriteCloser, error)
```

NewWriterLevel creates a new Writer that compresses data to the given Writer using the given level.

Level is any integer value between [lzma.BestSpeed](<#BestSpeed>) and [lzma.BestCompression](<#BestSpeed>).

Same as lzma.NewWriterSizeLevel\(w, \-1, level\).

<a name="NewWriterSize"></a>
## func [NewWriterSize](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/encoder.go#L81>)

```go
func NewWriterSize(w io.Writer, size int64) (io.WriteCloser, error)
```

NewWriterSize creates a new Writer that compresses data to the given Writer using the given size as the uncompressed data size.

If size is unknown, use \-1 instead.

Level is any integer value between [lzma.BestSpeed](<#BestSpeed>) and [lzma.BestCompression](<#BestSpeed>).

Parameter size and the size, [lzma.DefaultCompression](<#BestSpeed>), \(the lzma header\) are written to the passed in writer before any compressed data.

If size is \-1, last bytes are encoded in a different way to mark the end of the stream. The size of the compressed data will increase by 5 or 6 bytes.

Same as NewWriterSizeLevel\(w, size, lzma.DefaultCompression\).

<a name="NewWriterSizeLevel"></a>
## func [NewWriterSizeLevel](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/encoder.go#L40>)

```go
func NewWriterSizeLevel(w io.Writer, size int64, level int) (io.WriteCloser, error)
```

NewWriterSizeLevel writes to the given Writer the compressed version of data written to the returned [io.WriteCloser](<https://pkg.go.dev/io/#WriteCloser>). It is the caller's responsibility to call Close on the [io.WriteCloser](<https://pkg.go.dev/io/#WriteCloser>) when done.

Parameter size is the actual size of uncompressed data that's going to be written to [io.WriteCloser](<https://pkg.go.dev/io/#WriteCloser>). If size is unknown, use \-1 instead.

Parameter level is any integer value between [lzma.BestSpeed](<#BestSpeed>) and [lzma.BestCompression](<#BestSpeed>).

Arguments size and level \(the lzma header\) are written to the writer before any compressed data.

If size is \-1, last bytes are encoded in a different way to mark the end of the stream. The size of the compressed data will increase by 5 or 6 bytes.

The reason for which size is an argument is that, unlike gzip which appends the size and the checksum at the end of the stream, lzma stores the size before any compressed data. Thus, lzma can compute the size while reading data from pipe.

<a name="ArgumentValueError"></a>
## type [ArgumentValueError](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L39-L42>)

An ArgumentValueError reports an error encountered while parsing user provided arguments.

```go
type ArgumentValueError struct {
    // contains filtered or unexported fields
}
```

<a name="ArgumentValueError.Error"></a>
### func \(\*ArgumentValueError\) [Error](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L46>)

```go
func (e *ArgumentValueError) Error() string
```

Error returns the error message and implements the error interface on the ArgumentValueError type.

<a name="HeaderError"></a>
## type [HeaderError](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L6-L8>)

HeaderError is returned when the header is corrupt.

```go
type HeaderError struct {
    // contains filtered or unexported fields
}
```

<a name="HeaderError.Error"></a>
### func \(HeaderError\) [Error](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L12>)

```go
func (e HeaderError) Error() string
```

Error returns the error message and implements the error interface on the HeaderError type.

<a name="NWriteError"></a>
## type [NWriteError](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L28-L30>)

NWriteError is returned when the number of bytes returned by Writer.Write\(\) didn't meet expectances.

```go
type NWriteError struct {
    // contains filtered or unexported fields
}
```

<a name="NWriteError.Error"></a>
### func \(\*NWriteError\) [Error](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L34>)

```go
func (e *NWriteError) Error() string
```

Error returns the error message and implements the error interface on the NWriteError type.

<a name="Reader"></a>
## type [Reader](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/range.go#L24-L27>)

Reader is the actual read interface needed by \[NewDecoder\].

If the passed in io.Reader does not also have ReadByte, the \[NewDecoder\] will introduce its own buffering.

```go
type Reader interface {
    io.Reader
    ReadByte() (c byte, err error)
}
```

<a name="StreamError"></a>
## type [StreamError](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L17-L19>)

StreamError is returned when the stream is corrupt.

```go
type StreamError struct {
    // contains filtered or unexported fields
}
```

<a name="StreamError.Error"></a>
### func \(\*StreamError\) [Error](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/errors.go#L23>)

```go
func (e *StreamError) Error() string
```

Error returns the error message and implements the error interface on the StreamError type.

<a name="Writer"></a>
## type [Writer](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/lzma/range.go#L138-L142>)

Writer is the actual write interface needed by \[NewEncoder\].

If the passed in [io.Writer](<https://pkg.go.dev/io/#Writer>) does not also have WriteByte and Flush, the \[NewEncoder\] function will wrap it into a bufio.Writer.

```go
type Writer interface {
    io.Writer
    Flush() error
    WriteByte(c byte) error
}
```



  <type>markdown</type>


# web

```go
import "github.com/conneroisu/steroscopic-hardware/pkg/web"
```

Package web contains SVG templates and dom targets for the web UI.

SVG templates are used to render SVG icons and text in the web UI. Templates are embedded into the package using the go:embed directive.


## Index

- [Variables](<#variables>)
- [type Target](<#Target>)


## Variables

<a name="TargetLogContainer"></a>

```go
var (
    // TargetLogContainer is a target for the log container.
    // It is used to insert log entries into the DOM.
    TargetLogContainer = Target{
        ID:  "log-container",
        Sel: "#log-container",
    }

    TargetStatusContent = Target{
    }
)
```

<a name="CircleQuestion"></a>CircleQuestion is a template for the SVG circle\-question icon.

```go
var CircleQuestion = templ.Raw(circleQuestion)
```

<a name="CircleX"></a>CircleX is a template for the SVG circle\-x icon.

```go
var CircleX = templ.Raw(circleX)
```


```go
```

<a name="GreenUp"></a>GreenUp is a template for the SVG green\-up icon.

```go
var GreenUp = templ.Raw(greenUp)
```

<a name="LivePageTitle"></a>

```go
var (
    // LivePageTitle is the title of the live page.
    LivePageTitle = "Live Camera System"
)
```

<a name="RedDown"></a>RedDown is a template for the SVG red\-down icon.

```go
var RedDown = templ.Raw(redDown)
```

<a name="RefreshCw"></a>RefreshCw is a template for the SVG refresh\-cw icon.

```go
var RefreshCw = templ.Raw(refreshCw)
```

<a name="SettingsGear"></a>SettingsGear is a template for the SVG settings\-geat icon.

```go
var SettingsGear = templ.Raw(settingsGear)
```

<a name="Target"></a>
## type [Target](<https://github.com/conneroisu/steroscopic-hardware/blob/main/pkg/web/targets.go#L4-L7>)

Target is a struct representing a dom target.

```go
type Target struct {
    ID  string `json:"id"`
    Sel string `json:"sel"`
}
```



