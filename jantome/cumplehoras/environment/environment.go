package environment

import (
	//Librería de sistema para cargar las variables

	"log"
	"os"

	//Librería para las variables de entorno
	"github.com/joho/godotenv-master"
)

//Loadenvironment carga de las variables de sistema
func Loadenvironment(entorno string) {
	//Abrimos fichero
	//si variable entorno = local ponemos el nombre de la carpeta principal
	if entorno == "local" {
		err := godotenv.Load("cumplehoras.env")
		//controlamos error en la apertura del fichero
		if err != nil {
			log.Fatal("Error loading .env file", err.Error())
		}
	} else {
		// en caso de que la variable no este activada ponemos la ruta de linux
		err := godotenv.Load("/ejecutable/batch/env/cumplehoras.env")
		//controlamos error en la apertura del fichero
		if err != nil {
			log.Fatal("Error loading .env file", err.Error())
		}
	}

	//cargamos las variables
	//******DB2******//
	os.LookupEnv("DB_HOST")
	os.LookupEnv("DB_USER")
	os.LookupEnv("DB_PASSWORD")
	os.LookupEnv("DB_DATABASE")
	os.LookupEnv("FILE_HOUR")
	//******DB2******//
}
