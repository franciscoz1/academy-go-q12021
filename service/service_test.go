package service

import (
	"errors"
	"fmt"
	"testing"

	"Users/francisco.zamudio/projects/academy-go-q12021/model"

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

func (m *MockedRepository) GetExternalAPI(id string) (model.Pokemon, error) {
	args := m.Called(id)
	return args.Get(0).(model.Pokemon), args.Error(1)
}

func TestGetPokemonByID(t *testing.T) {
	tests := []struct {
		name                   string
		getAllPokemonDataValue [][]string
		input                  string
		expected               model.Pokemon
		err                    error
	}{
		{
			name:                   "happy path: returns pokemon given an id",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}},
			input:                  "2",
			expected:               model.Pokemon{ID: 2, Name: "ivysaur"},
			err:                    nil,
		},
		{
			name:                   "not happy path: the csv file is empty",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}},
			input:                  "2",
			expected:               model.Pokemon{},
			err:                    errors.New("The csv file is empty"),
		},
		{
			name:                   "not happy path: no pokemon with given id exits",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}},
			input:                  "7",
			expected:               model.Pokemon{},
			err:                    errors.New("No pokemon with given id exits"),
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		fmt.Println(test.name)
		dataRepository := new(MockedRepository)
		dataRepository.On("GetAllPokemonData").Return(test.getAllPokemonDataValue, nil)
		pokemonservice := New(dataRepository)
		pokemon, err := pokemonservice.GetPokemonByID(test.input)
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
		externalAPIResponse model.Pokemon
		externalAPIError    error
		savePokemonResponse error
		input               string
		expected            model.Pokemon
		err                 error
	}{
		{
			name:                "happy path: call external api. Get its response, save the response in csv, finally return response.",
			input:               "2",
			externalAPIResponse: model.Pokemon{ID: 9, Name: "blastoise"},
			externalAPIError:    nil,
			savePokemonResponse: nil,
			expected:            model.Pokemon{ID: 9, Name: "blastoise"},
			err:                 nil,
		},
		{
			name:                "not happy path: Called to external api failed.",
			input:               "10",
			externalAPIResponse: model.Pokemon{},
			externalAPIError:    errors.New("The HTTP request failed with error"),
			savePokemonResponse: nil,
			expected:            model.Pokemon{},
			err:                 errors.New("The HTTP request failed with error"),
		},
		{
			name:                "not happy path: Problem saving data to the csv",
			input:               "2",
			externalAPIResponse: model.Pokemon{ID: 9, Name: "blastoise"},
			externalAPIError:    nil,
			savePokemonResponse: errors.New("Problem saving data to the csv"),
			expected:            model.Pokemon{ID: 9, Name: "blastoise"},
			err:                 errors.New("Problem saving data to the csv"),
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		fmt.Println(test.name)
		dataRepository := new(MockedRepository)
		dataRepository.On("GetExternalAPI", test.input).Return(test.externalAPIResponse, test.externalAPIError)
		dataRepository.On("SavePokemon").Return(test.savePokemonResponse)
		pokemonservice := New(dataRepository)
		pokemon, err := pokemonservice.ExternalPokemon(test.input)
		if err != nil {
			assert.NotNil(test.err)
		} else {
			assert.Equal(test.expected, pokemon)
		}
	}
}

func TestConcurrentPokemonList(t *testing.T) {
	tests := []struct {
		name                   string
		getAllPokemonDataValue [][]string
		typeinput              string
		itemsinput             string
		itemsperworkerinput    string
		expected               int
		err                    error
	}{
		{
			name: "happy path: Get list of pokemon of items size.",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}, {"4", "bulbasaur"}, {"5", "ivysaur"}, {"6", "venusaur"},
				{"7", "bulbasaur"}, {"8", "ivysaur"}, {"9", "venusaur"}, {"10", "bulbasaur"}, {"11", "ivysaur"}, {"12", "venusaur"}, {"13", "bulbasaur"}, {"14", "ivysaur"}, {"15", "venusaur"},
				{"16", "ivysaur"}, {"17", "venusaur"}, {"18", "bulbasaur"}, {"19", "ivysaur"}, {"20", "venusaur"}},
			typeinput:           "odd",
			itemsinput:          "5",
			itemsperworkerinput: "5",
			expected:            5,
			err:                 nil,
		},
		{
			name: "happy path: The channel response has less posible values (EOF) then wanted items response",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}, {"4", "bulbasaur"}, {"5", "ivysaur"}, {"6", "venusaur"},
				{"7", "bulbasaur"}, {"8", "ivysaur"}, {"9", "venusaur"}, {"10", "bulbasaur"}, {"11", "ivysaur"}, {"12", "venusaur"}, {"13", "bulbasaur"}, {"14", "ivysaur"}, {"15", "venusaur"},
				{"16", "ivysaur"}, {"17", "venusaur"}, {"18", "bulbasaur"}, {"19", "ivysaur"}, {"20", "venusaur"}},
			typeinput:           "odd",
			itemsinput:          "11",
			itemsperworkerinput: "9",
			expected:            10,
			err:                 nil,
		},
		{
			name: "error path: Problem casting items.",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}, {"4", "bulbasaur"}, {"5", "ivysaur"}, {"6", "venusaur"},
				{"7", "bulbasaur"}, {"8", "ivysaur"}, {"9", "venusaur"}, {"10", "bulbasaur"}, {"11", "ivysaur"}, {"12", "venusaur"}, {"13", "bulbasaur"}, {"14", "ivysaur"}, {"15", "venusaur"},
				{"16", "ivysaur"}, {"17", "venusaur"}, {"18", "bulbasaur"}, {"19", "ivysaur"}, {"20", "venusaur"}},
			typeinput:           "odd",
			itemsinput:          "10.5",
			itemsperworkerinput: "15",
			expected:            10,
			err:                 errors.New("Problem casting items"),
		},
		{
			name: "error path: Problem casting items per worker.",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}, {"4", "bulbasaur"}, {"5", "ivysaur"}, {"6", "venusaur"},
				{"7", "bulbasaur"}, {"8", "ivysaur"}, {"9", "venusaur"}, {"10", "bulbasaur"}, {"11", "ivysaur"}, {"12", "venusaur"}, {"13", "bulbasaur"}, {"14", "ivysaur"}, {"15", "venusaur"},
				{"16", "ivysaur"}, {"17", "venusaur"}, {"18", "bulbasaur"}, {"19", "ivysaur"}, {"20", "venusaur"}},
			typeinput:           "odd",
			itemsinput:          "10",
			itemsperworkerinput: "15.5",
			expected:            10,
			err:                 errors.New("Problem casting items per worker"),
		},
		{
			name:                   "error path: The csv file is empty:",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}},
			typeinput:              "odd",
			itemsinput:             "5",
			itemsperworkerinput:    "5",
			expected:               5,
			err:                    errors.New("The csv is empty"),
		},
		{
			name:                   "error path: The csv file is empty:",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}},
			typeinput:              "odd",
			itemsinput:             "5",
			itemsperworkerinput:    "5",
			expected:               5,
			err:                    errors.New("The csv is empty"),
		},
		{
			name:                   "error path: Problem casting id:",
			getAllPokemonDataValue: [][]string{{"ID", "Name"}, {"1.3", "bulbasaur"}, {"2", "ivysaur"}, {"3", "venusaur"}},
			typeinput:              "odd",
			itemsinput:             "5",
			itemsperworkerinput:    "5",
			expected:               5,
			err:                    errors.New("Problem casting id"),
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		fmt.Println(test.name)
		dataRepository := new(MockedRepository)
		dataRepository.On("GetAllPokemonData").Return(test.getAllPokemonDataValue, nil)
		pokemonservice := New(dataRepository)
		pokemonList, err := pokemonservice.ConcurrentPokemonList(test.typeinput, test.itemsinput, test.itemsperworkerinput)
		if err != nil {
			assert.NotNil(test.err)
		} else {
			assert.Equal(test.expected, len(pokemonList))
		}
	}
}
