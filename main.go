package main

import (
	"log"
	"net/http"

	"Users/francisco.zamudio/projects/academy-go-q12021/controllers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/pokemon/{id}", controllers.HandlerGetPokemonByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/external/pokemon/{id}", controllers.HandlerGetExternalPokemonByID).Methods("GET", "OPTIONS")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
