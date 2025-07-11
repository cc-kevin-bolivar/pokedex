package repl

import "github.com/cc-kevin-bolivar/pokedex/internal/pokeapi"

type Config struct {
	Next     *string
	Previous *string
	Client   *pokeapi.Client
	Pokedex  map[string]pokeapi.Pokemon
}