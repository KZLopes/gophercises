package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	filename := flag.String("file", "gopher.json", "The json file that contains the book")
	flag.Parse()

	f, err := os.Open(*filename)
	checkError(err)
	defer f.Close()

	d := json.NewDecoder(f)
	var book Book
	err = d.Decode(&book)
	checkError(err)

	tmpl, err := template.New("arc-name").Parse("chapter")
	checkError(err)

	_ = tmpl

	startServer(book)
}

type Book map[string]chapter

type chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func startServer(b Book) {
	mux := http.NewServeMux()

	for k, v := range b {
		mux.HandleFunc(fmt.Sprintf("/%s", k), func(w http.ResponseWriter, r *http.Request) {
			t, err := loadTemplate()
			if err != nil {
				log.Printf("Falha ao carregar template para %s. %s\n", k, err)
				http.Redirect(w, r, r.Header.Get("Referer"), http.StatusNotFound)
				return
			}

			t.Execute(w, v)
		})
	}
	PORT := ":8080"
	log.Printf("Server running in port %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, mux))

}

func loadTemplate() (*template.Template, error) {
	t, err := template.ParseFiles("template.html")
	if err != nil {
		return nil, err
	}
	return t, nil
}

func checkError(e error) {
	if e != nil && e != io.EOF {
		log.Fatal(e)
	}

}
