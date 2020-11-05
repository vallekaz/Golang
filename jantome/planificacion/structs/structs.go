package structs

//Planificacion struct de la tabla de planificacions
type Planificacion struct {
	Numsec       int
	Nombre       string
	Ejecucion    string
	Condicionin  string
	Condicionout string
	Calendario   string
	Useralta     string
	Timalta      string
	UserModif    string
	Timesmod     string
}

//Ejecucion struct de la tabla de ejecuciones, para hacer select max
type Ejecucion struct {
	Numsec int
}

//Calendario struct para obtener el nombre del calendario
type Calendario struct {
	Namecalendar string
}
