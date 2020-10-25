//paquete inicial
package main

//importaciones
//importación para imprimir por pantalla
import (
	"fmt"
	"strings"
)

//constante
const holaMundo string = "Hola2 %s %s, bienvenido al curso de Go"

//funcion que debe de tener todos (funcion principal)
func main() {
	// declaramos variable llamando a la funcion
	//name := getName()
	//x, y, z := getMultiplesVariables()
	//imprimir por pantalla
	fmt.Print("Hola Mundo")

	//imprimir
	//espacio de linea ln
	//fmt.Println("Hola ", name)
	//con prinft nosotros damos formato (donde este %s es donde se pone la variable en los anteriores solo estaria al final)
	//fmt.Printf("Hola %s, bienvenido al curso de Go", name)
	//variable inicializada
	//var name string = "Nombre por defecto"
	//variable sin tipo pero tiene que ser inicializada para que sepa go de que tipo es
	//el := declar y asigan un valor a una variable
	//lastname := "lastname"
	//otra opcion es sin tipo pero iniizalizando
	//var miNumero = 100
	//varias variables a la vez
	//var (
	//		x = 1
	//		y = 2
	//		z = 3
	//	)
	fmt.Println("")
	//fmt.Printf(holaMundo, name, lastname)
	fmt.Println("")
	//definimos variable
	hola := "HOla"
	//mostramos la primera letra esta en ascii
	fmt.Println(hola[0])
	//ahora se imprime con nombre en vez de en ascci
	fmt.Println(string(hola[0]))
	//sacar longitud de la letra
	fmt.Println(len(hola))
	//arrays y slide
	//array se tiene que declarar el tamaño slide es una array dinamico
	imprimeArray()
	imprimeSlice()
	multiplo5()
	operacionesConString()
}

//Declaracion de funciones
func getName() string {
	//creacion de variables
	var name string
	//pedri ingresar nombre
	fmt.Print("Ingresa tu nombre: ")
	//%s lo capturado de teclado esta guardado en la variable name (El aspersan es xk vamos a cambiar el contenido de la variable)
	fmt.Scanf("%s", &name)
	return name
}

func getMultiplesVariables() (int, int, int) {
	return 1, 2, 3
}

func imprimeArray() {
	//nombre variable, el numero de ocurrencias, y el tipo de dato
	var array1 [2]string
	array1[0] = "Hola"
	array1[1] = "Mundo"
	fmt.Println(array1)
	//otra forma de declarar array
	array2 := [4]int{1, 2, 3, 4}
	fmt.Println(array2)

	// matriz

	var matriz [2][2]string
	matriz[0][0] = "Hola"
	matriz[0][1] = "Mundo"
	matriz[1][0] = "curso"
	matriz[1][1] = "go"
	fmt.Println(matriz)

	matriz2 := [2][2]int{{1, 2}, {3, 4}}
	fmt.Println(matriz2)

}

func imprimeSlice() {
	//array variable
	var slice1 []string
	//no tiene posición entonces se hace appedn
	slice1 = append(slice1, "Hola", "Slice")
	fmt.Println(slice1)
	fmt.Println(len(slice1))
	slice1 = append(slice1, "NuevoElemento")
	fmt.Println(slice1)
	fmt.Println(len(slice1))

	//Matrices en slide
	matriz := [][]string{{"hola", "Mundo"}, {"Curso", "Go"}}
	fmt.Println(matriz)

	//forma 2
	var matriz2 [][]string
	row1 := []string{"Hola2", "SliceMatriz2"}
	row2 := []string{"Curso2", "Go2"}
	matriz2 = append(matriz2, row1)
	matriz2 = append(matriz2, row2)
	fmt.Println(matriz2)

}

func multiplo5() {
	var numero = 0
	fmt.Println("Ingresa un número")
	//al ser un numero es %d
	fmt.Scanf("%d", &numero)
	//modulo 5 es dividir el número entre 5
	if numero%5 == 0 {
		fmt.Println("Es múltiplo de 5")
	} else {
		fmt.Println("No es múltiño de 5")
	}
}

func operacionesConString() {
	var texto = "Hola Go, Hola Antonio, Hola Mundo"
	var texto2 = "Prueba"
	fmt.Println(texto)
	fmt.Println(strings.ToUpper(texto))
	fmt.Println(strings.ToLower(texto))
	//reemplazar
	fmt.Println(strings.Replace(texto, "Hola", "Adiós", -1))
	fmt.Println(strings.ReplaceAll(texto, "Hola", "Adios"))
	//compara
	fmt.Println(strings.Compare(texto, texto2))
	//separacion de cadenaa
	fmt.Println(strings.Split(texto, ","))
	//buscar
	fmt.Println(strings.Contains(texto, "Hola"))
	fmt.Println(strings.Contains(texto, "coche"))

}
