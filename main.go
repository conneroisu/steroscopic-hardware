// Package main is the main entry point for the stereoscopic disparity algorithm web app.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic"
)

func main() {
	slog.SetDefault(DefaultLogger)
	err := steroscopic.Run(context.Background(), openBrowser)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func openBrowser() {
	var err error
	url := "http://localhost:8080"

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// DefaultLogger is a default logger.
var DefaultLogger = slog.New(
	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			if a.Key == "level" {
				return slog.Attr{}
			}
			if a.Key == slog.SourceKey {
				str := a.Value.String()
				split := strings.Split(str, "/")
				if len(split) > 2 {
					a.Value = slog.StringValue(
						strings.Join(split[len(split)-2:], "/"),
					)
					a.Value = slog.StringValue(
						strings.ReplaceAll(a.Value.String(), "}", ""),
					)
				}
			}
			return a
		}}),
)
