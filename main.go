package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/logan-bobo/pokedex-cli/internal/cache"
	"github.com/logan-bobo/pokedex-cli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(*config, *cahce.Cache) error
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
	}
}

func commandExit(conf *config, cache *cahce.Cache) error {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, cache *cahce.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func mapNext(conf *config, cache *cahce.Cache) error {
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

func mapPrevious(conf *config, cache *cahce.Cache) error {
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

func main() {
	cliCommands := buildCommandInterface()

	scanner := bufio.NewScanner(os.Stdin)

	cache := cahce.NewCache(60)

	for {
		fmt.Print("Pokedex -> ")

		scanner.Scan()

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		command, ok := cliCommands[scanner.Text()]

		if !ok {
			fmt.Println("Command not found")
			continue
		}

		command.callback(command.config, cache)
	}
}

