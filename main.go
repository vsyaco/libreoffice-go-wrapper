package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func changeFileExt(filename string, newExt string) string {
	ext := filepath.Ext(filename)
	newFilename := filename[0:len(filename)-len(ext)] + "." + newExt
	return newFilename
}

func main() {
	http.HandleFunc("/convert", convertHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the file from the request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file for the DOCX file
	tmpDocxFile, err := os.CreateTemp("", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tmpDocxFile.Close()

	// Write the DOCX file to the temporary file
	_, err = io.Copy(tmpDocxFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the DOCX file to PDF using LibreOffice
	cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf", "--outdir", os.TempDir(), tmpDocxFile.Name())
	err = cmd.Run()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/pdf
	w.Header().Set("Content-Type", "application/pdf")

	pdfFilename := changeFileExt(tmpDocxFile.Name(), "pdf")
	pdfFile, err := os.Open(pdfFilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer pdfFile.Close()

	_ = os.Remove(pdfFilename)
	_ = os.Remove(tmpDocxFile.Name())

	// Write the PDF file to the response
	_, err = io.Copy(w, pdfFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
