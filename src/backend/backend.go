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

const (
	paperUploadPageEndpoint   = "/paper-upload"
	paperUploadStatusEndpoint = paperUploadPageEndpoint + "/status"
)

const (
	paperTopicFormKey   = "paperTopic"
	uploaderNameFormKey = "uploaderName"
	reviewFileFormKey   = "reviewFile"
	paperFileFormKey    = "paperFile"
	paperIdKey          = "paperId"
)

type WebUIProcessor struct {
	centralWorker *central.Worker
}

func NewWebUIProcessor(centralWorker *central.Worker) *WebUIProcessor {
	processor := &WebUIProcessor{centralWorker: centralWorker}

	http.Handle("/", http.FileServer(http.Dir("web/html")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc(paperUploadPageEndpoint, processor.HandlePaperUploadRequest)
	http.HandleFunc(paperUploadStatusEndpoint, processor.HandlePaperUploadStatus)
	return processor
}

func (w *WebUIProcessor) HandlePaperUploadRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		http.ServeFile(writer, request, "web/html/paper_upload_page.html")
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

func (w *WebUIProcessor) HandlePaperUploadStatus(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		paperId := request.URL.Query().Get(paperIdKey)
		log.Println("New get paper status request. Paper id:", paperId)
		w.checkPaperStatus(paperId, writer)
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) ListenHttp() error {
	return http.ListenAndServe(":8888", nil)
}

func (w *WebUIProcessor) parsePaperUploadRequest(request *http.Request) (common.UploadedPaper, error) {
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		return common.UploadedPaper{}, fmt.Errorf("request form parse error: %s", err)
	}

	var uploadedPaper common.UploadedPaper
	uploadedPaper.Id = xid.New().String()
	uploadedPaper.Topic = request.Form.Get(paperTopicFormKey)
	uploadedPaper.Authors = append(uploadedPaper.Authors, common.Author{Name: request.Form.Get(uploaderNameFormKey)})
	uploadedPaper.PaperFilePath, err = storeFileFromRequest(request, uploadedPaper.Id, paperFileFormKey)
	if err != nil {
		return common.UploadedPaper{}, err
	}
	uploadedPaper.ReviewFilePath, err = storeFileFromRequest(request, uploadedPaper.Id, reviewFileFormKey)
	return uploadedPaper, err
}

func (w *WebUIProcessor) checkPaperStatus(paperId string, writer http.ResponseWriter) {
	returnCh := make(chan interface{})
	w.centralWorker.GetPaperStatusAction(paperId, returnCh)
	paperStatus, isOkay := (<-returnCh).(common.PaperProcessingResult)
	if !isOkay || paperStatus.Status == common.UnknownStatus {
		writer.WriteHeader(http.StatusBadRequest)
		errStr := fmt.Sprintf("Paper id %s is unkown.", paperStatus.Id)
		log.Println(errStr)
		if _, err := writer.Write([]byte(errStr)); err != nil {
			log.Println("Writing to http response writer failed. Error:", err)
		}
		return
	}

	paperStatusJson, err := json.Marshal(paperStatus)
	if err != nil {
		errStr := fmt.Sprintf("session status json marshal failed. Error: %s", err)
		log.Println(errStr)
		http.Error(writer, errStr, http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	if paperStatus.Status == common.SuccessStatus {
		writer.WriteHeader(http.StatusOK)
	} else {
		writer.WriteHeader(http.StatusAccepted)
	}
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

	var localFilePath string
	if formKey == paperFileFormKey {
		localFilePath = filepath.Join(common.LocalStoragePath, uploadId, common.PaperPdfFileName)
	} else if formKey == reviewFileFormKey {
		localFilePath = filepath.Join(common.LocalStoragePath, uploadId, common.ReviewPdfFileName)
	}

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
