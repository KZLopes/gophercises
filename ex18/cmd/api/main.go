package main

import (
	"ex18/internal/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	api.SetupRoutes(mux)
	log.Fatal(http.ListenAndServe(":2323", mux))
}
