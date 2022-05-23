package main

import (
	"github.com/volvinbur1/docs-chain/src/backend"
	"github.com/volvinbur1/docs-chain/src/central"
	"log"
)

func main() {
	centralWorker := central.NewWorker()
	defer centralWorker.Stop()
	go centralWorker.EnterMainLoop()

	webUIProcessor := backend.NewWebUIProcessor(centralWorker)
	log.Fatal(webUIProcessor.ListenHttp())
}
