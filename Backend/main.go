package main

import (
	"fmt"
	"log"
	"net/http"
	"paquetes/API"

	httplogger "github.com/jesseokeya/go-httplogger"
	"github.com/rs/cors"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	})
	router := API.NewRouter()
	fmt.Println("Sever running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", httplogger.Golog(c.Handler(router))))
	// analizador.InputHandler()
}