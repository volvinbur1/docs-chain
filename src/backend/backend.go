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

const localStoragePath = "bin/storage"

var closerHandler = func(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Print(err)
	}
}

type WebUIProcessor struct {
	centralWorker *central.Worker
}

func NewWebUIProcessor(centralWorker *central.Worker) *WebUIProcessor {
	processor := &WebUIProcessor{centralWorker: centralWorker}

	http.Handle("/", http.FileServer(http.Dir("web/static")))
	http.HandleFunc(paperUploadPageEndpoint, processor.HandlePaperUploadRequest)
	http.HandleFunc(paperUploadStatusEndpoint, processor.HandlePaperUploadStatus)
	return processor
}

func (w *WebUIProcessor) HandlePaperUploadRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		http.ServeFile(writer, request, "web/static/paper_upload_page.html")
	case http.MethodPost:
		newPaper, err := w.parsePaperUploadRequest(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		w.addPaperToProcessingQueue(newPaper, writer)
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) HandlePaperUploadStatus(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		paperId := request.URL.Query().Get(paperIdKey)
		responseStatus, err := w.checkPaperStatus(paperId, writer)
		if err != nil {
			http.Error(writer, err.Error(), responseStatus)
			return
		}

		writer.WriteHeader(responseStatus)
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
	uploadedPaper.CreatorName = request.Form.Get(uploaderNameFormKey)
	uploadedPaper.PaperPath, err = storeFileFromRequest(request, uploadedPaper.Id, paperFileFormKey)
	if err != nil {
		return common.UploadedPaper{}, err
	}
	uploadedPaper.ReviewPath, err = storeFileFromRequest(request, uploadedPaper.Id, reviewFileFormKey)
	return uploadedPaper, err
}

func (w *WebUIProcessor) addPaperToProcessingQueue(newPaper common.UploadedPaper, writer http.ResponseWriter) {
	w.centralWorker.AddNewPaperToQueue(newPaper)

	responseStatus, err := w.checkPaperStatus(newPaper.Id, writer)
	if err != nil {
		http.Error(writer, err.Error(), responseStatus)
		return
	}

	if responseStatus != http.StatusOK {
		writer.WriteHeader(http.StatusAccepted)
	} else {
		writer.WriteHeader(responseStatus)
	}
}

func (w *WebUIProcessor) checkPaperStatus(paperId string, writer http.ResponseWriter) (int, error) {
	paperStatus := w.centralWorker.GetSessionStatus(paperId)
	paperStatusJson, err := json.Marshal(paperStatus)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("session status json marshal failed. Error: %s", err)
	}
	_, err = writer.Write(paperStatusJson)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("writing to http response writer failed. Error: %s", err)
	}

	if paperStatus.Status == common.SuccessStatus {
		return http.StatusOK, nil
	}
	if paperStatus.Status == common.InProgressStatus {
		return http.StatusNoContent, nil
	}
	return http.StatusBadRequest, nil
}

func storeFileFromRequest(request *http.Request, uploadId, formKey string) (string, error) {
	requestFile, _, err := request.FormFile(formKey)
	if err != nil {
		return "", fmt.Errorf("getting from form failed: %s", err)
	}
	defer closerHandler(requestFile)

	if err = os.MkdirAll(filepath.Join(localStoragePath, uploadId), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create all file storega path subdirs: %s", err)
	}

	var localFilePath string
	if formKey == paperFileFormKey {
		localFilePath = filepath.Join(localStoragePath, uploadId, "paper.pdf")
	} else if formKey == reviewFileFormKey {
		localFilePath = filepath.Join(localStoragePath, uploadId, "review.pdf")
	}

	localFile, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s Error: %s", localFilePath, err)
	}
	defer closerHandler(localFile)

	_, err = io.Copy(localFile, requestFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy file from request to local one: %s", err)
	}

	return localFilePath, nil
}
