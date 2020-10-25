package environment

import (
	//Librería de sistema para cargar las variables
	"log"
	"os"
	//Librería para las variables de entorno
	"github.com/joho/godotenv-master"
)

//Loadenvironment carga de las variables de sistema
func Loadenvironment() {
	//Abrimos fichero
	err := godotenv.Load()
	//controlamos error en la apertura del fichero
	if err != nil {
		log.Fatal("Error loading .env file")
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
	//******SERVIDOR******//

}
