package central

import (
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage"
	"log"
)

type Worker struct {
	dbManager *storage.DatabaseManager
	centralCh chan common.UploadedPaper
}

func NewWorker() *Worker {
	return &Worker{
		dbManager: storage.NewDatabaseManager(),
		centralCh: make(chan common.UploadedPaper),
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

func (w *Worker) GetSessionStatus(sessionId string) common.ProcessingSession {
	return common.ProcessingSession{
		Id:     sessionId,
		Status: common.SuccessSessionStatus,
		NFT:    "test_nft",
	}
}

func (w *Worker) Stop() {
	close(w.centralCh)
	w.dbManager.Disconnect()
}

func (w *Worker) ProcessNewPaper(paper common.UploadedPaper) {
	w.centralCh <- paper
}
