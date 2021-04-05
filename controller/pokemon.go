package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"Users/francisco.zamudio/projects/academy-go-q12021/model"

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
	ConcurrentPokemonList(parity string, items string, ipw string) ([]model.Pokemon, error)
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

//HandlerGetConcurrentPokemonList  given a type, number of items and item per workers
//return a list of pokemon.
func (c Controller) HandlerGetConcurrentPokemonList(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	if vars["type"] != "even" && vars["type"] != "odd" {
		err := errors.New("The value of 'type' can only be 'even' or 'odd'")
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if vars["items"] == "0" {
		err := errors.New("The value of 'items' can not be 0")
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pokemon, err := c.service.ConcurrentPokemonList(vars["type"], vars["items"], vars["ipw"])
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
