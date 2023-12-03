package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrEmailTaken = errors.New("models: email address is already taken")
	ErrNotFound   = errors.New("models: resource not found")
)

type FileError struct {
	Issue string
}

func (e FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", e.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("could not seek to beginning of file: %v", err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, t := range allowedTypes {
		if contentType == t {
			return nil
		}
	}

	return FileError{Issue: "file type not allowed"}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if hasExtension(filename, allowedExtensions) {
		return nil
	}

	return FileError{Issue: "file extension not allowed"}
}
