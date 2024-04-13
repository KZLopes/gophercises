package main

import (
	"MyHtmlParser"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	exampleNumber := flag.String("ex", "1", "The number of the example to run. (1, 2, 3 or 4)")
	flag.Parse()

	filename := "ex" + *exampleNumber + ".html"
	exampleHtml, err := os.Open(filepath.Join("../../examples", filename))
	if err != nil {
		panic(err)
	}
	defer exampleHtml.Close()

	links, err := MyHtmlParser.Parse(exampleHtml)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
