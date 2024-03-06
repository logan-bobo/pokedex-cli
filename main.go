package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/logan-bobo/pokedex-cli/internal/cache"
	"github.com/logan-bobo/pokedex-cli/internal/pokeapi"
)

type pokedex struct {
	entities map[string]pokeapi.Pokemon
}

func newPokedex() *pokedex {
	p := pokedex{
		entities: map[string]pokeapi.Pokemon{},
	}

	return &p
}

type cliCommand struct {
	name        string
	description string
	config      *config
	callback    func(*config, *cache.Cache, *pokedex, string) error
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
			description: "Show all pokemon in an area",
			callback:    exploreLocation,
			config:      &conf,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    catchPokemon,
			config:      &conf,
		},
	}
}

func commandExit(conf *config, cache *cache.Cache, pokedex *pokedex, location string) error {
	fmt.Println("Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(conf *config, cache *cache.Cache, pokedex *pokedex, location string) error {
	fmt.Println("Welcome to the Pokedex!")

	return nil
}

func mapNext(conf *config, cache *cache.Cache, pokedex *pokedex, location string) error {
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

func mapPrevious(conf *config, cache *cache.Cache, pokedex *pokedex, location string) error {
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

func exploreLocation(conf *config, cache *cache.Cache, pokedex *pokedex, location string) error {
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

func catchPokemon(conf *config, cache *cache.Cache, pokedex *pokedex, name string) error {
	catch := false

	pokemon, err := pokeapi.GetPokemon(name, cache)

	if err != nil {
		return err
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	roll := rand.Intn(100)

	if pokemon.BaseExperience > 75 {
		if roll > 75 {
			catch = true
			fmt.Printf("Caught %v \n", pokemon.Name)
		} else {
			fmt.Printf("Failed to catch %v \n", pokemon.Name)
		}

	} else if pokemon.BaseExperience > 50 {
		if roll > 50 {
			catch = true
			fmt.Printf("Caught %v \n", pokemon.Name)
		} else {
			fmt.Printf("Failed to catch %v \n", pokemon.Name)
		}

	} else if pokemon.BaseExperience > 25 {
		if roll > 25 {
			catch = true
			fmt.Printf("Caught %v \n", pokemon.Name)
		} else {
			fmt.Printf("Failed to catch %v \n", pokemon.Name)
		}

	} else if pokemon.BaseExperience > 0 {
		catch = true
		fmt.Printf("Caught %v \n", pokemon.Name)
	}

	if catch {
		_, ok := pokedex.entities[pokemon.Name]

		if !ok {
			pokedex.entities[pokemon.Name] = pokemon
		} else {
			fmt.Println("Pokemon already registered in your pokedex")
		}
	}

	return nil
}

func main() {
	cliCommands := buildCommandInterface()

	scanner := bufio.NewScanner(os.Stdin)

	cache := cache.NewCache(60 * time.Second)

	pokedex := newPokedex()

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
			command.callback(command.config, cache, pokedex, "")
		} else {
			command.callback(command.config, cache, pokedex, inputSplit[1])
		}
	}
}
