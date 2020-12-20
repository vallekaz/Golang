package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/batcharq/batch"
)

var (
	//recuperamos el entorno de ejecucion mediante flags para saber las rutas
	entorno = flag.String("entorno", "", "entorno de ejecución")
)

func main() {
	flag.Parse()
	//entorno := "local"
	fmt.Println("inicio batch")
	batch.Start("prueba", *entorno)
	fmt.Println("fin inicio batc")
	fmt.Println("Hola")
	contador := 90
	cabe := "\n******Estadisticas******"
	esta1 := "\nRegistros leidos: " + strconv.Itoa(contador)
	esta2 := "\nRegistros tratados: " + strconv.Itoa(contador-10)
	esta3 := "\nRegistros no tratados: 10"
	esta4 := "\n******Fin Estadistidas******"
	inf := cabe + esta1 + esta2 + esta3 + esta4
	batch.Impr("prueba", inf, "w", *entorno)
	time.Sleep(60 * time.Second)
	batch.FinOk("prueba", *entorno)
	retorno := "100"
	descripcion := "Error en bla bla lba status.."
	batch.FinKo("prueba", retorno, descripcion, *entorno)
	os.Exit(0)
}
