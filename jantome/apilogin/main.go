package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	//Librería para cargar las variables de entorno
	"github.com/jantome/apilogin/environment"

	//Librería con las routas del servidor
	"github.com/jantome/apilogin/routes"
)

//main función principal
func main() {
	fmt.Println("Arrancando servidor...")
	//Cargarmos variables de entorno
	environment.Loadenvironment()
	port, _ := os.LookupEnv("SERV_PORT")
	portSSL, _ := os.LookupEnv("SERV_PORT_SSL")
	serSafe, _ := os.LookupEnv("SERV_SAFE")
	fmt.Println("Seguridad: ", serSafe)
	//Generamos los Endpoint
	routes.Routes()
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
