package main

import (
	"errors"
	"io"
	"net/http"
	"os"
)

func getReader(fileName string) (io.ReadCloser, error) {
	if fileName == "" {
		// pipe type
		return os.Stdin, nil
	} else if isValidURL(fileName) {
		// remote type
		response, err := http.Get(fileName)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != 200 {
			defer response.Body.Close()
			return nil, errors.New(response.Status)
		}
		return response.Body, nil
	} else {
		// local type
		file, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}

		return file, nil
	}
}
