package services

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/models"
	"Users/francisco.zamudio/projects/academy-go-q12021/repository"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//PokemonInteractor Struct used to have access to the different use case functions.
type PokemonInteractor struct {
	Repository            repository.InterfaceRepository
	ExternalapiRepository repository.InterfaceExternalAPIRepository
}

//GetPokemonByID given an id a pokemon is return
func (p PokemonInteractor) GetPokemonByID(id string) (models.Pokemon, error) {

	pokemons := []models.Pokemon{}
	var pokemon models.Pokemon

	records, err := p.Repository.GetAllPokemonData()
	if err != nil {
		err := errors.New("Problem retrieving data from the csv" + err.Error())
		return pokemon, err
	}

	if len(records) <= 2 {
		err := errors.New("The csv file is empty")
		return pokemon, err
	}
	for _, record := range records[1:] {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			err := errors.New("Problem casting id" + err.Error())
			return pokemon, err
		}
		data := models.Pokemon{ID: id, Name: record[1]}
		pokemons = append(pokemons, data)
	}

	for _, record := range pokemons {
		if strconv.Itoa(record.ID) == id {
			pokemon = record
			break
		}
	}
	if pokemon.Name == "" {
		err := errors.New("No pokemon with given id exits")
		return pokemon, err
	}
	return pokemon, nil
}

//ExternalPokemon given an id call an external API, then save and send its response
func (p PokemonInteractor) ExternalPokemon(id string) (models.Pokemon, error) {

	var pokemon models.Pokemon

	pokemon, err := p.ExternalapiRepository.GetExternalAPI(id)
	if err != nil {
		err := errors.New("The HTTP request failed with error" + err.Error())
		return pokemon, err
	}

	pokemonSlice := strucToSlice(pokemon)

	err = p.Repository.SavePokemon(pokemonSlice)
	if err != nil {
		err := errors.New("Problem saving data to the csv: " + err.Error())
		return pokemon, err
	}

	return pokemon, nil
}

func strucToSlice(pokemonStruct models.Pokemon) []string {
	v := reflect.ValueOf(pokemonStruct)
	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}
	pokemonSlice := make([]string, len(values))
	for i, v := range values {
		pokemonSlice[i] = fmt.Sprint(v)
	}
	return pokemonSlice
}
