package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/anyproto/anytype-heart/pb"
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
		apiAddr = "127.0.0.1:31009"
	}

	client, err := GetGRPCClient()
	if err != nil {
		return fmt.Errorf("error connecting to gRPC server: %w", err)
	}

	// Create a context for the initial calls.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set initial parameters.
	_, err = client.InitialSetParameters(ctx, &pb.RpcInitialSetParametersRequest{
		Platform: runtime.GOOS,
		Version:  Version,
		Workdir:  getDefaultWorkDir(),
	})
	if err != nil {
		return fmt.Errorf("failed to set initial parameters: %w", err)
	}

	// Recover the wallet.
	_, err = client.WalletRecover(ctx, &pb.RpcWalletRecoverRequest{
		Mnemonic: mnemonic,
		RootPath: rootPath,
	})
	if err != nil {
		return fmt.Errorf("wallet recovery failed: %w", err)
	}

	// Create a session.
	resp, err := client.WalletCreateSession(ctx, &pb.RpcWalletCreateSessionRequest{
		Auth: &pb.RpcWalletCreateSessionRequestAuthOfMnemonic{
			Mnemonic: mnemonic,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	sessionToken := resp.Token
	err = SaveToken(sessionToken)
	fmt.Println("ℹ Session token:", sessionToken)
	if err != nil {
		return fmt.Errorf("failed to save session token: %w", err)
	}

	// Start listening for session events.
	er, err := ListenForEvents(sessionToken)
	if err != nil {
		return fmt.Errorf("failed to start event listener: %w", err)
	}

	// Recover the account.
	ctx, cancel = ClientContextWithAuthTimeout(sessionToken, 5*time.Second)
	defer cancel()
	_, err = client.AccountRecover(ctx, &pb.RpcAccountRecoverRequest{})
	if err != nil {
		return fmt.Errorf("account recovery failed: %w", err)
	}

	// Wait for the account ID.
	accountID, err := WaitForAccountID(er)
	if err != nil {
		return fmt.Errorf("error waiting for account ID: %w", err)
	}
	fmt.Println("ℹ Account ID:", accountID)

	// Select the account.
	ctx2, cancel2 := ClientContextWithAuthTimeout(sessionToken, 5*time.Second)
	defer cancel2()
	_, err = client.AccountSelect(ctx2, &pb.RpcAccountSelectRequest{
		DisableLocalNetworkSync: false,
		Id:                      accountID,
		JsonApiListenAddr:       apiAddr,
		RootPath:                rootPath,
	})
	if err != nil {
		return fmt.Errorf("failed to select account: %w", err)
	}

	return nil
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
	client, err := GetGRPCClient()
	if err != nil {
		fmt.Println("Failed to connect to gRPC server:", err)
	}

	token, err := GetStoredToken()
	if err != nil {
		return fmt.Errorf("failed to get stored token: %w", err)
	}

	ctx, cancel := ClientContextWithAuthTimeout(token, 5*time.Second)
	defer cancel()

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

	if err := DeleteStoredMnemonic(); err != nil {
		return fmt.Errorf("failed to delete stored mnemonic: %w", err)
	}

	if err := DeleteStoredToken(); err != nil {
		return fmt.Errorf("failed to delete stored token: %w", err)
	}

	return nil
}
