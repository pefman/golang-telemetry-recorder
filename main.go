package main

import (
	"fmt"
	"os"

	"github.com/pefman/golang-telemetry-recorder/internal/menu"
)

func main() {
	fmt.Println("==============================================")
	fmt.Println("  F1 Telemetry Recorder - v1.0")
	fmt.Println("==============================================")
	fmt.Println()

	if err := menu.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
