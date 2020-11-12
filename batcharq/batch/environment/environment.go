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
		err := godotenv.Load("C:\\gopath\\src\\github.com\\batcharq\\batch\\batcharq.env")
		//controlamos error en la apertura del fichero
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		// en caso de que la variable no este activada ponemos la ruta de linux
		err := godotenv.Load("/ejecutable/online/env/batcharq.env")
		//controlamos error en la apertura del fichero
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	//******SERVIDOR******//
	os.LookupEnv("PATH_LOG")
	//******SERVIDOR******//

}
