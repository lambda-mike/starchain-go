package main

import (
	"github.com/starchain/api"
	"github.com/starchain/blockchain"
	"github.com/starchain/contracts"
	"github.com/starchain/proxy"
	"log"
	"net/http"
)

func main() {
	log.Println("Hello StarchainGo!")
	var (
		bchain          *blockchain.Blockchain
		clock           contracts.Clock
		blockchainProxy contracts.BlockchainOperator
	)
	clock = blockchain.BlockchainClock{}
	bchain = blockchain.New(clock)
	blockchainProxy = proxy.New(bchain)
	restApi := api.Create(&blockchainProxy)
	http.ListenAndServe(":8000", restApi)
}
