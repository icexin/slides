package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	buf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(buf)
	w, err := bodyWriter.CreateFormFile("uploadFile", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	io.Copy(w, f)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post("http://127.0.0.1:8080/upload", contentType, buf)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
}
