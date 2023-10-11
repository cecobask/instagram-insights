package main

import (
	"os"

	"github.com/cecobask/instagram-insights/cmd/root"
)

func main() {
	if err := root.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
