package config

const (
	// Default addresses
	LocalhostIP = "127.0.0.1"

	// Port configuration
	GRPCPort    = "31007"
	GRPCWebPort = "31008"
	APIPort     = "31009"

	// Full addresses
	DefaultGRPCAddress    = LocalhostIP + ":" + GRPCPort
	DefaultGRPCWebAddress = LocalhostIP + ":" + GRPCWebPort
	DefaultAPIAddress     = LocalhostIP + ":" + APIPort

	// URLs
	GRPCDNSAddress = "dns:///" + DefaultGRPCAddress

	// External URLs
	GitHubBaseURL    = "https://github.com/anyproto/anytype-cli"
	GitHubCommitURL  = GitHubBaseURL + "/commit/"
	GitHubReleaseURL = GitHubBaseURL + "/releases/tag/"

	// Anytype network address
	AnytypeNetworkAddress = "N83gJpVd9MuNRZAuJLZ7LiMntTThhPc6DtzWWVjb1M3PouVU"
)
