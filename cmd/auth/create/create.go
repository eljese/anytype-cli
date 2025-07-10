package create

import (
	"fmt"
	"strings"

	"github.com/anyproto/anytype-cli/internal"
	"github.com/spf13/cobra"
)

// NewCreateCmd creates the auth create command
func NewCreateCmd() *cobra.Command {
	var name string
	var rootPath string
	var apiAddr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Anytype account",
		Long:  "Create a new Anytype account with a generated mnemonic phrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("account name is required")
			}

			mnemonic, accountID, err := internal.CreateWallet(name, rootPath, apiAddr)
			if err != nil {
				return fmt.Errorf("failed to create account: %w", err)
			}

			// Success message first
			fmt.Println("✓ Account created successfully!")

			// Important warning
			fmt.Println("\n⚠️ IMPORTANT: Save your mnemonic phrase in a secure location.")
			fmt.Println("   This is the ONLY way to recover your account if you lose access.")

			// Print mnemonic in a box
			words := strings.Split(mnemonic, " ")
			fmt.Println("\n╔════════════════════════════════════════════════════════╗")
			fmt.Println("║                    MNEMONIC PHRASE                     ║")
			fmt.Println("╠════════════════════════════════════════════════════════╣")
			// Print 6 words per line
			fmt.Printf("║  %-52s  ║\n", strings.Join(words[0:6], " "))
			fmt.Printf("║  %-52s  ║\n", strings.Join(words[6:12], " "))
			fmt.Println("╚════════════════════════════════════════════════════════╝")

			// Account details
			fmt.Println("\n📋 Account Details:")
			fmt.Printf("   Name: %s\n", name)
			fmt.Printf("   ID: %s\n", accountID)

			// Final status
			fmt.Println("\n✓ You are now logged in to your new account.")
			fmt.Println("✓ Mnemonic saved to keychain.")

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Account name (required)")
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&rootPath, "root-path", "", "Custom root path for storing account data")
	cmd.Flags().StringVar(&apiAddr, "api-addr", "", "Custom API address (default: 127.0.0.1:31009)")

	return cmd
}
