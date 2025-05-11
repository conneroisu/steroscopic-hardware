// Package cmd implements the application's web server and HTTP API for stereoscopic image processing.
//
// The cmd package serves as the entry point for the application, providing:
//   - A web server with a UI for controlling the stereoscopic cameras
//   - API endpoints for camera configuration and image streaming
//   - Depth map generation from stereo image pairs
//   - Graceful shutdown handling
//
// The main packages are:
//   - Server: HTTP server implementation with proper timeouts
//   - Routes: API endpoint definitions for camera control and streaming
//   - Components: Templ-based UI components for web interface
//   - Handlers: HTTP handlers for API endpoints
package cmd

//go:generate gomarkdoc -o README.md -e .
