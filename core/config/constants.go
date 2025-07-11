package config

const (
	// Default addresses
	DefaultBindAddress = "0.0.0.0"

	// Port configuration
	GRPCPort    = "31007"
	GRPCWebPort = "31008"
	APIPort     = "31009"
	DaemonPort  = "31010"

	// Full addresses
	DefaultGRPCAddress    = DefaultBindAddress + ":" + GRPCPort
	DefaultGRPCWebAddress = DefaultBindAddress + ":" + GRPCWebPort
	DefaultAPIAddress     = DefaultBindAddress + ":" + APIPort
	DefaultDaemonAddress  = DefaultBindAddress + ":" + DaemonPort

	// URLs
	DaemonHTTPURL  = "http://" + DefaultDaemonAddress
	GRPCDNSAddress = "dns:///" + DefaultGRPCAddress

	// External URLs
	GitHubBaseURL    = "https://github.com/anyproto/anytype-cli"
	GitHubCommitURL  = GitHubBaseURL + "/commit/"
	GitHubReleaseURL = GitHubBaseURL + "/releases/tag/"

	// Environment variable names
	EnvGRPCAddr    = "ANYTYPE_GRPC_ADDR"
	EnvGRPCWebAddr = "ANYTYPE_GRPCWEB_ADDR"

	AnytypeNetworkAddress = "N83gJpVd9MuNRZAuJLZ7LiMntTThhPc6DtzWWVjb1M3PouVU"
)
