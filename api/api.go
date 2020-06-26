// api package provides handlers for http REST API
// Its goal is to provide an interface for the Blockchain
package api

import (
	"fmt"
	"log"
	"net/http"
)

func RegisterHandlers() {
	http.HandleFunc("/hello", hello)
	log.Println("INFO: Handlers registered successfully")
}

func hello(res http.ResponseWriter, req *http.Request) {
	log.Println("hello")
	fmt.Fprint(res, "hello")
}
