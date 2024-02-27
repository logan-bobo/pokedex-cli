package main

import (
	"fmt"
	"os"
	"bufio"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func buildCommandInterface() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandExit() error {
	os.Exit(0)	
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	return nil
}

func main() {
	cliCommands := buildCommandInterface()	
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("Pokedex -> ")
		
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
		
		command, ok :=  cliCommands[scanner.Text()]
		if !ok {
			fmt.Println("Command not found")
			continue
		}
		
		command.callback()
	}
}


