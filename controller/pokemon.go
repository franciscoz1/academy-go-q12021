package controller

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//Controller use to have access to the different use case functions.
type Controller struct {
	service Service
}

//Service Interface
type Service interface {
	GetPokemonByID(id string) (model.Pokemon, error)
	ExternalPokemon(id string) (model.Pokemon, error)
}

//New Create a instance of a Service
func New(s Service) *Controller {
	return &Controller{s}
}

//HandlerGetPokemonByID Connects to a csv repository, then calls GetPokemonByID service
//if correct it return the corresponding pokemon, else it return an error
func (c Controller) HandlerGetPokemonByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pokemon, err := c.service.GetPokemonByID(vars["id"])
	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	out, err := json.Marshal(pokemon)

	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(out))
}

//HandlerGetExternalPokemonByID Calls an external API,then saves the result in the csv
//finally returns the response.
func (c Controller) HandlerGetExternalPokemonByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	pokemon, err := c.service.ExternalPokemon(vars["id"])
	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	out, err := json.Marshal(pokemon)

	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(out))
}
