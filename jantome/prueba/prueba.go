package main

import (
	"fmt"
	"os"
)

func main() {
	//entorno := "local"
	batch.Start("prueba", entorno)*/
	fmt.Println("Hola")
	contador := 90
	cabe := "\n******Estadisticas******"
	esta1 := "\nRegistros leidos: " + strconv.Itoa(contador)
	esta2 := "\nRegistros tratados: " + strconv.Itoa(contador-10)
	esta3 := "\nRegistros no tratados: 10"
	esta4 := "\n******Fin Estadistidas******"
	inf := cabe + esta1 + esta2 + esta3 + esta4
	batch.Impr("prueba", inf, "w", entorno)
	time.Sleep(60 * time.Second)
	batch.FinOk("prueba", entorno)
	retorno := "100"
	descripcion := "Error en bla bla lba status.."
	batch.FinKo("prueba", retorno, descripcion, entorno)
	os.Exit(0)
}
