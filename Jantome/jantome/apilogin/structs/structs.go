package structs

import (
	//para los claims (token) y le ponemos un alias
	"github.com/dgrijalva/jwt-go-master"
)

//UserLogon para montar Json para las variables de entorno
/* Ya no se usa, ya que ponemos autorización basica no tenemos json de entrada
type UserLogon struct {
	Usuario  string `json:"user"`
	Password string `json:"password"`
}*/

//BatchUser estructura de la tabla
type BatchUser struct {
	Usuario  string `json:"user"`
	Password string `json:"password,omitempty"`
	Rol      string `json:"rol"`
}

//Claims Información dentro del paylod para el JWT
type Claims struct {
	User string `json:"user"`
	Rol  string `json:"rol"`
	//claims standar
	jwt.StandardClaims
}

//ResponseToken Estructura no obligatoria pero sirver para mostrar como json el tonken
type ResponseToken struct {
	Token string `json:"token"`
}

//EnvJSON con las variables de sistem para mostrar
type EnvJSON struct {
	Dbhost      string `json:"dbhost"`
	Dbuser      string `json:"dbuser"`
	Dbpassword  string `json:"dbpassword"`
	Dbdatabase  string `json:"dbdatabase"`
	Servport    string `json:"servport"`
	ServportSSL string `json:"servportssl"`
	Sersafe     string `json:"sersafe"`
	Logactivate string `json:"logactivate"`
	TokenLife   string `json:"tokenlife"`
	PathLog     string `json:"pathlog"`
}

//TestJSON para ver el estado de la api
type TestJSON struct {
	Status   string `json:"status"`
	ConexDb2 string `json:"conexdb2"`
}
