package main

import (
	"fmt"
	"time"

	"github.com/conneroisu/steroscopic-hardware/pkg/despair"
)

func main() {
	start := time.Now()
	if err := despair.RunSad("L_00001.png", "R_00001.png"); err != nil {
		panic(err)
	}
	fmt.Println("Done in", time.Since(start))
}
