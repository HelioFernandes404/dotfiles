package main

import (
	"os"

	"github.com/heliofernandes404/dotfiles/internal/dotfiles"
)

func main() {
	os.Exit(dotfiles.Run(os.Args[1:], os.Stdout, os.Stderr))
}
