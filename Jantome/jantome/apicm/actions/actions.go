package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/jantome/apicm/environment"

	"github.com/jantome/apicm/structs"
	"github.com/onlinearq/online"

	"github.com/jantome/apicm/db2"
)

//Definicion de variables que usara todo el programa
var (
	jsonerror       structs.Jsonerror
	usermessage     = ""
	internalmessage = ""
	nombreservicio  = "apicm"
	iderror         = ""
)

//HandlerEjecucion tabla ejecucion
func HandlerEjecucion(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para validar/usuario y pass y generar token en caso de que sea correcto
		getEjecucion(response, request)
	case "DELETE":
		//para la eliminacion de la tabla de ejecucion
		deleteEjecucion(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options1(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerEejecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getEjecucion, sacara el listado de la tabla de ejecucion
func getEjecucion(response http.ResponseWriter, request *http.Request) {
	//definicion de variables de la funcion necesarias
	sql := ""
	//control de paginación
	var ofsset int64
	var page int64
	var limit int64
	//sw para saber si metemos pie
	sipie := false
	//obtenemos las variables de la url
	pageurl, ok := request.URL.Query()["page"]
	limiturl, ok2 := request.URL.Query()["limit"]
	//comprobamos que extrae datos de la variable page
	if ok && len(pageurl[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		pageurl2 := pageurl[0]
		//convertimos a int
		page, _ = strconv.ParseInt(pageurl2, 10, 64)
	}
	//comprobamos que extrae datos de la variable limit
	if ok2 && len(limiturl[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		limiturl2 := limiturl[0]
		//convertimos a int
		limit, _ = strconv.ParseInt(limiturl2, 10, 64)
	}
	//si tenemos dato en page, restamos una a page y lo multiplicamos por limit para sacar el ofsset
	//esto es xk llegara siempre la pagina siguiente, es decir la 0 es la pagina 1, la 1 la 2 etc..
	if page > 0 {
		ofsset = page - 1
		ofsset = ofsset * limit
	}
	//montamos las querys para sacar los datos de lo que esta en ejecucion teniendo en cuenta los parametros
	//solo informado limit
	if limit > 0 && page == 0 {
		sql = fmt.Sprintf("SELECT nombre, fechaeje, estado FROM ejecucion WHERE numsec = 1 LIMIT %d", limit)
		sipie = true
	} else {
		//informado limit y pagina
		if limit > 0 && page > 0 {
			sql = fmt.Sprintf("SELECT nombre, fechaeje, estado FROM ejecucion WHERE numsec = 1 LIMIT %d OFFSET %d", limit, ofsset)
			sipie = true
		} else {
			//no esta informado nada
			sql = "SELECT nombre, fechaeje, estado FROM ejecucion WHERE numsec = 1"
		}
	}
	//Ejecutamos la query
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		//obtenemos id de error
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//jsonerror.InternalMessage = fmt.Sprintf("Error select ejecucion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror) //tiene que ser err2 si no machaca la info de err
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			//Mostramos error por pantalla y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default getEjecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := "Error select fechaeje getEjecucion id: " + iderror
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Creamos variable para aplantillar lectura
	var tabejecucion structs.Tabejecucion
	//Creamos la variable para devolver la salida con los link relacionados
	var ejecucionjson structs.Ejecucionjson
	//variable para acumular todas las lecutras
	acuejecucionjson := []structs.Ejecucionjson{}
	//Creamos sw para ver si devolvemos datos o no
	sidatos := false
	//Creamos bucle con la lectura
	for result.Next() {
		//activamos sw sidatos, para luego generar el json
		sidatos = true
		//hacemos un scan(aplantillar) por cada lectura
		err := result.Scan(&tabejecucion.Nombre, &tabejecucion.Fechaeje, &tabejecucion.Estado)
		//Controlar el error y devolver un 500
		if err != nil {
			//obtenemos id de error
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error default getEjecucion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			//Aunque saquemos mensaje de error, grabamos
			online.EjecutaError(nombreservicio, *structs.Entorno, "Error scan tabejecucion getejecucion", err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//aplantillamos la lectura en el formato de json que vamos a mostrar
		ejecucionjson.Name = tabejecucion.Nombre
		ejecucionjson.Fechaeje = tabejecucion.Fechaeje
		ejecucionjson.Estado = tabejecucion.Estado
		//Montamos las Url's que se pueden usar
		ejecucionjson.Links.Href = make(map[string]string)
		ejecucionjson.Links.Href["condicionin"] = "/cm/v1/ejecuciones/condicionin/" + tabejecucion.Nombre + "?fechaeje=" + tabejecucion.Fechaeje
		ejecucionjson.Links.Href["condicionout"] = "/cm/v1/ejecuciones/condicionout/" + tabejecucion.Nombre + "?fechaeje=" + tabejecucion.Fechaeje
		//añadimos a la variable de acumulados para luego poder montar correctamente el json
		acuejecucionjson = append(acuejecucionjson, ejecucionjson)
	}
	//Si no tenemos datos, sacamos error 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
	//Si si tenemos datos
	if sidatos {
		//Comprobamos si tenemos que sacar pie o no. (Esto es dependiendo de si tenemos paginación)
		if sipie {
			var pagmax float64
			//montamos url prev y next
			var pieejecucion structs.Pieejecucion
			pieejecucion.Pagdet.Links.Href = make(map[string]string)
			//Si tenemos page con valor que 1 el prev siempre sera el page menos 1
			if page > 1 {
				pieejecucion.Pagdet.Links.Href["prev"] = fmt.Sprintf("/cm/v1/ejecuciones?limit=%d&page=%d", limit, page-1)
				pieejecucion.Pagdet.Pagprev = page - 1
			}
			//para armar el next cuando corresponde necesitamos saber el número de paginas que tendra
			//hacemos un select count
			sql = "SELECT COUNT(*) FROM ejecucion WHERE numsec = 1"
			result, err = db2.EjecutaQuery(sql)
			//solo tendra un registro por lo que no hace falta que recorramos toda la tabla
			var ejecucioncount structs.Ejecucioncount
			result.Next()
			err = result.Scan(&ejecucioncount.Count)
			//Controlar el error y devolver un 500
			if err != nil {
				//obtenemos id de errro
				iderror = online.GeneraIDError()
				//json de error
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//	jsonerror.InternalMessage = fmt.Sprintf("Error select max. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si falla la generacion damos error grave
				if err2 != nil {
					http.Error(response, "Fatal Mistake", http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error default getEjecucion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				//Aunque saquemos mensaje de error, grabamos
				mensaje := fmt.Sprintf("Select Count ejecucion getejecucion id: %s ", iderror)
				online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//cerramos result para que no se quede la conexion abierta
			defer result.Close()
			//convertimos el resultado del count a float64
			countfloat := float64(ejecucioncount.Count)
			//convertimos limit en float64 para poder dividir
			limitfloat := float64(limit)
			//Dividmos
			divison := countfloat / limitfloat
			//redondeamos por funcion ya que quiero que siempre redonde hacia arriba
			t := math.Trunc(divison)
			//Si el resultado es 0 significa que es un número entero por lo la paginacion es exacta
			if math.Abs(divison-t) != 0 {
				if math.Abs(divison-t) >= 0.5 || math.Abs(divison-t) < 0.5 {
					pagmax = t + math.Copysign(1, divison)
				}
			} else {
				pagmax = divison
			}
			//convertimos pagmax a int64 para poder trabajar con el
			pagmax64 := int64(pagmax)
			//comparamos pero si page no esta informada le ponemos dos ya que empezmos en la 1
			if page == 0 {
				if 2 <= pagmax64 {
					pieejecucion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/ejecuciones?limit=%d&page=%d", limit, 2)
					pieejecucion.Pagdet.Pagnext = 2
				}
			} else {
				if page < pagmax64 {
					pieejecucion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/ejecuciones?limit=%d&page=%d", limit, page+1)
					pieejecucion.Pagdet.Pagnext = page + 1
				}
			}
			pieejecucion.Pagdet.Pagmax = pagmax64
			//Rellenamos el contenido acumulado de la select
			pieejecucion.Content = acuejecucionjson
			//creamos json
			JsResponser, err := json.Marshal(pieejecucion)
			//Controlar el error y devolver un 500
			if err != nil {
				//obtenemos id de error
				iderror = online.GeneraIDError()
				//Informamos el json
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//	jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si vuelve a fallar la generacion, ya grabamos en log
				if err2 != nil {
					mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
					http.Error(response, mensaje, http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error default getEjecucion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				//Aunque saquemos mensaje de error, grabamos
				mensaje := "Error genera json pieejecucion getejecucion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//creamos cabecera de respuesta
			response.Header().Set("Content-Type", "application/json")
			//devolvemos la respuesta
			response.Write(JsResponser)
		} else {
			JsResponser, err := json.Marshal(acuejecucionjson)
			if err != nil {
				//obtenemos id de error
				iderror = online.GeneraIDError()
				//Informamos el json
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//	jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si vuelve a fallar la generacion, ya grabamos en log
				if err2 != nil {
					mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
					http.Error(response, mensaje, http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error default getEjecucion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				//Aunque saquemos mensaje de error, grabamos
				mensaje := "Error generacion acueejecucionjson getejecucion" + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//creamos cabecera de respuesta
			response.Header().Set("Content-Type", "application/json")
			//grabamos cuerpo
			response.Write(JsResponser)
		}
	}
}

//deleteEjecucion, para eliminar si esta Ko por ejemplo, pero tiene que pasar por hold primero
func deleteEjecucion(response http.ResponseWriter, request *http.Request) {
	//variables inicializadas
	fechaeje2 := ""
	//Optenemos la Id de la urle
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//Recuperamos el parametro de fecha que llegara en la url
	fechaeje, ok := request.URL.Query()["fechaeje"]
	//comprobamos que extrae datos de la variable page
	if ok && len(fechaeje[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		fechaeje2 = fechaeje[0]
	} else {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		jsonerror.InternalMessage = fmt.Sprintf("Invalid Parameter url: fechaeje")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Fatal Mistake", http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default deleteEjecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := "Invalid Parameter url: fechaeje id: " + iderror
		online.EjecutaInfo(nombreservicio, *structs.Entorno, mensaje, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//comprobamos el estado ya que tiene que estar holdeado para poder elimintar
	sql := fmt.Sprintf("SELECT estado FROM ejecucion WHERE nombre = '%s' AND fechaeje ='%s' AND NUMSEC = 1", id, fechaeje2)
	//ejecutamos Query
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s ", iderror)
		//jsonerror.InternalMessage = fmt.Sprintf("Error select estado. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			http.Error(response, "Fatal Mistake", http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error deleteEjecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := "Error en la select ejecucion deleteejecuion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//solo tendra un registro pero aun asi tenemos que montar bucle, para que no falle
	//creamos variable donde aplantillar
	var estadoejecucion structs.Estadoejecucion
	//variable para indicar que no encontro datos
	datos := false
	for result.Next() {
		datos = true
		err = result.Scan(&estadoejecucion.Estado)
		//cerramos result
		defer result.Close()
		//controlamos el error
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
			//jsonerror.InternalMessage = fmt.Sprintf("Error scan estado. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				http.Error(response, "Fatal Mistake", http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error deleteEjecucion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			//grabamos en el log
			mensaje := "Error en la select ejecucion scan result select estado de ejecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Controlamos el estado y si es holdeado realizamos el delete
		if estadoejecucion.Estado == "ho" {
			sql := fmt.Sprintf("DELETE FROM ejecucion WHERE nombre = '%s' AND fechaeje ='%s'", id, fechaeje2)
			result, err := db2.EjecutaQuery(sql)
			defer result.Close()
			//Controlar el error para devolver un 500
			if err != nil {
				iderror = online.GeneraIDError()
				//json de error
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//jsonerror.InternalMessage = fmt.Sprintf("Error delete ejecucion. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si falla la generacion damos error grave
				if err2 != nil {
					http.Error(response, "Fatal Mistake", http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error deleteEjecucion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				//grabamos en el log
				mensaje := "Error en delete ejecucion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
		} else {
			//movemos 404 de no encontrado
			response.WriteHeader(http.StatusNotFound)
			return
		}
	}
	if !datos {
		//movemos 404 de no encontrado
		response.WriteHeader(http.StatusNotFound)
		return
	}
}

//options1 OPTIONS, GET, DELETE
func options1(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, DELETE")
	return
}

//HandlerCondicionin condiciones de entrada de la tabla ejecucion
func HandlerCondicionin(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para recuperar la info
		getCondicionin(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getCondicionnin condiciones de entrada en json de la tabla ejecucion
func getCondicionin(response http.ResponseWriter, request *http.Request) {
	//Creamos la variable necesaria
	fechaeje2 := ""
	//recuperamos el id de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//Recuperamos el parametro de fecha que llegara en la url
	fechaeje, ok := request.URL.Query()["fechaeje"]
	//comprobamos que extrae datos de la variable page
	if ok && len(fechaeje[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		fechaeje2 = fechaeje[0]
	} else {
		//obtenemos id de error
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Invalid Parameter url: fechaeje")
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			//Mostramos error por pantalla y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  getCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := jsonerror.InternalMessage + "getCondicionin"
		online.EjecutaInfo(nombreservicio, *structs.Entorno, mensaje, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Query para recuperar todas las condiciones de entrada
	sql := fmt.Sprintf("SELECT condicionin FROM ejecucion WHERE nombre ='%s' AND condicionin > '' and FECHAEJE ='%s' AND estado ='' ", id, fechaeje2)
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		//jsonerror.InternalMessage = fmt.Sprintf("Error select condicionin. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			//Mostramos error por pantalla y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  getCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := "Error select condicionin getCondicionin id: " + iderror
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos bucle para sacar las condiciones
	// Variable de lectura
	var condicionin structs.Condicionin
	//Variable para la acumulacion del json de salida
	jsoncondicionin := []structs.Condicionin{}
	//Sw para saber si sacamos datos o no
	sidatos := false
	for result.Next() {
		sidatos = true
		//aplantillamos en el struct de salida
		err = result.Scan(&condicionin.Condicionin)
		//Controlar el error y devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//jsonerror.InternalMessage = fmt.Sprintf("Error scan condicion. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				//Mostramos error por pantalla y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error  getCondicionin id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			mensaje := "Error scan condicionin getCondicionin id: " + iderror
			//Aunque saquemos mensaje de error, grabamos
			online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Acumulamos en el json
		jsoncondicionin = append(jsoncondicionin, condicionin)
	}
	//al salir del for es cuando creamos el json siempre y cuando tengamos algo en la lectura
	if sidatos {
		JsResponser, err := json.Marshal(jsoncondicionin)
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//		jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				//Mostramos error por pantalla y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error  getCondicionin id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			mensaje := "Error scan generajson getCondicionin id: " + iderror
			//Aunque saquemos mensaje de error, grabamos
			online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.Write(JsResponser)
	}
	//Si no tenemos datos sacamos 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
}

//HandlerCondicionout condiciones de salida
func HandlerCondicionout(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para recuperar la info
		getCondicionout(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getCondicionnout condiciones de salida en json
func getCondicionout(response http.ResponseWriter, request *http.Request) {
	//Creamos la variable necesaria
	fechaeje2 := ""
	//recuperamos el id de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//Recuperamos el parametro de fecha que llegara en la url
	fechaeje, ok := request.URL.Query()["fechaeje"]
	//comprobamos que extrae datos de la variable page
	if ok && len(fechaeje[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		fechaeje2 = fechaeje[0]
	} else {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Invalid Parameter url: fechaeje")
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			//Mostramos error por pantalla y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  getCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := jsonerror.InternalMessage + "getCondicionout"
		online.EjecutaInfo(nombreservicio, *structs.Entorno, mensaje, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Query para recuperar todas las condiciones de entrada
	sql := fmt.Sprintf("SELECT condicionout FROM ejecucion WHERE nombre ='%s' AND condicionout > '' and FECHAEJE ='%s' and estado =''", id, fechaeje2)
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		//jsonerror.InternalMessage = fmt.Sprintf("Error select condicionout. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			//Mostramos error por pantalla y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  getCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := "Error select condicionout from ejecucion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos bucle para sacar las condiciones
	//Variable para la lectura
	var condicionout structs.Condicionout
	//Variable para la acumulacion
	jsoncondicionout := []structs.Condicionout{}
	//sw para saber si tenemos datos o no
	sidatos := false
	for result.Next() {
		sidatos = true
		//aplantillamos en el struct de salida
		err = result.Scan(&condicionout.Condicionout)
		//Controlar el error y devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//jsonerror.InternalMessage = fmt.Sprintf("Error scan condicionout. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				//Mostramos error por pantalla y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error  getCondicionout id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			mensaje := "Error scan condicionout from ejecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//por cada lectura acumulamos
		jsoncondicionout = append(jsoncondicionout, condicionout)
	}
	//generamos json en caso de tener datos
	if sidatos {
		//creamos json con el valor acumulado
		JsResponser, err := json.Marshal(jsoncondicionout)
		//error en la generacion del json
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				//Mostramos error por pantalla y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error  getCondicionout id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			mensaje := "Error genera json from ejecucion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//devolvemos la respuesta
		response.Write(JsResponser)
	}
	//Si no tenemos datos sacamos 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
}

//HandlerPlanificacion para la tabla de planificacion
func HandlerPlanificacion(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para recuperar la info
		getPlanificacion(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options2(response, request)
	//Update
	case "PUT":
		putPlanificacion(response, request)
	//Insert
	case "POST":
		postPlanificacion(response, request)
	//Delete
	case "DELETE":
		deletePlanificacion(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getPlanificacion para recupera los datos de la tabla de planificacion (solo recuperar los que tengan el ejecucion informado
//que sera el primero de planificacion, del que dependera condiciones de entrada salida etc..)
func getPlanificacion(response http.ResponseWriter, request *http.Request) {
	//sw para saber si metemos pie
	sipie := false
	sql := ""
	//control de paginación
	var ofsset int64
	var page int64
	var limit int64
	//obtenemos las variables de la url
	pageurl, ok := request.URL.Query()["page"]
	limiturl, ok2 := request.URL.Query()["limit"]
	//comprobamos que extrae datos de la variable page
	if ok && len(pageurl[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		pageurl2 := pageurl[0]
		//convertimos a int
		page, _ = strconv.ParseInt(pageurl2, 10, 64)
	}
	//comprobamos que extrae datos de la variable limit
	if ok2 && len(limiturl[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		limiturl2 := limiturl[0]
		//convertimos a int
		limit, _ = strconv.ParseInt(limiturl2, 10, 64)
	}
	//si tenemos dato en page, restamos una a page y lo multiplicamos por limit
	//esto es xk llegara siempre la pagina siguiente, es decir la 0 es la pagina 1, la 1 la 2 etc..
	if page > 0 {
		ofsset = page - 1
		ofsset = ofsset * limit
	}
	//montamos query segun los parametros recuperados por URL
	if limit > 0 && page == 0 {
		sql = fmt.Sprintf("SELECT nombre, calendario, user_alta, timalta, user_modif, timesmod FROM planificacion WHERE ejecucion = 'n' LIMIT  %d", limit)
		sipie = true
	} else {
		if limit > 0 && page > 0 {
			sql = fmt.Sprintf("SELECT nombre, calendario, user_alta, timalta, user_modif, timesmod FROM planificacion WHERE ejecucion = 'n' LIMIT %d OFFSET %d", limit, ofsset)
			sipie = true
		} else {
			sql = "SELECT nombre, calendario, user_alta, timalta, user_modif, timesmod FROM planificacion WHERE ejecucion = 'n'"
		}
	}
	//Ejecutamos la query
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//jsonerror.InternalMessage = fmt.Sprintf("Error select ejecucion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default getPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		mensaje := "Error Select FROM planificacion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, mensaje, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Inicializacion de datos
	sidatos := false
	//Creamos var dond estara la lectura
	var tabplanificacion structs.Tabplanificacion
	//variable de acumulacion
	acuplanificacion := []structs.Tabplanificacion{}
	for result.Next() {
		sidatos = true
		//hacemos un scan(aplantillar) por cada lectura
		err := result.Scan(&tabplanificacion.Nombre, &tabplanificacion.Calendario, &tabplanificacion.Useralta, &tabplanificacion.Timalta, &tabplanificacion.Usermod, &tabplanificacion.Timesmod)
		//Controlar el error y devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
			//	jsonerror.InternalMessage = fmt.Sprintf("Error bucle planificacion. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error default getPlanificacion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error Scan Select FROM planificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Acumulamos
		acuplanificacion = append(acuplanificacion, tabplanificacion)
	}
	//Si no tenemos datos damos error con 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
	//Si tenemos datos
	if sidatos {
		//Comprobamos si tenemos que grabar o no pie
		if sipie {
			var pagmax float64
			//montamos url prev y next
			var pieplanificacion structs.Pieplanificacion
			pieplanificacion.Pagdet.Links.Href = make(map[string]string)
			//Si page es > 1 siempre restamos
			if page > 1 {
				pieplanificacion.Pagdet.Links.Href["prev"] = fmt.Sprintf("/cm/v1/planificacion?limit=%d&page=%d", limit, page-1)
				pieplanificacion.Pagdet.Pagprev = page - 1
			}
			//para armar el next cuando corresponde necesitamos saber el número de paginas que tendra
			//hacemos un select count
			sql = "SELECT COUNT(*) FROM planificacion WHERE ejecucion = 'n'"
			result, err = db2.EjecutaQuery(sql)
			//solo tendra un registro por lo que no hace falta que recorramos toda la tabla
			var planificacioncount structs.Planificacioncount
			result.Next()
			err = result.Scan(&planificacioncount.Count)
			//para query's que no son lista tenemos que cerrar result, si no se queda la conexion abierta al
			//ser procesos que no se acaban nunca.
			defer result.Close()
			//Controlar el error y devolver un 500
			if err != nil {
				iderror = online.GeneraIDError()
				//json de error
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//	jsonerror.InternalMessage = fmt.Sprintf("Error select max. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si falla la generacion damos error grave
				if err2 != nil {
					mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
					//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
					http.Error(response, mensaje, http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error default getPlanificacion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				descripcion := "Error select count from planificacoin id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//convertimos el resultado del count a float64
			countfloat := float64(planificacioncount.Count)
			//convertimos limit en float64 para poder dividir
			limitfloat := float64(limit)
			//Dividmos
			divison := countfloat / limitfloat
			//redondeamos por funcion ya que quiero que siempre redonde hacia arriba
			t := math.Trunc(divison)
			//Si el resultado es 0 significa que es un número entero por lo la paginacion es exacta
			if math.Abs(divison-t) != 0 {
				if math.Abs(divison-t) >= 0.5 || math.Abs(divison-t) < 0.5 {
					pagmax = t + math.Copysign(1, divison)
				}
			} else {
				pagmax = divison
			}
			//convertimos pagmax a int64 para poder trabajar con el
			pagmax64 := int64(pagmax)
			//comparamos pero si page no esta informada le ponemos dos ya que empezmos en la 1
			if page == 0 {
				if 2 <= pagmax64 {
					pieplanificacion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/planificacion?limit=%d&page=%d", limit, 2)
					pieplanificacion.Pagdet.Pagnext = 2
				}
			} else {
				if page < pagmax64 {
					pieplanificacion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/planificacion?limit=%d&page=%d", limit, page+1)
					pieplanificacion.Pagdet.Pagnext = page + 1
				}
			}
			pieplanificacion.Pagdet.Pagmax = pagmax64
			//informamos el content con lo acumulado en la lectura
			pieplanificacion.Content = acuplanificacion
			//creamos json
			JsResponser, err := json.Marshal(pieplanificacion)
			//Controlar el error y devolver un 500
			if err != nil {
				iderror = online.GeneraIDError()
				//Informamos el json
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//	jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si vuelve a fallar la generacion, ya grabamos en log
				if err2 != nil {
					mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
					//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
					http.Error(response, mensaje, http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error default getPlanificacion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				descripcion := "Error en la generacion del json id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//creamos cabecera de respuesta
			response.Header().Set("Content-Type", "application/json")
			//devolvemos la respuesta
			response.Write(JsResponser)
		} else {
			JsResponser, err := json.Marshal(acuplanificacion)
			if err != nil {
				iderror = online.GeneraIDError()
				//Informamos el json
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
				//	jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si vuelve a fallar la generacion, ya grabamos en log
				if err2 != nil {
					mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
					//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
					http.Error(response, mensaje, http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error default getPlanificacion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				descripcion := "Error en la generacion del json acuplanificacion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//creamos cabecera de respuesta
			response.Header().Set("Content-Type", "application/json")
			//grabamos cuerpo
			response.Write(JsResponser)
		}
	}
}

//options2 para los cors de esta api
func options2(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	//Get lista
	//Put update
	//Post create
	//Delete delete
	response.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	return
}

//options4 para los cors de esta api
func options4(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	//Get lista
	//Post create
	//Delete delete
	response.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	return
}

//putPlanificacion actualización de la tabla planificacion
func putPlanificacion(response http.ResponseWriter, request *http.Request) {
	//De entrada tendra un json en el que tendra como obligatorio el nombre (ya que sera con el que hacemos el update)
	//se puede actualizar unica y exclusivamente el calendario y los datos de auditoria (usermodif)
	//lo primero que hacemos es recuperar el json del body
	cuerpo := json.NewDecoder(request.Body)
	//ceramos la respuesta de body con defer para que se ejecute al final
	defer request.Body.Close()
	//creamos variable donde aplantillaremos
	var putplanificacion structs.Putplanificacion
	//decodificamos en el struc correspondiente
	err := cuerpo.Decode(&putplanificacion)
	//controlamos el error
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error decode putplanificacion. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error en el decode del json id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//validamos campos obligatorios
	if putplanificacion.Name == "" {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Mandatory data nombre")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Dato obligatorio nombre putplanificacion"
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//calendario obligatorio
	if putplanificacion.Calendar == "" {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Mandatory data Calendario")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Dato obligatorio Calendario putplanificacion"
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//usuario modif obligatorio
	if putplanificacion.Usermod == "" {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Mandatory data Usermodif")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Dato obligatorio Usermodif putplanificacion"
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//hacemos una query para verificar si existe
	//comprobamos que exista haciendo un count
	sql := fmt.Sprintf("SELECT COUNT(*) FROM planificacion WHERE nombre = '%s'", putplanificacion.Name)
	result2, err := db2.EjecutaQuery(sql)
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		//	jsonerror.InternalMessage = fmt.Sprintf("Error select count planificacion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error en count from planificacion putplanificacion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos variable donde leermos
	var planificacioncount structs.Planificacioncount
	//solo tendremos una ocurrencia
	result2.Next()
	defer result2.Close()
	//aplantillamos
	err = result2.Scan(&planificacioncount.Count)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		//	jsonerror.InternalMessage = fmt.Sprintf("Error scan count. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error en scan del count putPlanificacion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//si tiene contenido realiza el update
	if planificacioncount.Count > 0 {
		//Una vez validado los datos obligatorios montamos la query
		sql = fmt.Sprintf("UPDATE planificacion SET calendario = '%s', user_modif = '%s' WHERE nombre = '%s'", putplanificacion.Calendar, putplanificacion.Usermod, putplanificacion.Name)
		//ejecutamos la query
		result, err := db2.EjecutaQuery(sql)
		//este defer de resultado, es para los put, ya que si no se queda la conexión abierta con mysql
		//Tambien pasa en los post
		defer result.Close()
		//Controlar el error para devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error update planificacion. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				iderror = online.GeneraIDError()
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error default putPlanificacion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error en el update planificacion putplanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
	} else {
		//en caso de no existir lo que hacemos es mostrar un 404
		response.WriteHeader(http.StatusNotFound)
	}
}

//postPlanificacion, insert en la tabla de planificacion con ejecucion = 'n' (sin condiciones)
func postPlanificacion(response http.ResponseWriter, request *http.Request) {
	//De entrada tendra un json en el que tendra como obligatorio el nombre, el calendario y el usuario de alta
	//lo primero que hacemos es recuperar el json del body
	cuerpo := json.NewDecoder(request.Body)
	//ceramos la respuesta de body con defer para que se ejecute al final
	defer request.Body.Close()
	//creamos variable donde aplantillaremos
	var postplanificacion structs.Postplanificacion
	//decodificamos en el struc correspondiente
	err := cuerpo.Decode(&postplanificacion)
	//controlamos el error
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s ", iderror)
		//jsonerror.InternalMessage = fmt.Sprintf("Error decode putplanificacion. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error postPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error decode cuerpo postplanificacion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Comprobación de los datos obligatorios
	//Nombre
	if postplanificacion.Name == "" {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Mandatory data nombre")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error postPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Dato obligatorio nombre postplanificacion "
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Calendario
	if postplanificacion.Calendar == "" {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Mandatory data Calendar")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error postPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Dato obligatorio Calendar postplanificacion "
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Usuario de alta
	if postplanificacion.Useralt == "" {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Mandatory data Useralt")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error postPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Dato obligatorio Useralt postplanificacion "
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 400 al error
		response.WriteHeader(http.StatusBadRequest)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//generamos el sql
	sql := fmt.Sprintf("INSERT INTO planificacion VALUES(NULL, '%s', 'n', '','','%s','%s',CURRENT_TIMESTAMP, '%s',CURRENT_TIMESTAMP)", postplanificacion.Name, postplanificacion.Calendar, postplanificacion.Useralt, postplanificacion.Useralt)
	//ejecutamos query
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		//	jsonerror.InternalMessage = fmt.Sprintf("Error update planificacion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error postPlanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error en el insert planificacion id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//este defer de resultado, es igual que en los put para que cierre la conexión con mysql
	defer result.Close()
	//devolvemos 201 creado
	response.WriteHeader(http.StatusCreated)
}

//deletePlanificacion eliminar de la tabla. (Elimina todo condiciones incluidas)
func deletePlanificacion(response http.ResponseWriter, request *http.Request) {
	//recuperamos el id de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//Comprobamos que el id esta informado con algo distinto de planificacion, que eso indicara que viene algo
	//informado
	if id != "planificacion" {
		//comprobamos que exista haciendo un count
		sql := fmt.Sprintf("SELECT COUNT(*) FROM planificacion WHERE nombre = '%s'", id)
		result2, err := db2.EjecutaQuery(sql)
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//		jsonerror.InternalMessage = fmt.Sprintf("Error select count planificacion. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error deletePlanificacion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error select count deleteplanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//creamos variable donde leermos
		var planificacioncount structs.Planificacioncount
		//solo tendremos una ocurrencia
		result2.Next()
		defer result2.Close()
		//aplantillamos
		err = result2.Scan(&planificacioncount.Count)
		//Controlar el error para devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//		jsonerror.InternalMessage = fmt.Sprintf("Error scan count. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error deleteplanificacion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error scan count deleteplanificacion id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Solo hacemos delete si el count es > 0
		if planificacioncount.Count > 0 {
			//montamos la query para el delete
			sql = fmt.Sprintf("DELETE FROM planificacion WHERE nombre = '%s'", id)
			result, err := db2.EjecutaQuery(sql)
			//este defer de resultado, es igual que en los put para que cierre la conexión con mysql
			defer result.Close()
			//Controlar el error para devolver un 500
			if err != nil {
				iderror = online.GeneraIDError()
				//json de error
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
				//		jsonerror.InternalMessage = fmt.Sprintf("Error delete planificacion. Descripción: %s", err.Error())
				JsResponser, err2 := json.Marshal(jsonerror)
				//si falla la generacion damos error grave
				if err2 != nil {
					mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
					//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
					http.Error(response, mensaje, http.StatusInternalServerError)
					descripcion := "Error en la generacion del json de error deleteplanificacion id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
					//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
					response.WriteHeader(http.StatusInternalServerError)
					return
				}
				descripcion := "Error en el delete deleteplanificacion id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
		} else {
			//en caso de no existir lo que hacemos es mostrar un 404
			response.WriteHeader(http.StatusNotFound)
		}
	} else {
		//devolvemos 400 en caso de que no este informada el id
		response.WriteHeader(http.StatusBadRequest)
	}
}

//HandlerPlanifCondicionin para las condiciones de entrada en planificacion
func HandlerPlanifCondicionin(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	//para recuperar las condiciones de entrada de la tabla planif
	case "GET":
		getcondicionin2(response, request)
	case "POST":
		//ejecutamos la funcion para recuperar la info
		postCondicionin(response, request)
	case "DELETE":
		delCondicionin(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options4(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerPlanifCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getcondicionin2 para sacar las condiciones de la tabla de planificacion
func getcondicionin2(response http.ResponseWriter, request *http.Request) {
	//recuperamos el id de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//Query para recuperar todas las condiciones de entrada
	sql := fmt.Sprintf("SELECT condicionin FROM planificacion WHERE nombre ='%s' AND condicionin > '' ", id)
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		//	jsonerror.InternalMessage = fmt.Sprintf("Error select condicionin. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default getcondicionin2 id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error select condicionin from planificacion getcondicionin2 id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos bucle para sacar las condiciones
	// Variable de lectura
	var condicionin structs.Condicionin
	//Variable para la acumulacion del json de salida
	jsoncondicionin := []structs.Condicionin{}
	//Sw para saber si sacamos datos o no
	sidatos := false
	for result.Next() {
		sidatos = true
		//aplantillamos en el struct de salida
		err = result.Scan(&condicionin.Condicionin)
		//Controlar el error y devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s ", iderror)
			//	jsonerror.InternalMessage = fmt.Sprintf("Error scan condicion. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
				iderror = online.GeneraIDError()
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error default getcondicionin2 id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error scan select condicion2 getcondicionin2 id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Acumulamos en el json
		jsoncondicionin = append(jsoncondicionin, condicionin)
	}
	//al salir del for es cuando creamos el json siempre y cuando tengamos algo en la lectura
	if sidatos {
		JsResponser, err := json.Marshal(jsoncondicionin)
		if err != nil {
			iderror = online.GeneraIDError()
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err2 != nil {
				//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
				iderror = online.GeneraIDError()
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error default getcondicionin2 id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error en la generacion del json de salida getcondicionin2 id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.Write(JsResponser)
	}
	//Si no tenemos datos sacamos 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
}

//postCondicionin para insertar las condiciones de entrada en la planificacion
func postCondicionin(response http.ResponseWriter, request *http.Request) {
	//De entrada tendra un json en el que tendra como obligatorio el nombre, el calendario y el usuario de alta
	//lo primero que hacemos es recuperar el json del body
	cuerpo := json.NewDecoder(request.Body)
	//ceramos la respuesta de body con defer para que se ejecute al final
	defer request.Body.Close()
	//creamos variable donde aplantillaremos
	var postcondicionin structs.Postcondicionin
	//decodificamos en el struc correspondiente
	err := cuerpo.Decode(&postcondicionin)
	//controlamos el error
	if err != nil {
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error decode postcondicionin. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default postCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error dedoce body postCondicionin id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//comprobamos que los datos de entrada estan informados
	if postcondicionin.Name != "" {
		if postcondicionin.Condicionin != "" {
			if postcondicionin.Useralt != "" {
				//si todos los datos estan informados primero recuperamos el calendario del programa principal
				sql := fmt.Sprintf("SELECT calendario FROM planificacion WHERE nombre = '%s' AND ejecucion = 'n'", postcondicionin.Name)
				result, err := db2.EjecutaQuery(sql)
				//controlamos error
				if err != nil {
					iderror = online.GeneraIDError()
					//json de error
					jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
					//		jsonerror.InternalMessage = fmt.Sprintf("Error select calendario planificacion. Descripción: %s", err.Error())
					JsResponser, err2 := json.Marshal(jsonerror)
					//si falla la generacion damos error grave
					if err2 != nil {
						mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
						//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
						http.Error(response, mensaje, http.StatusInternalServerError)
						descripcion := "Error en la generacion del json de error default postCondicionin id: " + iderror
						online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
						//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
						response.WriteHeader(http.StatusInternalServerError)
						return
					}
					descripcion := "Error en la select calendario from planifcacion postcondicionin: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
					//Creamos cabecera
					response.Header().Set("Content-Type", "application/json")
					//movemos 500 al error
					response.WriteHeader(http.StatusInternalServerError)
					//grabamos el json de error
					response.Write(JsResponser)
					return
				}
				//solo tendra un registro, pero tenemos que montar bucle igual para que no falle
				var calendarplanificacion structs.Calendarplanificacion
				datos := false
				for result.Next() {
					datos = true
					//aplantillamos
					err = result.Scan(&calendarplanificacion.Calendario)
					defer result.Close()
					if err != nil {
						iderror = online.GeneraIDError()
						//json de error
						jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
						//			jsonerror.InternalMessage = fmt.Sprintf("Error scan calendarplanificacion. Descripción: %s", err.Error())
						JsResponser, err2 := json.Marshal(jsonerror)
						//si falla la generacion damos error grave
						if err2 != nil {
							mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
							//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
							http.Error(response, mensaje, http.StatusInternalServerError)
							descripcion := "Error en la generacion del json de error default postCondicionin id: " + iderror
							online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
							//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
							response.WriteHeader(http.StatusInternalServerError)
							return
						}
						descripcion := "Error en la scan select postCondicionin id: " + iderror
						online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
						//Creamos cabecera
						response.Header().Set("Content-Type", "application/json")
						//movemos 500 al error
						response.WriteHeader(http.StatusInternalServerError)
						//grabamos el json de error
						response.Write(JsResponser)
						return
					}
					//cerramos result
					defer result.Close()
					//realizamos insert
					sql = fmt.Sprintf("INSERT INTO planificacion VALUES( NULL, '%s','', '%s','','%s','%s',CURRENT_TIMESTAMP,'%s',CURRENT_TIMESTAMP)", postcondicionin.Name, postcondicionin.Condicionin, calendarplanificacion.Calendario, postcondicionin.Useralt, postcondicionin.Useralt)
					//ejecutamos la query
					result, err = db2.EjecutaQuery(sql)
					if err != nil {
						iderror = online.GeneraIDError()
						//json de error
						jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
						//	jsonerror.InternalMessage = fmt.Sprintf("Error insert condicionin planificacion. Descripción: %s", err.Error())
						JsResponser, err2 := json.Marshal(jsonerror)
						//si falla la generacion damos error grave
						if err2 != nil {
							mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
							//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
							http.Error(response, mensaje, http.StatusInternalServerError)
							descripcion := "Error en la generacion del json de error default postCondicionin id: " + iderror
							online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
							//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
							response.WriteHeader(http.StatusInternalServerError)
							return
						}
						descripcion := "Error insert planificacion postCondicionin id: " + iderror
						online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
						//Creamos cabecera
						response.Header().Set("Content-Type", "application/json")
						//movemos 500 al error
						response.WriteHeader(http.StatusInternalServerError)
						//grabamos el json de error
						response.Write(JsResponser)
						return
					}
					//mostramos el mensaje de creado
					response.WriteHeader(http.StatusCreated)
					defer result.Close()
				}
				if !datos {
					descripcion := "Dato no encontrado postCondicionin "
					online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
					//movemos 204 al error
					response.WriteHeader(http.StatusNotFound)
					return
				}
			} else {
				descripcion := "Dato obligatorio Useralt postCondicionin "
				online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
				response.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			descripcion := "Dato obligatorio Condicionin postCondicionin "
			online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		descripcion := "Dato obligatorio Name postCondicionin "
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, nil)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

//delCondicionin para la elmiminación de la condicion de entrada de la tabla planificación
func delCondicionin(response http.ResponseWriter, request *http.Request) {
	condicion2 := ""
	//recuperamos el id(name) de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//obtenemos las variables de la url (condición)
	condicion, ok := request.URL.Query()["condicion"]
	//comprobamos que extrae datos de la variable page
	if ok && len(condicion[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		condicion2 = condicion[0]
	}
	//con el name recuperado y la condicion a borrar, montamos la query, lo primero que hacemos es comprobar que
	//existe
	sql := fmt.Sprintf("SELECT COUNT(*) FROM planificacion WHERE nombre = '%s' AND condicionin = '%s'", id, condicion2)
	result, err := db2.EjecutaQuery(sql)
	if err != nil {
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//jsonerror.InternalMessage = fmt.Sprintf("Error select count planificacion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error delCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error eselect count delCondicionin id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos variable donde leermos
	var planificacioncount structs.Planificacioncount
	//solo tendremos una ocurrencia
	result.Next()
	defer result.Close()
	//aplantillamos
	err = result.Scan(&planificacioncount.Count)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error scan count. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error delCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error scan delCondicionin id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Solo hacemos delete si el count es > 0
	if planificacioncount.Count > 0 {
		sql = fmt.Sprintf("DELETE FROM planificacion WHERE nombre = '%s' AND condicionin = '%s'", id, condicion2)
		result, err = db2.EjecutaQuery(sql)
		//este defer de resultado, es igual que en los put para que cierre la conexión con mysql
		defer result.Close()
		//Controlar el error para devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//		jsonerror.InternalMessage = fmt.Sprintf("Error delete planificacion condicionin. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error delCondicionin id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error delete delCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
	} else {
		//en caso de no existir lo que hacemos es mostrar un 404
		response.WriteHeader(http.StatusNotFound)
	}
}

//HandlerCalendar para recuperar el nombre de los calendarios que existen
func HandlerCalendar(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para recuperar la info
		getCalendar(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerPlanifCondicionin id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getCalendar, recuperacion de los calendarios
func getCalendar(response http.ResponseWriter, request *http.Request) {
	//recueramos año
	date := time.Now()
	anno := fmt.Sprintf("%d", date.Year())
	//montamos la query
	sql := fmt.Sprintf("SELECT nombre FROM calendarios WHERE year = %s", anno)
	result, err := db2.EjecutaQuery(sql)
	if err != nil {
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//jsonerror.InternalMessage = fmt.Sprintf("Error select ejecucion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default getCalendar id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error select getCalendar id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Inicializacion de datos
	sidatos := false
	//Creamos var dond estara la lectura
	var calendar structs.Calendar
	//variable de acumulacion
	acucalendar := []structs.Calendar{}
	for result.Next() {
		sidatos = true
		//hacemos un scan(aplantillar) por cada lectura
		err := result.Scan(&calendar.Name)
		//Controlar el error y devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//		jsonerror.InternalMessage = fmt.Sprintf("Error bucle calendar. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error getCalendar id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error scan getCalendar id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Acumulamos
		acucalendar = append(acucalendar, calendar)
	}
	//Si no tenemos datos damos error con 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
	if sidatos {
		JsResponser, err := json.Marshal(acucalendar)
		if err != nil {
			iderror = online.GeneraIDError()
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
			//	jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error getCalendar id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error en la generacion del json  getCalendar id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//grabamos cuerpo
		response.Write(JsResponser)
	}
}

//HandlerPlanifCondicionout para las condiciones de entrada en planificacion
func HandlerPlanifCondicionout(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	//para recuperar las condiciones de entrada de la tabla planif
	case "GET":
		getcondicionout2(response, request)
	case "POST":
		//ejecutamos la funcion para recuperar la info
		postCondicionout(response, request)
	case "DELETE":
		delCondicionout(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options4(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error default HandlerPlanifCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getcondicionout2 para sacar las condiciones de la tabla de planificacion
func getcondicionout2(response http.ResponseWriter, request *http.Request) {
	//recuperamos el id de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//Query para recuperar todas las condiciones de entrada
	sql := fmt.Sprintf("SELECT condicionout FROM planificacion WHERE nombre ='%s' AND condicionout > '' ", id)
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error select condicionin. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error getcondicionout2 id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error select from planificacion getcondicionout2 id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos bucle para sacar las condiciones
	// Variable de lectura
	var condicionout structs.Condicionout
	//Variable para la acumulacion del json de salida
	jsoncondicionout := []structs.Condicionout{}
	//Sw para saber si sacamos datos o no
	sidatos := false
	for result.Next() {
		sidatos = true
		//aplantillamos en el struct de salida
		err = result.Scan(&condicionout.Condicionout)
		//Controlar el error y devolver un 500
		if err != nil {
			iderror = online.GeneraIDError()
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
			//		jsonerror.InternalMessage = fmt.Sprintf("Error scan condicion. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error getcondicionout2 id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error scan select getcondicionout2 id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Acumulamos en el json
		jsoncondicionout = append(jsoncondicionout, condicionout)
	}
	//al salir del for es cuando creamos el json siempre y cuando tengamos algo en la lectura
	if sidatos {
		JsResponser, err := json.Marshal(jsoncondicionout)
		if err != nil {
			iderror = online.GeneraIDError()
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
			//		jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error getcondicionout2 id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error en la generacion del json  getcondicionout2 id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.Write(JsResponser)
	}
	//Si no tenemos datos sacamos 204
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
}

//postCondicionout para insertar las condiciones de entrada en la planificacion
func postCondicionout(response http.ResponseWriter, request *http.Request) {
	//De entrada tendra un json en el que tendra como obligatorio el nombre, el calendario y el usuario de alta
	//lo primero que hacemos es recuperar el json del body
	cuerpo := json.NewDecoder(request.Body)
	//ceramos la respuesta de body con defer para que se ejecute al final
	defer request.Body.Close()
	//creamos variable donde aplantillaremos
	var postcondicionout structs.Postcondicionout
	//decodificamos en el struc correspondiente
	err := cuerpo.Decode(&postcondicionout)
	//controlamos el error
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error decode postcondicionout. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error postCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error decode body postCondicionout id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//comprobamos que los datos de entrada estan informados
	if postcondicionout.Name != "" {
		if postcondicionout.Condicionout != "" {
			if postcondicionout.Useralt != "" {
				//si todos los datos estan informados primero recuperamos el calendario del programa principal
				sql := fmt.Sprintf("SELECT calendario FROM planificacion WHERE nombre = '%s' AND ejecucion = 'n'", postcondicionout.Name)
				result, err := db2.EjecutaQuery(sql)
				//controlamos error
				if err != nil {
					iderror = online.GeneraIDError()
					//json de error
					jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
					//	jsonerror.InternalMessage = fmt.Sprintf("Error select calendario planificacion. Descripción: %s", err.Error())
					JsResponser, err2 := json.Marshal(jsonerror)
					//si falla la generacion damos error grave
					if err2 != nil {
						mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
						//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
						http.Error(response, mensaje, http.StatusInternalServerError)
						descripcion := "Error en la generacion del json de error postCondicionout id: " + iderror
						online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
						//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
						response.WriteHeader(http.StatusInternalServerError)
						return
					}
					descripcion := "Error select calendario from planificacion postCondicionout id: " + iderror
					online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
					//Creamos cabecera
					response.Header().Set("Content-Type", "application/json")
					//movemos 500 al error
					response.WriteHeader(http.StatusInternalServerError)
					//grabamos el json de error
					response.Write(JsResponser)
					return
				}
				//solo tendra un registro, pero tenemos que montar bucle igual para que no falle
				var calendarplanificacion structs.Calendarplanificacion
				datos := false
				for result.Next() {
					datos = true
					//aplantillamos
					err = result.Scan(&calendarplanificacion.Calendario)
					defer result.Close()
					if err != nil {
						iderror = online.GeneraIDError()
						//json de error
						jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
						//		jsonerror.InternalMessage = fmt.Sprintf("Error scan calendarplanificacion. Descripción: %s", err.Error())
						JsResponser, err2 := json.Marshal(jsonerror)
						//si falla la generacion damos error grave
						if err2 != nil {
							mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
							//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
							http.Error(response, mensaje, http.StatusInternalServerError)
							descripcion := "Error en la generacion del json de error postCondicionout id: " + iderror
							online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
							//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
							response.WriteHeader(http.StatusInternalServerError)
							return
						}
						descripcion := "Error scan postCondicionout id: " + iderror
						online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
						//Creamos cabecera
						response.Header().Set("Content-Type", "application/json")
						//movemos 500 al error
						response.WriteHeader(http.StatusInternalServerError)
						//grabamos el json de error
						response.Write(JsResponser)
						return
					}
					//cerramos result
					defer result.Close()
					//realizamos insert
					sql = fmt.Sprintf("INSERT INTO planificacion VALUES( NULL, '%s','', '','%s','%s','%s',CURRENT_TIMESTAMP,'%s',CURRENT_TIMESTAMP)", postcondicionout.Name, postcondicionout.Condicionout, calendarplanificacion.Calendario, postcondicionout.Useralt, postcondicionout.Useralt)
					//ejecutamos la query
					result, err = db2.EjecutaQuery(sql)
					if err != nil {
						//json de error
						jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
						//		jsonerror.InternalMessage = fmt.Sprintf("Error insert condicionin planificacion. Descripción: %s", err.Error())
						JsResponser, err2 := json.Marshal(jsonerror)
						//si falla la generacion damos error grave
						if err2 != nil {
							mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
							//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
							http.Error(response, mensaje, http.StatusInternalServerError)
							descripcion := "Error en la generacion del json de error postCondicionout id: " + iderror
							online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
							//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
							response.WriteHeader(http.StatusInternalServerError)
							return
						}
						descripcion := "Error einsert postCondicionout id: " + iderror
						online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
						//Creamos cabecera
						response.Header().Set("Content-Type", "application/json")
						//movemos 500 al error
						response.WriteHeader(http.StatusInternalServerError)
						//grabamos el json de error
						response.Write(JsResponser)
						return
					}
					//mostramos el mensaje de creado
					response.WriteHeader(http.StatusCreated)
					defer result.Close()
				}
				if !datos {
					descripcion := "Error no encontrado postCondicionout id: " + iderror
					online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, err)
					//movemos 204 al error
					response.WriteHeader(http.StatusNotFound)
					return
				}
			} else {
				descripcion := "Falta dato obligatorio Useralt postCondicionout id: " + iderror
				online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, err)
				response.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			descripcion := "Falta dato obligatorio Condicionout postCondicionout id: " + iderror
			online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, err)
			response.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		descripcion := "Falta dato obligatorio Name postCondicionout id: " + iderror
		online.EjecutaInfo(nombreservicio, *structs.Entorno, descripcion, err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}
}

//delCondicionout para la elmiminación de la condicion de entrada de la tabla planificación
func delCondicionout(response http.ResponseWriter, request *http.Request) {
	condicion2 := ""
	//recuperamos el id(name) de la url
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	//obtenemos las variables de la url (condición)
	condicion, ok := request.URL.Query()["condicion"]
	//comprobamos que extrae datos de la variable page
	if ok && len(condicion[0]) > 0 {
		//nos quedamos con la primera ocurrencia por si existiese alguna más
		condicion2 = condicion[0]
	}
	//con el name recuperado y la condicion a borrar, montamos la query, lo primero que hacemos es comprobar que
	//existe
	sql := fmt.Sprintf("SELECT COUNT(*) FROM planificacion WHERE nombre = '%s' AND condicionout = '%s'", id, condicion2)
	result, err := db2.EjecutaQuery(sql)
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error select count planificacion. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error delCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error select count delCondicionout id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//creamos variable donde leermos
	var planificacioncount structs.Planificacioncount
	//solo tendremos una ocurrencia
	result.Next()
	defer result.Close()
	//aplantillamos
	err = result.Scan(&planificacioncount.Count)
	//Controlar el error para devolver un 500
	if err != nil {
		iderror = online.GeneraIDError()
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error scan count. Descripción: %s", err.Error())
		JsResponser, err2 := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err2 != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error delCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error scan select delCondicionout id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Solo hacemos delete si el count es > 0
	if planificacioncount.Count > 0 {
		sql = fmt.Sprintf("DELETE FROM planificacion WHERE nombre = '%s' AND condicionout = '%s'", id, condicion2)
		result, err = db2.EjecutaQuery(sql)
		//este defer de resultado, es igual que en los put para que cierre la conexión con mysql
		defer result.Close()
		//Controlar el error para devolver un 500
		if err != nil {
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			//	jsonerror.InternalMessage = fmt.Sprintf("Error delete planificacion condicionin. Descripción: %s", err.Error())
			JsResponser, err2 := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err2 != nil {
				mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
				//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
				http.Error(response, mensaje, http.StatusInternalServerError)
				descripcion := "Error en la generacion del json de error  delCondicionout id: " + iderror
				online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err2)
				//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			descripcion := "Error delete delCondicionout id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
	} else {
		//en caso de no existir lo que hacemos es mostrar un 404
		response.WriteHeader(http.StatusNotFound)
	}
}

//HandlerLog donde devolver la salida de la ejecucion
func HandlerLog(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	//para recuperar las condiciones de entrada de la tabla planif
	case "GET":
		getLog(response, request)
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  HandlerLog id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getLog donde recuperar y devolver en formato txt el archivo de log
func getLog(response http.ResponseWriter, request *http.Request) {
	urlpath := request.URL.Path
	id := path.Base(urlpath)
	pathLog := "C:\\gopath\\src\\github.com\\log\\" + id
	//abrimos el fichero
	file, err := os.Open(pathLog)
	if err != nil {
		//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
		iderror = online.GeneraIDError()
		if os.IsNotExist(err) {
			//movemos 404
			response.WriteHeader(http.StatusNotFound)
			return
		}
		//Informamos el json
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support id: %s", iderror)
		//	jsonerror.InternalMessage = fmt.Sprintf("Error sysout. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si vuelve a fallar la generacion, ya grabamos en log
		if err != nil {
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  getLog id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		descripcion := "Error Open log getLog id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	result, _ := ioutil.ReadAll(file)
	//Creamos cabecera
	response.Header().Set("Content-Type", "text/plain")
	//grabamos el json de error
	response.Write(result)
	//fmt.Println(string(result))
}

//options3 OPTIONS, GET
func options3(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	return
}

//HandlerEnv para recuperar las variables de entorno
func HandlerEnv(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	//para recuperar las condiciones de entrada de la tabla planif
	case "GET":
		getEnv(response, request)
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  HandlerEnv id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getEnv recupera las variables de entorno
func getEnv(response http.ResponseWriter, request *http.Request) {
	//Creamos la variable donde mapearemos los datos
	var envjson structs.EnvJSON
	//Recuperamos los datos
	envjson.Dbhost, _ = os.LookupEnv("DB_HOST")
	envjson.Dbuser, _ = os.LookupEnv("DB_USER")
	envjson.Dbpassword, _ = os.LookupEnv("DB_PASSWORD")
	envjson.Dbdatabase, _ = os.LookupEnv("DB_DATABASE")
	envjson.Servport, _ = os.LookupEnv("SERV_PORT")
	envjson.ServportSSL, _ = os.LookupEnv("SERV_PORT_SSL")
	envjson.Sersafe, _ = os.LookupEnv("SERV_SAFE")
	envjson.PathCert, _ = os.LookupEnv("PATH_CERT")
	envjson.PathKey, _ = os.LookupEnv("PATH_KEY")
	//Generamos el json
	JsResponser, err := json.Marshal(envjson)
	//controlamos el error de json
	if err != nil {
		iderror = online.GeneraIDError()
		mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
		//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
		http.Error(response, mensaje, http.StatusInternalServerError)
		descripcion := "Error en la generacion del json de error  getEnv id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	//creamos cabecera de respuesta
	response.Header().Set("Content-Type", "application/json")
	//devolvemos la respuesta
	response.Write(JsResponser)
}

//HandlerEnvRefresh para recuperar las variables de entorno
func HandlerEnvRefresh(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	//para recuperar las condiciones de entrada de la tabla planif
	case "GET":
		getEnvRefresh(response, request)
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  HandlerEnvRefresh id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

//getEnvRefresh para actualizar la variable de entorno en tiempo de ejecucion
func getEnvRefresh(response http.ResponseWriter, request *http.Request) {
	//Limpiamos las variable de entorno
	os.Clearenv()
	//Volvemos a cargar las variables
	environment.Loadenvironment(*structs.Entorno)
	//Ejecutamos getEnv para que muestre el Json con la info de las variables de entorno
	getEnv(response, request)
}

//HandlerTest para recuperar las variables de entorno
func HandlerTest(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	//para recuperar las condiciones de entrada de la tabla planif
	case "GET":
		getTest(response, request)
	case "OPTIONS":
		options3(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y grabar en log
		if err != nil {
			//ejecutamos la funcion para generar el ID de error para mostrar y grabar en el log y poder localizarlo más rápidamente
			iderror = online.GeneraIDError()
			mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
			//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
			http.Error(response, mensaje, http.StatusInternalServerError)
			descripcion := "Error en la generacion del json de error  HandlerTest id: " + iderror
			online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
			//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
			response.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Aunque saquemos mensaje de error, grabamos
		online.EjecutaInfo(nombreservicio, *structs.Entorno, jsonerror.UserMessage, nil)
		//para que funcione correctamente el orden tiene que ser este. Grabar cabecera, escribir cabecera, escribir cuerpo(json)
		//creamos cabecera de respuesta
		response.Header().Set("Content-Type", "application/json")
		//movemos 405 al error
		response.WriteHeader(http.StatusMethodNotAllowed)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
}

func getTest(response http.ResponseWriter, request *http.Request) {
	var result *sql.Rows
	var err error
	//Creamos variable
	var testJSON structs.TestJSON
	testJSON.Status = "OK"
	//realizamos un select a alguna tabla para comprobar si tenemos conexión con db2
	sql := "SELECT * FROM calendarios"
	result, err = db2.EjecutaQuery(sql)
	//defer result.Close()
	//defer result.Close()
	if err != nil {
		testJSON.ConexDb2 = "KO"

	} else {
		testJSON.ConexDb2 = "OK"
		defer result.Close()
	}
	//Generamos el json
	JsResponser, err := json.Marshal(testJSON)
	//controlamos el error de json
	if err != nil {
		iderror = online.GeneraIDError()
		mensaje := fmt.Sprintf("Fatal Mistake id: %s", iderror)
		//Mostramos error por pantalla para que puedan localizar y tambien lo guardamos en el log
		http.Error(response, mensaje, http.StatusInternalServerError)
		descripcion := "Error en la generacion del json de error  getTest id: " + iderror
		online.EjecutaError(nombreservicio, *structs.Entorno, descripcion, err)
		//movemos 500 al error y no grabamos ni tipo ni json, ya que esto se guardara en log
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	//creamos cabecera de respuesta
	response.Header().Set("Content-Type", "application/json")
	//devolvemos la respuesta
	response.Write(JsResponser)
}
