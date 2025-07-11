package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cc-kevin-bolivar/pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/"

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
					Name string `json:"name"`
			} `json:"stat"`
	} `json:"stats"`
	Types []struct {
			Type struct {
					Name string `json:"name"`
			} `json:"type"`
	} `json:"types"`
}

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

type LocationAreasResp struct{
	Count    int `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetail struct {
	Name          string `json:"name"`
	PokemonEncounters []struct {
			Pokemon struct {
					Name string `json:"name"`
			} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetLocationAreas(url *string) (LocationAreasResp, error) {
	endpoint := "/location-area"
	fullURL := baseURL + endpoint
	
	if url != nil {
		fullURL = *url
	}

	// Verificar cache primero
	if data, ok := c.cache.Get(fullURL); ok {
		fmt.Println("Using cached data...")
		locationAreasResp := LocationAreasResp{}
		err := json.Unmarshal(data, &locationAreasResp)
		if err != nil {
			return LocationAreasResp{}, err
		}
		return locationAreasResp, nil
	}

	// Hacer request si no está en cache
	res, err := c.httpClient.Get(fullURL)
	if err != nil {
		return LocationAreasResp{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return LocationAreasResp{}, fmt.Errorf("bad status code: %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasResp{}, err
	}

	// Guardar en cache antes de retornar
	c.cache.Add(fullURL, data)

	locationAreasResp := LocationAreasResp{}
	err = json.Unmarshal(data, &locationAreasResp)
	if err != nil {
		return LocationAreasResp{}, err
	}

	return locationAreasResp, nil
}

func (c *Client) GetLocationAreaDetails(areaName string) (LocationAreaDetail, error) {
	endpoint := "/location-area/" + areaName
	fullURL := baseURL + endpoint

	// Verificar cache primero
	if data, ok := c.cache.Get(fullURL); ok {
			fmt.Println("Using cached data...")
			var area LocationAreaDetail
			err := json.Unmarshal(data, &area)
			if err != nil {
					return LocationAreaDetail{}, err
			}
			return area, nil
	}

	res, err := c.httpClient.Get(fullURL)
	if err != nil {
			return LocationAreaDetail{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
			return LocationAreaDetail{}, err
	}

	// Guardar en cache
	c.cache.Add(fullURL, data)

	var area LocationAreaDetail
	err = json.Unmarshal(data, &area)
	if err != nil {
			return LocationAreaDetail{}, err
	}

	return area, nil
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	if name == "" {
			return Pokemon{}, errors.New("pokemon name cannot be empty")
	}

	endpoint := "/pokemon/" + strings.ToLower(strings.TrimSpace(name))
	fullURL := baseURL + endpoint

	// Verificar cache primero
	if data, ok := c.cache.Get(fullURL); ok {
			var pokemon Pokemon
			if err := json.Unmarshal(data, &pokemon); err != nil {
					return Pokemon{}, fmt.Errorf("error decoding cached data: %w", err)
			}
			return pokemon, nil
	}

	// Validar URL
	if _, err := url.ParseRequestURI(fullURL); err != nil {
			return Pokemon{}, fmt.Errorf("invalid URL formed: %w", err)
	}

	res, err := c.httpClient.Get(fullURL)
	if err != nil {
			return Pokemon{}, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
			return Pokemon{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
			return Pokemon{}, fmt.Errorf("error reading response: %w", err)
	}

	// Guardar en cache
	c.cache.Add(fullURL, data)

	var pokemon Pokemon
	if err := json.Unmarshal(data, &pokemon); err != nil {
			return Pokemon{}, fmt.Errorf("error decoding response: %w", err)
	}

	return pokemon, nil
}