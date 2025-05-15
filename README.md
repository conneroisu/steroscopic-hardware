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
- `Vivado/`, `Vitis/` – FPGA/embedded hardware design files
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
