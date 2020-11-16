package routes

import (
	"net/http"

	"github.com/jantome/apicm/actions"
)

//Routes rutas permitidas cuando arranque el servidor
func Routes() {
	//Endpoint para la tabla de ejecucion
	http.HandleFunc("/cm/v1/ejecuciones", actions.HandlerEjecucion)
	//Endpoint para la tabla de condiciones de entrada  (ponemos la última / ya que recibiremos el ID)
	http.HandleFunc("/cm/v1/ejecuciones/condicionin/", actions.HandlerCondicionin)
	//Endpoint para la tabla de condiciones de salida  (ponemos la última / ya que recibiremos el ID)
	http.HandleFunc("/cm/v1/ejecuciones/condicionout/", actions.HandlerCondicionout)
	//Endpoint para la tabla de planificacion
	http.HandleFunc("/cm/v1/planificacion", actions.HandlerPlanificacion)
}
