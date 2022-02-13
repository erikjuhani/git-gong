package main

import (
	"github.com/erikjuhani/git-gong/cmd"
	"github.com/erikjuhani/git-gong/config"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	cmd.Execute()
}
