package backend

const (
	apiBaseEndpoint        = "/api/v1"
	addPaperEndpoint       = apiBaseEndpoint + "/addPaper"
	getPaperStatusEndpoint = apiBaseEndpoint + "/getPaperStatus"
	getPaperInfoEndpoint   = apiBaseEndpoint + "/getPaperInfo"
	searchForPaperEndpoint = addPaperEndpoint + "/searchForPaper"
)

const (
	paperFileFormKey        = "paperFile"
	paperTopicFormKey       = "paperTopic"
	paperDescriptionFormKey = "paperDescription"
	authorNameFormKey       = "authorName"
	authorSurnameFormKey    = "authorSurname"
	authorDegreeFormKey     = "authorDegree"
)

const (
	paperIdKey       = "paperId"
	paperNftKey      = "paperNft"
	searchPayloadKey = "searchPayload"
)
