package structs

//Ejecucion estructura de la tabla
type Ejecucion struct {
	Nombre   string `json:"nombre"`
	Fechaeje string `json:"fechaEje"`
	Estado   string `json:"estado"`
}

//Ejecucionjson estructura de la tabla
type Ejecucionjson struct {
	Nombre   string `json:"nombre"`
	Fechaeje string `json:"fechaEje"`
	Estado   string `json:"estado"`
	Links    struct {
		Href map[string]string `json:"href,omitempty"`
	} `json:"links,omitempty"`
}

//Jsonerror estructura del json de error
type Jsonerror struct {
	UserMessage     string `json:"userMessage"`
	InternalMessage string `json:"internalMessage,omitempty"`
}

//Ejecucioncount structura para el select count
type Ejecucioncount struct {
	Count int64
}

//Piepaginacion pie de ejecucion que contendra la información necesaria para la paginacion
type Piepaginacion struct {
	Pagdet struct {
		Links struct {
			Href map[string]string `json:"href,omitempty"`
		}
		Pagmax  int64 `json:"pagcount,omitempty"`
		Pagnext int64 `json:"pagnext,omitempty"`
		Pagprev int64 `json:"pagprev,omitempty"`
	} `json:"pagination_details"`
}

//Condicionin struct para el json con las condiciones de entrada
type Condicionin struct {
	Condicionin string `json:"condicionin"`
}

//Condicionout struct para el json con las condiciones de entrada
type Condicionout struct {
	Condicionout string `json:"condicionout"`
}

//Planificacion struct para la lectura de la tabla y para mostrar en el json
type Planificacion struct {
	Nombre     string `json:"name"`
	Calendario string `json:"calendar,omitempty"`
	Useralta   string `json:"useralt,omitempty"`
	Timalta    string `json:"timalta,omitempty"`
	Usermod    string `json:"usermod,omitempty"`
	Timesmod   string `json:"timesmod,omitempty"`
}

//Planificacionout struct para el json con el número de planificacion
type Planificacionout struct {
	Count int64 `json:"condicionout"`
}
