package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/auth"
	"github.com/anyproto/anytype-cli/cmd/config"
	"github.com/anyproto/anytype-cli/cmd/serve"
	"github.com/anyproto/anytype-cli/cmd/service"
	"github.com/anyproto/anytype-cli/cmd/shell"
	"github.com/anyproto/anytype-cli/cmd/space"
	"github.com/anyproto/anytype-cli/cmd/update"
	"github.com/anyproto/anytype-cli/cmd/version"
	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
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
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Show version information")
	rootCmd.Flags().BoolP("help", "h", false, "Show help for command")

	rootCmd.AddCommand(
		auth.NewAuthCmd(),
		config.NewConfigCmd(),
		serve.NewServeCmd(),
		service.NewServiceCmd(),
		shell.NewShellCmd(rootCmd),
		space.NewSpaceCmd(),
		update.NewUpdateCmd(),
		version.NewVersionCmd(),
	)
}

func printVersion() {
	output.Print(core.GetVersionBrief())
}
