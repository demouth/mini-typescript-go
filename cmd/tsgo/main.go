package main

import (
	"fmt"
	"os"

	ts "github.com/demouth/mini-typescript-go/internal/compiler"
	"github.com/demouth/mini-typescript-go/internal/vfs"
)

func main() {
	currentDirectory, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	fs := vfs.FromOS()

	host := ts.NewCompilerHost(currentDirectory, fs)
	program := ts.NewProgram(
		ts.ProgramOptions{
			Host: host,
		},
	)
	program.Emit()
}
