package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/auth"
	"github.com/anyproto/anytype-cli/cmd/daemon"
	"github.com/anyproto/anytype-cli/cmd/server"
	"github.com/anyproto/anytype-cli/cmd/shell"
	"github.com/anyproto/anytype-cli/cmd/space"
	"github.com/anyproto/anytype-cli/cmd/update"
	"github.com/anyproto/anytype-cli/cmd/version"
	"github.com/anyproto/anytype-cli/core"
)

var (
	versionFlag bool
	rootCmd     = &cobra.Command{
		Use:   "anytype <command> <subcommand> [flags]",
		Short: "Anytype CLI",
		Long:  "Seamlessly interact with Anytype from the command line",
		Run: func(cmd *cobra.Command, args []string) {
			if versionFlag {
				printVersion()
				return
			}
			cmd.Help()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show version information")

	rootCmd.AddCommand(
		auth.NewAuthCmd(),
		daemon.NewDaemonCmd(),
		server.NewServerCmd(),
		shell.NewShellCmd(rootCmd),
		space.NewSpaceCmd(),
		update.NewUpdateCmd(),
		version.NewVersionCmd(),
	)
}

func printVersion() {
	fmt.Println(core.GetVersionBrief())
}
