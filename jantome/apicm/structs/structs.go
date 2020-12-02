package structs

//Tabejecucion estructura de la tabla
type Tabejecucion struct {
	Nombre   string `json:"nombre"`
	Fechaeje string `json:"fechaEje"`
	Estado   string `json:"estado"`
}

//Ejecucionjson json de salida con la información de la tabla de ejecución
type Ejecucionjson struct {
	Name     string `json:"name"`
	Fechaeje string `json:"fechaEje"`
	Estado   string `json:"estado"`
	Links    struct {
		Href map[string]string `json:"href,omitempty"`
	} `json:"links,omitempty"`
}

//Pieejecucion pie de ejecucion que contendra la información necesaria para la paginacion
type Pieejecucion struct {
	Content []Ejecucionjson `json:"data"`
	Pagdet  struct {
		Links struct {
			Href map[string]string `json:"href,omitempty"`
		}
		Pagmax  int64 `json:"pagcount,omitempty"`
		Pagnext int64 `json:"pagnext,omitempty"`
		Pagprev int64 `json:"pagprev,omitempty"`
	} `json:"pagination_details"`
}

//Ejecucioncount structura para el select count
type Ejecucioncount struct {
	Count int64
}

//Estadoejecucion lectura del estado de la tabla de ejecucion
type Estadoejecucion struct {
	Estado string
}

//Jsonerror estructura del json de error
type Jsonerror struct {
	UserMessage     string `json:"userMessage"`
	InternalMessage string `json:"internalMessage,omitempty"`
}

//Condicionin struct para el json con las condiciones de entrada
type Condicionin struct {
	Condicionin string `json:"condicionin"`
}

//Condicionout struct para el json con las condiciones de entrada
type Condicionout struct {
	Condicionout string `json:"condicionout"`
}

//Tabplanificacion struct para la lectura de la tabla y para mostrar en el json
//La lectura de la tabla tambien sirve como Json de salida. Ya que no tendra link's de salida, por lo que no creamos
//ningun struct más como en el caso de ejecucion
type Tabplanificacion struct {
	Nombre     string `json:"name"`
	Calendario string `json:"calendar,omitempty"`
	Useralta   string `json:"useralt,omitempty"`
	Timalta    string `json:"timalta,omitempty"`
	Usermod    string `json:"usermod,omitempty"`
	Timesmod   string `json:"timesmod,omitempty"`
}

//Pieplanificacion pie de ejecucion que contendra la información necesaria para la paginacion
type Pieplanificacion struct {
	Content []Tabplanificacion `json:"data"`
	Pagdet  struct {
		Links struct {
			Href map[string]string `json:"href,omitempty"`
		}
		Pagmax  int64 `json:"pagcount,omitempty"`
		Pagnext int64 `json:"pagnext,omitempty"`
		Pagprev int64 `json:"pagprev,omitempty"`
	} `json:"pagination_details"`
}

//Planificacioncount struct para el json con el número de planificacion
type Planificacioncount struct {
	Count int64 `json:"condicionout"`
}

//Putplanificacion update de la tabla planificacion (se actualizara calendario y usuermodifi)
type Putplanificacion struct {
	Name     string `json:"name"`
	Calendar string `json:"calendar"`
	Usermod  string `json:"usermod"`
}

//Postplanificacion alta en la tabla de planificacion
type Postplanificacion struct {
	Name     string `json:"name"`
	Calendar string `json:"calendar"`
	Useralt  string `json:"useralt"`
}

//Postcondicionin alta en la tabla de planificacion de las condiciones de entrada
type Postcondicionin struct {
	Name        string `json:"name"`
	Condicionin string `json:"condicionin"`
	Useralt     string `json:"useralt"`
}

//Calendarplanificacion para recuperar el calendario de la tabla de planificacion
type Calendarplanificacion struct {
	Calendario string
}

//Calendar struct de la tabla y para sacar el json con el listado de calendarios
type Calendar struct {
	Name string `json:"name"`
}

//Postcondicionout alta en la tabla de planificacion de las condiciones de entrada
type Postcondicionout struct {
	Name         string `json:"name"`
	Condicionout string `json:"condicionout"`
	Useralt      string `json:"useralt"`
}
