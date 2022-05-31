package analyzer

import (
	"code.sajari.com/docconv"
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/jdkato/prose/v2"
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage"
	"hash/fnv"
	"log"
	"regexp"
	"strings"
)

type PaperPdfProcessor struct {
	paperId         string
	filePath        string
	canonizedText   string
	paperShingles   common.PaperShingles
	database        storage.DatabaseInterface
	compareResultCh chan CompareResult
	dispatcher      *Dispatcher
}

func NewPaperPdfProcessor(newPaper common.UploadedPaper, database storage.DatabaseInterface, dispatcher *Dispatcher) *PaperPdfProcessor {
	return &PaperPdfProcessor{
		paperId:         newPaper.Id,
		filePath:        newPaper.FilePath,
		database:        database,
		compareResultCh: make(chan CompareResult, workersCount),
		dispatcher:      dispatcher,
	}
}

func (p *PaperPdfProcessor) PrepareFile() error {
	readResult, err := docconv.ConvertPath(p.filePath)
	if err != nil {
		return fmt.Errorf("error when reading pdf for %s. Error: %s", p.paperId, err)
	}
	log.Printf("Paper %s plain text read successfully.", p.paperId)

	paperPlainText, err := removePunctuation(readResult.Body)
	if err != nil {
		return err
	}
	log.Printf("Paper %s punctuation from text has been removed.", p.paperId)

	paperPlainText, err = performPosTaggingAnalyze(paperPlainText)
	if err != nil {
		return err
	}
	log.Printf("Paper %s text pos-tagging analyze performed successfully.", p.paperId)

	p.canonizedText = removePlural(paperPlainText)
	return nil
}

func (p *PaperPdfProcessor) MakeShingles() error {
	if strings.ContainsAny(p.canonizedText, specialChars) {
		return fmt.Errorf("canonized text for paper %s is not preapared", p.paperId)
	}

	words := strings.Fields(p.canonizedText)
	var shinglesList []uint32
	for idx := 0; idx < len(words)-shinglesCount; idx++ {
		shingle := strings.Join(words[idx:idx+shinglesCount], "")
		fnv32a := fnv.New32a()
		_, err := fnv32a.Write([]byte(shingle))
		if err != nil {
			return fmt.Errorf("hashing shingle for %s failed with error: %s", p.paperId, err)
		}
		shinglesList = append(shinglesList, fnv32a.Sum32())
	}

	shinglesList = removeDuplicate(shinglesList)
	p.paperShingles = common.PaperShingles{
		Id:                p.paperId,
		Shingles:          shinglesList,
		WordsInShingleCnt: shinglesCount,
		HashAlgorithm:     shingleHashAlgorithm,
	}
	log.Printf("Shingles hashes for paper %s has been created.", p.paperId)

	return p.database.AddPaperShingles(p.paperShingles)
}

func (p *PaperPdfProcessor) PerformAnalyze() (common.AnalysisResult, error) {
	log.Printf("Paper %s analysis started...", p.paperId)

	papersShinglesList, err := p.database.GetAllPapersShingles()
	if err != nil {
		return common.AnalysisResult{}, err
	}
	if len(papersShinglesList) == 0 {
		return common.AnalysisResult{}, fmt.Errorf("comparison dataset is empty")
	}

	isCurrentPaperInDb := false
	taskQueue := p.dispatcher.GetTaskQueue()
	for _, paperShingles := range papersShinglesList {
		if paperShingles.Id == p.paperId {
			isCurrentPaperInDb = true
			continue
		}
		taskQueue <- DocTask{
			TargetPaperShingles: paperShingles,
			Comparator:          p,
		}
	}

	return p.calculateAnalysisResult(len(papersShinglesList) - boolToInt[isCurrentPaperInDb])
}

func (p *PaperPdfProcessor) CompareToDoc(targetPaperShingles common.PaperShingles) {
	mapB := make(map[uint32]interface{}, len(targetPaperShingles.Shingles))
	for _, x := range targetPaperShingles.Shingles {
		mapB[x] = nil
	}

	diffCnt := len(targetPaperShingles.Shingles)
	for _, x := range p.paperShingles.Shingles {
		if _, exist := mapB[x]; !exist {
			diffCnt++
		}
	}

	similarityCnt := len(targetPaperShingles.Shingles) + len(p.paperShingles.Shingles) - diffCnt

	p.compareResultCh <- CompareResult{
		TargetPaperId:  targetPaperShingles.Id,
		SimilarityRate: float64(similarityCnt) / float64(diffCnt),
	}
}

func (p *PaperPdfProcessor) calculateAnalysisResult(papersToCompareCnt int) (common.AnalysisResult, error) {
	analysisResult := common.AnalysisResult{Id: p.paperId}
	minUniqueness := 100.0
	resRetrievedCnt := 0
	for resRetrievedCnt < papersToCompareCnt {
		res, isOkay := <-p.compareResultCh
		if !isOkay {
			return common.AnalysisResult{}, fmt.Errorf("compare result for paper %s retrieving from channel failed", p.paperId)
		}
		resRetrievedCnt++

		uniqueness := (1 - res.SimilarityRate) * 100
		if uniqueness < minUniqueness {
			minUniqueness = uniqueness
		}

		if uniqueness < UniquenessThresholdValue {
			analysisResult.SimilarPapersId = append(analysisResult.SimilarPapersId, res.TargetPaperId)
		}
	}
	analysisResult.Uniqueness = minUniqueness
	log.Printf("Paper %s analysis finiched. Uniqueness: %.2f", p.paperId, analysisResult.Uniqueness)
	return analysisResult, nil
}

func removePunctuation(text string) (string, error) {
	text = strings.ToLower(text)

	reg, err := regexp.Compile(`[^a-zA-Z\d ]+`)
	if err != nil {
		return "", fmt.Errorf("alpha-numerical regural expression instance creation failed. Error %s", err)
	}
	newText := reg.ReplaceAllString(text, "")

	reg, err = regexp.Compile(`\s+`)
	if err != nil {
		return "", fmt.Errorf("space-removal regural expression instance creation failed. Error %s", err)
	}
	newText = reg.ReplaceAllString(newText, " ")

	return newText, nil
}

func performPosTaggingAnalyze(text string) (string, error) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return "", fmt.Errorf("prose document instance from paper pdf plain text. Error: %s", err)
	}

	docTokens := make([]prose.Token, len(doc.Tokens()))
	copy(docTokens, doc.Tokens())

	for idx := 0; idx < len(docTokens); idx++ {
		for _, removeTag := range tagsToRemoveList {
			if docTokens[idx].Tag == removeTag {
				docTokens = append(docTokens[:idx], docTokens[idx+1:]...)
				idx--
				break
			}
		}
	}

	var sb strings.Builder
	for _, entity := range docTokens {
		sb.WriteString(entity.Text)
		sb.WriteString(" ")
	}
	return sb.String(), nil
}

func removePlural(text string) string {
	words := strings.Fields(text)
	singleMaker := pluralize.NewClient()
	for idx, word := range words {
		words[idx] = singleMaker.Singular(word)
	}

	return strings.Join(words, " ")
}

func removeDuplicate(slice []uint32) []uint32 {
	allKeys := make(map[uint32]interface{})
	var list []uint32
	for _, item := range slice {
		if _, exist := allKeys[item]; !exist {
			allKeys[item] = nil
			list = append(list, item)
		}
	}
	return list
}
