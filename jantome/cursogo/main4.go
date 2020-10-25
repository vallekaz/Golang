//paquete inicial
package main

//importaciones
//importación para imprimir por pantalla
import (
	"fmt"
	/*"github.com/jantome/cursogo/maps"*/)

//propia variable
type Galisteo string

//struc
type Persona struct {
	Nombre             string
	Apellidos          string
	DocumentoIdentidad string
	Telefono           []string
	Direccion          string
	Edad               int
}

//una struc dentro de otra
type Casa struct {
	NumeroCasa int
	Personas   []Persona
}

//funcion que debe de tener todos (funcion principal)
func main() {

	/*fmt.Println(maps.GetMap())
	fmt.Println(maps.GetKeyMap("Maria"))*/
	var miVariable Galisteo = "Mi propia variable"
	fmt.Println(miVariable)

	antonio := Persona{
		Nombre:             "Antonio",
		Apellidos:          "Galisteo",
		DocumentoIdentidad: "5555",
		Telefono:           []string{"111", "2322"},
		Direccion:          "yyy",
		Edad:               30,
	}

	fmt.Println(antonio)

	maria := Persona{
		Nombre:             "Maria",
		Apellidos:          "Galisteo",
		DocumentoIdentidad: "5555",
		Telefono:           []string{"111", "2322"},
		Direccion:          "yyy",
		Edad:               55,
	}

	//otra manera
	jorge := new(Persona)
	jorge.Nombre = "jorge"
	jorge.Apellidos = "apellidos"
	jorge.DocumentoIdentidad = "333"
	jorge.Telefono = []string{"3232", "2343"}
	jorge.Direccion = "yyy"
	jorge.Edad = 44

	// cuando es puntero
	casa := Casa{
		NumeroCasa: 1,
		Personas:   []Persona{antonio, maria, *jorge},
	}

	fmt.Println(casa)

	casa.GetNumeroCasa()
	casa.GetPersonasCasa()

}

//para recuperar algo de un struct
func (c Casa) GetNumeroCasa() {
	fmt.Println("El número de la casa es: ", c.NumeroCasa)
}

func (c Casa) GetPersonasCasa() {
	fmt.Println("Las personas de la casa son: ", c.Personas)
}
