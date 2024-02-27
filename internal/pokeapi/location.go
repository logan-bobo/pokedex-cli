package pokeapi

import (
	"net/http"
	"io"
	"fmt"
	"errors"
)

const endpoint = "https://pokeapi.co/api/v2/"

func getAPIEndpoint(path string) ([]byte, error) {
	requestPath := fmt.Sprintf("%v%v", endpoint, path)

	resp, err := http.Get(requestPath)

	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	resp.Body.Close()

	if resp.StatusCode > 299 {
		err := fmt.Sprintf("Non 200 status code, got %v on path %v", resp.StatusCode, requestPath)
		return []byte{}, errors.New(err)
	}

	return body, nil
}

func GetNextLocations() ([]byte, error) {
	body, err := getAPIEndpoint("location")

	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
