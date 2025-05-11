// Package main is the main entry point for the stereoscopic disparity algorithm web app.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/conneroisu/steroscopic-hardware/cmd"
)

func main() {
	err := cmd.Run(context.Background(), openBrowser)
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
		err = errors.New("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
