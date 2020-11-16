package main

import (
	"log"
	"net/http"

	"github.com/jantome/apicm/routes"
)

func init() {
	//cargamos las rutas permitidas
	routes.Routes()
}
func main() {
	port := ":8585"
	//Ponemos a escuchar el servidor, el nil es por que las routas ya las arrancamos
	err := http.ListenAndServe(port, nil)
	//controlamos el error
	if err != nil {
		log.Fatal(err.Error())
	}
}
