package main

import (
	"fmt"
	"os"

	"github.com/Angel-Sechar/skill-craft/internal/tui"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("skill-craft v" + version)
		return
	}

	if err := tui.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
