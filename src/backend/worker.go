package backend

import (
	"fmt"
	"net/http"
)

const (
	mainPageEndpoint    = "/main"
	uploadPaperEndpoint = "/upload-paper"
)

type Worker struct {
}

func NewWorker() *Worker {
	worker := &Worker{}
	http.Handle("/", http.FileServer(http.Dir("./web/static")))
	http.HandleFunc(uploadPaperEndpoint, worker.HandlePaperUploadRequest)
	return worker
}

func (w *Worker) HandlePaperUploadRequest(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		fmt.Fprint(writer, err)
	}
}

func (w *Worker) ListenHttp() error {
	return http.ListenAndServe(":8888", nil)
}
