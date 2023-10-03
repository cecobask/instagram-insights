package main

import (
	"github.com/cecobask/instagram-insights/cmd/root"
	"os"
)

func main() {
	err := root.NewRootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}
