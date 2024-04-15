package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/yudaph/s3-integration/s3client"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Post("/", HandleUpload)
	http.ListenAndServe(":8080", r)
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// TODO: Do Validation

	// Upload file to S3
	filename := "suffix_" + header.Filename
	url, err := s3client.UploadFile(r.Context(), "qwera-upload-test", filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(url))
}
