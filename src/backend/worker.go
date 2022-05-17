package backend

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	paperUploadPageEndpoint = "/paper-upload"
)

const (
	paperFileFormKey = "paperFile"
)

type Worker struct {
}

func NewWorker() *Worker {
	worker := &Worker{}
	http.Handle("/", http.FileServer(http.Dir("web/static")))
	http.HandleFunc(paperUploadPageEndpoint, worker.HandlePaperUploadRequest)
	return worker
}

func (w *Worker) HandlePaperUploadRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		http.ServeFile(writer, request, "web/static/paper_upload_page.html")
	} else {
		if err := request.ParseMultipartForm(32 << 20); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		requestFile, header, err := request.FormFile(paperFileFormKey)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer requestFile.Close()

		fileStoragePath := filepath.Join("bin", "storage", "paper")
		if err = os.MkdirAll(fileStoragePath, os.ModePerm); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		localFile, err := os.OpenFile(filepath.Join(fileStoragePath, header.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer localFile.Close()
		_, err = io.Copy(localFile, requestFile)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (w *Worker) ListenHttp() error {
	return http.ListenAndServe(":8888", nil)
}
