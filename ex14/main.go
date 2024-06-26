package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", recoverMw(mux, true)))
}

func recoverMw(app http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				stack := debug.Stack()
				log.Println(string(stack))
				if dev {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "<h1>panic: %v</h1><pre>%s<pre>", err, string(stack))
					return
				}
				http.Error(w, "Someting Went Wrong", http.StatusInternalServerError)
			}
		}()
		nw := &respWriter{ResponseWriter: w}
		app.ServeHTTP(nw, r)
		nw.flush()
	}
}

// type ResponseWriter interface {
// 	Header() Header
// 	Write([]byte)(int,error)
// 	WriteHeader(statusCOde int)
// }

type respWriter struct {
	http.ResponseWriter
	writes [][]byte
	status int
}

func (rw *respWriter) Write(b []byte) (int, error) {
	rw.writes = append(rw.writes, b)
	return len(b), nil
}

func (rw *respWriter) Writeheader(statusCode int) {
	rw.status = statusCode
}

func (rw *respWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}
	for _, write := range rw.writes {
		_, err := rw.ResponseWriter.Write(write)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rw *respWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("does not impement Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *respWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	flusher.Flush()
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
