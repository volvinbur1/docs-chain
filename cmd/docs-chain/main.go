package main

import (
	"github.com/volvinbur1/docs-chain/src/backend"
	"log"
)

func main() {
	worker := backend.NewWorker()
	log.Fatal(worker.ListenHttp())
}
