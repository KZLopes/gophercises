package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", debugHandler)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", devMw(mux)))
}

func devMw(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s</pre>", err, makeLinks(string(stack)))
			}
		}()
		app.ServeHTTP(w, r)
	}
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 	u := r.
	// 	fmt.Println("URI:", u)
	// 	p, _ := url.Parse(u.String())
	//
	// 	frag := p.Fragment
	// 	fmt.Println("FRAG:", frag)
	// 	lineNum, _ := strconv.Atoi(strings.TrimPrefix(frag, "bug"))
	contents, _ := io.ReadAll(file)
	writeWithChroma(w, string(contents))

}

func writeWithChroma(w http.ResponseWriter, contents string) {
	lexer := lexers.Get("go")
	style := styles.Get("monokai")
	formatter := html.New(html.Standalone(true), html.WithClasses(true), html.WithLineNumbers(true), html.WithLinkableLineNumbers(true, "bug") /*, html.HighlightLines([][2]int{{lineNum, lineNum}})*/)
	iterator, _ := lexer.Tokenise(nil, contents)

	formatter.Format(w, style, iterator)
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func makeLinks(stack string) string {
	lines := strings.Split(stack, "\n")

	for idx, line := range lines {
		if strings.HasPrefix(line, "\t") {
			collon := strings.LastIndex(line, ":")
			fp := line[1:collon]
			afterCollon := line[collon+1:]
			// Line Number
			lineNum := strings.Split(afterCollon, " ")[0]

			v := url.Values{}
			v.Set("path", fp)
			lines[idx] = "\t<a href=\"/debug/?" + v.Encode() + "#bug" + lineNum + "\">" + fp + "</a>" + line[collon:]

		}
	}

	return strings.Join(lines, "\n")
}
