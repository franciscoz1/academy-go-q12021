package repository

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"Users/francisco.zamudio/projects/academy-go-q12021/model"
)

//Repository Struct that has all needed functions to interact with the
//cvs file
type Repository struct {
}

//New Create a instance of a Repository
func New() *Repository {
	return &Repository{}
}

//GetAllPokemonData reads a csv file and extracts all its data
func (r Repository) GetAllPokemonData() ([][]string, error) {
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
func (r Repository) SavePokemon(pokemon []string) error {
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

//GetExternalAPI calls an external api
func (r Repository) GetExternalAPI(id string) (model.Pokemon, error) {
	var pokemon model.Pokemon
	response, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + id)
	if err != nil {
		return pokemon, err
	}
	data, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(data, &pokemon)
	return pokemon, err
}
