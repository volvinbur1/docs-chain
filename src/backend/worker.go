package backend

import "net/http"

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) HandleUserLogin() {

}

func (w *Worker) ListenHttp() {
	http.ListenAndServe("*:8080", nil)
}
