package main

import (
	"github.com/starchain/api"
	"log"
	"net/http"
)

func main() {
	log.Println("Hello StarchainGo!")
	api.RegisterHandlers()
	http.ListenAndServe(":8000", nil)
}
