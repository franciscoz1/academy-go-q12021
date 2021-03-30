package main

import (
	"Users/francisco.zamudio/projects/academy-go-q12021/controller"
	"Users/francisco.zamudio/projects/academy-go-q12021/repository"
	"Users/francisco.zamudio/projects/academy-go-q12021/router"
	"Users/francisco.zamudio/projects/academy-go-q12021/service"
	"log"
	"net/http"
)

func main() {
	pokemonRepository := repository.New()
	pokemonservice := service.New(pokemonRepository)
	pokemonController := controller.New(pokemonservice)
	router := router.New(pokemonController)
	log.Fatal(http.ListenAndServe(":8080", router))
}
