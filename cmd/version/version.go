package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
)

func NewVersionCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				fmt.Println(core.GetVersionVerbose())
			} else {
				fmt.Println(core.GetVersionBrief())
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed version information")

	return cmd
}
