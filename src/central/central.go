package central

import (
	"fmt"
	"github.com/volvinbur1/docs-chain/src/analyzer"
	"github.com/volvinbur1/docs-chain/src/blockchain"
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/ipfs"
	"github.com/volvinbur1/docs-chain/src/storage"
	"log"
	"sync"
	"time"
)

const sizeOfProcessingQueue = 1024

type Worker struct {
	dbManager              *storage.DatabaseManager
	centralCh              chan common.ServiceTask
	dispatcher             *analyzer.Dispatcher
	blockChain             *blockchain.BlockChain
	processingPapersStatus sync.Map // key - paper id; value status
}

func NewWorker() *Worker {
	return &Worker{
		dbManager:              storage.NewDatabaseManager(),
		centralCh:              make(chan common.ServiceTask, sizeOfProcessingQueue),
		dispatcher:             analyzer.NewDispatcher(),
		blockChain:             blockchain.New(),
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
		res.Status = status.(int)
	}

	returnCh <- res
}

func (w *Worker) handleNewPaperUpload(newPaper common.UploadedPaper) {
	log.Printf("New paper %s processing started.", newPaper.Id)
	w.processingPapersStatus.Store(newPaper.Id, common.InProgressStatus)

	analysisResult, err := w.analyzePaperPdf(newPaper)
	if err != nil {
		log.Println(err)
		w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
		return
	}

	if analysisResult.Uniqueness < analyzer.UniquenessThresholdValue {
		log.Println("Not enough uniqueness for paper ", newPaper.Id)
		w.processingPapersStatus.Store(newPaper.Id, common.NotEnoughUniquenessStatus)
		return
	}

	err = w.savePaperInSystem(newPaper, analysisResult)
	if err != nil {
		log.Println(err)
		w.processingPapersStatus.Store(newPaper.Id, common.ProcessingFailedStatus)
		return
	}
	w.processingPapersStatus.Store(newPaper.Id, common.SuccessStatus)
}

func (w *Worker) analyzePaperPdf(newPaper common.UploadedPaper) (common.AnalysisResult, error) {
	pdfProcessor := analyzer.NewPaperPdfProcessor(newPaper, w.dbManager, w.dispatcher)
	if err := pdfProcessor.PrepareFile(); err != nil {
		return common.AnalysisResult{}, err
	}
	if err := pdfProcessor.PrepareFile(); err != nil {
		return common.AnalysisResult{}, err
	}
	return pdfProcessor.PerformAnalyze()
}

func (w *Worker) savePaperInSystem(newPaper common.UploadedPaper, analysisResult common.AnalysisResult) error {
	qrCodePath, err := w.createPaperQrCode(newPaper, analysisResult)
	if err != nil {
		return err
	}

	nftImageUrl, err := blockchain.GetImageUrl(qrCodePath)
	if err != nil {
		return err
	}

	return w.createNft(newPaper, nftImageUrl)
}

func (w *Worker) createPaperQrCode(newPaper common.UploadedPaper, analysisResult common.AnalysisResult) (string, error) {
	ipfsCid, err := ipfs.AddFileToIpfs(newPaper.FilePath)
	if err != nil {
		return "", fmt.Errorf("adding paper %s pdf to ipfs failed. Error: %s", newPaper.Id, err)
	}

	var similarPapersNft []string
	for _, id := range analysisResult.SimilarPapersId {
		nft, err := w.dbManager.GetPaperNftById(id)
		if err != nil {
			log.Printf("Getting nft for paper %s failed. Error: %s", id, err)
			continue
		}
		similarPapersNft = append(similarPapersNft, nft)
	}

	paperMetadata := common.PaperMetadata{
		Id:               newPaper.Id,
		Topic:            newPaper.Topic,
		UploadDate:       time.Now().Format(time.RFC850),
		Authors:          newPaper.Authors,
		PaperIpfsHash:    ipfsCid,
		PaperUniqueness:  fmt.Sprintf("%.2f", analysisResult.Uniqueness),
		SimilarPapersNfr: similarPapersNft,
	}
	if err = w.dbManager.AddNewPaper(paperMetadata); err != nil {
		return "", err
	}

	return CreateNFTQRCode(paperMetadata)
}

func (w *Worker) createNft(newPaper common.UploadedPaper, nftImageUrl string) error {
	nftMetadata := blockchain.NftMetaData{
		Name:        newPaper.NftName,
		Symbol:      newPaper.NftSymbol,
		Description: fmt.Sprint(newPaper.Id, "|", newPaper.Topic),
		ImageUrl:    nftImageUrl,
	}

	nftResponse, err := w.blockChain.AddToSolana(nftMetadata)
	if err != nil {
		return err
	}

	nftResponse.Id = newPaper.Id
	return w.dbManager.AddPaperNft(nftResponse)
}
