//go:build run
// +build run

// Package main is an quick way to run the stereoscopic disparity algorithm.
package main

import (
	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

func main() {

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	// Load left and right images
	leftImg, err := despair.LoadPNG("./testdata/L_00001.png")
	if err != nil {
		return err
	}

	rightImg, err := despair.LoadPNG("./testdata/R_00001.png")
	if err != nil {
		return err
	}
	got := despair.RunSad(leftImg, rightImg, 16, 64)
	err = despair.SavePNG("output.png", got)
	if err != nil {
		return err
	}
	return nil
}
