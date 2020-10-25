package routes

import (
	//Librería con las acciones por cada ruta
	"github.com/jantome/apilogin/actions"
	//librería de http
	"net/http"
)

//Routes funcion para arrancar el servidor
func Routes() {
	//Endpoint GET para el Login de usuario
	http.HandleFunc("/bboo/v1/login", actions.HandlerLogin)
	//Endpoint GET para comprobar el token
	http.HandleFunc("/bboo/v1/token", actions.HandlerViewToken)
	//Endpoint GET para las variables de sistema (autentificado por token)
	http.HandleFunc("/bboo/v1/server/env", actions.HandlerEnv)
	//Endpoint GET para la actualización en tiempo de ejecución de las variables de sistema
	http.HandleFunc("/bboo/v1/server/env/refresh", actions.HandlerEnvRefresh)
	//Endpoint GET para comprobar el estado del servidor
	http.HandleFunc("/bboo/v1/server/test", actions.HandlerTest)
}
