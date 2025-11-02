package core

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"

	"github.com/anyproto/anytype-cli/core/config"
	"github.com/anyproto/anytype-cli/core/output"
)

// Authenticate performs the full authentication flow for a bot account using an account key.
// This includes wallet recovery, session creation, account recovery, account selection, and config persistence.
func Authenticate(accountKey, rootPath, apiAddr string) error {
	if rootPath == "" {
		rootPath = config.GetDataDir()
	}
	if apiAddr == "" {
		apiAddr = config.DefaultAPIAddress
	}

	var sessionToken string
	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
			Platform: runtime.GOOS,
			Version:  Version,
			Workdir:  config.GetWorkDir(),
		})
		if err != nil {
			return fmt.Errorf("failed to set initial parameters: %w", err)
		}
		if resp.Error.Code != pb.RpcInitialSetParametersResponseError_NULL {
			return fmt.Errorf("failed to set initial parameters: %s", resp.Error.Description)
		}

		resp2, err := client.WalletRecover(ctx, &pb.RpcWalletRecoverRequest{
			AccountKey: accountKey,
			RootPath:   rootPath,
		})
		if err != nil {
			return fmt.Errorf("wallet recovery failed: %w", err)
		}
		if resp2.Error.Code != pb.RpcWalletRecoverResponseError_NULL {
			return fmt.Errorf("wallet recovery failed: %s", resp2.Error.Description)
		}

		resp3, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
			Auth: &pb.RpcWalletCreateSessionRequestAuthOfAccountKey{
				AccountKey: accountKey,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		if resp3.Error.Code != pb.RpcWalletCreateSessionResponseError_NULL {
			return fmt.Errorf("failed to create session: %s", resp3.Error.Description)
		}
		sessionToken = resp3.Token
		return nil
	})
	if err != nil {
		return err
	}

	err = SaveToken(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to save session token: %w", err)
	}

	er, err := ListenForEvents(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to start event listener: %w", err)
	}

	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.AccountRecover(ctx, &pb.RpcAccountRecoverRequest{})
		if err != nil {
			return fmt.Errorf("account recovery failed: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	accountId, err := WaitForAccountId(er)
	if err != nil {
		return fmt.Errorf("error waiting for account Id: %w", err)
	}
	output.Info("Account Id: %s", accountId)

	var techSpaceId string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
			Id:                accountId,
			JsonApiListenAddr: apiAddr,
			RootPath:          rootPath,
		})
		if err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		if resp.Account != nil && resp.Account.Info != nil {
			techSpaceId = resp.Account.Info.TechSpaceId
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := config.SetAccountIdToConfig(accountId); err != nil {
		output.Warning("failed to save account Id: %v", err)
	}
	if techSpaceId != "" {
		if err := config.SetTechSpaceIdToConfig(techSpaceId); err != nil {
			output.Warning("failed to save tech space Id: %v", err)
		}
	}

	return nil
}

// ValidateAccountKey checks if the provided account key is valid.
func ValidateAccountKey(accountKey string) error {
	if accountKey == "" {
		return fmt.Errorf("account key cannot be empty")
	}

	words := strings.Fields(accountKey)
	if len(words) >= 12 {
		return fmt.Errorf("this appears to be a mnemonic phrase, not an account key - the CLI only supports bot accounts created via 'anytype auth create'")
	}

	decoded, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return fmt.Errorf("invalid account key format: must be valid base64")
	}

	if len(decoded) < 32 {
		return fmt.Errorf("invalid account key format: insufficient key material")
	}

	return nil
}

// Login handles user interaction for login by prompting for account key if not provided,
// validating it, performing authentication, and saving the key to keychain.
func Login(accountKey, rootPath, apiAddr string) error {
	if accountKey == "" {
		output.Print("Enter account key: ")
		reader := bufio.NewReader(os.Stdin)
		accountKey, _ = reader.ReadString('\n')
		accountKey = strings.TrimSpace(accountKey)
	}

	if err := ValidateAccountKey(accountKey); err != nil {
		return err
	}

	if err := Authenticate(accountKey, rootPath, apiAddr); err != nil {
		return err
	}

	if err := SaveAccountKey(accountKey); err != nil {
		output.Warning("failed to save account key in keychain: %v", err)
	} else {
		output.Success("Account key saved to keychain.")
	}

	return nil
}

// Logout logs out the current user by stopping the account, closing the wallet session,
// deleting stored credentials, and clearing the config.
func Logout() error {
	token, err := GetStoredToken()
	if err != nil {
		return fmt.Errorf("failed to get stored token: %w", err)
	}

	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountStop(ctx, &pb.RpcAccountStopRequest{
			RemoveData: false,
		})
		if err != nil {
			return fmt.Errorf("failed to log out: %w", err)
		}
		if resp.Error.Code != pb.RpcAccountStopResponseError_NULL {
			output.Warning("Failed to log out: %s", resp.Error.Description)
		}

		resp2, err := client.WalletCloseSession(ctx, &pb.RpcWalletCloseSessionRequest{Token: token})
		if err != nil {
			return fmt.Errorf("failed to close session: %w", err)
		}
		if resp2.Error.Code != pb.RpcWalletCloseSessionResponseError_NULL {
			output.Warning("Failed to close session: %s", resp2.Error.Description)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := DeleteStoredAccountKey(); err != nil {
		return fmt.Errorf("failed to delete stored account key: %w", err)
	}

	if err := DeleteStoredToken(); err != nil {
		return fmt.Errorf("failed to delete stored token: %w", err)
	}

	configMgr := config.GetConfigManager()
	if err := configMgr.Delete(); err != nil {
		output.Warning("failed to clear config: %v", err)
	}

	CloseEventReceiver()

	return nil
}

// CreateWallet creates a new wallet and account, establishes a session,
// saves credentials, and returns the account key and account ID.
func CreateWallet(name, rootPath, apiAddr string) (string, string, error) {
	if rootPath == "" {
		rootPath = config.GetDataDir()
	}
	if apiAddr == "" {
		apiAddr = config.DefaultAPIAddress
	}

	var sessionToken string
	var accountKey string

	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
			Platform: runtime.GOOS,
			Version:  Version,
			Workdir:  config.GetWorkDir(),
		})
		if err != nil {
			return fmt.Errorf("failed to set initial parameters: %w", err)
		}

		createResp, err := client.WalletCreate(ctx, &pb.RpcWalletCreateRequest{
			RootPath: rootPath,
		})
		if err != nil {
			return fmt.Errorf("wallet creation failed: %w", err)
		}
		accountKey = createResp.AccountKey

		sessionResp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
			Auth: &pb.RpcWalletCreateSessionRequestAuthOfAccountKey{
				AccountKey: accountKey,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		sessionToken = sessionResp.Token
		return nil
	})

	if err != nil {
		return "", "", err
	}

	err = SaveToken(sessionToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to save session token: %w", err)
	}

	_, err = ListenForEvents(sessionToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to start event listener: %w", err)
	}

	var accountId string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountCreate(ctx, &pb.RpcAccountCreateRequest{
			Name:              name,
			StorePath:         rootPath,
			JsonApiListenAddr: apiAddr,
		})
		if err != nil {
			return fmt.Errorf("account creation failed: %w", err)
		}
		accountId = resp.Account.Id
		return nil
	})
	if err != nil {
		return "", "", err
	}

	var techSpaceId string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
			Id:                accountId,
			JsonApiListenAddr: apiAddr,
			RootPath:          rootPath,
		})
		if err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		if resp.Account != nil && resp.Account.Info != nil {
			techSpaceId = resp.Account.Info.TechSpaceId
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	if err := SaveAccountKey(accountKey); err != nil {
		output.Warning("failed to save account key: %v", err)
	}

	if err := config.SetAccountIdToConfig(accountId); err != nil {
		output.Warning("failed to save account Id: %v", err)
	}
	if techSpaceId != "" {
		if err := config.SetTechSpaceIdToConfig(techSpaceId); err != nil {
			output.Warning("failed to save tech space Id: %v", err)
		}
	}

	return accountKey, accountId, nil
}
