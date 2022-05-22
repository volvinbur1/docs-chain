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
	UnknownSessionStatus    = "unknown"
	InProgressSessionStatus = "inProgress"
	SuccessSessionStatus    = "success"
)

type ProcessingSession struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	NFT    string `json:"NFT"`
}
