package analyzer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/jdkato/prose/v2"
	"github.com/ledongthuc/pdf"
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage"
	"hash/fnv"
	"log"
	"regexp"
	"strings"
)

type PaperPdfProcessor struct {
	paperFilePath string
	canonizedText string
	dbManager     *storage.DatabaseManager
	compareResult chan CompareResult
	dispatcher    *Dispatcher
}

func NewPaperPdfProcessor(paperFilePath string, dbManager *storage.DatabaseManager, dispatcher *Dispatcher) *PaperPdfProcessor {
	return &PaperPdfProcessor{
		paperFilePath: paperFilePath,
		dbManager:     dbManager,
		compareResult: make(chan CompareResult, workersCount),
		dispatcher:    dispatcher,
	}
}

func (p *PaperPdfProcessor) PrepareFile(paperId string) error {
	log.Printf("Paper preanalyze processing started")
	paperPlainText, err := readPaperPdf(p.paperFilePath)
	if err != nil {
		return err
	}
	log.Printf("Paper %s plain text read successfully.", paperId)

	paperPlainText, err = removePunctuation(paperPlainText)
	if err != nil {
		return err
	}
	log.Printf("Paper %s punctuation from text has been removed.", paperId)

	paperPlainText, err = performPosTaggingAnalyze(paperPlainText)
	if err != nil {
		return err
	}
	log.Printf("Paper %s text pos-tagging analyze performed successfully.", paperId)

	p.canonizedText = removePlural(paperPlainText)
	return nil
}

func (p *PaperPdfProcessor) MakeShingles(paperId string) error {
	if strings.ContainsAny(p.canonizedText, specialChars) {
		return fmt.Errorf("canonized text for paper %s is not preapared", paperId)
	}

	words := strings.Fields(p.canonizedText)
	var shinglesList []uint32
	for idx := 0; idx < len(words)-shinglesCount; idx++ {
		shingle := strings.Join(words[idx:idx+shinglesCount], "")
		fnv32a := fnv.New32a()
		_, err := fnv32a.Write([]byte(shingle))
		if err != nil {
			return fmt.Errorf("hashing shingle for %s failed with error: %s", paperId, err)
		}
		shinglesList = append(shinglesList, fnv32a.Sum32())
	}

	log.Printf("Shingles hashes for paper %s has been created.", paperId)
	return p.dbManager.AddPaperShingles(common.PaperShingles{
		Id:                paperId,
		Shingles:          shinglesList,
		WordsInShingleCnt: shinglesCount,
		HashAlgorithm:     shingleHashAlgorithm,
	})
}

func (p *PaperPdfProcessor) PerformAnalyze() (common.AnalysisResult, error) {
	_, err := p.dbManager.GetAllPapersShingles()
	if err != nil {
		return common.AnalysisResult{}, err
	}

	return common.AnalysisResult{}, errors.New("not implemented")
}

func (p *PaperPdfProcessor) compareToOtherDoc(otherPaperShingles common.PaperShingles) {
	//TODO: implement
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

	return buffer.String(), nil
}

func removePunctuation(text string) (string, error) {
	text = strings.ToLower(text)

	reg, err := regexp.Compile("[^a-zA-Z\\d ]+ ")
	if err != nil {
		return "", fmt.Errorf("alpha-numerical regural expression instance creation failed. Error %s", err)
	}
	newText := reg.ReplaceAllString(text, "")

	reg, err = regexp.Compile(" +(?= )")
	if err != nil {
		return "", fmt.Errorf("space-removal regural expression instance creation failed. Error %s", err)
	}
	newText = reg.ReplaceAllString(newText, "")

	return newText, nil
}

func performPosTaggingAnalyze(text string) (string, error) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		return "", fmt.Errorf("prose document instance from paper pdf plain text. Error: %s", err)
	}

	docEntities := make([]prose.Entity, len(doc.Entities()))
	copy(docEntities, doc.Entities())

	for idx := 0; idx < len(docEntities); idx++ {
		for _, removeTag := range tagsToRemoveList {
			if docEntities[idx].Label == removeTag {
				docEntities = append(docEntities[:idx], docEntities[idx+1:]...)
				idx--
				break
			}
		}
	}

	var sb strings.Builder
	for _, entity := range docEntities {
		sb.WriteString(entity.Text)
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
