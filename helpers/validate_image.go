package helpers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func ValidateFile(file *multipart.FileHeader) error {
	// Periksa ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("Invalid file format. Only JPG, JPEG, and PNG are allowed.")
	}

	// Baca tipe MIME file untuk memastikan keamanan tambahan
	fileType, err := getFileType(file)
	if err != nil {
		return err
	}

	// Periksa tipe MIME file
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if !allowedMimeTypes[fileType] {
		return fmt.Errorf("Invalid file format. Only JPG, JPEG, and PNG are allowed.")
	}

	return nil
}

func getFileType(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Baca beberapa byte pertama untuk mendeteksi tipe MIME
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	// Tentukan tipe MIME dari byte pertama
	fileType := http.DetectContentType(buffer)

	return fileType, nil
}
