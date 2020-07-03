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
	"strconv"
	"strings"
)

type AddressDto struct {
	Address string `json:"address"`
}

type BlockDto struct {
	Body              string `json:"body"`
	Hash              string `json:"hash"`
	Height            int    `json:"height"`
	Owner             string `json:"owner"`
	PreviousBlockHash string `json:"previousBlockHash"`
	Time              int64  `json:"time"`
}

type StarDto struct {
	Address   string          `json:"address"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"star"`
	Signature string          `json:"signature"`
}

var blockchain *contracts.BlockchainOperator

func Create(b *contracts.BlockchainOperator) http.Handler {
	blockchain = b
	api := newRestApi()
	api.Add("GET /hello", hello)
	api.Add("GET /block/\\d+", getBlockByHeight)
	api.Add("GET /block/hash/\\w+", getBlockByHash)
	api.Add("GET /blocks/\\w+", getBlocks)
	api.Add("POST /requestValidation", requestValidation)
	api.Add("POST /submitStar", submitStar)
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
	var addr AddressDto
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
	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, msg)
}

func getBlockByHeight(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: getBlockByHeight")
	var parts []string
	if parts = strings.Split(req.URL.Path, "/"); len(parts) != 3 {
		log.Println("ERR: getBlockByHeight: wrong url format", req.URL.Path)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Could not fetch block by height: bad request URL")
		return
	}
	heightStr := parts[2]
	height, err := strconv.Atoi(heightStr)
	if err != nil {
		log.Println("ERR: getBlockByHeight: could not parse block height param: ", heightStr)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Could not parse block height param: "+heightStr)
		return
	}
	block, err := (*blockchain).GetBlockByHeight(height)
	respondWithBlock(res, req, &block, err)
}

func getBlockByHash(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: getBlockByHash")
	var parts []string
	if parts = strings.Split(req.URL.Path, "/"); len(parts) != 4 {
		log.Println("ERR: getBlockByHash: wrong url format", req.URL.Path)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Could not fetch block by hash: bad request URL")
		return
	}
	hash := parts[3]
	block, err := (*blockchain).GetBlockByHash(hash)
	respondWithBlock(res, req, &block, err)
}

func getBlocks(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: getBlocks")
	var parts []string
	if parts = strings.Split(req.URL.Path, "/"); len(parts) != 3 {
		log.Println("ERR: getBlockByHash: wrong url format", req.URL.Path)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Could not fetch block by hash: bad request URL")
		return
	}
	addr := parts[2]
	blocksData := (*blockchain).GetStarsByWalletAddress(addr)
	json, err := json.Marshal(blocksData)
	if err != nil {
		log.Println("ERR: getBlocks failed to marshal blocks: ", err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Failed to serialize blocks data into JSON")
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, string(json))
}

func submitStar(res http.ResponseWriter, req *http.Request) {
	log.Println("INFO: submitStar")
	if req.Body == nil {
		log.Println("ERR: submitStar: request body is nil")
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Error occurred when decoding star data: empty body")
		return
	}
	var starDto StarDto
	if err := json.NewDecoder(req.Body).Decode(&starDto); err != nil {
		log.Println("ERR: submitStar: ", err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Error occurred when decoding star from JSON")
		return
	}
	star := contracts.StarData{
		Address:   starDto.Address,
		Message:   starDto.Message,
		Data:      starDto.Data,
		Signature: starDto.Signature,
	}
	block, err := (*blockchain).SubmitStar(star)
	if err != nil {
		log.Println("ERR: submitStar: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, "Star submition failed: "+err.Error())
		return
	}
	blockDto := BlockDto{
		Body:              block.Body,
		Hash:              block.Hash,
		Height:            block.Height,
		Owner:             block.Owner,
		PreviousBlockHash: block.PreviousBlockHash,
		Time:              block.Time,
	}
	blockJson, err := json.Marshal(blockDto)
	if err != nil {
		log.Println("ERR: submitStar failed to marshal block: ", err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Failed to serialize block into JSON")
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, string(blockJson))
}

func respondWithBlock(res http.ResponseWriter, req *http.Request, block *contracts.Block, err error) {
	if err != nil {
		log.Println("ERR: respondWithBlock: block not found: ", err)
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprint(res, "Block not found")
		return
	}
	blockDto := BlockDto{
		Body:              block.Body,
		Hash:              block.Hash,
		Height:            block.Height,
		Owner:             block.Owner,
		PreviousBlockHash: block.PreviousBlockHash,
		Time:              block.Time,
	}
	blockJson, err := json.Marshal(blockDto)
	if err != nil {
		log.Println("ERR: respondWithBlock failed to marshal block: ", err)
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "Failed to serialize block into JSON")
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	fmt.Fprint(res, string(blockJson))
}
