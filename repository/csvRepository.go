package repository

import (
	"encoding/csv"
	"os"
)

//DataRepository Struct that has all needed functions to interact with the
//cvs file
type DataRepository struct {
}

//InterfaceRepository Interface with all functions available for the csv
type InterfaceRepository interface {
	GetAllPokemonData() ([][]string, error)
	SavePokemon([]string) error
}

//GetAllPokemonData reads a csv file and extracts all its data
func (r DataRepository) GetAllPokemonData() ([][]string, error) {
	recordFile, err := os.Open("./pokemon.csv")
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(recordFile)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, err
}

//SavePokemon Saves a pokemon on a file
func (r DataRepository) SavePokemon(pokemon []string) error {
	recordFile, err := os.OpenFile("./pokemon.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	w := csv.NewWriter(recordFile)

	if err := w.Write(pokemon); err != nil {
		return err
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	return err
}
