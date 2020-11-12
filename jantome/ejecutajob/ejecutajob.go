package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/jantome/ejecutajob/db2"
	"github.com/jantome/ejecutajob/environment"
	"github.com/jantome/ejecutajob/structs"
)

var (
	// inicializado a true para que entre la primera vez
	repite = true
	//recuperamos el entorno de ejecucion mediante flags para saber las rutas
	entorno = flag.String("entorno", "", "entorno de ejecución")
	date    = time.Now()
	anno    = fmt.Sprintf("%d", date.Year())
	mes     = fmt.Sprintf("%d", date.Month())
	//Convertido a String
	dia = fmt.Sprintf("%d", date.Day())
	//montamos la fecha de controlm formato YYMMDD (Este formato para MYSQL como date lo interpreta correctamente)
	fechacm = fmt.Sprintf("%s%s%s", anno[2:4], mes, dia)
)

func main() {
	flag.Parse()
	//cargamos variables de entorno
	environment.Loadenvironment(*entorno)
	for repite == true {
		ejecuta()
	}

}

func ejecuta() {
	//variable inicializada
	comando := ""
	consola := ""
	letra := ""
	//nada más entrar ponemos a false el SW y solo se activa en caso de actualizar condiciones de salida
	repite = false
	//recuperamos primero los numsec 1 que estan en la tabla con estado pl para ver recorrer luego las condiciones y ver que
	//se puede ejecutar
	sql := fmt.Sprintf("SELECT nombre FROM ejecucion WHERE estado ='pl' AND fechaeje = '%s'", fechacm)
	result, err := db2.EjecutaQuery(sql)

	if err != nil {
		fmt.Println("Error select nombre", err.Error())
		os.Exit(1)
		return
	}

	//Montamos bucle, por cada registro leido, comprobaremos si tiene conciones de entra sin cumplir
	//todos los struct los genero fuero para tener solo 1
	var ejecucion structs.Ejecucion

	//creamos las variables qui para que no las cree mil veces
	for result.Next() {
		//leemos y aplantillamos
		err = result.Scan(&ejecucion.Nombre)
		if err != nil {
			fmt.Println("Error en la lectura nombre", err.Error())
			os.Exit(1)
			return
		}

		//vemos si tiene alguna condicion de entrada sin cumplir
		sql2 := fmt.Sprintf("SELECT count(*) FROM ejecucion WHERE nombre = '%s' and condicionin > ' ' AND estado <> 'ok' AND fechaeje = '%s' ", ejecucion.Nombre, fechacm)
		result2, err := db2.EjecutaQuery(sql2)

		if err != nil {
			fmt.Println("Error select count", err.Error())
			os.Exit(1)
			return
		}
		//Result solo tiene que tener un resultado con el numero de filas, en caso de ser mayor de 0 es que todavía tiene
		//condiciones pendientes. No se monta bucle al tener solo una fila
		result2.Next()
		var countEje structs.CountEje
		//aplantillamos
		err = result2.Scan(&countEje.Numero)
		//comprobamos si tiene o no condiciones pendientes
		//if countEje.Numero > 0 {
		//	fmt.Println("condiciones pendientes")
		//}

		if countEje.Numero == 0 {
			//	fmt.Println("Ninguna condicion pendiente")
			//como no tiene ninguna condicion esperando ejecutamos job pero antes actualizmos a en ejecucion
			sql3 := fmt.Sprintf("UPDATE ejecucion SET estado = 'ej' WHERE nombre = '%s' AND numsec = 1 AND fechaeje = '%s'", ejecucion.Nombre, fechacm)
			_, err = db2.EjecutaQuery(sql3)
			if err != nil {
				fmt.Println("Error update KO", err.Error())
				os.Exit(1)
				return
			}
			//Ejecuta y saca la salida directamente teniendo en cuenta el entorno de ejecucion
			if *entorno == "local" {
				comando = "go run c:\\gopath\\src\\github.com\\jantome\\" + ejecucion.Nombre + "\\" + ejecucion.Nombre + ".go"
				consola = "cmd"
				letra = "/C"
			} else {
				//ruta linux
				comando = "cd /ejecutable/batch/app; ./" + ejecucion.Nombre
				consola = "bash"
				letra = "-c"
			}
			//ejecuta := "go run c:\\gopath\\src\\github.com\\jantome\\" + ejecucion.Nombre + "\\" + ejecucion.Nombre + ".go"
			c := exec.Command(consola, letra, comando)
			_, err := c.Output()
			//controlamos el error, y si falla actualizamos la ejecucion como fallida
			if err != nil {
				//Este print le dejo por el momento, pero sera el job el que deje un log de salida
				fmt.Println(err.Error())
				//actualizamos a fallido la ejecucion
				sql3 = fmt.Sprintf("UPDATE ejecucion SET estado = 'ko' WHERE nombre = '%s' AND numsec = 1 AND fechaeje = '%s'", ejecucion.Nombre, fechacm)
				_, err = db2.EjecutaQuery(sql3)
				if err != nil {
					fmt.Println("Error update KO", err.Error())
					os.Exit(1)
					return
				}
			} else {
				//En caso de no tener error, grabamos la marca de OK
				sql3 := fmt.Sprintf("UPDATE ejecucion SET estado= 'ok' WHERE nombre = '%s' AND fechaeje = '%s'", ejecucion.Nombre, fechacm)
				_, err = db2.EjecutaQuery(sql3)
				if err != nil {
					fmt.Println("Error update OK", err.Error())
					os.Exit(1)
					return
				}
				//Una vez actualizado a Ok, buscamos si tiene condiciones de salida, en caso de que las tenga
				//las cumplimos

				sql3 = fmt.Sprintf("SELECT condicionout FROM ejecucion WHERE nombre ='%s' AND condicionout > '' AND fechaeje = '%s' ", ejecucion.Nombre, fechacm)
				result3, err := db2.EjecutaQuery(sql3)

				if err != nil {
					fmt.Println("Error select condicionout", err.Error())
					os.Exit(1)
					return
				}

				//leemos
				var ejecucion2 structs.Ejecucion2
				for result3.Next() {
					err = result3.Scan(&ejecucion2.Condicionout)
					//realizmos update de la condicion de salida en toda la tabla
					sql4 := fmt.Sprintf("UPDATE ejecucion SET estado='ok' WHERE condicionin = '%s' AND fechaeje = '%s'", ejecucion2.Condicionout, fechacm)
					_, err = db2.EjecutaQuery(sql4)

					if err != nil {
						fmt.Println("Error udate Ok condicion out", err.Error())
						os.Exit(1)
						return
					}
					//si hemos actualizado condicionin, significa que tenemos que volver a lanzar el bucle, por si
					//algun proceso más se puede volver a ejecutar.
					repite = true
				}
			}
			//Este print le dejamos por el momento, pero lo eliminaremos.
			//fmt.Println(string(cout))
		}
	}
}
