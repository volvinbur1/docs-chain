package central

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/ledongthuc/pdf"
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
		log.Println("New paper was taken in processing")
		w.processingPapersStatus.Store(newPaper.Id, common.InProgressStatus)
		err := w.dbManager.AddNewPaper(newPaper)
		if err != nil {
			log.Println(err)
			w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
			continue
		}

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

// readPaperPdf reads a paper pdf plain text starting from 5 and to (n-2) pages
func readPaperPdf(path string) (string, error) {
	file, pdfReader, err := pdf.Open(path)
	defer common.CloserHandler(file)
	if err != nil {
		return "", fmt.Errorf("%s file oped failed. Error: %s", path, err)
	}

	var buffer bytes.Buffer
	for pageNumber := 5; pageNumber < pdfReader.NumPage()-2; pageNumber++ {
		page := pdfReader.Page(pageNumber)
		if page.V.IsNull() {
			log.Printf("Page %d from %s reading failed.", pageNumber, path)
			continue
		}

		plainTextStr, err := page.GetPlainText(nil)
		if err != nil {
			log.Printf("Getting plain text from page %d from %s failed. Error: %s", pageNumber, path, err)
			continue
		}
		buffer.WriteString(plainTextStr)
	}

	return base64.StdEncoding.EncodeToString([]byte(buffer.String())), nil
}
