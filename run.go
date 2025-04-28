//go:build run
// +build run

// Package main is an quick way to run the stereoscopic disparity algorithm.
package main

import (
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

func main() {
	if err := despair.RunSadPaths("L_00001.png", "R_00001.png", 16, 64); err != nil {
		panic(err)
	}
}
