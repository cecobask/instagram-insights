package main

import (
	"os"

	"github.com/cecobask/instagram-insights/cmd/root"
)

func main() {
	err := root.NewRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
