package utils

import "os"

func LoadChainConfigs() []ChainConfig {
	var chains []ChainConfig

	if cfg, ok := evmConfigFromEnv("Ethereum", "ETHEREUM"); ok {
		chains = append(chains, cfg)
	}
	if cfg, ok := evmConfigFromEnv("Base", "BASE"); ok {
		chains = append(chains, cfg)
	}
	if cfg, ok := svmConfigFromEnv("Solana", "SOLANA"); ok {
		chains = append(chains, cfg)
	}

	return chains
}

func evmConfigFromEnv(name, prefix string) (ChainConfig, bool) {
	wsURL := os.Getenv(prefix + "_WSS_URL")
	contractAddress := os.Getenv(prefix + "_CONTRACT_ADDRESS")
	if wsURL == "" || contractAddress == "" {
		return ChainConfig{}, false
	}

	return ChainConfig{
		Name:            name,
		Runtime:         RuntimeEVM,
		WSURL:           wsURL,
		ContractAddress: contractAddress,
		EventTypes:      []string{"Transfer"},
	}, true
}

func svmConfigFromEnv(name, prefix string) (ChainConfig, bool) {
	wsURL := os.Getenv(prefix + "_WSS_URL")
	programID := os.Getenv(prefix + "_PROGRAM_ID")
	if wsURL == "" || programID == "" {
		return ChainConfig{}, false
	}

	return ChainConfig{
		Name:       name,
		Runtime:    RuntimeSVM,
		WSURL:      wsURL,
		ProgramID:  programID,
		EventTypes: []string{"Transfer", "Swap", "Mint", "Burn"},
	}, true
}
