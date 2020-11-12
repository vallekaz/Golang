package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/jantome/mqarq/environment"

	"github.com/jantome/mqarq/db2"

	"github.com/jantome/mqarq/structs"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func holdea(programa string, fechaeje string) {
	//comprobamos que esta informado el nombre del programa
	if programa != "" {
		if fechaeje != "" {
			//antes de holdear comprobamos si no esta OK, en caso de estarlo no podemos holdear
			sql := fmt.Sprintf("SELECT estado WHERE numsec = 1 AND nombre = '%s' AND fechaeje = '%s'", programa, fechaeje)
			result, err := db2.EjecutaQuery(sql)
			if err != nil {
				log.Println("Error select estado", err.Error())
			}
			//solo tendra una linea, por lo que no es necesario el bucle
			result.Next()
			var ejecucion2 structs.Ejecucion2
			//aplantillamos
			err = result.Scan(&ejecucion2.Estado)
			//comprobamos el estado para ver si podemos holdea
			if ejecucion2.Estado != "ok" {
				//Actualizamos el estado a ho en la tabla ejecucion y con numsec 1(ya que las condiciones no las queremos tocar)
				sql = fmt.Sprintf("UPDATE ejecucion SET estado = 'ho' WHERE numsec = 1 AND nombre = '%s' AND fechaeje = '%s'", programa, fechaeje)
				_, err = db2.EjecutaQuery(sql)
				if err != nil {
					log.Println("Error Update estado", err.Error())
				}
			}

		} else {
			log.Println("nombre de programa no informado")
		}
	} else {
		log.Println("nombre de programa no informado")
	}
}

func free(programa string, fechaeje string) {
	comando := ""
	consola := ""
	letra := ""
	//comprobamos que esta informado el nombre del programa
	if programa != "" {
		if fechaeje != "" {
			//Actualizamos el estado a pl en la tabla ejecucion y con numsec 1(ya que las condiciones no las queremos tocar)
			//y ejecutamos ejecutajobs
			sql := fmt.Sprintf("UPDATE ejecucion SET estado = 'pl' WHERE numsec = 1 AND nombre = '%s' AND fechaeje = '%s'", programa, fechaeje)
			_, err := db2.EjecutaQuery(sql)
			if err != nil {
				log.Println("Error Update estado", err.Error())
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
			//de esta manera no casca en caso de que falle el ejecutajob, y seguira ejecutandose el programa de MQ
			c := exec.Command(consola, letra, comando)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Run()
		} else {
			log.Println("nombre de programa no informado")
		}
	} else {
		log.Println("nombre de programa no informado")
	}
}

func rerun(programa string, fechaeje string) {
	comando := ""
	consola := ""
	letra := ""
	//comprobamos que esta informado el nombre del programa
	if programa != "" {
		if fechaeje != "" {
			//poner el job que se pide hacer rerun como eje
			sql := fmt.Sprintf("UPDATE ejecucion SET estado = 'ej' WHERE nombre = '%s' AND numsec = 1 AND fechaeje = '%s'", programa, fechaeje)
			_, err := db2.EjecutaQuery(sql)
			if err != nil {
				log.Println("Update Ko error")
			}
			//ejecutar el job indicado controlando el error
			if *entorno == "local" {
				comando = "go run c:\\gopath\\src\\github.com\\jantome\\" + programa + "\\" + programa + ".go"
				consola = "cmd"
				letra = "/C"
			} else {
				//ruta linux
				comando = "cd /ejecutable/batch/app; ./" + programa
				consola = "bash"
				letra = "-c"
			}
			c := exec.Command(consola, letra, comando)
			_, err = c.Output()

			//controlamos el error, y si falla actualizamos la ejecucion como fallida
			if err != nil {
				//Este print le dejo por el momento, pero sera el job el que deje un log de salida
				fmt.Println(err.Error())
				//actualizamos a fallido la ejecucion
				sql = fmt.Sprintf("UPDATE ejecucion SET estado = 'ko' WHERE nombre = '%s' AND numsec = 1 AND fechaeje = '%s'", programa, fechaeje)
				_, err := db2.EjecutaQuery(sql)
				if err != nil {
					log.Println("Update Ko error")
				}
			} else {
				//em caso de que finalice OK actualizmaos
				sql = fmt.Sprintf("UPDATE ejecucion SET estado= 'ok' WHERE nombre = '%s' AND fechaeje = '%s'", programa, fechaeje)
				_, err = db2.EjecutaQuery(sql)
				//controlamos error
				if err != nil {
					log.Println("Update OK error")
				} else {
					//Finaliza Ok el update y actualizamos las condiciones de salida
					sql = fmt.Sprintf("SELECT condicionout FROM ejecucion WHERE nombre ='%s' AND condicionout > '' AND fechaeje = '%s' ", programa, fechaeje)
					result, err := db2.EjecutaQuery(sql)
					if err != nil {
						log.Println("Error Select condicionout")
					}
					//leemos
					var ejecucion structs.Ejecucion
					for result.Next() {
						err = result.Scan(&ejecucion.Condicionout)
						//realizmos update de la condicion de salida en toda la tabla
						sql2 := fmt.Sprintf("UPDATE ejecucion SET estado='ok' WHERE condicionin = '%s' AND fechaeje = '%s'", ejecucion.Condicionout, fechaeje)
						_, err = db2.EjecutaQuery(sql2)
						if err != nil {
							log.Println("Error update condicionin")
						}
					}
				}
			}
		} else {
			log.Println("fechaeje no informada")
		}
	} else {
		log.Println("nombre de programa no informado")
	}
}

func petic(programa string, fechaeje string) {
	comando := ""
	consola := ""
	letra := ""
	//comprobamos que esta informado el nombre del programa
	if programa != "" {
		if fechaeje != "" {
			//Realizmos un insert con el nombre del programa y la fecha
			sql := fmt.Sprintf("INSERT INTO ejecucion VALUES('%s', 1, '%s', '','','pl')", programa, fechaeje)
			log.Println(sql)
			_, err := db2.EjecutaQuery(sql)
			if err != nil {
				log.Println("Error Inser Petic", err.Error())
			}
			//para saber que es a petición cremos una condicion de peticion
			sql = fmt.Sprintf("INSERT INTO ejecucion VALUES('%s', 2, '%s', 'petic-%s','','ok')", programa, fechaeje, programa)
			_, err = db2.EjecutaQuery(sql)
			if err != nil {
				log.Println("Error Inser Petci", err.Error())
			} else {
				//Una vez planificado, ejecutamos el ejecutajob, que ejecutara todo lo que tenga la condicion cumplida como es el caso
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
				//de esta manera no casca en caso de que falle el ejecutajob, y seguira ejecutandose el job
				c := exec.Command(consola, letra, comando)
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Run()
			}
		} else {
			log.Println("nombre de programa no informado")
		}
	} else {
		log.Println("nombre de programa no informado")
	}
}

var (
	//recuperamos el entorno de ejecucion mediante flags para saber las rutas
	entorno = flag.String("entorno", "", "entorno de ejecución")
)

func main() {
	flag.Parse()
	//cargamos variables de entorno
	environment.Loadenvironment(*entorno)
	parrabbitmq, _ := os.LookupEnv("RABBITMQ")
	conn, err := amqp.Dial(parrabbitmq)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cmqueue", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	//Creamos las variables de aplantillamiento de CM
	var jsoncm structs.Jsoncm

	go func() {
		for d := range msgs {
			//log.Printf("Received a message: %s", d.Body)
			//Nos creamos el buffer
			buf := bytes.NewBuffer(d.Body)
			//Lo ponermos en el decodificador json
			decoder := json.NewDecoder(buf)
			//Decodificamos aplantillando
			err := decoder.Decode(&jsoncm)
			//OJOO controlar el error (guardando en cola mq o algo cuando no decodifique)
			if err != nil {
				log.Print("Error conversion json ", err)
				log.Printf("Received a message: %s", d.Body)
			}
			//realizamos la acción que hemos recibido en el json
			//holdear
			if jsoncm.Accion == "hold" {
				//llamamos a la funcion para realizar el hold
				holdea(jsoncm.Programa, jsoncm.Fechaeje)
			}
			//liberar
			if jsoncm.Accion == "free" {
				free(jsoncm.Programa, jsoncm.Fechaeje)
			}
			//rerun
			if jsoncm.Accion == "rerun" {
				rerun(jsoncm.Programa, jsoncm.Fechaeje)
			}
			//ejecucion a peticion
			if jsoncm.Accion == "peticion" {
				petic(jsoncm.Programa, jsoncm.Fechaeje)
			}
			//bypaas (Saltarse las condiciones de entrada) es lo mismo que rerun (ejecuta un job concreto)
			if jsoncm.Accion == "bypass" {
				rerun(jsoncm.Programa, jsoncm.Fechaeje)
			}
			//Borrar el conetenido para que no salga en unacked (tratados en error)
			contenido := bytes.Count(d.Body, []byte("."))
			t := time.Duration(contenido)
			time.Sleep(t * time.Second)
			//log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
