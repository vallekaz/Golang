package main

//usaremos 2 libreriras http e IO para cuando se hace llamadas al servidor

import (
	"encoding/json"
	"io"
	"net/http"
)

//ruebas
type Persona struct {
	//nombre de la persona
	Nombre string
	//apellidos de la persona
	Apellidos string
}

func main() {
	//creamos enpoind se define enpoind + funcion que se quiere ejecutar
	http.HandleFunc("/", handlerRaiz)
	//segunda ruta
	//creamos enpoind se define enpoind + funcion que se quiere ejecutar
	http.HandleFunc("/usuarios", handlerUsuarios)

	//creamos enpoind se define enpoind + funcion que se quiere ejecutar
	http.HandleFunc("/personas", handlerPersonas)

	//creamos el servidor
	http.ListenAndServe(":8000", nil)

}

// la funcion recibe como parametro una rsponse y una request, si optimizamos se puede pasar rques como puntero (*)
func handlerRaiz(response http.ResponseWriter, request *http.Request) {
	//con io scribimos el json o el string en este caso
	io.WriteString(response, "Hola Mundo Api!")
}

// la funcion recibe como parametro una rsponse y una request, si optimizamos se puede pasar rques como puntero (*)
func handlerUsuarios(response http.ResponseWriter, request *http.Request) {
	//con io scribimos el json o el string en este caso
	io.WriteString(response, "Hola Usuarios!")
}

// la funcion recibe como parametro una rsponse y una request, si optimizamos se puede pasar rques como puntero (*)
func handlerPersonas(response http.ResponseWriter, request *http.Request) {
	persona := Persona{"Antonio", "Mesa"}
	//convertimos en json las variables on el json convertido + error
	jsResponse, err := json.Marshal(persona)
	//si al convertir tenemos error
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	//si no tenemos error
	response.Header().Set("Content-Type", "application/json")
	response.Write(jsResponse)

}
