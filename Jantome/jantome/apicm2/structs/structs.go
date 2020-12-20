package structs

//Calendario tabla // Ejemplo para emitir empty Nombre        string `json:"nombre,omitempty"`
//lo que esta dentro de la etiqueta Json es como saldra en los JSON de entrada/salida
type Calendario struct {
	Nombre        string `json:"nombre"`
	Mes           string `json:"meses"`
	Dia           string `json:"dia"`
	Usuarioalta   string `json:"usuarioalta"`
	Fechacreacion string `json:"fechacreacion"`
	Usuariomodif  string `json:"usuariomodificacion"`
	Fechamodif    string `json:"fechamodificacion"`
}

//EnvJSON para montar Json para las variables de entorno
type EnvJSON struct {
	Dbhost      string `json:"dbhost"`
	Dbuser      string `json:"dbuser"`
	Dbpassword  string `json:"dbpassword"`
	Dbdatabase  string `json:"dbdatabase"`
	Servport    string `json:"servport"`
	Logactivate string `json:"logactivate"`
}
