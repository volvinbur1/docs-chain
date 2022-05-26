package central

import (
	"github.com/volvinbur1/docs-chain/src/analyzer"
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage"
	"log"
	"sync"
)

const sizeOfProcessingQueue = 1024

type Worker struct {
	dbManager              *storage.DatabaseManager
	centralCh              chan common.UploadedPaper
	processingPapersStatus sync.Map // key - paper id; value status
}

func NewWorker() *Worker {
	return &Worker{
		dbManager:              storage.NewDatabaseManager(),
		centralCh:              make(chan common.UploadedPaper, sizeOfProcessingQueue),
		processingPapersStatus: sync.Map{},
	}
}

func (w *Worker) EnterMainLoop() {
	for newPaper := range w.centralCh {
		log.Printf("New paper %s processing started.", newPaper.Id)
		w.processingPapersStatus.Store(newPaper.Id, common.InProgressStatus)
		err := w.dbManager.AddNewPaper(newPaper)
		if err != nil {
			log.Println(err)
			w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
			continue
		}

		_ = analyzer.NewPaperPdfProcessor(newPaper.PaperPath, w.dbManager)
		//w.processingPapersStatus.Store(newPaper.Id, common.SuccessStatus)
	}
}

func (w *Worker) GetSessionStatus(paperId string) common.PaperProcessResult {
	res := common.PaperProcessResult{
		Id:  paperId,
		NFT: "test_nft",
	}
	status, exist := w.processingPapersStatus.Load(paperId)
	if !exist {
		res.Status = common.UnknownStatus
	} else {
		res.Status = status.(string)
	}
	return res
}

func (w *Worker) Stop() {
	close(w.centralCh)
	w.dbManager.Disconnect()
}

func (w *Worker) AddNewPaperToQueue(newPaper common.UploadedPaper) {
	w.centralCh <- newPaper
	w.processingPapersStatus.Store(newPaper.Id, common.IsReadyForProcessingStatus)
	log.Printf("New paper added to processing queue. Paper id: %s", newPaper.Id)
}
