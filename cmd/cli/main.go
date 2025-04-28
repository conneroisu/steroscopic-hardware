package main

import (
	"fmt"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

func main() {
	start := time.Now()
	if err := despair.RunSad("L_00001.png", "R_00001.png", &despair.Parameters{
		BlockSize:    16,
		MaxDisparity: 64,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Done in", time.Since(start))
}
