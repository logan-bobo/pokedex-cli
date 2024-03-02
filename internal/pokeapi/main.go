package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/logan-bobo/pokedex-cli/internal/cache"
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

func getAPIEndpoint(path string, cache *cache.Cache) ([]byte, error) {
	requestURL := fmt.Sprintf("%v%v", endpoint, path)

	cacheObj := cache.Get(requestURL)

	if cacheObj {
		return cache.Data[requestURL].Val, nil

	} else {
		resp, err := http.Get(requestURL)

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
					requestURL,
				),
			)
		}
		
		cache.Add(requestURL, body)

		return body, nil
	}
}

func GetLocations(offset string, cache *cache.Cache) (Locations, error) {
	loc := Locations{}

	path := fmt.Sprintf("location-area/?offset=%v", offset)

	body, err := getAPIEndpoint(path, cache)

	if err != nil {
		return loc, err
	}

	err = json.Unmarshal(body, &loc)

	if err != nil {
		return loc, err
	}

	return loc, nil
}
