package structs

//Jsoncm para el json de lectura del CM
type Jsoncm struct {
	Accion   string `json:"accion"`
	Programa string `json:"programa"`
	Fechaeje string `json:"fechaeje"`
}

//Ejecucion struct de la tabla de ejecucion para recuperar el nombre
type Ejecucion struct {
	Condicionout string
}

//Ejecucion2 struct de la tabla de ejecución para recuperar el estado
type Ejecucion2 struct {
	Estado string
}
