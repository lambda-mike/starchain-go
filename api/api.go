// api package provides handlers for http REST API
// Its goal is to provide an interface for the Blockchain
package api

import (
	"fmt"
	"github.com/starchain/contracts"
	"log"
	"net/http"
)

var blockchain *contracts.Blockchain

func RegisterHandlers(blockchain *contracts.Blockchain) {
	blockchain = blockchain
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/requestValidation", requestValidation)
	log.Println("INFO: Handlers registered successfully")
}

func hello(res http.ResponseWriter, req *http.Request) {
	log.Println("hello")
	fmt.Fprint(res, "hello")
}

func requestValidation(res http.ResponseWriter, req *http.Request) {
	log.Println("requestValidation")
	fmt.Fprint(res, "TODO ...")
}
