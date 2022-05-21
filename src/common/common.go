package common

type UploadedPaper struct {
	Id          string `json:"uploadId" bson:"uploadId"`
	Topic       string `json:"topic" bson:"topic"`
	CreatorName string `json:"creatorName" bson:"creatorName"`
	IpfsHash    string `json:"ipfsHash" bson:"ipfsHash"`
	PaperPath   string `json:"paperPath" bson:"paperPath"`
	ReviewPath  string `json:"reviewPath" bson:"reviewPath"`
}
