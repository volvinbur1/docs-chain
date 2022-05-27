package common

type UploadedPaper struct {
	Id          string `json:"uploadId" bson:"uploadId"`
	Topic       string `json:"topic" bson:"topic"`
	CreatorName string `json:"creatorName" bson:"creatorName"`
	IpfsHash    string `json:"ipfsHash" bson:"ipfsHash"`
	PaperPath   string `json:"paperPath" bson:"paperPath"`
	ReviewPath  string `json:"reviewPath" bson:"reviewPath"`
}

const (
	UnknownStatus              = "unknown"
	IsReadyForProcessingStatus = "isReadyForProcessing"
	InProgressStatus           = "inProgress"
	SuccessStatus              = "success"
)

type PaperProcessResult struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	NFT    string `json:"NFT"`
}

//blockchain
const (
	Network         = "network"
	DevNetwork      = "devnet"
	MainNetwork     = "mainnet-beta"
	Description     = "description"
	Mnemonic        = "secret_recovery_phrase"
	DerivationPath  = "derivation_path"
	NfrName         = "nft_name"
	NftSymbol       = "nft_symbol"
	NftUrl          = "nft_url"
	NftUploadMethod = "nft_upload_method"
	ApiKey          = "APIKeyID"
	ApiSecret       = "APISecretKey"
	Link            = "LINK"
)

//blockchain endpoint
const (
	NftBaseUrl   = "https://api.theblockchainapi.com/"
	MintEndpoint = "v1/solana/nft"
)

//nft image creator
const (
	ImageBaseUrl = "https://api.imgbb.com/1/upload"
	Image        = "image"
	Key          = "key"
	Expiration   = "expiration"
)

type NftMintResponse struct {
	NftMint,
	MintRecoveryPhrase,
	TransactionSignature string
}
