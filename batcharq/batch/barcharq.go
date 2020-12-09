package batch

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/batcharq/batch/environment"
)

var (
	//formato de los logs
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	//obtenemos la hora actual que para poder formatearla
	date = time.Now()
	//formateamos para guardar el log con formato nombre fecha-hora (XXXDxxxxxxHXX:XX:XX)
	anno       = fmt.Sprintf("%d", date.Year())
	mes        = fmt.Sprintf("%d", date.Month())
	dia        = fmt.Sprintf("%d", date.Day())
	hora       = time.Now().Format("1504")
	second     = fmt.Sprintf("%d", date.Second())
	formateado = fmt.Sprintf("D%s%s%sH%s%s", anno[2:4], mes, dia, hora, second)
)

//Start funcion que se utilizara para generar el fichero de log nada más empezar un proceso batch
func Start(nombre string, entorno string) {
	//apertura del fichero, nos devuelve file para que podamos cerrarlo cuando finalicemos la funcion
	file := openFicher(nombre, entorno, true)
	//Comprobamos el tipo de error para ver que tenemos que guardar
	infoLogger.Println("Start", nombre)
	//Cerramos fichero
	defer file.Close()
}

//OpenFicher funcion que se invocara en todos los sitios para la apertura del fichero
func openFicher(nombre string, entorno string, primera bool) (file *os.File) {
	//obtenemos la ruta donde se guardaran los logs
	//pathLog := "./" + nombre + formateado
	environment.Loadenvironment(entorno)
	pathLog, _ := os.LookupEnv("PATH_LOG")
	//OJO de manera temporal, vamos a tener el nombre del fichero sin año mes día. Es decir vamos a tener solo
	//un fichero por ejecucion, hasta que tengamos un historico de ejecuciones. Por lo que el nombre por el momento
	//sera siempre el mismo. Por lo que tenemos que hacer un delete primero
	//pathLog = pathLog + nombre + formateado
	pathLog = pathLog + nombre
	//si es la primera vez borramos el fichero
	if primera == true {
		err := os.Remove(pathLog)
		if err != nil {
			if os.IsNotExist(err) == false {
				log.Fatal(err)
			}
		}
	}
	//pathLog = pathLog + nombre + formateado
	/*pathLog, _ := os.LookupEnv("PATH_LOG")*/
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

//FinOk funcion para la finalización correcta del proceso batch
func FinOk(nombre string, entorno string) {
	//apertura del fichero, nos devuelve file para que podamos cerrarlo cuando finalicemos la funcion
	file := openFicher(nombre, entorno, false)
	infoLogger.Println("Finish",
		nombre,
		"Tiempo de ejecucion: ",
		time.Since(date))
	//Cerramos fichero
	defer file.Close()
}

// Impr funcion que servirar para imprimir toda la info que creamos necesaria. Por ejempo estadisticas o display's
func Impr(nombre string, inf string, tipo string, entorno string) {
	file := openFicher(nombre, entorno, false)
	switch tipo {
	//warning
	case "w":
		warningLogger.Println(inf)
	//info
	case "i":
		infoLogger.Println(inf)
	//error
	/*case "e":
	errorLogger.Println(inf)*/
	//default es info
	default:
		infoLogger.Println(inf)
	}
	defer file.Close()
}

//FinKo funcion para imprimir error cuando esta KO
func FinKo(nombre string, retorno string, descripcion string, entorno string) {
	file := openFicher(nombre, entorno, false)
	errorLogger.Println("\n¡¡¡¡¡¡Error!!!!!!\n",
		"Retorno: ",
		retorno,
		"\nDescripcion: ",
		descripcion,
		"\n¡¡¡¡¡¡Error!!!!!!")
	defer file.Close()

}
