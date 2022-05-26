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
	LocalStoragePath  = "bin/storage"
	PaperPdfFileName  = "paper.pdf"
	ReviewPdfFileName = "review.pdf"
)

type UploadedPaper struct {
	Id          string `json:"uploadId" bson:"uploadId"`
	Topic       string `json:"topic" bson:"topic"`
	CreatorName string `json:"creatorName" bson:"creatorName"`
	IpfsHash    string `json:"ipfsHash" bson:"ipfsHash"`
	PaperPath   string `json:"paperPath" bson:"paperPath"`
	ReviewPath  string `json:"reviewPath" bson:"reviewPath"`
}

type PaperAnalysis struct {
	Id string `json:"id" bson:"id"`
}

type PaperNftMetadata struct {
	Id               string   `json:"id" bson:"id"`
	Topic            string   `json:"topic" bson:"topic"`
	UploadDate       string   `json:"uploadDate" bson:"uploadDate"`
	Authors          []Author `json:"authors" bson:"authors"`
	PaperIpfsHash    string   `json:"paperIpfsHash" bson:"paperIpfsHash"`
	ReviewRating     string   `json:"reviewRating" bson:"reviewRating"`
	PaperUniqueness  string   `json:"paperUniqueness" bson:"paperUniqueness"`
	SimilarPapersNfr []string `json:"similarPapersNfr,omitempty" bson:"similarPapersNfr,omitempty"`
}

type Author struct {
	Name          string `json:"name" bson:"name"`
	Surname       string `json:"surname" bson:"surname"`
	MiddleName    string `json:"middleName,omitempty" bson:"middleName,omitempty"`
	ScienceDegree string `json:"scienceDegree" bson:"scienceDegree"`
}

type PaperShingles struct {
	Id                string   `json:"id" bson:"id"`
	Shingles          []uint32 `json:"shingles" bson:"shingles"`
	WordsInShingleCnt int      `json:"wordsInShingleCnt" bson:"wordsInShingleCnt"`
	HashAlgorithm     string   `json:"hashAlgorithm" bson:"hashAlgorithm"`
}

const (
	UnknownStatus              = "unknown"
	IsReadyForProcessingStatus = "isReadyForProcessing"
	InProgressStatus           = "inProgress"
	SuccessStatus              = "success"
	ProcessingFailedStatus     = "fail"
)

type PaperProcessingResult struct {
	Id     string `json:"id"`
	Status string `json:"status"`
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
