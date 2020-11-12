package main

import (
	"fmt"
	"os"
)

func main() {
	//batch.Start("prueba")
	fmt.Println("Hola")
	/*contador := 90
	cabe := "\n******Estadisticas******"
	esta1 := "\nRegistros leidos: " + strconv.Itoa(contador)
	esta2 := "\nRegistros tratados: " + strconv.Itoa(contador-10)
	esta3 := "\nRegistros no tratados: 10"
	esta4 := "\n******Fin Estadistidas******"
	inf := cabe + esta1 + esta2 + esta3 + esta4
	batch.Impr("prueba", inf, "w")
	time.Sleep(60 * time.Second)
	batch.FinOk("prueba")
	retorno := "100"
	descripcion := "Error en bla bla lba status.."
	batch.FinKo("prueba", retorno, descripcion)*/
	os.Exit(0)
}
