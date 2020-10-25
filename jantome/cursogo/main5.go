//paquete inicial
package main

import (
	"fmt"
	"strconv"
)

func main() {
	//para declarar canal
	canal := make(chan string)
	/*punteros()*/
	lanzaHilos(200, canal)

	//mientras tenga valores el canal
	for valor := range canal {
		fmt.Println(valor)
	}

}

/*func punteros() {
	x := 100
	//con el * indica que es un puntero (x almacena el dato e Y almacena la dirección de memoria)
	var y *int
	//devuelve la dirección  de memoria
	y = &x
	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(*y)

	//modificar el valor atras dle puntero
	*y = 500
	fmt.Println(*y)

}*/

func holaMundo(i int, canal chan<- string) {
	//fmt.Println("Hola Mundo Numero: ", i)
	//gurdamos en canal --> Transforma a string strconv
	canal <- "Hola Mundo Número:" + strconv.Itoa(i)
}

// en chan solo puede escribir tipos string, si la fecha fuese al reves seria de lectura
func lanzaHilos(numHilos int, canal chan<- string) {
	for i := 0; i < numHilos; i++ {
		//mismo hilo sin go, hilos diferentes con go
		go holaMundo(i, canal)
	}
}
