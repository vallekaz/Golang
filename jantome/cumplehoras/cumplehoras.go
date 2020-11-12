package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/jantome/cumplehoras/structs"

	"github.com/jantome/cumplehoras/db2"

	"github.com/jantome/cumplehoras/environment"
)

var (
	//recuperamos la fecha de ejecucion y la formateamos separada
	date = time.Now()
	anno = fmt.Sprintf("%d", date.Year())
	mes  = fmt.Sprintf("%d", date.Month())
	//Convertido a String
	dia = fmt.Sprintf("%d", date.Day())
	//montamos la fecha de controlm formato YYMMDD
	fechacm = fmt.Sprintf("%s%s%s", anno[2:4], mes, dia)
	entorno = flag.String("entorno", "", "entorno de ejecución")
)

func main() {
	//Ejemplo ejecuciones para lanzar objetos
	//Ejecuta y saca la salida directamente
	/*c := exec.Command("cmd", "/C", "go run c:\\gopath\\src\\github.com\\jantome\\prueba\\main.go")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()*/

	//Ejecuta y saca la salida segun lo imprimamos
	/*c := exec.Command("cmd", "/C", "go run c:\\gopath\\src\\github.com\\jantome\\prueba\\main.go")
	cout, err := c.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(cout))*/

	fmt.Println("Comienza cumplehoras...")
	//recuperamos el entorno de ejecucion mediante flags para saber las rutas
	//entorno := flag.String("entorno", "", "entorno de ejecución")
	flag.Parse()
	//variables inizialidas
	comando := ""
	consola := ""
	letra := ""
	//cargamos variables de entorno
	environment.Loadenvironment(*entorno)
	//Leer fichero
	pathfile, _ := os.LookupEnv("FILE_HOUR")
	//al nombre del fichero le añadimos el día formato cm
	pathfile = pathfile + fechacm
	file, err := os.Open(pathfile)
	if err != nil {
		fmt.Println("Error lectura fichero horas", err.Error())
		os.Exit(1)
		return
	}
	//cerramos el fichero con defer, para que sea al final
	defer file.Close()
	//recuperamos la hora actual
	hora := time.Now().Format("15:04")
	//Sw para entrar al segundo bucle
	t := false
	//realizamos un scanner del fichero
	scanner := bufio.NewScanner(file)
	sql := ""
	//montamos un bucle por cada linea leida hasta el final del fichero
	for scanner.Scan() {
		t = false
		//segundo bucle donde se comprobara la hora
		for t == false {
			//recuperamos la hora actual lo que esta entre "" es para indicar el formato que necesitamos
			hora = time.Now().Format("15:04")
			//comparamos la hora actual con la del fichero
			if hora == scanner.Text() {
				//buscamos el nombre del que tiene la condicion de entrada de la hora
				sql = fmt.Sprintf("SELECT nombre FROM ejecucion WHERE condicionin = '%s' AND fechaeje ='%s'", scanner.Text(), fechacm)
				result, err := db2.EjecutaQuery(sql)
				if err != nil {
					fmt.Println("Error select nombre", err.Error())
					os.Exit(1)
					return
				}
				//solo tendra un registro por lo que no montamos bucle de lectura
				result.Next()
				//creamos variable de aplantillamento de la lectura
				var ejecucion structs.Ejecucion
				err = result.Scan(&ejecucion.Nombre)

				if err != nil {
					fmt.Println("Error lectura select horas", err.Error())
					os.Exit(1)
					return
				}
				//una vez tenemos el nombre, realizamos una query para poner todos los estados a Ok, de la hora leida en el fichero
				sql = fmt.Sprintf("UPDATE ejecucion SET estado ='ok' WHERE nombre = '%s' AND fechaeje ='%s'", ejecucion.Nombre, fechacm)
				_, err = db2.EjecutaQuery(sql)

				if err != nil {
					fmt.Println("Error update estado", err.Error())
					os.Exit(1)
					return
				}
				//Buscamos la condición de salida que tiene que dejar la hora leida en el fichero
				sql = fmt.Sprintf("SELECT condicionout FROM ejecucion WHERE nombre ='%s' and condicionout > ''", ejecucion.Nombre)
				result, err = db2.EjecutaQuery(sql)

				if err != nil {
					fmt.Println("Error select condicionout", err.Error())
					os.Exit(1)
					return
				}
				//solo tendra un resultado por lo que no creamos bucle de lectura
				result.Next()
				//aplantillamos
				var ejecucion2 structs.Ejecucion2
				err = result.Scan(&ejecucion2.Condicionout)
				if err != nil {
					fmt.Println("Error lectura  select condicionout", err.Error())
					os.Exit(1)
					return
				}
				//Actualizamos todas las lineas de la condicion de entrada que estan esperando los jobs (se actualiza con la condición de salida de la hora)
				sql = fmt.Sprintf("UPDATE ejecucion SET estado ='ok' WHERE condicionin = '%s' AND fechaeje ='%s'", ejecucion2.Condicionout, fechacm)
				_, err = db2.EjecutaQuery(sql)

				if err != nil {
					fmt.Println("Error update estado condicionin", err.Error())
					os.Exit(1)
					return
				}
				//ejecutar ejecutajob, para que se lancen los job's con las condiciones de entradas cumplidas
				//ademas lo hacemos de la manera que no para la ejecucion en caso de que falle
				if *entorno == "local" {
					comando = "go run c:\\gopath\\src\\github.com\\jantome\\ejecutajob\\ejecutajob.go -entorno=local"
					consola = "cmd"
					letra = "/C"
				} else {
					comando = "cd /ejecutable/batch/app; ./ejecutajob "
					consola = "bash"
					letra = "-c"
				}
				//de esta manera no casca en caso de que falle el ejecutajob, y seguira ejecutandose el cumple horas
				c := exec.Command(consola, letra, comando)
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
				t = true
			}
			// si la hora actual es menor que la del fichero, tenemos que seguir con esta, ya que significa que todavía
			//no llego la hora
			if hora < scanner.Text() {
				//esperamos un minuto antes de hacer la siguiente comprobación
				time.Sleep(60 * time.Second)
			}
			// si la hora actual es mayor, significa que ya paso la hora, por lo que leemos el siguiente registro
			if hora > scanner.Text() {
				//activamos sw para salir del segundo bucle y leer el siguiente registro
				t = true
			}
		}
	}
	fmt.Println("Fin cumplehoras...")
}
