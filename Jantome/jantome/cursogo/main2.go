//paquete inicial
package main

//importaciones
//importación para imprimir por pantalla
import (
	"fmt"
)

//constante
const holaMundo string = "Hola2 %s %s, bienvenido al curso de Go"

//funcion que debe de tener todos (funcion principal)
func main() {
	//imprimir por pantalla
	fmt.Print("Hola Mundo")
	//creacion de variables
	var name string
	//pedri ingresar nombre
	fmt.Print("Ingresa tu nombre: ")
	//%s lo capturado de teclado esta guardado en la variable name (El aspersan es xk vamos a cambiar el contenido de la variable)
	fmt.Scanf("%s", &name)
	//imprimir
	//espacio de linea ln
	fmt.Println("Hola ", name)
	//con prinft nosotros damos formato (donde este %s es donde se pone la variable en los anteriores solo estaria al final)
	fmt.Printf("Hola %s, bienvenido al curso de Go", name)
	//variable inicializada
	//var name string = "Nombre por defecto"
	//variable sin tipo pero tiene que ser inicializada para que sepa go de que tipo es
	//el := declar y asigan un valor a una variable
	lastname := "lastname"
	//otra opcion es sin tipo pero iniizalizando
	//var miNumero = 100
	//varias variables a la vez
	//var (
	//		x = 1
	//		y = 2
	//		z = 3
	//	)
	fmt.Println("")
	fmt.Printf(holaMundo, name, lastname)
}
