package services

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/models"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedRepository struct {
	mock.Mock
}

func (m *MockedRepository) GetAllPokemonData() ([][]string, error) {
	args := m.Called()
	return args.Get(0).([][]string), args.Error(1)
}

func (m *MockedRepository) SavePokemon(pokemon []string) error {

	return nil

}

type MockedExternalAPIRepository struct {
	mock.Mock
}

func (m *MockedExternalAPIRepository) GetExternalAPI(id string) (models.Pokemon, error) {
	args := m.Called(id)
	return args.Get(0).(models.Pokemon), args.Error(1)
}

func TestGetPokemonByID(t *testing.T) {
	tests := []struct {
		name                   string
		getAllPokemonDataValue [][]string
		input                  string
		expected               models.Pokemon
		err                    error
	}{
		{
			name:                   "happy path: returns pokemon given an id",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}},
			input:                  "2",
			expected:               models.Pokemon{ID: 2, Name: "ivysaur"},
			err:                    nil,
		},
		{
			name:                   "not happy path: the csv file is empty",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}},
			input:                  "2",
			expected:               models.Pokemon{},
			err:                    errors.New("The csv file is empty"),
		},
		{
			name:                   "not happy path: no pokemon with given id exits",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}},
			input:                  "7",
			expected:               models.Pokemon{},
			err:                    errors.New("No pokemon with given id exits"),
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		fmt.Println(test.name)
		dataRepository := new(MockedRepository)
		externalapiRepository := new(MockedExternalAPIRepository)
		dataRepository.On("GetAllPokemonData").Return(test.getAllPokemonDataValue, nil)
		pokemonInteractor := PokemonInteractor{dataRepository, externalapiRepository}
		pokemon, err := pokemonInteractor.GetPokemonByID(test.input)
		if err != nil {
			assert.NotNil(test.err)
		} else {
			assert.Equal(test.expected, pokemon)
		}
	}
}

func TestExternalPokemon(t *testing.T) {
	tests := []struct {
		name                string
		externalAPIResponse models.Pokemon
		externalAPIError    error
		savePokemonResponse error
		input               string
		expected            models.Pokemon
		err                 error
	}{
		{
			name:                "happy path: call external api. Get its response, save the response in csv, finally return response.",
			input:               "2",
			externalAPIResponse: models.Pokemon{ID: 9, Name: "blastoise"},
			externalAPIError:    nil,
			savePokemonResponse: nil,
			expected:            models.Pokemon{ID: 9, Name: "blastoise"},
			err:                 nil,
		},
		{
			name:                "not happy path: Called to external api failed.",
			input:               "10",
			externalAPIResponse: models.Pokemon{},
			externalAPIError:    errors.New("The HTTP request failed with error"),
			savePokemonResponse: nil,
			expected:            models.Pokemon{},
			err:                 errors.New("The HTTP request failed with error"),
		},
		{
			name:                "not happy path: Problem saving data to the csv",
			input:               "2",
			externalAPIResponse: models.Pokemon{ID: 9, Name: "blastoise"},
			externalAPIError:    nil,
			savePokemonResponse: errors.New("Problem saving data to the csv"),
			expected:            models.Pokemon{ID: 9, Name: "blastoise"},
			err:                 errors.New("Problem saving data to the csv"),
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		fmt.Println(test.name)
		dataRepository := new(MockedRepository)
		externalapiRepository := new(MockedExternalAPIRepository)
		externalapiRepository.On("GetExternalAPI", test.input).Return(test.externalAPIResponse, test.externalAPIError)
		dataRepository.On("SavePokemon").Return(test.savePokemonResponse)
		pokemonInteractor := PokemonInteractor{dataRepository, externalapiRepository}
		pokemon, err := pokemonInteractor.ExternalPokemon(test.input)
		if err != nil {
			assert.NotNil(test.err)
		} else {
			assert.Equal(test.expected, pokemon)
		}
	}
}
