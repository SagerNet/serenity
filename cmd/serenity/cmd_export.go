package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sagernet/serenity/common/metadata"
	"github.com/sagernet/serenity/common/semver"
	"github.com/sagernet/serenity/server"
	"github.com/sagernet/sing-box/log"
	boxOption "github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"

	"github.com/spf13/cobra"
)

var (
	commandExportFlagPlatform string
	commandExportFlagVersion  string
)

var commandExport = &cobra.Command{
	Use:   "export [profile]",
	Short: "Export configuration without running HTTP services",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var profileName string
		if len(args) == 1 {
			profileName = args[0]
		}
		err := export(profileName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	mainCommand.AddCommand(commandExport)
	commandExport.Flags().StringVarP(&commandExportFlagPlatform, "platform", "p", "", "platform: ios, macos, tvos, android (empty by default)")
	commandExport.Flags().StringVarP(&commandExportFlagVersion, "version", "v", "", "sing-box version (latest by default)")
}

func export(profileName string) error {
	var (
		platform metadata.Platform
		version  *semver.Version
		err      error
	)
	if commandExportFlagPlatform != "" {
		platform, err = metadata.ParsePlatform(commandExportFlagPlatform)
		if err != nil {
			return err
		}
	}
	if commandExportFlagVersion != "" {
		version = common.Ptr(semver.ParseVersion(commandExportFlagVersion))
	}

	options, err := readConfigAndMerge()
	if err != nil {
		return err
	}
	if disableColor {
		if options.Log == nil {
			options.Log = &boxOption.LogOptions{}
		}
		options.Log.DisableColor = true
	}
	ctx, cancel := context.WithCancel(globalCtx)
	instance, err := server.New(ctx, options)
	if err != nil {
		cancel()
		return E.Cause(err, "create service")
	}
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer func() {
		signal.Stop(osSignals)
		close(osSignals)
	}()
	startCtx, finishStart := context.WithCancel(context.Background())
	go func() {
		_, loaded := <-osSignals
		if loaded {
			cancel()
			closeMonitor(startCtx)
		}
	}()
	err = instance.StartHeadless()
	finishStart()
	if err != nil {
		cancel()
		return E.Cause(err, "start service")
	}
	boxOptions, err := instance.RenderHeadless(profileName, metadata.Metadata{
		Platform: platform,
		Version:  version,
	})
	if err != nil {
		return err
	}
	encoder := json.NewEncoderContext(globalCtx, os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(boxOptions)
	if err != nil {
		return E.Cause(err, "encode config")
	}
	return nil
}
