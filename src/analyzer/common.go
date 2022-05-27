package analyzer

import "github.com/volvinbur1/docs-chain/src/common"

var tagsToRemoveList = [...]string{"CC", "DT", "EX", "IN", "PDT", "TO", "UH", "WDT", "WP", "WP$", "WRB"}

const specialChars = ".,!?@#$&+-*/=^%~(){}[]<>'`|\"\\"
const shinglesCount = 7
const shingleHashAlgorithm = "fnv32a"

const taskQueueSize = 1024
const workersCount = 4

type DocsComparator interface {
	CompareToDoc(targetPaperShingles common.PaperShingles)
}

type DocTask struct {
	TargetPaperShingles common.PaperShingles
	Comparator          DocsComparator
}

type CompareResult struct {
	TargetPaperId  string
	SimilarityRate float64
}
