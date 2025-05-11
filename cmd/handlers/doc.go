// Package handlers contains functions for handling API requests.
//
// This Go package `handlers` is part of a stereoscopic hardware system project that manages HTTP requests for a web UI controlling stereo cameras (for 3D vision). It handles the communication between the web interface and the physical camera hardware.
//
// ## Core Components
//
// ### API Handling Structure
//   - `APIFn` is the fundamental type - a function signature that processes HTTP requests and returns errors
//   - `Make()` converts these API functions into standard HTTP handlers, with built-in error handling
//   - `ErrorHandler()` wraps API functions to provide formatted HTML error responses (using color-coded success/failure messages)
//
// ### Key Handlers
//
//  1. **Camera Configuration (`ConfigureCamera`):**
//     - Processes form data for camera setup (port, baud rate, compression)
//     - Configures either left or right camera streams based on parameters
//     - Creates new output streams after successful configuration
//     - Includes validation and error handling for all input parameters
//
//  2. **Parameters Management (`ParametersHandler`):**
//     - Handles changes to disparity map generator parameters
//     - Processes form data for block size and maximum disparity values
//     - Uses mutex locking to ensure thread safety when updating shared parameters
//     - Logs parameter changes for debugging
//
//  3. **Port Discovery (`GetPorts`):**
//     - Enumerates and returns available serial ports as HTML options
//     - Implements retry logic (up to 10 attempts) if ports aren't initially found
//     - Formats port information for direct use in form select elements
//
//  4. **Image Streaming (`StreamHandlerFn`):**
//     - Sets up MJPEG streaming with multipart boundaries
//     - Manages client registration and connection lifecycle
//     - Implements performance optimizations:
//     - Buffer pooling to minimize memory allocation
//     - JPEG quality control and compression
//     - Frame rate limiting (10 FPS)
//     - Connection timeouts (30 minutes)
//     - Efficient image encoding with reusable buffers
//
// ### UI Integration
//
//   - `MorphableHandler()` supports HTMX integration by detecting the presence of HX-Request headers
//   - Serves either full page or partial content based on request type, enabling dynamic UI updates without full page reloads
//
// ## Technical Design Highlights
//
//   - Thread safety with mutex locks for parameter updates
//   - Memory efficiency through object pooling (JPEG options)
//   - Graceful error handling with formatted responses
//   - Efficient image streaming with buffer reuse
//   - Robust port detection with retry mechanisms
//   - Context-aware logging throughout the system
//
// This package serves as the interface layer between the web UI and the underlying stereoscopic hardware, providing both configuration management and real-time image streaming capabilities.
package handlers

//go:generate gomarkdoc -o README.md -e .
