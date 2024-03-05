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

type LocationData struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func getAPIEndpoint(path string, cache *cache.Cache) ([]byte, error) {
	requestURL := fmt.Sprintf("%v%v", endpoint, path)

	data, cacheObj := cache.Get(requestURL)

	if cacheObj {
		return data, nil

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

func ExploreLocation(location string, cache *cache.Cache) (LocationData, error) {
	loc := LocationData{}

	path := fmt.Sprintf("location-area/%v", location)

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
