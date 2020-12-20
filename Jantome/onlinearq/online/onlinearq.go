package online

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/onlinearq/online/environment"
)

var (
	//formato de los logs
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

//Start para generar el fichero de log segun arranca el servidor online
func Start(nombre string, entorno string, puerto string, descripcion string) {
	file := openFicher(nombre, entorno)
	//Comprobamos el tipo de error para ver que tenemos que guardar
	infoLogger.Println("Start server:", nombre, "Port:", puerto[1:5], descripcion)
	//Cerramos fichero
	defer file.Close()
}

//OpenFicher apertura del fichero
func openFicher(nombre string, entorno string) (file *os.File) {
	//cargamos la variables de entorno, que nos indicaran donde guardar los log's
	environment.Loadenvironment(entorno)
	pathLog, _ := os.LookupEnv("PATH_LOG")
	//Tendremos un fichero de log por cada día, por lo que obtenemos la fecha del día
	date := time.Now()
	//Formateamos la fecha para el fichero con el formato Nombre/año/mes/dia
	anno := fmt.Sprintf("%d", date.Year())
	mes := fmt.Sprintf("%d", date.Month())
	dia := fmt.Sprintf("%d", date.Day())
	fechaformateada := fmt.Sprintf("%s%s%s", anno[2:4], mes, dia)
	//juntamos el nombre del servidro + la fecha formateada
	nombre = nombre + fechaformateada
	pathLog = pathLog + nombre
	//apertura del fichero o creacion del mismo en caso de no existir
	file, err := os.OpenFile(pathLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//Estructura de como se mostrar en el fichero
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	return file
}

//EjecutaError para cuando queramos grabar un Error
func EjecutaError(nombre string, entorno string, descripcion string, err error) {
	//abrimos el fichero
	file := openFicher(nombre, entorno)
	if err != nil {
		errorLogger.Println("Error:", descripcion, " ", err.Error())
	} else {
		errorLogger.Println("Error:", descripcion)
	}

	defer file.Close()
}

//EjecutaInfo para cuando queramos grabar un Error
func EjecutaInfo(nombre string, entorno string, descripcion string, err error) {
	//abrimos el fichero
	file := openFicher(nombre, entorno)
	if err != nil {
		infoLogger.Println(descripcion, " ", err.Error())
	} else {
		infoLogger.Println(descripcion)
	}

	defer file.Close()
}

//GeneraIDError para la generacion del id del error con la fecha/timestamp
func GeneraIDError() (id string) {
	date := time.Now()
	dia := fmt.Sprintf("%d", date.Day())
	nanose := fmt.Sprintf("%d", date.Nanosecond())
	id = dia + nanose
	return id
}
