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
	paperUploadPageEndpoint          = "/paper-upload"
	paperUploadSessionStatusEndpoint = paperUploadPageEndpoint + "/status"
)

const (
	paperTopicFormKey   = "paperTopic"
	uploaderNameFormKey = "uploaderName"
	reviewFileFormKey   = "reviewFile"
	paperFileFormKey    = "paperFile"
	sessionIdKey        = "sessionId"
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
	http.HandleFunc(paperUploadSessionStatusEndpoint, processor.HandlePaperUploadSessionStatus)
	return processor
}

func (w *WebUIProcessor) HandlePaperUploadRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		http.ServeFile(writer, request, "web/static/paper_upload_page.html")
	case http.MethodPost:
		paper, err := w.parsePaperUploadRequest(request)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		w.centralWorker.ProcessNewPaper(paper)
		idJson, err := json.Marshal(paper.Id)
		if err != nil {
			http.Error(writer, fmt.Sprint("Session id json marshal failed. Error: ", err), http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(idJson)
		if err != nil {
			http.Error(writer, fmt.Sprint("Writing to http response writer failed. Error: ", err), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusAccepted)
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *WebUIProcessor) HandlePaperUploadSessionStatus(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		sessionId := request.URL.Query().Get(sessionIdKey)
		sessionStatus := w.centralWorker.GetSessionStatus(sessionId)
		sessionStatusJson, err := json.Marshal(sessionStatus)
		if err != nil {
			http.Error(writer, fmt.Sprint("Session status json marshal failed. Error: ", err), http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(sessionStatusJson)
		if err != nil {
			http.Error(writer, fmt.Sprint("Writing to http response writer failed. Error: ", err), http.StatusInternalServerError)
			return
		}
		if sessionStatus.Status == common.SuccessSessionStatus {
			writer.WriteHeader(http.StatusOK)
		} else if sessionStatus.Status == common.InProgressSessionStatus {
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusBadRequest)
		}
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
