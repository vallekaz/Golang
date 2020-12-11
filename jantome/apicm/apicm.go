package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jantome/apicm/environment"
	"github.com/jantome/apicm/routes"
	"github.com/jantome/apicm/structs"
	"github.com/onlinearq/online"
)

//Variables globales
var (
	//recuperamos el entorno de ejecucion mediante flags para saber las rutas
	//entorno    = flag.String("entorno", "", "entorno de ejecución")
	nameserver = "apicm"
)

//Init al arrancar
func init() {
	//cargamos las rutas permitidas
	routes.Routes()
}

//funcion principal
func main() {
	//parseamos los flag que podamos recibir
	flag.Parse()
	//Cargamos la variables de entorno
	environment.Loadenvironment(*structs.Entorno)
	port, _ := os.LookupEnv("SERV_PORT")
	portSSL, _ := os.LookupEnv("SERV_PORT_SSL")
	serSafe, _ := os.LookupEnv("SERV_SAFE")
	certpem, _ := os.LookupEnv("PATH_CERT")
	keypem, _ := os.LookupEnv("PATH_KEY")
	switch serSafe {
	//servidor seguro
	case "S":
		fmt.Println("Port safe: ", portSSL)
		fmt.Println("Starting secure server")
		//grabamos en el log el arranque del servidor
		online.Start(nameserver, *structs.Entorno, portSSL, "Secure Server")
		//Arrancamos con seguridad SSL
		err := http.ListenAndServeTLS(portSSL, certpem, keypem, nil)
		//controlamos el error de listener
		if err != nil {
			log.Fatal(err)
		}
	//Sin servidor seguro
	case "N":
		fmt.Println("Port no safe: ", port)
		fmt.Println("Starting unsecure server")
		//grabamos en el log el arranque del servidor
		online.Start(nameserver, *structs.Entorno, portSSL, "Unsecure Server")
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
			//grabamos en el log el arranque del servidor
			online.Start(nameserver, *structs.Entorno, portSSL, "Secure Server")

			//Arrancamos con seguridad SSL
			err := http.ListenAndServeTLS(portSSL, certpem, keypem, nil)
			//Arracamos sin serguridad SSL
			if err != nil {
				log.Fatal(err)
			}

		}()
		fmt.Println("Port no safe: ", port)
		fmt.Println("Starting unsecure server")
		//grabamos en el log el arranque del servidor
		online.Start(nameserver, *structs.Entorno, portSSL, "Unsecure Server")
		//Ponemos a escuchar el servidor (Arrancamos, con las rutas ya montadas)
		err := http.ListenAndServe(port, nil)
		//Arracamos sin serguridad SSL
		if err != nil {
			log.Fatal(err)
		}
	}
}
