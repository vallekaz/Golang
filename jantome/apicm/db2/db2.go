package db2

import (
	//para darle formato a la variable
	"fmt"

	//Librerías para base de datos
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//EjecutaQuery funcion para que con el string de la query ejecute todo conexion query y desconexión
func EjecutaQuery(query string) (result *sql.Rows, e error) {
	//se conectara a la base de datos
	db, err := connectDb2()

	if err != nil {
		return nil, err
	}

	//con el string que tenemos de entrada hacemos la query
	result, err = realizaQuery(query, db)

	if err != nil {
		return nil, err
	}
	// cerramos con defer para que no se nos olvide
	defer db.Close()

	return result, nil

}

//ConnectDb2 FUNCION PARA LA CONEXIÓN CON LA BASE DE DATOS (Devuelve un DB y error)
func connectDb2() (db *sql.DB, e error) {
	//Datos para la conexión
	usuario := "root"
	pass := ""
	host := "tcp(127.0.0.1:3306)"
	nombreBaseDeDatos := "cm"
	//el formato tiene que ser  sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/batch")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
	//si es error devolvemos el error
	if err != nil {
		return nil, err
	}
	return db, nil
}

//RealizaQuery FUNCION PARA REALIZAR LA QUERY (De entrada tendra la query y el db, de salida las filas y el error)
func realizaQuery(query string, db *sql.DB) (result *sql.Rows, e error) {
	//result, err := db.Query("SELECT * FROM batch_calendario")
	result, err := db.Query(query)
	//controlamos error
	if err != nil {
		return nil, err
	}
	return result, nil
}
