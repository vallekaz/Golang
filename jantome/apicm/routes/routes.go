package routes

import (
	"net/http"

	"github.com/jantome/apicm/actions"
)

//Routes rutas permitidas cuando arranque el servidor
func Routes() {
	//Endpoint para la tabla de ejecucion
	http.HandleFunc("/cm/v1/ejecuciones", actions.HandlerEjecucion)
	//Endpoint para la tabla de ejecucion acabado en / para recibir la Id para borrar
	http.HandleFunc("/cm/v1/ejecuciones/", actions.HandlerEjecucion)
	//Endpoint para la tabla de condiciones de entrada  (ponemos la última / ya que recibiremos el ID)
	http.HandleFunc("/cm/v1/ejecuciones/condicionin/", actions.HandlerCondicionin)
	//Endpoint para la tabla de condiciones de salida  (ponemos la última / ya que recibiremos el ID)
	http.HandleFunc("/cm/v1/ejecuciones/condicionout/", actions.HandlerCondicionout)
	//Endpoint para la tabla de planificacion
	http.HandleFunc("/cm/v1/planificacion", actions.HandlerPlanificacion)
	//Endpoint para la tabla de planificacion con / para el delete ya que recibira el ID en la ULR de esa manera
	http.HandleFunc("/cm/v1/planificacion/", actions.HandlerPlanificacion)
	//Endpoint para el el resto de metodos para las condiciones de entrada de planificacion (los que no tienen por url /(dato))
	http.HandleFunc("/cm/v1/planificacion/condicionin", actions.HandlerPlanifCondicionin)
	//Endpoint para el get para las condiciones de entrada de planificacion
	http.HandleFunc("/cm/v1/planificacion/condicionin/", actions.HandlerPlanifCondicionin)
	//Endpoint para recuperar los calendarios
	http.HandleFunc("/cm/v1/calendar", actions.HandlerCalendar)
	//Endpoint para el el resto de metodos para las condiciones de salida de planificacion (los que no tienen por url /(dato))
	http.HandleFunc("/cm/v1/planificacion/condicionout", actions.HandlerPlanifCondicionout)
	//Endpoint para el get para las condiciones de salida de planificacion
	http.HandleFunc("/cm/v1/planificacion/condicionout/", actions.HandlerPlanifCondicionout)
	//Endpoint para recuperar el log de la ejecucion tiene / ya que ira el nombre
	http.HandleFunc("/cm/v1/ejecuciones/logs/", actions.HandlerLog)
}
