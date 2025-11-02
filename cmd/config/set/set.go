package set

import (
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/cmd/cmdutil"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long:  `Set a specific configuration value`,
		Args:  cmdutil.ExactArgs(2, "cannot set config: key and value arguments required"),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			value := args[1]

			configMgr := config.GetConfigManager()
			if err := configMgr.Load(); err != nil {
				return output.Error("failed to load config: %w", err)
			}

			switch key {
			case "accountId":
				if err := configMgr.SetAccountId(value); err != nil {
					return output.Error("failed to set account Id: %w", err)
				}
			case "techSpaceId":
				if err := configMgr.SetTechSpaceId(value); err != nil {
					return output.Error("failed to set tech space Id: %w", err)
				}
			default:
				return output.Error("unknown config key: %s", key)
			}

			output.Success("Set %s = %s", key, value)
			return nil
		},
	}
}
