package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// Librerías con las funciones de authentication
	"github.com/jantome/apilogin/authentication"
	"github.com/jantome/apilogin/environment"

	//Librería db2
	"github.com/jantome/apilogin/db2"
	//Librería de Log's
	"github.com/jantome/apilogin/logs"
	//Libreria con los struct (copy)
	"github.com/jantome/apilogin/structs"
)

//creamos la variable accesible para todos para los log's
var (
	erroLog       string = ""
	terrorwarning string = "w"
	terrorinfo    string = "i"
	terrorerro    string = "e"
)

//HandlerLogin función principal de login que recibira user y pass
func HandlerLogin(response http.ResponseWriter, request *http.Request) {
	//Evaluamos el tipo de petición que se esta realizando solo se permite GET
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para validar/usuario y pass y generar token en caso de que sea correcto
		getLogin(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		optionsLogin(response, request)
	default:
		erroLog = fmt.Sprintf("Not implemented Method %s", request.Method)
		logs.GrabaLog(nil, erroLog, terrorinfo)
		http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return

	}
}

//getLogin, Get con los datos con los que se  quiere logar
func getLogin(response http.ResponseWriter, request *http.Request) {
	enableCors(&response)
	//como por el momento solo postman admite en get json, el usuario y pass lo metemos como
	//seguridad basica.
	usuario, password, _ := request.BasicAuth()

	//Comprobamos que estan informados el usuario y la password
	if usuario == "" || usuario == " " {
		logs.GrabaLog(nil, "Mandatory data User", terrorinfo)
		http.Error(response, "Mandatory data User", http.StatusBadRequest)
		return
	}

	//Comprobamos que estan informados el usuario y la password
	if password == "" || password == " " {
		logs.GrabaLog(nil, "Mandatory data Password", terrorinfo)
		http.Error(response, "Mandatory data User", http.StatusBadRequest)
		return
	}
	//montamos la query
	query := fmt.Sprintf("SELECT * FROM batch_usuarios WHERE Usuario='%s' AND Password='%s'", usuario, password)
	result, err := db2.EjecutaQuery(query)
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	//Como la estructura de la tabla es la misma que la del json de entrada creamos otra variable para aplantilalr la taba
	var batchUser structs.BatchUser
	//Realizamos lectura del registro (Solo debe de existir 1)
	result.Next()
	err = result.Scan(&batchUser.Usuario, &batchUser.Password, &batchUser.Rol)
	//Controlamos el error de montar el dato en la structura
	if err != nil {
		logs.GrabaLog(err, "Unregistered User", terrorwarning)
		http.Error(response, "Unregistered User", http.StatusUnauthorized)
		return
	}
	//Una vez el usuario esta autentificado, generamos el Token
	valorToken, err := authentication.GenerateJWT(batchUser.Usuario, batchUser.Rol)
	//Comprobamos error en la generación del token
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInsufficientStorage)
		return
	}
	//creamos una structura de respuesta del token para mostrar de salida un json con el token
	var token structs.ResponseToken
	//movemos el valor de token a la structur del json
	token.Token = valorToken
	//generamos el json
	JsResponser, err := json.Marshal(token)

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

//HandlerViewToken función para comprobar el token
func HandlerViewToken(response http.ResponseWriter, request *http.Request) {
	//Evaluamos el tipo de petición que se esta realizando solo se permite GET
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para validar/usuario y pass y generar token en caso de que sea correcto
		getToken(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		optionsViewToken(response, request)
	default:
		erroLog = fmt.Sprintf("Not implemented Method %s", request.Method)
		logs.GrabaLog(nil, erroLog, terrorinfo)
		http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return

	}
}

//getToken, accion del get para endpoint Token (comprobación si es valido y si esta activo)
func getToken(response http.ResponseWriter, request *http.Request) {
	//habilitamos cors
	enableCors(&response)
	//Recuperamos de la cabecera el token
	reqToken := request.Header.Get("Authorization")
	//Comprobamos si tiene valor
	if reqToken == "" {
		logs.GrabaLog(nil, "Token Empty", terrorerro)
		http.Error(response, "Mandatory Token", http.StatusNetworkAuthenticationRequired)
		return
	}
	//LLamamos a la funcion que esta en authentication para comprobar el token y devolvera una structura
	//claims para tener el rol del usuario
	rol, descError, err := authentication.CompruebaToken(reqToken)
	//comprobamos si tenemos error
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	//comprobamos si tuviesemos algun error en descripcion (validaciones token)
	if descError != "" {
		logs.GrabaLog(nil, descError, terrorinfo)
		http.Error(response, descError, http.StatusBadRequest)
		return
	}
	//si todo es correcto, creamos un json con la misma estructura de la tabla (ya que los vacios los omite)
	//montamos json
	JsResponser, err := json.Marshal(rol)

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

//HandlerEnv para recuperar las variables de sistema (solo autorizados tiene que venir un token)
func HandlerEnv(response http.ResponseWriter, request *http.Request) {
	//Evaluamos la petición para saber que estamos haciendo
	switch request.Method {
	case "GET":
		//hacemos funciones para que el código sea legible
		getEnv(response, request)
	default:
		http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return
	}
}

//getEnv recupera las variables de entorno, siempre y cuando el usuarios este autorizado
func getEnv(response http.ResponseWriter, request *http.Request) {
	//Recuperamos de la cabecera el token
	reqToken := request.Header.Get("Authorization")
	//Comprobamos si tiene valor
	if reqToken == "" {
		logs.GrabaLog(nil, "Token Empty", terrorerro)
		http.Error(response, "Mandatory Token", http.StatusNetworkAuthenticationRequired)
		return
	}
	//LLamamos a la funcion que esta en authentication para comprobar el token y devolvera una structura
	//claims del cual obtendremos el rol
	datos, descError, err := authentication.CompruebaToken(reqToken)
	//comprobamos si tenemos error
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//Cualquier error damos no autorizado (401)
	if descError != "" {
		logs.GrabaLog(nil, descError, terrorinfo)
		http.Error(response, descError, http.StatusUnauthorized)
		return
	}

	//Si tenemos permisos recuperamos las variables de sistema y montamos el json
	if datos.Rol == "admin" {
		//recuperamos las variables de sistema poniendolas en la copy de la cual montara el json
		var envJSON structs.EnvJSON
		envJSON.Dbhost, _ = os.LookupEnv("DB_HOST")
		envJSON.Dbuser, _ = os.LookupEnv("DB_USER")
		envJSON.Dbpassword, _ = os.LookupEnv("DB_PASSWORD")
		envJSON.Dbdatabase, _ = os.LookupEnv("DB_DATABASE")
		envJSON.Servport, _ = os.LookupEnv("SERV_PORT")
		envJSON.ServportSSL, _ = os.LookupEnv("SERV_PORT_SSL")
		envJSON.Sersafe, _ = os.LookupEnv("SERV_SAFE")
		envJSON.Logactivate, _ = os.LookupEnv("LOG_ACTIVATE")
		envJSON.TokenLife, _ = os.LookupEnv("TOKEN_LIFE")

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

	//si no es admin damos usuario no autorizado
	if datos.Rol != "admin" {
		erroLog = fmt.Sprintf("Usuario no permitido BBOO/Server/ENV. User %s, tipo %s", datos.User, datos.Rol)
		logs.GrabaLog(nil, erroLog, terrorinfo)
		http.Error(response, "", http.StatusUnauthorized)
		return

	}
}

//HandlerEnvRefresh para actualización de variables en tipo de ejecuión
func HandlerEnvRefresh(response http.ResponseWriter, request *http.Request) {
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

//GetRefresh para acutalizar en tiempo de ejecución las variables siempre y cuando el usuario tenga permisos
func getRefresh(response http.ResponseWriter, request *http.Request) {
	//Recuperamos de la cabecera el token
	reqToken := request.Header.Get("Authorization")
	//Comprobamos si tiene valor
	if reqToken == "" {
		logs.GrabaLog(nil, "Token Empty", terrorerro)
		http.Error(response, "Mandatory Token", http.StatusNetworkAuthenticationRequired)
		return
	}
	//LLamamos a la funcion que esta en authentication para comprobar el token y devolvera una structura
	//claims del cual obtendremos el rol
	datos, descError, err := authentication.CompruebaToken(reqToken)
	//comprobamos si tenemos error
	if err != nil {
		logs.GrabaLog(err, "", terrorerro)
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	//Cualquier error damos no autorizado (401)
	if descError != "" {
		logs.GrabaLog(nil, descError, terrorinfo)
		http.Error(response, descError, http.StatusUnauthorized)
		return
	}

	//Comprobamos si puede ejecutar refresh
	if datos.Rol == "admin" {
		info := fmt.Sprintf("Refresh environment for user: %s", datos.User)
		//muestra en consola, y graba en log, es una información lo suficientemente importante,
		//como para que no se pierda
		log.Println(info)
		logs.GrabaLog(nil, info, terrorinfo)
		//Limpiamos las variables de entorno, para volver a cargarlas de nuevo
		os.Clearenv()
		//vuelve a ejecutar environment para cargar las variables de entorno
		environment.Loadenvironment()
		//respondemos que se realizo correctamente
		response.WriteHeader(http.StatusAccepted)
		return
	}

	//si no es admin damos usuario no autorizado
	if datos.Rol != "admin" {
		erroLog = fmt.Sprintf("Usuario no permitido BBOO/Server/ENV. User %s, tipo %s", datos.User, datos.Rol)
		logs.GrabaLog(nil, erroLog, terrorinfo)
		http.Error(response, "", http.StatusUnauthorized)
		return

	}
}

//HandlerTest comprobación del estado de la api (No contiene seguridad)
func HandlerTest(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		var testJSON structs.TestJSON
		//comprobamos conexión con una select normal
		query := "SELECT * FROM batch_usuarios"
		_, err := db2.EjecutaQuery(query)
		//Controlamos el error
		if err != nil {
			testJSON.Status = "KO"
			testJSON.ConexDb2 = err.Error()
		} else {
			//si no tiene error
			testJSON.Status = "OK"
			testJSON.ConexDb2 = "OK"
		}

		//generamos el Json
		JsResponser, err := json.Marshal(testJSON)
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
	default:
		http.Error(response, "Not Implemented", http.StatusNotImplemented)
		return
	}
}

//optionsLogin para mostrar lo que se acepta
func optionsLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Authorization")
	return
}

//enableCors, para habilitar los cors
func enableCors(response *http.ResponseWriter) {
	(*response).Header().Set("Access-Control-Allow-Origin", "*")
}

//optionsViewToken para mostrar lo que se acepta
func optionsViewToken(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	response.Header().Set("Access-Control-Allow-Headers", "Authorization")
	return
}
