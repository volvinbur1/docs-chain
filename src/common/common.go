package common

import (
	"io"
	"log"
)

var CloserHandler = func(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Print(err)
	}
}

const (
	LocalStoragePath = "bin/storage"
	PaperPdfFileName = "paper.pdf"
)

type UploadedPaper struct {
	Id          string
	NftName     string
	NftSymbol   string
	Topic       string
	Description string
	Authors     []Author
	FilePath    string
}

type AnalysisResult struct {
	Id              string
	Uniqueness      float64
	SimilarPapersId []string
}

type PaperMetadata struct {
	Topic            string   `json:"topic"`
	Description      string   `json:"description"`
	Authors          []Author `json:"authors"`
	UploadDate       string   `json:"uploadDate"`
	Uniqueness       string   `json:"uniqueness"`
	IpfsHash         string   `json:"ipfsHash,omitempty"`
	SimilarPapersNft []string `json:"similarPapersNft,omitempty"`
}

type NftMetadata struct {
	Address     string `json:"address"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Transaction string `json:"transaction"`
	Image       string `json:"image"`
}

const (
	Okay          = "okay"
	Processing    = "processing"
	Failed        = "failed"
	NoResults     = "noResults"
	LowUniqueness = "lowUniqueness"
)

type ApiResponse struct {
	Status  string `json:"state"`
	Message string `json:"message,omitempty"`
}

type AddPaperResponse struct {
	ApiResponse
	Id                string      `json:"id"`
	Uniqueness        string      `json:"uniqueness,omitempty"`
	IpfsHash          string      `json:"ipfsHash,omitempty"`
	SimilarPapersNft  []string    `json:"similarPapersNft,omitempty"`
	Nft               NftMetadata `json:"nft,omitempty"`
	NftRecoveryPhrase string      `json:"nftRecoveryPhrase,omitempty"`
}

type GetPaperResponse struct {
	ApiResponse
	Nft      string        `json:"nft"`
	Metadata PaperMetadata `json:"metadata,omitempty"`
}

type SearchForPaperResponse struct {
	ApiResponse
	Payload       string          `json:"payload"`
	NftMetadata   []NftMetadata   `json:"nftMetadata,omitempty"`
	PaperMetadata []PaperMetadata `json:"paperMetadata,omitempty"`
}

type Author struct {
	Name          string `json:"name" bson:"name"`
	Surname       string `json:"surname" bson:"surname"`
	ScienceDegree string `json:"scienceDegree" bson:"scienceDegree"`
}

type PaperShingles struct {
	Id                string   `json:"id" bson:"id"`
	Shingles          []uint32 `json:"shingles" bson:"shingles"`
	WordsInShingleCnt int      `json:"wordsInShingleCnt" bson:"wordsInShingleCnt"`
	HashAlgorithm     string   `json:"hashAlgorithm" bson:"hashAlgorithm"`
}

const (
	UnknownStatus = iota
	ProcessingFailedStatus
	NotEnoughUniquenessStatus
	ReadyForProcessingStatus
	InProgressStatus
	ProcessedStatus
)

type ServiceAction = int

const (
	NewPaperUploadAction ServiceAction = iota
	GetPaperProcessingStatusAction
	GetPaperByHashAction
	GetPaperInfoByNFTAction
)

type ServiceTask struct {
	Action   ServiceAction
	Payload  interface{}
	ReturnCh chan<- interface{}
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

type NftResponse struct {
	Id                   string `json:"id" bson:"id"`
	Mint                 string `json:"mint" bson:"mint"`
	MintRecoveryPhrase   string `json:"mintRecoveryPhrase" bson:"mintRecoveryPhrase"`
	TransactionSignature string `json:"transactionSignature" bson:"transactionSignature"`
}

//IPFS
const IpfsUrl = "https://ipfs.io/ipfs/"
