package main

import (
	"github.com/starchain/api"
	"github.com/starchain/contracts"
	"log"
	"net/http"
)

func main() {
	log.Println("Hello StarchainGo!")
	// TODO assign proper implementation!!
	var blockchain contracts.Blockchain
	restApi := api.Create(&blockchain)
	http.ListenAndServe(":8000", restApi)
}
