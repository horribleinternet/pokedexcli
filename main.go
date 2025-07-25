package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextURL string
	prevURL string
}

var commands map[string]cliCommand

func initCommands() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func main() {
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)
	prompt := "Pokedex > "
	fmt.Print(prompt)
	context := config{nextURL: pokeapi.LocationStartURL, prevURL: ""}
	for scanner.Scan() {
		input := cleanInput(scanner.Text())
		comm, ok := commands[input[0]]
		if ok {
			if err := comm.callback(&context); err != nil {
				fmt.Printf("Error: %v", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		fmt.Print(prompt)
	}
}

func commandExit(context *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(context *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandMap(context *config) error {
	if context.nextURL == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	return printMap(context.nextURL, context)
}

func commandMapb(context *config) error {
	if context.prevURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	return printMap(context.prevURL, context)
}

func printMap(url string, context *config) error {
	locations, nextURL, prevURL, err := pokeapi.LocationPage(url)
	if err != nil {
		return err
	}
	context.nextURL = nextURL
	context.prevURL = prevURL
	for _, location := range locations {
		fmt.Println(location)
	}
	return nil
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	out := strings.Fields(text)
	if len(out) == 0 {
		return []string{""}
	}
	return out
}
