package analyzer

import (
	"sync"
)

type Worker struct {
	Id int

	WorkerPool  chan chan DocTask
	TaskChannel chan DocTask
	quit        chan bool

	processedTaskMutex sync.Mutex
	processedTaskCount int
}

func NewWorker(workerId int, workerPool chan chan DocTask) *Worker {
	return &Worker{
		Id:          workerId,
		WorkerPool:  workerPool,
		TaskChannel: make(chan DocTask),
		quit:        make(chan bool),
	}
}

func (w *Worker) GetTotallyProcessedTasks() int {
	w.processedTaskMutex.Lock()
	defer w.processedTaskMutex.Unlock()
	return w.processedTaskCount
}

func (w *Worker) Start() {
	//fmt.Println("Worker ", w.Id, " started")
	go w.threadController()
}

func (w *Worker) Stop() {
	//fmt.Println("Worker ", w.Id, " stopping...")
	w.quit <- true
	//fmt.Println("Worker ", w.Id, " stopped")
}

func (w *Worker) threadController() {
	for {
		w.WorkerPool <- w.TaskChannel
		select {
		case task := <-w.TaskChannel:
			task.Comparator.CompareToDoc(task.TargetPaperShingles)

			w.processedTaskMutex.Lock()
			w.processedTaskCount++
			w.processedTaskMutex.Unlock()
		case <-w.quit:
			return
		}
	}
}
