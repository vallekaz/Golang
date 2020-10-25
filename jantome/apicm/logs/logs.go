package logs

import (
	"log"
	"os"
)

/*Definimos los posible serrores que vamos a tener*/
var (
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

//GrabaLog recibe de entrada el error para grabar en el ficheor de log
func GrabaLog(err2 error, descripcion string, tipo string) {
	//Solo grabamos si el log esta activado
	logActivate, _ := os.LookupEnv("LOG_ACTIVATE")
	switch logActivate {
	case "S":
		//apertura del fichero de log's en caso de no existir lo crea
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		//Estructura de como se mostrar en el fichero
		infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
		errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

		//Comprobamos el tipo de error para ver que tenemos que guardar
		//Warning
		if tipo == "w" {
			if err2 != nil && descripcion != "" {
				warningLogger.Println(err2, descripcion)
			}

			if err2 != nil {
				warningLogger.Println(err2)
			}

			if descripcion != "" {
				warningLogger.Println(descripcion)
			}

		}

		//Info
		if tipo == "i" {
			if err2 != nil && descripcion != "" {
				infoLogger.Println(err2, descripcion)
			}

			if err2 != nil {
				infoLogger.Println(err2)
			}

			if descripcion != "" {
				infoLogger.Println(descripcion)
			}
		}

		//Error
		if tipo == "e" {
			if err2 != nil && descripcion != "" {
				errorLogger.Println(err2, descripcion)
			}

			if err2 != nil {
				errorLogger.Println(err2)
			}

			if descripcion != "" {
				errorLogger.Println(descripcion)
			}
		}

	}
}
