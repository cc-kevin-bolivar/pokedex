package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cc-kevin-bolivar/pokedex/internal/pokeapi"
	"github.com/cc-kevin-bolivar/pokedex/repl"
)

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func main() {
	// Crear el cliente con cache de 5 minutos
	pokeClient := pokeapi.NewClient(5 * time.Minute)
	
	// Configuración inicial con el cliente
	cfg := &repl.Config{
    Client:  &pokeClient,
    Pokedex: make(map[string]pokeapi.Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}

		commandName := input[0]
		args := []string{}
		if len(input) > 1 {
			args = input[1:]
		}

		commands := repl.GetCommands(cfg)

		if command, exists := commands[commandName]; exists {
			err := command.Callback(cfg, args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}