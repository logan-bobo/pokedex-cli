package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/logan-bobo/pokedex-cli/internal/cache"
	"github.com/logan-bobo/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(*config, *cache.Cache, string) error
}

type config struct {
	next     string
	previous string
}

func buildCommandInterface() map[string]cliCommand {
	conf := config{}
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      &conf,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      &conf,
		},
		"map": {
			name:        "map",
			description: "Show the next 20 locations",
			callback:    mapNext,
			config:      &conf,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 locations",
			callback:    mapPrevious,
			config:      &conf,
		},
		"explore": {
			name:        "explore",
			description: "show all pokemon in an area",
			callback:    exploreLocation,
			config:      &conf,
		},
	}
}

func commandExit(conf *config, cache *cache.Cache, location string) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, cache *cache.Cache, location string) error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func mapNext(conf *config, cache *cache.Cache, location string) error {
	var locations pokeapi.Locations
	var err error

	if conf.next == "" {
		locations, err = pokeapi.GetLocations("0", cache)

		if err != nil {
			return err
		}

	} else {
		u, err := url.Parse(conf.next)

		if err != nil {
			return err
		}

		query := u.Query()

		offset, ok := query["offset"]

		if !ok {
			return errors.New(fmt.Sprintf("Offset not found in URL: %v", conf.next))
		}

		locations, err = pokeapi.GetLocations(offset[0], cache)

		if err != nil {
			return err
		}
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	conf.next = locations.Next

	previous, ok := locations.Previous.(string)

	if ok {
		conf.previous = previous
	}

	return nil
}

func mapPrevious(conf *config, cache *cache.Cache, location string) error {
	var locations pokeapi.Locations
	var err error

	if conf.previous == "" {
		fmt.Println("No location to go back to...")
		return err

	} else {
		u, err := url.Parse(conf.previous)

		if err != nil {
			return err
		}

		query := u.Query()

		offset, ok := query["offset"]

		if !ok {
			return errors.New(fmt.Sprintf("Offset not in URL: %v", conf.previous))
		}

		locations, err = pokeapi.GetLocations(offset[0], cache)

		if err != nil {
			return err
		}
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	previous, ok := locations.Previous.(string)

	if ok {
		conf.previous = previous
	}

	conf.next = locations.Next

	return nil
}

func exploreLocation(conf *config, cache *cache.Cache, location string) error {
	locations, err := pokeapi.ExploreLocation(location, cache)

	if err != nil {
	return err
	}
	
	fmt.Println("Found Pokemon...")
	for _, pokemon := range locations.PokemonEncounters {
		fmt.Printf("- %v \n", pokemon.Pokemon.Name)
	}

	return nil
}

func main() {
	cliCommands := buildCommandInterface()

	scanner := bufio.NewScanner(os.Stdin)

	cache := cache.NewCache(60 * time.Second)

	for {
		fmt.Print("Pokedex -> ")

		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		inputRaw := scanner.Text()

		inputSplit := strings.Split(inputRaw, " ")

		command, ok := cliCommands[inputSplit[0]]

		if !ok {
			fmt.Println("Command not found")
			continue
		}

		if len(inputSplit) == 1 {
			command.callback(command.config, cache, "")
		} else {
			command.callback(command.config, cache, inputSplit[1])
		}
	}
}
