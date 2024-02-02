package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func NewTempFile() (*os.File, error) {
	now := time.Now().Unix()
	return os.CreateTemp("/tmp", fmt.Sprintf("%d-*.jpg", now))
}

func NewTempFileName() (string, error) {
	f, err := NewTempFile()
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Name(), nil
}

func DownloadImage(url string) (string, error) {
	r, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	// Create a file to store the image
	f, err := NewTempFile()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Write the image to the file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
