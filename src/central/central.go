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
	centralCh              chan common.ServiceTask
	processingPapersStatus sync.Map // key - paper id; value status
}

func NewWorker() *Worker {
	return &Worker{
		dbManager:              storage.NewDatabaseManager(),
		centralCh:              make(chan common.ServiceTask, sizeOfProcessingQueue),
		processingPapersStatus: sync.Map{},
	}
}

func (w *Worker) EnterMainLoop() {
	for task := range w.centralCh {
		switch task.Action {
		case common.NewPaperUploadAction:
			newPaper, isOkay := task.Payload.(common.UploadedPaper)
			if !isOkay {
				log.Println("Interface to new paper converting failed.", task.Payload)
				continue
			}
			go w.handleNewPaperUpload(newPaper)
		case common.GetPaperProcessingStatusAction:
			paperId, isOkay := task.Payload.(string)
			if !isOkay {
				log.Println("Interface to string converting failed.", task.Payload)
			}
			go w.handlePaperStatusRequest(paperId, task.ReturnCh)
		default:
			log.Println("Unknown action:", task.Action)
		}
	}
}

func (w *Worker) Stop() {
	close(w.centralCh)
	w.dbManager.Disconnect()
}

func (w *Worker) NewPaperAction(newPaper common.UploadedPaper) {
	w.centralCh <- common.ServiceTask{
		Action:  common.NewPaperUploadAction,
		Payload: newPaper,
	}
	w.processingPapersStatus.Store(newPaper.Id, common.IsReadyForProcessingStatus)
	log.Printf("New paper added to processing queue. Paper id: %s", newPaper.Id)
}

func (w *Worker) GetPaperStatusAction(paperId string, returnCh chan<- interface{}) {
	w.centralCh <- common.ServiceTask{
		Action:   common.GetPaperProcessingStatusAction,
		Payload:  paperId,
		ReturnCh: returnCh,
	}
}

func (w *Worker) handlePaperStatusRequest(paperId string, returnCh chan<- interface{}) {
	res := common.PaperProcessingResult{
		Id:  paperId,
		NFT: "test_nft", //TODO: replace; this one just for tests
	}
	status, exist := w.processingPapersStatus.Load(paperId)
	if !exist {
		res.Status = common.UnknownStatus
	} else {
		res.Status = status.(string)
	}

	returnCh <- res
}

func (w *Worker) handleNewPaperUpload(newPaper common.UploadedPaper) {
	log.Printf("New paper %s processing started.", newPaper.Id)
	w.processingPapersStatus.Store(newPaper.Id, common.InProgressStatus)

	if err := w.dbManager.AddNewPaper(newPaper); err != nil {
		log.Println(err)
		w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
		return
	}

	pdfProcessor := analyzer.NewPaperPdfProcessor(newPaper.PaperPath, w.dbManager)
	if err := pdfProcessor.PrepareFile(newPaper.Id); err != nil {
		log.Println(err)
		w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
		return
	}
	if err := pdfProcessor.PrepareFile(newPaper.Id); err != nil {
		log.Println(err)
		w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
		return
	}
	w.processingPapersStatus.Store(newPaper.Id, common.SuccessStatus)
}
