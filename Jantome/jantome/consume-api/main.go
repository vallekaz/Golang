package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Tarea struct {
	UserId    int    `json:userId`
	Id        int    `json:id`
	Title     string `json:title`
	Completed bool   `json:completed`
}

func main() {
	var urlApi = "https://jsonplaceholder.typicode.com/todos"

	var cliente = &http.Client{Timeout: 10 * time.Second}
	//llamar cleinte
	response, err := cliente.Get(urlApi)

	if err != nil {
		panic(err.Error())
	}
	//creamos una variable tareas que sera con la estructrua de Tarea
	var tareas []Tarea
	//decodficmoas
	json.NewDecoder(response.Body).Decode(&tareas)
	//si huiese error nos lo miestra
	err = json.NewDecoder(response.Body).Decode(&tareas)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(tareas)

}
