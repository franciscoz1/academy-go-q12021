package service

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/model"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//Service use to have access to the different use case functions.
type Service struct {
	repository Repository
}

//Repository Interface
type Repository interface {
	GetAllPokemonData() ([][]string, error)
	SavePokemon([]string) error
	GetExternalAPI(id string) (model.Pokemon, error)
}

//New Create a instance of a Service
func New(s Repository) *Service {
	return &Service{s}
}

//GetPokemonByID given an id a pokemon is return
func (s Service) GetPokemonByID(id string) (model.Pokemon, error) {

	pokemons := []model.Pokemon{}
	var pokemon model.Pokemon

	records, err := s.repository.GetAllPokemonData()
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
		data := model.Pokemon{ID: id, Name: record[1]}
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
func (s Service) ExternalPokemon(id string) (model.Pokemon, error) {

	var pokemon model.Pokemon

	pokemon, err := s.repository.GetExternalAPI(id)
	if err != nil {
		err := errors.New("The HTTP request failed with error" + err.Error())
		return pokemon, err
	}

	pokemonSlice := strucToSlice(pokemon)

	err = s.repository.SavePokemon(pokemonSlice)
	if err != nil {
		err := errors.New("Problem saving data to the csv: " + err.Error())
		return pokemon, err
	}

	return pokemon, nil
}

func strucToSlice(pokemonStruct model.Pokemon) []string {
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
