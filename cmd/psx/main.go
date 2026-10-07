package main

import (
	"os"

	"github.com/m-mdy-m/psx/internal/command"
)

func main() {
	os.Exit(command.Execute())
}
