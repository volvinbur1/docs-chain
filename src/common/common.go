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
	Id            string
	NftName       string
	NftSymbol     string
	Topic         string
	Description   string
	Authors       []Author
	PaperFilePath string
}

type AnalysisResult struct {
	Id              string
	Uniqueness      float64
	SimilarPapersId []string
}

type PaperMetadata struct {
	Id               string   `json:"id" bson:"id"`
	Topic            string   `json:"topic" bson:"topic"`
	UploadDate       string   `json:"uploadDate" bson:"uploadDate"`
	Authors          []Author `json:"authors" bson:"authors"`
	PaperIpfsHash    string   `json:"paperIpfsHash,omitempty" bson:"paperIpfsHash,omitempty"`
	ReviewRating     string   `json:"reviewRating,omitempty" bson:"reviewRating,omitempty"`
	PaperUniqueness  string   `json:"paperUniqueness,omitempty" bson:"paperUniqueness,omitempty"`
	SimilarPapersNfr []string `json:"similarPapersNfr,omitempty" bson:"similarPapersNfr,omitempty"`
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
	IsReadyForProcessingStatus
	InProgressStatus
	SuccessStatus
)

type PaperProcessingResult struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	NFT    string `json:"NFT"`
}

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
