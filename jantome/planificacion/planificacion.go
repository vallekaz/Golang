package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jantome/planificacion/environment"

	"github.com/jantome/planificacion/db2"
	"github.com/jantome/planificacion/structs"
)

var (
	//recuperamos la fecha de ejecucion y la formateamos separada
	date = time.Now()
	anno = fmt.Sprintf("%d", date.Year())
	mes  = fmt.Sprintf("%d", date.Month())
	//Convertido a String
	dia = fmt.Sprintf("%d", date.Day())
	//Sin convertir a string
	//dia := date.Day()
	//montamos la fecha de controlm formato YYMMDD
	fechacm = fmt.Sprintf("%s%s%s", anno[0:2], mes, dia)
	//numero día
	numdia = fmt.Sprintf("%d", date.YearDay())
)

func planificar() {
	//recuperamos el entorno de ejecuion para saber las rutas
	entorno := flag.String("entorno", "", "entorno de ejecución")
	flag.Parse()
	//cargamos variables de entorno
	environment.Loadenvironment(*entorno)
	//recuperamos los calendarios que se tienen que planificar
	//PONEMOS ` en el and, ya que es como lo reconoce al ser numeros en MYSQL
	query := fmt.Sprintf("SELECT nombre FROM calendarios WHERE year = '%s' AND `%s` = 'Y'", anno, numdia)
	result, err := db2.EjecutaQuery(query)
	//comprobamos error
	if err != nil {
		fmt.Println("Error seleccion calendario", err.Error())
		os.Exit(1)
		return
	}
	//montamos where de calendario
	var calendario structs.Calendario
	//Switch para saber si es la primera vez o no
	primeravez := true
	//variable where
	valor := ""
	for result.Next() {
		err = result.Scan(&calendario.Namecalendar)
		if err != nil {
			fmt.Println("Error en la lectura de la tabla calendario", err.Error())
			os.Exit(1)
			return
		}
		if primeravez {
			valor = fmt.Sprintf("WHERE calendario = '%s'", calendario.Namecalendar)
			primeravez = false
		} else {
			valor = valor + fmt.Sprintf("OR calendario = '%s'", calendario.Namecalendar)
		}
	}

	//recuperamos los datos de la tabla de planificación
	query = "SELECT * FROM planificacion " + valor + " ORDER BY nombre, numsec"
	result, err = db2.EjecutaQuery(query)
	//comprobamos error
	if err != nil {
		fmt.Println("Error en la lectura de planificacion: ", query, err.Error())
		os.Exit(1)
		return
	}
	//creamos la variable donde formatearemos la lectura
	var planifi structs.Planificacion
	//Variable para numero secuencial
	newnumsec := 0
	//creamos la variable de nombre anterior para comprar
	nombreAnt := ""
	//contador con los jobs planificados
	jobplanif := 0
	//recoremos la tabla y vamos planificando
	for result.Next() {
		//realizamos la lectura
		err = result.Scan(&planifi.Numsec, &planifi.Nombre, &planifi.Ejecucion, &planifi.Condicionin, &planifi.Condicionout, &planifi.Calendario, &planifi.Useralta, &planifi.Timalta, &planifi.UserModif, &planifi.Timesmod)
		if err != nil {
			fmt.Println("Error en la lectura de la tabla planificacion", err.Error())
			os.Exit(1)
			return
		}
		//Si el nombre no es el mismo al de la lectura anterior, grabamos en la BBDD, el numsec 1
		//que sera la planificacion (las condiciones seran guardados en numsec > 1)
		if nombreAnt != planifi.Nombre {
			query = fmt.Sprintf("INSERT INTO ejecucion VALUES('%s', 1, '%s', '','','pl')", planifi.Nombre, fechacm)
			_, err = db2.EjecutaQuery(query)
			if err != nil {
				fmt.Println("Error insert en ejecucion 1: ", err.Error())
				os.Exit(1)
				return
			}
			nombreAnt = planifi.Nombre
			//ponemos 1 al numsec
			newnumsec = 1
			jobplanif = jobplanif + 1
		}
		//comprobamos que tiene informado y damos de alta
		if planifi.Condicionin != "" {
			//Sumamos 1 al numsec que sera lo que se grabe
			newnumsec = newnumsec + 1
			//realizamos el insert
			query = fmt.Sprintf("INSERT INTO ejecucion VALUES('%s', %d, '%s', '%s','','')", planifi.Nombre, newnumsec, fechacm, planifi.Condicionin)
			_, err = db2.EjecutaQuery(query)
			if err != nil {
				fmt.Println("Error insert en ejecucion condicionin: ", err.Error())
				os.Exit(1)
				return
			}
		}

		if planifi.Condicionout != "" {
			//Sumamos 1 al numsec que sera lo que se grabe
			newnumsec = newnumsec + 1
			//realizamos el insert
			query = fmt.Sprintf("INSERT INTO ejecucion VALUES('%s', %d, '%s', '','%s','')", planifi.Nombre, newnumsec, fechacm, planifi.Condicionout)
			_, err = db2.EjecutaQuery(query)
			if err != nil {
				fmt.Println("Error insert en ejecucion condicionout: ", err.Error())
				os.Exit(1)
				return
			}
		}

	}
	fmt.Println("Numero de Jobs planficados: ", jobplanif)
}

func limpia() {
	fmt.Println("Limpieza CM los ejecutados OK...")
	query := fmt.Sprintf("DELETE FROM ejecucion WHERE estado ='ok'")
	_, err := db2.EjecutaQuery(query)
	//Comprobamos Error
	if err != nil {
		fmt.Println("Error delete limpia: ", err.Error())
		os.Exit(1)
		return
	}
	fmt.Println("Fin limpieza...")
}

func main() {
	fmt.Println("Comienza planificación...")
	//primero limpia CM de la ejecucion anterior
	limpia()
	//planifica
	planificar()
	fmt.Println("Fin planificacion...")
}