package controllers

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/repository"
	"Users/francisco.zamudio/projects/academy-go-q12021/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//HandlerGetPokemonByID Connects to a csv repository, then calls GetPokemonByID service
//if correct it return the corresponding pokemon, else it return an error
func HandlerGetPokemonByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	dataRepository := repository.DataRepository{}
	pokemonInteractor := services.PokemonInteractor{Repository: dataRepository}
	pokemon, err := pokemonInteractor.GetPokemonByID(vars["id"])
	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	out, err := json.Marshal(pokemon)

	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

//HandlerGetExternalPokemonByID Calls an external API,then saves the result in the csv
//finally returns the response.
func HandlerGetExternalPokemonByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	dataRepository := repository.DataRepository{}
	externalAPIRepository := repository.ExternalAPIRepository{}
	pokemonInteractor := services.PokemonInteractor{Repository: dataRepository, ExternalapiRepository: externalAPIRepository}
	pokemon, err := pokemonInteractor.ExternalPokemon(vars["id"])
	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	out, err := json.Marshal(pokemon)

	if err != nil {
		fmt.Println("An error encountered ::", err)
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}
