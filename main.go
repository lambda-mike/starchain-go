package main

import (
	"github.com/starchain/api"
	"github.com/starchain/contracts"
	"log"
	"net/http"
)

func main() {
	var blockchain contracts.Blockchain
	log.Println("Hello StarchainGo!")
	api.RegisterHandlers(&blockchain)
	http.ListenAndServe(":8000", nil)
}
