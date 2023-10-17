package main

import (
	"github.com/cecobask/instagram-insights/cmd/root"
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/spf13/cobra/doc"
)

func main() {
	if err := doc.GenMarkdownTree(root.NewRootCommand(), instagram.PathDocs); err != nil {
		panic(err)
	}
}
