package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/sagernet/serenity"
	"github.com/sagernet/sing-box/log"

	"github.com/spf13/cobra"
)

var configPath string

var commandRun = &cobra.Command{
	Use:   "run",
	Short: "Run serenity",
	Run: func(cmd *cobra.Command, args []string) {
		err := run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	commandRun.Flags().StringVarP(&configPath, "config", "c", "config.json", "set configuration file path")
	command.AddCommand(commandRun)
}

func run() error {
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	var options serenity.Options
	err = options.UnmarshalJSON(configContent)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	server, err := serenity.NewServer(ctx, options)
	if err != nil {
		return err
	}
	err = server.Start()
	if err != nil {
		return err
	}
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(osSignals)
	debug.FreeOSMemory()
	<-osSignals
	cancel()
	err = server.Close()
	if err != nil {
		return err
	}
	return nil
}
