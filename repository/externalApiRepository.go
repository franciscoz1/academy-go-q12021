package repository

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//ExternalAPIRepository Struct that has all needed functions to interact with an
//external API
type ExternalAPIRepository struct {
}

//InterfaceExternalAPIRepository Interface with all functions available when calling external API
type InterfaceExternalAPIRepository interface {
	GetExternalAPI(id string) (models.Pokemon, error)
}

//GetExternalAPI calls an external api
func (e ExternalAPIRepository) GetExternalAPI(id string) (models.Pokemon, error) {
	var pokemon models.Pokemon
	response, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + id)
	if err != nil {
		return pokemon, err
	}
	data, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(data, &pokemon)
	return pokemon, err
}
