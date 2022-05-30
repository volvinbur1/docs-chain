package main

import (
	"fmt"
	"github.com/volvinbur1/docs-chain/src/blockchain"
)

func main() {
	//centralWorker := central.NewWorker()
	//defer centralWorker.Stop()
	//go centralWorker.EnterMainLoop()
	//
	//webUIProcessor := backend.NewWebUIProcessor(centralWorker)
	//log.Fatal(webUIProcessor.ListenHttp())
	Url, err := blockchain.GetImageUrl("nebin/storage/dsdsdsdsfdde3r3e3/QmapHwRHhzJPtX79oXit3gNs4gAmF3LX8FNTjNGiR1syrY.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(Url)
}
