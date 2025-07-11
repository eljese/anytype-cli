package core

import (
	"bufio"
	"context"
	"fmt"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anyproto/anytype-cli/core/config"
)

// getDefaultDataPath returns the default data path for Anytype based on the operating system
func getDefaultDataPath() string {
	if dataPath := os.Getenv("DATA_PATH"); dataPath != "" {
		return dataPath
	}

	baseDir := getDefaultWorkDir()
	return filepath.Join(baseDir, "data")
}

// getDefaultWorkDir returns the default work directory for Anytype based on the operating system
func getDefaultWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "anytype")
	case "windows":
		return filepath.Join(homeDir, "AppData", "Roaming", "anytype")
	default:
		return filepath.Join(homeDir, ".config", "anytype")
	}
}

// LoginAccount performs the common steps for logging in with a given mnemonic and root path.
func LoginAccount(mnemonic, rootPath, apiAddr string) error {
	if rootPath == "" {
		rootPath = getDefaultDataPath()
	}
	if apiAddr == "" {
		apiAddr = config.DefaultAPIAddress
	}

	var sessionToken string

	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
			Platform: runtime.GOOS,
			Version:  Version,
			Workdir:  getDefaultWorkDir(),
		})
		if err != nil {
			return fmt.Errorf("failed to set initial parameters: %w", err)
		}

		_, err = client.WalletRecover(ctx, &pb.RpcWalletRecoverRequest{
			Mnemonic: mnemonic,
			RootPath: rootPath,
		})
		if err != nil {
			return fmt.Errorf("wallet recovery failed: %w", err)
		}

		resp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
			Auth: &pb.RpcWalletCreateSessionRequestAuthOfMnemonic{
				Mnemonic: mnemonic,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		sessionToken = resp.Token
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

	accountID, err := WaitForAccountID(er)
	if err != nil {
		return fmt.Errorf("error waiting for account ID: %w", err)
	}
	fmt.Println("ℹ Account ID:", accountID)

	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
			Id:                accountID,
			JsonApiListenAddr: apiAddr,
			RootPath:          rootPath,
		})
		if err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		return nil
	})

	return err
}

func Login(mnemonic, rootPath, apiAddr string) error {
	usedStoredMnemonic := false
	if mnemonic == "" {
		storedMnemonic, err := GetStoredMnemonic()
		if err == nil && storedMnemonic != "" {
			mnemonic = storedMnemonic
			fmt.Println("Using stored mnemonic from keychain.")
			usedStoredMnemonic = true
		} else {
			fmt.Print("Enter mnemonic (12 words): ")
			reader := bufio.NewReader(os.Stdin)
			mnemonic, _ = reader.ReadString('\n')
			mnemonic = strings.TrimSpace(mnemonic)
		}
	}

	if len(strings.Split(mnemonic, " ")) != 12 {
		return fmt.Errorf("mnemonic must be 12 words")
	}

	err := LoginAccount(mnemonic, rootPath, apiAddr)
	if err != nil {
		return fmt.Errorf("failed to log in: %w", err)
	}

	if !usedStoredMnemonic {
		if err := SaveMnemonic(mnemonic); err != nil {
			fmt.Println("Warning: failed to save mnemonic in keychain:", err)
		} else {
			fmt.Println("✓ Mnemonic saved to keychain.")
		}
	}

	return nil
}

func Logout() error {
	// Need to get token for WalletCloseSession request parameter
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
			fmt.Println("Failed to log out:", resp.Error.Description)
		}

		resp2, err := client.WalletCloseSession(ctx, &pb.RpcWalletCloseSessionRequest{Token: token})
		if err != nil {
			return fmt.Errorf("failed to close session: %w", err)
		}
		if resp2.Error.Code != pb.RpcWalletCloseSessionResponseError_NULL {
			fmt.Println("Failed to close session:", resp2.Error.Description)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := DeleteStoredMnemonic(); err != nil {
		return fmt.Errorf("failed to delete stored mnemonic: %w", err)
	}

	if err := DeleteStoredToken(); err != nil {
		return fmt.Errorf("failed to delete stored token: %w", err)
	}

	CloseEventReceiver()

	return nil
}

// CreateWallet creates a new wallet with the given root path and returns the mnemonic and account ID
func CreateWallet(name, rootPath, apiAddr string) (string, string, error) {
	if rootPath == "" {
		rootPath = getDefaultDataPath()
	}
	if apiAddr == "" {
		apiAddr = config.DefaultAPIAddress
	}

	var mnemonic string
	var sessionToken string

	err := GRPCCallNoAuth(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
			Platform: runtime.GOOS,
			Version:  Version,
			Workdir:  getDefaultWorkDir(),
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
		mnemonic = createResp.Mnemonic

		sessionResp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
			Auth: &pb.RpcWalletCreateSessionRequestAuthOfMnemonic{
				Mnemonic: mnemonic,
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

	var accountID string
	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		resp, err := client.AccountCreate(ctx, &pb.RpcAccountCreateRequest{
			Name:              name,
			StorePath:         rootPath,
			JsonApiListenAddr: apiAddr,
		})
		if err != nil {
			return fmt.Errorf("account creation failed: %w", err)
		}
		accountID = resp.Account.Id
		return nil
	})
	if err != nil {
		return "", "", err
	}

	err = GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
		_, err := client.AccountSelect(ctx, &pb.RpcAccountSelectRequest{
			Id:                accountID,
			JsonApiListenAddr: apiAddr,
			RootPath:          rootPath,
		})
		if err != nil {
			return fmt.Errorf("failed to select account: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	_ = SaveMnemonic(mnemonic)

	return mnemonic, accountID, nil
}
