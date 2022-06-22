package backend

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"github.com/volvinbur1/docs-chain/src/central"
	"github.com/volvinbur1/docs-chain/src/common"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//TODO: remove after tests
var reqCount = 0

type WebUIProcessor struct {
	centralWorker *central.Worker
}

func NewWebUIProcessor(centralWorker *central.Worker) *WebUIProcessor {
	processor := &WebUIProcessor{centralWorker: centralWorker}

	http.Handle("/", http.FileServer(http.Dir("web/html")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc(addPaperEndpoint, processor.HandleAddPaperRequest)
	http.HandleFunc(getPaperStatusEndpoint, processor.HandleGetPaperStatusRequest)
	http.HandleFunc(getPaperInfoEndpoint, processor.HandleGetPaperInfoRequest)
	http.HandleFunc(searchForPaperEndpoint, processor.HandleSearchForPaperRequest)
	return processor
}

func (w *WebUIProcessor) HandleAddPaperRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		newPaper, err := w.parsePaperUploadRequest(request)
		if err != nil {
			log.Println(err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		w.centralWorker.NewPaperAction(newPaper)
		w.checkPaperStatus(newPaper.Id, writer)
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) HandleGetPaperStatusRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		reqCount++
		paperId := request.URL.Query().Get(paperIdKey)
		log.Println("New get paper status request. Paper id:", paperId)
		w.checkPaperStatus(paperId, writer)
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) HandleGetPaperInfoRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		paperNft := request.URL.Query().Get(paperNftKey)
		log.Println("New get paper info request. Paper nft:", paperNft)
		response := common.GetPaperResponse{
			ApiResponse: common.ApiResponse{
				Status: common.Okay,
			},
			Nft: paperNft,
			Metadata: common.PaperMetadata{
				Topic:       "Blockchain Technology in Finance",
				Description: "The banking and financial-services industry has taken notice of blockchain technology’s many advantages. This special issue explores its unlikely origins, tremendous impact, implementation challenges, and enormous potential.",
				Authors: []common.Author{{
					Name:          "Harsh",
					Surname:       "Sheth",
					ScienceDegree: "professor",
				}},
				UploadDate: "2022/18/05",
				Uniqueness: "95%",
				IpfsHash:   "QmXoypizjW3WknFiJnKLwsCnL72vedxjQkDdP1mXWo6uco",
			},
		}

		paperStatusJson, err := json.Marshal(response)
		if err != nil {
			errStr := fmt.Sprintf("session status json marshal failed. Error: %s", err)
			log.Println(errStr)
			http.Error(writer, errStr, http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if _, err = writer.Write(paperStatusJson); err != nil {
			log.Println("Writing to http response writer failed. Error:", err)
		}
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) HandleSearchForPaperRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		payload := request.URL.Query().Get(searchPayloadKey)
		log.Println("New search for paper request. Search payload:", payload)
		response := common.SearchForPaperResponse{
			ApiResponse: common.ApiResponse{
				Status: common.Okay,
			},
			Payload: payload,

			NftMetadata: []common.NftMetadata{
				{
					Address:     "503315395963",
					Symbol:      "OBT",
					Name:        "OverviewOfBlockchainTechnology",
					Transaction: "0x1dDB2C0897daD1843E362EE2fdD2dbDC0eB3a9Ec",
				},
				{
					Address:     "100311325963",
					Symbol:      "BTF",
					Name:        "BlockchainTechnologyInFinance",
					Transaction: "0x1dDB2C0897daF18432362E72fdD2dbDC0eB3a9Ec",
				},
			},
			PaperMetadata: []common.PaperMetadata{
				{
					Topic:       "Overview of Blockchain Technology",
					Description: " A blockchain is distributed, decentralized database of records which enables fast reliable transactions without a centralized management overlooking it",
					Authors: []common.Author{{
						Name:          "Janvi",
						Surname:       "Dattani",
						ScienceDegree: "phd",
					}},
					UploadDate: "2022/18/05",
					Uniqueness: "97%",
					IpfsHash:   "QmXorpiznW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco",
				},
				{
					Topic:       "Blockchain Technology in Finance",
					Description: "The banking and financial-services industry has taken notice of blockchain technology’s many advantages.",
					Authors: []common.Author{{
						Name:          "Harsh",
						Surname:       "Sheth",
						ScienceDegree: "professor",
					}},
					UploadDate: "2022/18/05",
					Uniqueness: "95%",
					IpfsHash:   "QmXoypizjW3WknFiJnKLwsCnL72vedxjQkDdP1mXWo6uco",
				},
			},
		}

		paperStatusJson, err := json.Marshal(response)
		if err != nil {
			errStr := fmt.Sprintf("session status json marshal failed. Error: %s", err)
			log.Println(errStr)
			http.Error(writer, errStr, http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if _, err = writer.Write(paperStatusJson); err != nil {
			log.Println("Writing to http response writer failed. Error:", err)
		}
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) ListenHttp() error {
	return http.ListenAndServe(":8888", nil)
}

func (w *WebUIProcessor) parsePaperUploadRequest(request *http.Request) (common.UploadedPaper, error) {
	reqCount = 0
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		return common.UploadedPaper{}, fmt.Errorf("request form parse error: %s", err)
	}

	var uploadedPaper common.UploadedPaper
	uploadedPaper.Id = xid.New().String()
	uploadedPaper.Topic = request.Form.Get(paperTopicFormKey)
	uploadedPaper.Description = request.Form.Get(paperDescriptionFormKey)

	for i := 1; i <= 3; i++ {
		name := request.Form.Get(fmt.Sprint(authorNameFormKey, i))
		surname := request.Form.Get(fmt.Sprint(authorSurnameFormKey, i))
		degree := request.Form.Get(fmt.Sprint(authorDegreeFormKey, i))

		if len(name) == 0 && len(surname) == 0 && len(degree) == 0 {
			continue
		}

		uploadedPaper.Authors = append(uploadedPaper.Authors, common.Author{
			Name:          name,
			Surname:       surname,
			ScienceDegree: degree,
		})
	}

	uploadedPaper.FilePath, err = storeFileFromRequest(request, uploadedPaper.Id, paperFileFormKey)
	return uploadedPaper, err
}

func (w *WebUIProcessor) checkPaperStatus(paperId string, writer http.ResponseWriter) {
	//returnCh := make(chan interface{})
	//w.centralWorker.GetPaperStatusAction(paperId, returnCh)
	//paperStatus, isOkay := (<-returnCh).(common.PaperProcessingResult)
	//if !isOkay || paperStatus.Status == common.UnknownStatus {
	//	writer.WriteHeader(http.StatusBadRequest)
	//	errStr := fmt.Sprintf("Paper id %s is unkown.", paperStatus.Id)
	//	log.Println(errStr)
	//	if _, err := writer.Write([]byte(errStr)); err != nil {
	//		log.Println("Writing to http response writer failed. Error:", err)
	//	}
	//	return
	//}

	var response common.AddPaperResponse
	if reqCount > 3 {
		response = common.AddPaperResponse{
			ApiResponse: common.ApiResponse{
				Status: common.Okay,
			},
			Id:         paperId,
			Uniqueness: "95%",
			IpfsHash:   "QmXoypizjW3WknFiJnKLwsCnL72vedxjQkDdP1mXWo6uco",
			Nft: common.NftMetadata{
				Address:     "100311325963",
				Symbol:      "BTF",
				Name:        "BlockchainTechnologyInFinance",
				Transaction: "0x1dDB2C0897daF18432362E72fdD2dbDC0eB3a9Ec",
				Image:       "",
			},
			NftRecoveryPhrase: "fdEvyu6BVer34DFGwsa",
		}

		//response = common.AddPaperResponse{
		//	ApiResponse: common.ApiResponse{
		//		Status:  common.LowUniqueness,
		//		Message: "Your uploaded paper is not unique enough",
		//	},
		//	Id:               paperId,
		//	Uniqueness:       "34%",
		//	SimilarPapersNft: []string{"100311587963", "506312321966", "935161325569"},
		//}
	} else {
		response = common.AddPaperResponse{
			ApiResponse: common.ApiResponse{
				Status: common.Processing,
			},
			Id: paperId,
		}
	}

	paperStatusJson, err := json.Marshal(response)
	if err != nil {
		errStr := fmt.Sprintf("session status json marshal failed. Error: %s", err)
		log.Println(errStr)
		http.Error(writer, errStr, http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	//if paperStatus.Status == common.ProcessedStatus {
	//	writer.WriteHeader(http.StatusOK)
	//} else {
	//	writer.WriteHeader(http.StatusAccepted)
	//}
	if _, err = writer.Write(paperStatusJson); err != nil {
		log.Println("Writing to http response writer failed. Error:", err)
	}
}

func storeFileFromRequest(request *http.Request, uploadId, formKey string) (string, error) {
	requestFile, _, err := request.FormFile(formKey)
	if err != nil {
		return "", fmt.Errorf("getting from form failed: %s", err)
	}
	defer common.CloserHandler(requestFile)

	if err = os.MkdirAll(filepath.Join(common.LocalStoragePath, uploadId), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create all file storega path subdirs: %s", err)
	}

	localFilePath := filepath.Join(common.LocalStoragePath, uploadId, common.PaperPdfFileName)
	localFile, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s Error: %s", localFilePath, err)
	}
	defer common.CloserHandler(localFile)

	_, err = io.Copy(localFile, requestFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy file from request to local one: %s", err)
	}

	return localFilePath, nil
}
