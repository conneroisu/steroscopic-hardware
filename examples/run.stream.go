//go:build run && stream
// +build run,stream

// Package main is an quick way to run the stereoscopic disparity algorithm in streaming mode.
package main

import (
	"fmt"
	"image"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

const (
	blockSize  = 16
	maxDespair = 32
	numWorkers = 32
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	fmt.Printf("Using %d workers (green threads)\n", numWorkers)
	inps, outs := despair.SetupConcurrentSAD(numWorkers)
	defer close(inps)

	for i := 0; i < 10; i++ {
		leftImg, err := despair.LoadPNG("./testdata/L_00001.png")
		if err != nil {
			return err
		}

		rightImg, err := despair.LoadPNG("./testdata/R_00001.png")
		if err != nil {
			return err
		}
		chunkSize := max(1, leftImg.Rect.Dy()/(numWorkers*4))
		numChunks := (leftImg.Rect.Dy() + chunkSize - 1) / chunkSize

		start := time.Now()
		for y := leftImg.Rect.Min.Y; y < leftImg.Rect.Max.Y; y += chunkSize {
			inps <- despair.InputChunk{
				Left:  leftImg,
				Right: rightImg,
				Region: image.Rect(
					leftImg.Rect.Min.X,
					y,
					leftImg.Rect.Max.X,
					min(y+chunkSize, leftImg.Rect.Max.Y),
				),
			}
		}
		got := despair.AssembleDisparityMap(outs, leftImg.Rect, numChunks)
		end := time.Now()

		fmt.Printf("Elapsed time: %v\n", end.Sub(start))
		err = despair.SavePNG("output.png", got)
		if err != nil {
			return err
		}
	}
	return nil
}
