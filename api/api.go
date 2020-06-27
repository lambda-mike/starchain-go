// api package provides handlers for http REST API
// Its goal is to provide an interface for the Blockchain
package api

import (
	"encoding/json"
	"fmt"
	"github.com/starchain/contracts"
	"log"
	"net/http"
)

type Address = struct {
	Address string `json:"address"`
}

var blockchain *contracts.Blockchain

func RegisterHandlers(b *contracts.Blockchain) {
	blockchain = b
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/requestValidation", requestValidation)
	log.Println("INFO: Handlers registered successfully")
}

func hello(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: hello")
	fmt.Fprint(res, "hello")
}

func requestValidation(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: requestValidation")
	if req.Body == nil {
		log.Println("ERR: requestValidation: request body is nil")
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Error occurred when decoding address from JSON: empty body")
		return
	}
	var addr Address
	if err := json.NewDecoder(req.Body).Decode(&addr); err != nil {
		log.Println("ERR: requestValidation: ", err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Error occurred when decoding address from JSON")
		return
	}
	if addr.Address == "" {
		log.Println("ERR: requestValidation: empty address field")
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "address is required")
		return
	}
	msg, err := (*blockchain).RequestMessageOwnershipVerification(addr.Address)
	if err != nil {
		log.Println("ERR: requestValidation: ", err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Error occurred when calling blockchain for the validation msg")
		return
	}
	fmt.Fprint(res, msg)
}
