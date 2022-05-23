package central

import (
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage"
	"log"
	"time"
)

const sizeOfProcessingQueue = 1024

type Worker struct {
	dbManager              *storage.DatabaseManager
	centralCh              chan common.UploadedPaper
	processingPapersStatus map[string]string // key - paper id; value status
}

func NewWorker() *Worker {
	return &Worker{
		dbManager:              storage.NewDatabaseManager(),
		centralCh:              make(chan common.UploadedPaper, sizeOfProcessingQueue),
		processingPapersStatus: map[string]string{},
	}
}

func (w *Worker) EnterMainLoop() {
	for newPaper := range w.centralCh {
		log.Println("New paper was taken in processing")
		w.processingPapersStatus[newPaper.Id] = common.InProgressStatus
		err := w.dbManager.AddNewPaper(newPaper)
		if err != nil {
			log.Println(err)
		}

		//TODO: section below is just for web testing
		{
			time.Sleep(10 * time.Second)
			w.processingPapersStatus[newPaper.Id] = common.SuccessStatus
		}
	}
}

func (w *Worker) GetSessionStatus(paperId string) common.PaperProcessResult {
	res := common.PaperProcessResult{
		Id:  paperId,
		NFT: "test_nft",
	}
	status, exist := w.processingPapersStatus[paperId]
	if !exist {
		res.Status = common.UnknownStatus
	} else {
		res.Status = status
	}
	return res
}

func (w *Worker) Stop() {
	close(w.centralCh)
	w.dbManager.Disconnect()
}

func (w *Worker) AddNewPaperToQueue(newPaper common.UploadedPaper) {
	w.centralCh <- newPaper
	w.processingPapersStatus[newPaper.Id] = common.IsReadyForProcessingStatus
	log.Printf("New paper added to processing queue. Paper id: %s", newPaper.Id)
}
