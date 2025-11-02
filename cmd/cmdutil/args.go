package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

func ExactArgs(n int, msg string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > n {
			return fmt.Errorf("too many arguments")
		}
		if len(args) < n {
			return fmt.Errorf("%s", msg)
		}
		return nil
	}
}
