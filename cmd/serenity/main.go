package main

import (
	"github.com/sagernet/sing-box/log"

	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serenity",
	Short: "The sing-box configuration generator.",
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
