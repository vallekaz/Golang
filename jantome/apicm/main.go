package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jantome/apicm/environment"
	"github.com/jantome/apicm/routes"
)

//Variables globales
var (
	//recuperamos el entorno de ejecucion mediante flags para saber las rutas
	entorno = flag.String("entorno", "", "entorno de ejecución")
)

//Init al arrancar
func init() {
	//cargamos las rutas permitidas
	routes.Routes()
}

//funcion principal
func main() {
	flag.Parse()
	//fmt.Println("entorno", *entorno)
	fmt.Println("Arrancando servidor...")
	//Cargamos la variables de entorno
	environment.Loadenvironment(*entorno)
	port, _ := os.LookupEnv("SERV_PORT")
	portSSL, _ := os.LookupEnv("SERV_PORT_SSL")
	serSafe, _ := os.LookupEnv("SERV_SAFE")
	switch serSafe {
	//servidor seguro
	case "S":
		fmt.Println("Port safe: ", portSSL)
		fmt.Println("Starting secure server")
		//Arrancamos con seguridad SSL
		err := http.ListenAndServeTLS(portSSL, "cert.pem", "key.pem", nil)
		//controlamos el error de listener
		if err != nil {
			log.Fatal(err)
		}
	//Sin servidor seguro
	case "N":
		fmt.Println("Port no safe: ", port)
		fmt.Println("Starting unsecure server")
		//Ponemos a escuchar el servidor (Arrancamos, con las rutas ya montadas)
		err := http.ListenAndServe(port, nil)
		//Arracamos sin serguridad SSL
		if err != nil {
			log.Fatal(err)
		}
	//Arrancar ambos servidores seguro/no seguro
	case "A":
		//Arrancamos con go para primero arrancar sin seguridad y luego funcion con seguridad
		go func() {
			fmt.Println("Port safe: ", portSSL)
			fmt.Println("Starting secure server")
			//Arrancamos con seguridad SSL
			err := http.ListenAndServeTLS(portSSL, "cert.pem", "key.pem", nil)
			//Arracamos sin serguridad SSL
			if err != nil {
				log.Fatal(err)
			}

		}()
		fmt.Println("Port no safe: ", port)
		fmt.Println("Starting unsecure server")
		//Ponemos a escuchar el servidor (Arrancamos, con las rutas ya montadas)
		err := http.ListenAndServe(port, nil)
		//Arracamos sin serguridad SSL
		if err != nil {
			log.Fatal(err)
		}
	}
}
