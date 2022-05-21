package backend

import (
	"fmt"
	"github.com/rs/xid"
	"github.com/volvinbur1/docs-chain/src/common"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	paperUploadPageEndpoint = "/paper-upload"
)

const (
	paperTopicFormKey   = "paperTopic"
	uploaderNameFormKey = "uploaderName"
	reviewFileFormKey   = "reviewFile"
	paperFileFormKey    = "paperFile"
)

const localStoragePath = "bin/storage"

var closerHandler = func(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Print(err)
	}
}

type Worker struct{}

func NewWorker() *Worker {
	worker := &Worker{}
	http.Handle("/", http.FileServer(http.Dir("web/static")))
	http.HandleFunc(paperUploadPageEndpoint, worker.HandlePaperUploadRequest)
	return worker
}

func (w *Worker) HandlePaperUploadRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		http.ServeFile(writer, request, "web/static/paper_upload_page.html")
	case http.MethodPost:
		if err := w.parsePaperUploadRequest(request); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(writer, request, "/", http.StatusFound)
		}
	default:
		http.Error(writer, fmt.Sprintf("Http is method %s is not supported", request.Method), http.StatusNotImplemented)
	}
}

func (w *Worker) ListenHttp() error {
	return http.ListenAndServe(":8888", nil)
}

func (w *Worker) parsePaperUploadRequest(request *http.Request) error {
	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		return fmt.Errorf("request form parse error: %s", err)
	}

	var uploadedPaper common.UploadedPaper
	uploadedPaper.Id = xid.New().String()
	uploadedPaper.Topic = request.Form.Get(paperTopicFormKey)
	uploadedPaper.CreatorName = request.Form.Get(uploaderNameFormKey)
	uploadedPaper.PaperPath, err = storeFileFromRequest(request, uploadedPaper.Id, paperFileFormKey)
	if err != nil {
		return err
	}
	uploadedPaper.ReviewPath, err = storeFileFromRequest(request, uploadedPaper.Id, reviewFileFormKey)
	fmt.Printf("%+v", uploadedPaper)
	return err
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
