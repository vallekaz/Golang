package environment

import (
	//Librería de sistema para cargar las variables

	"fmt"
	"log"
	"os"

	//Librería para las variables de entorno
	"github.com/joho/godotenv-master"
)

//Loadenvironment carga de las variables de sistema
func Loadenvironment(entorno string) {
	fmt.Println("entorno", entorno)
	//Abrimos fichero
	//si variable entorno = local ponemos el nombre de la carpeta principal
	if entorno == "local" {
		err := godotenv.Load("apilogin.env")
		//controlamos error en la apertura del fichero
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		// en caso de que la variable no este activada ponemos la ruta de linux
		err := godotenv.Load("/ejecutable/online/env/apilogin.env")
		//controlamos error en la apertura del fichero
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	//cargamos las variables
	//******DB2******//
	os.LookupEnv("DB_HOST")
	os.LookupEnv("DB_USER")
	os.LookupEnv("DB_PASSWORD")
	os.LookupEnv("DB_DATABASE")
	//******DB2******//

	//******SERVIDOR******//
	os.LookupEnv("SERV_PORT")
	os.LookupEnv("LOG_ACTIVATE")
	os.LookupEnv("PATH_LOG")
	//******SERVIDOR******//

}
