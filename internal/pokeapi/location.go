package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const endpoint = "https://pokeapi.co/api/v2/"

type Locations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

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
		return []byte{}, errors.New(
			fmt.Sprintf(
				"Non 200 status code, got %v on path %v", 
				resp.StatusCode, 
				requestPath,
			),
		)
	}

	return body, nil
}

func GetLocations(offset string) (Locations, error) {
	loc := Locations{}

	path := fmt.Sprintf("location-area/?offset=%v", offset)

	body, err := getAPIEndpoint(path)

	if err != nil {
		return loc, err
	}

	err = json.Unmarshal(body, &loc)

	if err != nil {
		return loc, err
	}

	return loc, nil
}
