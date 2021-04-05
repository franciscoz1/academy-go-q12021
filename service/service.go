package service

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"

	"Users/francisco.zamudio/projects/academy-go-q12021/model"
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

//ConcurrentPokemonList given a type, items and items per worker  return a list of pokemon
func (s Service) ConcurrentPokemonList(parity string, items string, ipw string) ([]model.Pokemon, error) {

	maximumItemsInResponse, err := strconv.Atoi(items)
	if err != nil {
		err := errors.New("Problem casting items" + err.Error())
		return nil, err
	}

	itemsPerWorker, err := strconv.Atoi(ipw)
	if err != nil {
		err := errors.New("Problem casting items per worker" + err.Error())
		return nil, err
	}

	records, err := s.repository.GetAllPokemonData()
	if err != nil {
		err := errors.New("Problem retrieving data from the csv" + err.Error())
		return nil, err
	}

	if len(records) <= 2 {
		err := errors.New("The csv file is empty")
		return nil, err
	}

	pokemons := []model.Pokemon{}
	pokemonResponse := []model.Pokemon{}

	for _, record := range records[1:] {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			err := errors.New("Problem casting id" + err.Error())
			return nil, err
		}
		data := model.Pokemon{ID: id, Name: record[1]}
		pokemons = append(pokemons, data)
	}

	poolsize := runtime.GOMAXPROCS(0)

	jobs := make(chan model.Pokemon, len(pokemons))
	results := make(chan model.Pokemon, len(pokemons))

	wg := new(sync.WaitGroup)

	wg.Add(poolsize)

	for w := 1; w <= poolsize; w++ {
		go worker(parity, itemsPerWorker, jobs, results, wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < len(pokemons); i++ {
		jobs <- pokemons[i]
	}
	close(jobs)

	for result := range results {
		pokemonResponse = append(pokemonResponse, result)
		if len(pokemonResponse) == maximumItemsInResponse {
			break
		}
	}
	return pokemonResponse, nil
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

func worker(parity string, itemsPerWorker int, jobs <-chan model.Pokemon, results chan<- model.Pokemon, wg *sync.WaitGroup) {
	itemsCounter := 0
	for item := range jobs {
		if itemsCounter == itemsPerWorker {
			break
		}
		if confirmParity(parity, item) {
			results <- item
		}
		itemsCounter++
	}
	defer wg.Done()
}

func confirmParity(parity string, pokemon model.Pokemon) bool {
	if (parity == "even" && pokemon.ID%2 == 0) || (parity == "odd" && pokemon.ID%2 != 0) {
		return true
	}
	return false
}
