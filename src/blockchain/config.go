package blockchain

import "os"

type BlockChain struct {
	Mnemonic,
	BlockApiKey,
	BlockApiSecret,
	ImgUrlApiKey string
}

func New() *BlockChain {
	return &BlockChain{
		Mnemonic:       getEnv("MNEMONIC", ""),
		BlockApiKey:    getEnv("BLOCK_API_KEY", ""),
		BlockApiSecret: getEnv("BLOCK_API_SECRET", ""),
		ImgUrlApiKey:   getEnv("IMG_URL_API_KEY", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
