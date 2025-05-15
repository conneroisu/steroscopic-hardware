// Package camera provides interfaces and implementations for different camera types
// (e.g., static, serial, output) and their management in a stereoscopic hardware system.
//
// The camera package defines a common Camera interface and several concrete camera types:
//
//   - StaticCamera: Loads images from files for testing and simulation.
//   - SerialCamera: Communicates with hardware cameras over a serial port.
//   - OutputCamera: Processes stereo images to generate a depth map.
//
// The package also provides a Manager interface and default manager implementation for
// orchestrating multiple cameras and their data channels.
//
// Example usage:
//
//	leftCam := camera.NewStaticCamera(ctx, "./testdata/L_00001.png", camera.LeftCameraType)
//	rightCam := camera.NewSerialCamera(ctx, camera.RightCameraType, "/dev/ttyUSB0", 115200, 0)
//	outputCam := camera.NewOutputCamera(ctx)
//
//	camera.SetCamera(ctx, camera.LeftCameraType, leftCam)
//	camera.SetCamera(ctx, camera.RightCameraType, rightCam)
//	camera.SetCamera(ctx, camera.OutputCameraType, outputCam)
//
// See the documentation for each type for more details.
package camera
