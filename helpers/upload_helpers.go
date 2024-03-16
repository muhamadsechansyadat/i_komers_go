package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func GenerateUniqueFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	hash := md5.Sum([]byte(originalName))
	hashStr := hex.EncodeToString(hash[:])
	ext := filepath.Ext(originalName)
	newFileName := fmt.Sprintf("%d_%s%s", timestamp, hashStr, ext)
	return newFileName
}

func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func SaveUploadedFile(file *multipart.FileHeader, filePath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
