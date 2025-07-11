package repl

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cc-kevin-bolivar/pokedex/internal/pokeapi"
)

type command struct {
	Name 			  string
	Description string
	Callback    func(*Config, []string) error
}

func GetCommands(cfg *Config) map[string]command {
	return map[string]command{
		"help": {
			Name:				 "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name: 			 "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:				 "map",
			Description: "Displays next 20 location areas",
			Callback:		 commandMap,
		},
		"mapb": {
			Name:				 "mapb",
			Description: "Displays previous 20 location areas",
			Callback:		 commandMapb,
		},
		"explore": {
    	Name:        "explore",
    	Description: "Explora un área para ver los Pokémon disponibles",
    	Callback:    commandExplore,
		},
		"catch": {
    	Name:        "catch",
    	Description: "Intenta atrapar un Pokémon",
    	Callback:    commandCatch,
		},
		"inspect": {
			Name:					"inspect",
			Description:  "Muestra detalles de un Pokémon capturado",
			Callback:		  commandInspect,
		},
		"pokedex": {
			Name:					"pokedex",
			Description: "Muestra los Pokémon capturados",
			Callback:    commandPokedex,
		},
	}
}

func commandHelp(cfg *Config, args []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
  fmt.Println("Usage: ")

	for _, cmd := range GetCommands(cfg) {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}

	fmt.Println()
	return nil
}

func commandExit(cfg *Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, args []string) error {
	// Usar cfg.Client en lugar de crear un nuevo cliente
	locationAreas, err := cfg.Client.GetLocationAreas(cfg.Next)
	if err != nil {
			return err
	}

	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous

	for _, area := range locationAreas.Results {
			fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *Config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}

	pokeClient := pokeapi.NewClient(5 * time.Minute)
	locationAreas, err := pokeClient.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	cfg.Next = locationAreas.Next
	cfg.Previous = locationAreas.Previous

	for _, area := range locationAreas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(cfg *Config, args []string) error {
	if len(args) == 0 {
			return errors.New("Debes especificar un área (ej: explore pastoria-city-area)")
	}

	areaName := args[0]
	fmt.Printf("Exploring %s...\n", areaName)

	area, err := cfg.Client.GetLocationAreaDetails(areaName)
	if err != nil {
			return fmt.Errorf("error fetching area details: %w", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range area.PokemonEncounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *Config, args []string) error {
	if len(args) == 0 {
			return errors.New("Debes especificar un Pokémon (ej: catch pikachu)")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := cfg.Client.GetPokemon(pokemonName)
	if err != nil {
			return fmt.Errorf("error fetching Pokémon: %w", err)
	}

	// Calcular probabilidad de captura (menor para Pokémon con más experiencia)
	catchProbability := 50.0 // Base 50%
	catchProbability -= float64(pokemon.BaseExperience) * 0.1

	// Asegurar probabilidad mínima
	if catchProbability < 10 {
			catchProbability = 10
	}

	// Generar número aleatorio
	rand.Seed(time.Now().UnixNano())
	roll := rand.Float64() * 100

	if roll <= catchProbability {
			fmt.Printf("%s was caught!\n", pokemonName)
			fmt.Println("You may now inspect it with the inspect command.")
			
			// Inicializar Pokédex si es necesario
			if cfg.Pokedex == nil {
					cfg.Pokedex = make(map[string]pokeapi.Pokemon)
			}
			
			// Añadir al Pokédex
			cfg.Pokedex[pokemonName] = pokemon
	} else {
			fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(cfg *Config, args []string) error {
	if len(args) == 0 {
		return errors.New("Debes especificar un Pokémon (ej: inspect pikachu)")
	}
	pokemonName := strings.ToLower(args[0])

	// verificar si el Pokémon esta en el Pokédex
	pokemon, exists := cfg.Pokedex[pokemonName]
	if !exists {
		fmt.Println("Tu no has capturado ese pokemon")
		return nil
	}

	// Mostrar información detallada
	fmt.Printf("Nombre: %s\n", pokemon.Name)
	fmt.Printf("Altura: %d\n", pokemon.Height)
	fmt.Printf("Peso: %d\n", pokemon.Weight)

	// Mostrar stats
	fmt.Printf("stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	// Mostrar tipos
	fmt.Printf("Tipos:\n")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf(" - %s\n", pokemonType.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *Config, args []string) error {
	fmt.Println("Your Pokedex:")
	
	// Verificar si el Pokédex está vacío
	if len(cfg.Pokedex) == 0 {
			fmt.Println("  You haven't caught any Pokémon yet!")
			return nil
	}
	
	// Mostrar todos los Pokémon capturados en orden alfabético
	names := make([]string, 0, len(cfg.Pokedex))
	for name := range cfg.Pokedex {
			names = append(names, name)
	}
	
	sort.Strings(names) // Ordenar alfabéticamente
	
	for _, name := range names {
			fmt.Printf(" - %s\n", name)
	}
	
	return nil
}