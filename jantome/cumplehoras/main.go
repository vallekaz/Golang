package main

import (
	"fmt"
	"time"
)

var (
	//recuperamos la fecha de ejecucion y la formateamos separada
	date = time.Now()
	anno = fmt.Sprintf("%d", date.Year())
	mes  = fmt.Sprintf("%d", date.Month())
	//Convertido a String
	dia = fmt.Sprintf("%d", date.Day())
	//montamos la fecha de controlm formato YYMMDD
	fechacm = fmt.Sprintf("%s%s%s", anno[0:2], mes, dia)
)

func main() {

}
