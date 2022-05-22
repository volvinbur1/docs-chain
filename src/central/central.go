package central

import (
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage"
	"log"
)

const sizeOfProcessingQueue = 1024

type Worker struct {
	dbManager *storage.DatabaseManager
	centralCh chan common.UploadedPaper
}

func NewWorker() *Worker {
	return &Worker{
		dbManager: storage.NewDatabaseManager(),
		centralCh: make(chan common.UploadedPaper, sizeOfProcessingQueue),
	}
}

func (w *Worker) EnterMainLoop() {
	for newPaper := range w.centralCh {
		log.Println("New paper took in processing")
		err := w.dbManager.AddNewPaper(newPaper)
		if err != nil {
			log.Println(err)
		}
	}
}

func (w *Worker) GetSessionStatus(sessionId string) common.PaperProcessResult {
	return common.PaperProcessResult{
		Id:     sessionId,
		Status: common.InProgressStatus,
		NFT:    "test_nft",
	}
}

func (w *Worker) Stop() {
	close(w.centralCh)
	w.dbManager.Disconnect()
}

func (w *Worker) AddNewPaperToQueue(newPaper common.UploadedPaper) {
	w.centralCh <- newPaper
	log.Printf("New paper added to processing queue. Paper id: %s", newPaper.Id)
}
