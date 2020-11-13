package routes

import (
	//librería de http
	"net/http"
	//libreria con las acciones que realizara
	"github.com/jantome/apicm/actions"
)

//Routes funcion para arrancar el servidor
func Routes() {
	//Endpoint GET calendario listado
	http.HandleFunc("/calendario", actions.HandlerListCalendarios)
	/*Endpoint GET calendario listado usamos la misma funcion para cuando tenemos ID, ya que lo recuperas
	de la url y lo montamos en el where */
	http.HandleFunc("/calendario/", actions.HandlerListCalendarios)
	/*Endpoint GET /server/refresh para actualizar las variables de entorno*/
	http.HandleFunc("/server/refresh/", actions.HandlerRefresh)
	/*Endpoint GET /server/env para recuperar las variables de entorno*/
	http.HandleFunc("/server/env/", actions.HandlerEnv)
}
