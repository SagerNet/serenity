package main

import (
	"os"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var Version = "unknown"

var commandVersion = &cobra.Command{
	Use:   "version",
	Short: "Print current version of serenity",
	Run:   printVersion,
	Args:  cobra.NoArgs,
}

var nameOnly bool

func init() {
	commandVersion.Flags().BoolVarP(&nameOnly, "name", "n", false, "print version name only")
	mainCommand.AddCommand(commandVersion)
}

func printVersion(cmd *cobra.Command, args []string) {
	if nameOnly {
		os.Stdout.WriteString(Version + "\n")
		return
	}
	version := "serenity version " + Version + "\n\n"
	version += "Environment: " + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH + "\n"

	var revision string

	debugInfo, loaded := debug.ReadBuildInfo()
	if loaded {
		for _, setting := range debugInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			}
		}
	}

	if revision != "" {
		version += "Revision: " + revision + "\n"
	}

	os.Stdout.WriteString(version)
}
