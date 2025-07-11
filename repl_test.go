package main

import (
	// "reflect"
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	// definimos los casos de prueba
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: " Charmander  Bulbasaur  PIKACHU ",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input: "   ",
			expected: []string{},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: " ONLY",
			expected: []string{"only"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		fmt.Println("================	§ =======================")
		fmt.Println("Input:", actual)
		// Primero verificamos la longitud
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: expected %v, got %v", c.expected, actual)
			continue
		}
		// Luego verificamos cada palabra
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("at index %d: expected %s, got %s", i, c.expected[i], actual[i])
				break
			}
		}

	}
}
