package analyzer

type Dispatcher struct {
	workerPool chan chan DocTask
	taskQueue  chan DocTask
	quit       chan bool

	workersList []*Worker
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		taskQueue:  make(chan DocTask, taskQueueSize),
		workerPool: make(chan chan DocTask, workersCount),
		quit:       make(chan bool),
	}
}

func (d *Dispatcher) GetProcessedTaskCnt() int {
	totallyProcessedTaskCnt := 0
	for _, worker := range d.workersList {
		totallyProcessedTaskCnt += worker.GetTotallyProcessedTasks()
	}
	return totallyProcessedTaskCnt
}

func (d *Dispatcher) Run() {
	//fmt.Println("Dispatcher started.")
	for i := 0; i < workersCount; i++ {
		worker := NewWorker(i, d.workerPool)
		worker.Start()
		d.workersList = append(d.workersList, worker)
	}

	go d.dispatch()
}

func (d *Dispatcher) Stop() {
	//fmt.Println("Dispatcher stopping...")
	d.quit <- true
	for _, worker := range d.workersList {
		worker.Stop()
	}
	//fmt.Println("Dispatcher stopped.")
}

func (d *Dispatcher) GetTaskQueue() chan DocTask {
	return d.taskQueue
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case taskCh := <-d.workerPool:
			task := <-d.taskQueue
			taskCh <- task
		case <-d.quit:
			return
		}
	}
}
