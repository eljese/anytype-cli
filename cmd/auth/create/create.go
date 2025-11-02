package create

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the auth create command
func NewCreateCmd() *cobra.Command {
	var name string
	var rootPath string
	var listenAddress string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new bot account",
		Long:  "Create a new Anytype bot account with a generated account key. The account key is your credential for bot authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				output.Print("Enter account name: ")
				reader := bufio.NewReader(os.Stdin)
				name, _ = reader.ReadString('\n')
				name = strings.TrimSpace(name)

				if name == "" {
					return output.Error("account name is required")
				}
			}

			accountKey, accountId, err := core.CreateBotWallet(name, rootPath, listenAddress)
			if err != nil {
				return output.Error("failed to create account: %w", err)
			}

			output.Success("Bot account created successfully!")

			output.Warning("IMPORTANT: Save your account key in a secure location.")
			output.Info("   This is the ONLY way to authenticate your bot account.")

			output.Print("")
			keyLen := len(accountKey)
			boxWidth := keyLen + 4
			if boxWidth < 24 {
				boxWidth = 24
			}

			topBorder := "╔" + strings.Repeat("═", boxWidth) + "╗"
			midBorder := "╠" + strings.Repeat("═", boxWidth) + "╣"
			botBorder := "╚" + strings.Repeat("═", boxWidth) + "╝"

			title := "BOT ACCOUNT KEY"
			titlePadding := (boxWidth - len(title)) / 2
			titleLine := "║" + strings.Repeat(" ", titlePadding) + title + strings.Repeat(" ", boxWidth-titlePadding-len(title)) + "║"

			keyLine := fmt.Sprintf("║  %s  ║", accountKey)

			output.Print(topBorder)
			output.Print(titleLine)
			output.Print(midBorder)
			output.Print(keyLine)
			output.Print(botBorder)

			output.Print("")
			output.Print("📋 Bot Account Details:")
			output.Print("   Name: %s", name)
			output.Print("   Account Id: %s", accountId)

			output.Print("")
			output.Success("You are now logged in to your new bot account.")
			output.Success("Account key saved to keychain.")

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Account name")
	cmd.Flags().StringVar(&rootPath, "root-path", "", "Root path for account data")
	cmd.Flags().StringVar(&listenAddress, "listen-address", config.DefaultAPIAddress, "API listen address in `host:port` format")

	return cmd
}
