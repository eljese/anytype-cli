package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/internal"
)

func NewVersionCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				fmt.Println(internal.GetVersionVerbose())
			} else {
				fmt.Println(internal.GetVersionBrief())
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed version information")

	return cmd
}
