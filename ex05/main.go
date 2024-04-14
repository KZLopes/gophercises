package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	url := "https://www.calhoun.io"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}
