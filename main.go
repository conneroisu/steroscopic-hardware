// Package main is the main entry point for the stereoscopic disparity algorithm web app.
package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/conneroisu/steroscopic-hardware/cmd/steroscopic"
)

func main() {
	err := steroscopic.Run(context.Background(), func() {})
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
