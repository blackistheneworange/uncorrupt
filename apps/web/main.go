package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/blackistheneworange/gohttprouter"
	"github.com/blackistheneworange/uncorrupt"
)

var pageHtml []byte
var sslEnabled = flag.Bool("ssl", true, "")

func handleRender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(pageHtml)
}

func handleCorrupt(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bytes = []byte{}
	fileInput := true

	parseErr := r.ParseMultipartForm(50 << 20)
	if parseErr != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	file, fileHeader, fileErr := r.FormFile("inputFile")
	if fileErr != nil {
		text := r.FormValue("inputText")
		if len(text) == 0 {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		fileInput = false
		bytes = []byte(text)
	} else {
		bytes, _ = io.ReadAll(file)
		defer file.Close()
	}

	key := r.FormValue("key")

	corrupted := uncorrupt.Corrupt(bytes, key)

	w.Header().Set("Content-Length", strconv.Itoa(len(corrupted)))
	if fileInput {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+fileHeader.Filename)
	} else {
		w.Header().Set("Content-Type", "application/base-64")
	}

	w.Write(corrupted)
}

func handleUncorrupt(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var bytes = []byte{}
	fileInput := true

	parseErr := r.ParseMultipartForm(50 << 20)
	if parseErr != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	file, fileHeader, fileErr := r.FormFile("inputFile")
	if fileErr != nil {
		text := r.FormValue("inputText")
		if len(text) == 0 {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		fileInput = false
		bytes = []byte(text)
	} else {
		bytes, _ = io.ReadAll(file)
		defer file.Close()
	}

	key := r.FormValue("key")

	uncorrupted := uncorrupt.Uncorrupt(bytes, key)

	w.Header().Set("Content-Length", strconv.Itoa(len(uncorrupted)))
	if fileInput {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+fileHeader.Filename)
	} else {
		w.Header().Set("Content-Type", "application/base-64")
	}

	w.Write(uncorrupted)
}

func main() {
	flag.Parse()

	port := os.Getenv("PORT")
	router := gohttprouter.NewRouter()

	fs := http.FileServer(http.Dir("./page/static"))
	router.GET("/static/:any", http.StripPrefix("/static/", fs).ServeHTTP)

	router.GET("/", handleRender)
	router.POST("/corrupt", handleCorrupt)
	router.POST("/uncorrupt", handleUncorrupt)

	log.Println("Server started in port", port)

	if *sslEnabled {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/certificate.crt", "certs/private.key", router))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, router))
	}
}

func init() {
	bytes, err := os.ReadFile("./page/index.html")
	if err != nil {
		panic(err)
	}

	pageHtml = bytes
}
