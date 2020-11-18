package actions

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path"
	"strconv"

	"github.com/jantome/apicm/structs"

	"github.com/jantome/apicm/db2"
)

//Definicion de variables que usara todo el programa
var (
	jsonerror structs.Jsonerror
)

//HandlerEjecucion tabla ejecucion
func HandlerEjecucion(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para validar/usuario y pass y generar token en caso de que sea correcto
		getEjecucion(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options1(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y devolver un 500
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json1. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
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
	//definicion de variables necesarias
	sql := ""
	//control de paginación
	var ofsset int64
	var page int64
	var limit int64
	//obtenemos las variables de la url
	pageurl, ok := request.URL.Query()["page"]
	limiturl, ok2 := request.URL.Query()["limit"]
	//sw para saber si metemos pie
	sipie := false
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
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Error select ejecucion. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
			return
		}
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Creamos variable para aplantillar lectura
	var ejecucion structs.Ejecucion
	//Creamos la variable para devolver la salida con los link relacionados
	var ejecucionjson structs.Ejecucionjson
	//Creamos sw para ver si devolvemos datos o no
	sidatos := false
	//Creamos bucle con la lectura
	for result.Next() {
		sidatos = true
		//hacemos un scan(aplantillar) por cada lectura
		err := result.Scan(&ejecucion.Nombre, &ejecucion.Fechaeje, &ejecucion.Estado)
		//Controlar el error y devolver un 500
		if err != nil {
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error bucle ejecucion. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//aplantillamos
		ejecucionjson.Nombre = ejecucion.Nombre
		ejecucionjson.Fechaeje = ejecucion.Fechaeje
		ejecucionjson.Estado = ejecucion.Estado
		//Montamos las Url's que se pueden usar
		ejecucionjson.Links.Href = make(map[string]string)
		ejecucionjson.Links.Href["condicionin"] = "/cm/v1/ejecuciones/condicionin/" + ejecucion.Nombre + "?fechaeje=" + ejecucion.Fechaeje
		ejecucionjson.Links.Href["condicionout"] = "/cm/v1/ejecuciones/condicionout/" + ejecucion.Nombre + "?fechaeje=" + ejecucion.Fechaeje //Montamos
		JsResponser, err := json.Marshal(ejecucionjson)
		//Controlar el error y devolver un 500
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
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
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}

	if sidatos {
		if sipie {
			var pagmax float64
			//montamos url prev y next
			var piepaginacion structs.Piepaginacion
			piepaginacion.Pagdet.Links.Href = make(map[string]string)
			if page > 1 {
				piepaginacion.Pagdet.Links.Href["prev"] = fmt.Sprintf("/cm/v1/ejecuciones?limit=%d&page=%d", limit, page-1)
				piepaginacion.Pagdet.Pagprev = page - 1
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
				//json de error
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
				jsonerror.InternalMessage = fmt.Sprintf("Error select max. Descripción: %s", err.Error())
				JsResponser, err := json.Marshal(jsonerror)
				//si falla la generacion damos error grave
				if err != nil {
					http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
					return
				}
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
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
					piepaginacion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/ejecuciones?limit=%d&page=%d", limit, 2)
					piepaginacion.Pagdet.Pagnext = 2
				}
			} else {
				if page < pagmax64 {
					piepaginacion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/ejecuciones?limit=%d&page=%d", limit, page+1)
					piepaginacion.Pagdet.Pagnext = page + 1
				}
			}
			piepaginacion.Pagdet.Pagmax = pagmax64
			//creamos json
			JsResponser, err := json.Marshal(piepaginacion)
			//Controlar el error y devolver un 500
			if err != nil {
				//Informamos el json
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
				jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
				JsResponser, err := json.Marshal(jsonerror)
				//si vuelve a fallar la generacion, ya grabamos en log
				if err != nil {
					http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
					return
				}
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
	}
}

//options1 para los cors de esta api
func options1(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	return
}

//HandlerCondicionin condiciones de entrada --> OK
func HandlerCondicionin(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para recuperar la info
		getCondicionin(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options1(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y devolver un 500
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json1. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
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

//getCondicionnin condiciones de entrada en json --> Arreglado json de salida --> OK
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
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Invalid Parameter url: fechaeje")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
			return
		}
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Query para recuperar todas las condiciones de entrada
	sql := fmt.Sprintf("SELECT condicionin FROM ejecucion WHERE nombre ='%s' AND condicionin > '' and FECHAEJE ='%s'", id, fechaeje2)
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Error select condicionin. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
			return
		}
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
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error scan condicion. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
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
			jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
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

//HandlerCondicionout condiciones de salida --> OK
func HandlerCondicionout(response http.ResponseWriter, request *http.Request) {
	//Methodos permitidos GET-OPTIONS
	switch request.Method {
	case "GET":
		//ejecutamos la funcion para recuperar la info
		getCondicionout(response, request)
	//tenemos que habilitar el metodo options, para que se puedan verificar los cors
	case "OPTIONS":
		options1(response, request)
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y devolver un 500
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json1. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
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

//getCondicionnout condiciones de salida en json --> Arreglado json de salida --> Ok
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
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Invalid Parameter url: fechaeje")
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
			return
		}
		//Creamos cabecera
		response.Header().Set("Content-Type", "application/json")
		//movemos 500 al error
		response.WriteHeader(http.StatusInternalServerError)
		//grabamos el json de error
		response.Write(JsResponser)
		return
	}
	//Query para recuperar todas las condiciones de entrada
	sql := fmt.Sprintf("SELECT condicionout FROM ejecucion WHERE nombre ='%s' AND condicionout > '' and FECHAEJE ='%s'", id, fechaeje2)
	result, err := db2.EjecutaQuery(sql)
	//Controlar el error para devolver un 500
	if err != nil {
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Error select condicionout. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
			return
		}
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
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error scan condicionout. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
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
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
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
	default:
		jsonerror.UserMessage = fmt.Sprintf("Not implemented Method %s", request.Method)
		//Montamos el json de error
		JsResponser, err := json.Marshal(jsonerror)
		//Controlar el error y devolver un 500
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json1. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
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
		//json de error
		jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
		jsonerror.InternalMessage = fmt.Sprintf("Error select ejecucion. Descripción: %s", err.Error())
		JsResponser, err := json.Marshal(jsonerror)
		//si falla la generacion damos error grave
		if err != nil {
			http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
			return
		}
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
	var planificacion structs.Planificacion
	for result.Next() {
		sidatos = true
		//hacemos un scan(aplantillar) por cada lectura
		err := result.Scan(&planificacion.Nombre, &planificacion.Calendario, &planificacion.Useralta, &planificacion.Timalta, &planificacion.Usermod, &planificacion.Timesmod)
		//Controlar el error y devolver un 500
		if err != nil {
			//json de error
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error bucle planificacion. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si falla la generacion damos error grave
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
			//Creamos cabecera
			response.Header().Set("Content-Type", "application/json")
			//movemos 500 al error
			response.WriteHeader(http.StatusInternalServerError)
			//grabamos el json de error
			response.Write(JsResponser)
			return
		}
		//Como no hemos tenido error creamos el json de salida
		JsResponser, err := json.Marshal(planificacion)
		//Controlar el error y devolver un 500
		if err != nil {
			//Informamos el json
			jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
			jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
			JsResponser, err := json.Marshal(jsonerror)
			//si vuelve a fallar la generacion, ya grabamos en log
			if err != nil {
				http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
				return
			}
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
	if !sidatos {
		//movemos 204 al error
		response.WriteHeader(http.StatusNoContent)
		return
	}
	if sidatos {
		if sipie {
			var pagmax float64
			//montamos url prev y next
			var piepaginacion structs.Piepaginacion
			piepaginacion.Pagdet.Links.Href = make(map[string]string)
			if page > 1 {
				piepaginacion.Pagdet.Links.Href["prev"] = fmt.Sprintf("/cm/v1/planificacion?limit=%d&page=%d", limit, page-1)
				piepaginacion.Pagdet.Pagprev = page - 1
			}
			//para armar el next cuando corresponde necesitamos saber el número de paginas que tendra
			//hacemos un select count
			sql = "SELECT COUNT(*) FROM planificacion WHERE ejecucion = 'n'"
			result, err = db2.EjecutaQuery(sql)
			//solo tendra un registro por lo que no hace falta que recorramos toda la tabla
			var planificacionout structs.Planificacionout
			result.Next()
			err = result.Scan(&planificacionout.Count)
			//Controlar el error y devolver un 500
			if err != nil {
				//json de error
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
				jsonerror.InternalMessage = fmt.Sprintf("Error select max. Descripción: %s", err.Error())
				JsResponser, err := json.Marshal(jsonerror)
				//si falla la generacion damos error grave
				if err != nil {
					http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
					return
				}
				//Creamos cabecera
				response.Header().Set("Content-Type", "application/json")
				//movemos 500 al error
				response.WriteHeader(http.StatusInternalServerError)
				//grabamos el json de error
				response.Write(JsResponser)
				return
			}
			//convertimos el resultado del count a float64
			countfloat := float64(planificacionout.Count)
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
					piepaginacion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/planificacion?limit=%d&page=%d", limit, 2)
					piepaginacion.Pagdet.Pagnext = 2
				}
			} else {
				if page < pagmax64 {
					piepaginacion.Pagdet.Links.Href["next"] = fmt.Sprintf("/cm/v1/planificacion?limit=%d&page=%d", limit, page+1)
					piepaginacion.Pagdet.Pagnext = page + 1
				}
			}
			piepaginacion.Pagdet.Pagmax = pagmax64
			//creamos json
			JsResponser, err := json.Marshal(piepaginacion)
			//Controlar el error y devolver un 500
			if err != nil {
				//Informamos el json
				jsonerror.UserMessage = fmt.Sprintf("Internal error, contact support")
				jsonerror.InternalMessage = fmt.Sprintf("Error json2. Descripción: %s", err.Error())
				JsResponser, err := json.Marshal(jsonerror)
				//si vuelve a fallar la generacion, ya grabamos en log
				if err != nil {
					http.Error(response, "Error Grave generacion Json de error", http.StatusInternalServerError)
					return
				}
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
	}
}

//options2 para los cors de esta api
func options2(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Origin", "*")
	//Get lista
	//Put update
	//Post create
	//Delete delete
	response.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
	return
}
