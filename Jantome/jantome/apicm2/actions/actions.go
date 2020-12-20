package actions

import (
	//Librería de log
	"log"
	//Librería para fichero de log
	"github.com/jantome/apicm2/logs"
	//Librería de sistema
	"os"
	//Librería para las variables de entorno
	"github.com/jantome/apicm2/environment"
	//Librería con las opciones de db2
	"github.com/jantome/apicm2/db2"
	//Libreria FMT para el formateo de la query
	"fmt"
	//Libería de structs
	"github.com/jantome/apicm2/structs"
	//Librería para json
	"encoding/json"
	//Librería para http
	"net/http"
	//Librerías para SQL
	_ "github.com/go-sql-driver/mysql"
)

/*creamos las variables aqui para que todas las funciones, para que se pueda grabar la descripcion del erro en caso de llamar
a logs*/
var (
	descerror     string = ""
	terrorwarning string = "w"
	terrorinfo    string = "i"
	terrorerro    string = "e"
)

//HandlerListCalendarios Función para endpoint /prueba
func HandlerListCalendarios(response http.ResponseWriter, request *http.Request) {
	//Evaluamos la petición para saber que estamos haciendo
	switch request.Method {
	case "GET":
		//hacemos funciones para que el código sea legible
		getCalendario(response, request)
	case "POST":
		postCalendario(response, request)
	/*	http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return*/
	case "DELETE":
		delCalendario(response, request)
	/*	http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return*/
	case "PUT":
		putCalendario(response, request)
		/*	http.Error(response, "Not Implemented", http.StatusNotImplemented)
			return*/
	default:
		descerror = fmt.Sprintf("Not implemented Method %s", request.Method)
		logs.GrabaLog(nil, descerror, terrorinfo)
		http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return
	}

}

//getCalendario recupera de la tabla batch_calendario para el get ya se con id o sin el
func getCalendario(response http.ResponseWriter, request *http.Request) {

	//recuperamos variable de URL para el posible where en la select
	variable := variURL(request.RequestURI)

	//creamos la variable de query inicializada
	query := ""
	//Montamos la select dependiendo si tiene Where o no
	if variable != "" {
		query = fmt.Sprintf("SELECT * FROM batch_calendario WHERE nombre_calendario = '%s'", variable)
	} else {
		query = "SELECT * FROM batch_calendario"
	}

	//LLamamos a la funcion con la query que queremos realizar
	result, err := db2.EjecutaQuery(query)

	//Controlamos el error de la ejecución de la query
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//declaramos SW para saber si tenemos datos o no, dependiente de si entra en el for
	siDatos := false
	//con Next recorremos toda la query para mostrar los datos
	for result.Next() {
		siDatos = true
		//Creamos la variable de tipo calendario donde vamos a leer el resultaod
		var calendario structs.Calendario

		err = result.Scan(&calendario.Nombre, &calendario.Mes, &calendario.Dia, &calendario.Usuarioalta, &calendario.Fechacreacion, &calendario.Usuariomodif, &calendario.Fechamodif)
		//Controlamos el erro de montar el dato en la structura
		if err != nil {
			logs.GrabaLog(err, "", terrorerro)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		//montamos json
		JsResponser, err := json.Marshal(calendario)
		//controlamos el error de json
		if err != nil {
			logs.GrabaLog(err, "", terrorerro)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//devolvemos la respuesta
		response.Write(JsResponser)

	}
	//si no tenemos datos devolvemos error 404 (no encontrado)
	if siDatos == false {
		http.Error(response, "404 not found", http.StatusNotFound)
		return
	}
}

//variURL con la URL de entrada buscamos la parte variable desde el primer /
func variURL(url string) string {
	//primero sacamos el tamaño
	lenght := len(url)
	//quitamos la primera posición
	url2 := url[1:lenght]
	//volvemos a sacar el tamaño
	lenght = len(url2)
	//inicializamos a 1 ya que empieza por 0 la cuenta y si no luego devuelve "/" y variable
	lenght2 := 1
	//creamos un SW para no pasar dos veces por /
	encontrado := false
	//buscamos la posición de la "/"
	for i := 0; i < len(url2); i++ {
		//si la encontramos nos guardamos el valor de I
		if string(url2[i]) == "/" && encontrado == false {
			lenght2 = lenght2 + i
			encontrado = true
		}
	}
	//si la longitud de lenght2 no es mayor que 0, significa que no tiene variables
	if lenght2 != 1 {
		//Utilizamos desde lenght2 que es la "/" hasta el final lenght(que es lo que calculamos al principio)
		url2 = url2[lenght2:lenght]
	} else {
		// si es 1 significa que no encontro valor, por lo que la variable se devuelve vacia
		url2 = ""
	}
	//devolvemos la parte variable para poder montar el where
	return url2
}

//postCalendario post de la tabla calendario
func postCalendario(response http.ResponseWriter, request *http.Request) {

	//Utilizamos lo que viene en el cuerpo como json
	cuerpo := json.NewDecoder(request.Body)
	//ceramos la respuesta de body con defer para que se ejecute al final
	defer request.Body.Close()

	//creamos variable donde almacenaremos los datos recuperados en el json
	var calendario structs.Calendario

	//decodificamos
	err := cuerpo.Decode(&calendario)

	//comprobamos si tenemos error
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}
	// comprobamos los datos del body
	pasaValida, calendario := compruebaPostCalendar(calendario)

	if pasaValida != "" {
		http.Error(response, pasaValida, http.StatusBadRequest)
		return
	}
	//Montamos la query
	query := fmt.Sprintf("INSERT INTO batch_calendario VALUES('%s','%s','%s','%s',CURRENT_TIMESTAMP,'%s',CURRENT_TIMESTAMP)", calendario.Nombre, calendario.Mes, calendario.Dia, calendario.Usuarioalta, calendario.Usuariomodif)

	//ejecutamos query ponemos _ ya que no queremos la variable result
	_, err = db2.EjecutaQuery(query)

	//comprobamos error al realizar el insert
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

//compruebaPostCalendar realiza la comprobación previa al post en la tabla calendario
func compruebaPostCalendar(calendario structs.Calendario) (string, structs.Calendario) {
	pasaValida := ""
	calendario2 := calendario

	if calendario.Nombre == " " || calendario.Nombre == "" || &calendario.Nombre == nil {
		pasaValida = "Mandatory field Name"
		return pasaValida, calendario2
	}

	if calendario.Mes == " " || calendario.Mes == "" || &calendario.Mes == nil {
		calendario2.Mes = "*"
	}

	if calendario.Dia == " " || calendario.Dia == "" || &calendario.Dia == nil {
		calendario2.Dia = "*"
	}

	if calendario.Usuarioalta == " " || calendario.Usuarioalta == "" || &calendario.Usuarioalta == nil {
		pasaValida = "Mandatory field usuarioalta"
		return pasaValida, calendario2
	}

	if calendario.Usuariomodif == " " || calendario.Usuariomodif == "" || &calendario.Usuariomodif == nil {
		calendario2.Usuariomodif = calendario.Usuarioalta
	}

	return pasaValida, calendario2
}

//delCalendario realiza delete en la tabla calendario
func delCalendario(response http.ResponseWriter, request *http.Request) {

	//recuperamos de la url el ID que queremos eliminar
	variable := variURL(request.RequestURI)

	//comprobamos la variable que recupere y si esta vacia fallamos
	if variable == "" {
		http.Error(response, "Mandatory ID", http.StatusBadRequest)
		return
	}

	//en caso contrario montamos la query para el delete
	query := fmt.Sprintf("DELETE FROM batch_calendario WHERE nombre_calendario = '%s'", variable)
	_, err := db2.EjecutaQuery(query)
	//comprobamos el error
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//actualizamos el estado a 201 que es creado
	response.WriteHeader(http.StatusCreated)
}

//putCalendario realiza el put en la tabla calendario
func putCalendario(response http.ResponseWriter, request *http.Request) {
	//lo primero que hacemos es recuperar el json del body
	cuerpo := json.NewDecoder(request.Body)
	//ceramos la respuesta de body con defer para que se ejecute al final
	defer request.Body.Close()

	//creamos variable donde almacenaremos los datos recuperados en el json
	var calendario structs.Calendario

	//decodificamos
	err := cuerpo.Decode(&calendario)

	//comprobamos si tenemos error
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	//comprobara si existe el dato y que datos se van a actualizar
	valida, calendario := compruebaPutCalendario(calendario)

	//comprobamos si existiese un error y de que tipo es
	if valida != "" && valida != "Internal Server Error" {
		http.Error(response, valida, http.StatusBadRequest)
		return
	}

	if valida == "Internal Server Error" {
		logs.GrabaLog(err, valida, terrorerro)
		http.Error(response, valida, http.StatusInternalServerError)
		return
	}

	//una vez devuelto y sin errores realizamos la query

	query := fmt.Sprintf("UPDATE batch_calendario SET mes_ejecucion ='%s', dia_ejecucion ='%s', usuario_modif ='%s' WHERE nombre_calendario = '%s'", calendario.Mes, calendario.Dia, calendario.Usuariomodif, calendario.Nombre)

	//ejecutamos la query
	_, err = db2.EjecutaQuery(query)

	//comprobamos el error
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//actualizamos el estado de la petición a 204 que es actualizamos pero no tiene contenido
	response.WriteHeader(http.StatusNoContent)

}

//compruebaPutCalendario comprobación de lo recibido en el put para la tabla calendario
func compruebaPutCalendario(calendario structs.Calendario) (string, structs.Calendario) {
	valida := ""
	//comprobamos que nos llega el id, y que este existe
	if calendario.Nombre == " " || calendario.Nombre == "" || &calendario.Nombre == nil {
		valida = "Mandatory field Nombre"
		return valida, calendario
	}

	if calendario.Usuariomodif == " " || calendario.Usuariomodif == "" || &calendario.Usuariomodif == nil {
		valida = "Mandatory field UsuarioModif"
		return valida, calendario
	}

	//Una vez pasado los obligatorios, lo que hacemos es comprobar si el dato existe

	query := fmt.Sprintf("SELECT * FROM batch_calendario WHERE nombre_calendario = '%s'", calendario.Nombre)

	//ejecuamtos la query
	result, err := db2.EjecutaQuery(query)

	//comprobamos error
	if err != nil {
		valida = "Internal Server Error"
		return valida, calendario
	}

	// en result hacemos next, solo debe de existir un registro
	result.Next()
	//Creamos la variable de tipo calendario donde vamos a leer el resultaod
	var calendario2 structs.Calendario

	err = result.Scan(&calendario2.Nombre, &calendario2.Mes, &calendario2.Dia, &calendario2.Usuarioalta, &calendario2.Fechacreacion, &calendario2.Usuariomodif, &calendario2.Fechamodif)
	//Controlamos el erro de montar el dato en la structura
	if err != nil {
		valida = "Internal Server Error"
		return valida, calendario
	}

	//Ahora comprobamos del primer calendario de los datos que podemos actualizar si los hemos recibido para actualizar
	if calendario.Mes != " " && calendario.Mes != "" && &calendario.Mes != nil {
		calendario2.Mes = calendario.Mes
	}

	if calendario.Dia != " " && calendario.Dia != "" && &calendario.Dia != nil {
		calendario2.Dia = calendario.Dia
	}

	//el usuario ya lo hemos comprobado antes, por lo que directamente lo movemos
	calendario2.Usuariomodif = calendario.Usuariomodif

	//devolvemos el segundo calendario
	return valida, calendario2
}

//HandlerRefresh Función para endpoint /prueba
func HandlerRefresh(response http.ResponseWriter, request *http.Request) {

	//Evaluamos la petición para saber que estamos haciendo
	switch request.Method {
	case "GET":
		//hacemos funciones para que el código sea legible
		getRefresh(response, request)
	default:
		http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return
	}
}

//getRefresh para actualizar las variables de entorno
func getRefresh(response http.ResponseWriter, request *http.Request) {
	log.Println("Refresh environment")
	//Limpiamos las variables de entorno, para volver a cargarlas de nuevo
	os.Clearenv()
	//vuelve a ejecutar environment para cargar las variables de entorno
	environment.Loadenvironment()
	//respondemos que se realizo correctamente
	response.WriteHeader(http.StatusAccepted)
	return
}

//HandlerEnv para recuperar en un json la variable de entorno y su valor
func HandlerEnv(response http.ResponseWriter, request *http.Request) {
	//Creamos el campo
	var envJSON structs.EnvJSON
	//Rellenamos con la información
	envJSON.Dbhost, _ = os.LookupEnv("DB_HOST")
	envJSON.Dbuser, _ = os.LookupEnv("DB_USER")
	envJSON.Dbpassword, _ = os.LookupEnv("DB_PASSWORD")
	envJSON.Dbdatabase, _ = os.LookupEnv("DB_DATABASE")
	envJSON.Servport, _ = os.LookupEnv("SERV_PORT")
	envJSON.Logactivate, _ = os.LookupEnv("LOG_ACTIVATE")

	//montamos json
	JsResponser, err := json.Marshal(envJSON)

	//controlamos el error de json
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//creamos cabecera de respuesta
	response.Header().Set("Content-Type", "application/json")
	//devolvemos la respuesta
	response.Write(JsResponser)
}
