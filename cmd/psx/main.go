package main

import (
	"os"

	"github.com/m-mdy-m/psx/internal/command"
)

var Version = "development"

func main() {
	if err := command.Execute(Version); err != nil {
		os.Exit(1)
	}
}
