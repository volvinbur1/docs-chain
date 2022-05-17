package backend

import (
	"fmt"
	"log"
	"net/http"
)

const (
	paperUploadPageEndpoint = "/paper-upload"
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
		err := request.ParseForm()
		if err != nil {
			log.Print(fmt.Fprint(writer, err))
		}
	}
}

func (w *Worker) ListenHttp() error {
	return http.ListenAndServe(":8888", nil)
}
