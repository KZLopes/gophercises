package api

import (
	"ex18/internal/primitive"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var fs = http.FileServer(http.Dir("./storage/"))

func SetupRoutes(mux *http.ServeMux) {

	mux.Handle("GET /files/", http.StripPrefix("/files/", fs))

	mux.HandleFunc("POST /primitive/", submitHandler)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/static/form.html")
	})

}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	f, fh, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "invalid file format", http.StatusBadRequest)
		return
	}
	ext := filepath.Ext(fh.Filename)
	switch ext {
	case ".jpg", "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	default:
		http.Error(w, "invalid file format", http.StatusBadRequest)
		return
	}

	outFile, err := primitive.Transform(f, ext)
	if err != nil {
		http.Error(w, "invalid file format", http.StatusBadRequest)
		return
	}

	file, err := os.CreateTemp("storage", "*"+ext)
	if err != nil {
		http.Error(w, "invalid file format", http.StatusBadRequest)
		return
	}
	defer file.Close()
	io.Copy(file, outFile)

	redUrl := fmt.Sprintf("/files/%s", filepath.Base(file.Name()))
	http.Redirect(w, r, redUrl, http.StatusFound)
}

/*
	Save uploaded image to system
*/
