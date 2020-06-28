// api package provides handlers for http REST API
// Its goal is to provide an interface for the Blockchain
package api

import (
	"encoding/json"
	"fmt"
	"github.com/starchain/contracts"
	"log"
	"net/http"
	"regexp"
)

type BlockDto struct {
	Body              string `json:"body"`
	Hash              string `json:"hash"`
	Height            int    `json:"height"`
	Owner             string `json:"owner"`
	PreviousBlockHash string `json:"previousBlockHash"`
	Time              int64  `json:"time"`
}

// TODO postfix Dto
type Address struct {
	Address string `json:"address"`
}

var blockchain *contracts.Blockchain

func Create(b *contracts.Blockchain) http.Handler {
	blockchain = b
	api := newRestApi()
	api.Add("GET /hello", hello)
	api.Add("POST /requestValidation", requestValidation)
	// TODO /block
	//log.Println("INFO: Handlers registered successfully")
	log.Println("INFO: REST API created successfully")
	return api
}

type restApi struct {
	handlers map[string]http.HandlerFunc
	cache    map[string]*regexp.Regexp
}

func newRestApi() *restApi {
	return &restApi{
		handlers: make(map[string]http.HandlerFunc),
		cache:    make(map[string]*regexp.Regexp),
	}
}

func (a *restApi) Add(regex string, handler http.HandlerFunc) {
	a.handlers[regex] = handler
	compiled, err := regexp.Compile(regex)
	if err != nil {
		log.Panicln("Could not compile regexp: ", regex)
	}
	a.cache[regex] = compiled
}

func (a *restApi) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	pattern := req.Method + " " + req.URL.Path
	log.Println("DBG: ServeHTTP: ", pattern)
	for regex, handler := range a.handlers {
		if a.cache[regex].MatchString(pattern) {
			handler(res, req)
			return
		}
	}
	http.NotFound(res, req)
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

func getBlock(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: getBlock")
	log.Println("DBG: getBlock:url", req.URL)
	/*if err := json.NewDecoder(req.Body).Decode(&addr); err != nil {
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
	*/
	fmt.Fprint(res, "TODO")
}
