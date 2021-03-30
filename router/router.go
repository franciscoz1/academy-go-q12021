package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

//ControllerHandlers list of handler methods
type ControllerHandlers interface {
	HandlerGetPokemonByID(w http.ResponseWriter, r *http.Request)
	HandlerGetExternalPokemonByID(w http.ResponseWriter, r *http.Request)
}

//New create a router object
func New(c ControllerHandlers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id}", c.HandlerGetPokemonByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/external/pokemon/{id}", c.HandlerGetExternalPokemonByID).Methods("GET", "OPTIONS")
	http.Handle("/", r)
	return r
}
