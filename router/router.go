package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

//ControllerHandlers list of handler methods
type ControllerHandlers interface {
	HandlerGetPokemonByID(w http.ResponseWriter, r *http.Request)
	HandlerGetExternalPokemonByID(w http.ResponseWriter, r *http.Request)
	HandlerGetConcurrentPokemonList(w http.ResponseWriter, r *http.Request)
}

//New create a router object
func New(c ControllerHandlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id}", c.HandlerGetPokemonByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/external/pokemon/{id}", c.HandlerGetExternalPokemonByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/pokemon", c.HandlerGetConcurrentPokemonList).Queries("type", "{type:[a-z]+}", "items", "{items:[0-9]+}", "ipw", "{ipw:[0-9]+}").Methods("GET", "OPTIONS")
	http.Handle("/", r)
	return r
}
