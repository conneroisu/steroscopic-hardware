//go:build run
// +build run

// Package main is an quick way to run the stereoscopic disparity algorithm.
package main

import (
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

func main() {
	if err := despair.RunSadPaths("./testdata/im0.png", "./testdata/im1.png", 16, 64); err != nil {
		panic(err)
	}
}
