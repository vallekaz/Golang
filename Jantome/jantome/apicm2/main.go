package main

import (
	//Librería para las variables de entorno
	"os"
	//libreria variables de entorno
	"github.com/jantome/apicm2/environment"
	//Libería de Log's
	"log"
	//Librería para HTTP
	"net/http"
	//Librería para las rutas
	"github.com/jantome/apicm2/routes"
)

//función principal
func main() {
	log.Println("Servidor Arrancado")
	//Cargar variables de entorno
	environment.Loadenvironment()
	port, _ := os.LookupEnv("SERV_PORT")
	//cargamos las rutas permitidas
	routes.Routes()
	//Ponemos a escuchar el servidor (Arrancamos, con las routas ya montadas)
	err := http.ListenAndServe(port, nil)
	//controlamos el error de listener
	if err != nil {
		log.Fatal(err)
	}

}
