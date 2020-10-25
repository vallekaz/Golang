package authentication

import (
	//Librería de logs

	"github.com/jantome/apilogin/logs"

	//Liberia para las llaves RSA
	"crypto/rsa"
	"os"
	"strconv"

	//Librería para la lectura de los ficheores
	"io/ioutil"
	"time"

	//Libreria para el jwt
	"github.com/dgrijalva/jwt-go-master"
	//Librería para los structs
	"github.com/jantome/apilogin/structs"
)

/* Creación del token
creamos con clave privada, y consultamos con clave publica */

//Creamos las variables para las claves
var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

var (
	erroLog       string = ""
	terrorwarning string = "w"
	terrorinfo    string = "i"
	terrorerro    string = "e"
)

//GenerateJWT recibe de entrada el usuario devuelve string con el token o error si tuviese
func GenerateJWT(userLogon string, rol string) (string, error) {
	//Leemos los archivos de claves
	//Lectura bytes de la clave privada ./ donde esta el archivo main
	privateBytes, err := ioutil.ReadFile("./keys/private.rsa")
	if err != nil {
		logs.GrabaLog(err, "Open Private Key", terrorerro)
		return "", err
		//log.Fatal("NO se puede leer el archivo privado")
	}

	//Convertimos los bytes en las llaves
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateBytes)

	if err != nil {
		logs.GrabaLog(err, "Conver Private Key Byte", terrorerro)
		return "", err
		//log.Fatal("No puede hacer el parse de private")
	}
	//de   las variables de entorno, cogemos el tiempo que durara el token
	duracion, _ := os.LookupEnv("TOKEN_LIFE")
	//lo pasamos numerico
	duracionnum, _ := strconv.Atoi(duracion)
	//lo convertimos en time
	prueba := time.Duration(duracionnum)
	//Creamos una variable Claims con el contenido del payload del token
	claims := structs.Claims{
		//usuario de entrada
		User: userLogon,
		//Rol del claim
		Rol: rol,
		//standarclaims estructura de jwt
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * prueba).Unix(),
			Issuer:    "Login User",
		},
	}
	//Convertimos en token
	//NewWithClaims (metodo de firma, claims que queremos que convierta en payload)
	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	//Lo convertimos a base64
	//convertimos a string firmando con la llave privada
	result, err := token.SignedString(privateKey)

	if err != nil {
		logs.GrabaLog(err, "Token signing failed", terrorerro)
		//log.Fatal("Error al firmar el token")
		return "", err
	}

	return result, nil
}

//CompruebaToken de si el token es valido
func CompruebaToken(reqToken string) (structs.Claims, string, error) {
	//creamos una estructura como la tabla para devolver el rol
	var contenido structs.Claims
	descError := ""
	//lectura bytes de publica la hacemos cuando comprobamos el token, ya que para generarlo hace falta la clave privada
	publicBytes, err := ioutil.ReadFile("./keys/public.rsa.pub")
	//comprobamos el error
	if err != nil {
		logs.GrabaLog(err, "Open public Key", terrorerro)
		return contenido, "", err
		//log.Fatal("NO se puede leer el archivo publico")
	}
	//Convertimos los bytes en las llaves
	//publicKey, err := jwt.ParseECPublicKeyFromPEM(publicBytes)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	//fmt.Println(publicKey)

	if err != nil {
		logs.GrabaLog(err, "Read public Key", terrorerro)
		return contenido, "", err
		//log.Fatal("NO se puede leer el archivo publico")
	}

	//comprobamos si es correcto
	//fmt.Println(reqToken)

	//Comprobamos el Token mediante ParseWithClaims
	token, err := jwt.ParseWithClaims(reqToken, &structs.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil

	})
	//para formatear el token si es correcto y sacar datos
	// prueba := token.Claims.(*structs.Claims)
	//comprobamos si existe error.
	if err != nil {
		//comprobamos el tipo de erro que devuelve
		switch err.(type) {
		//error de validacion
		case *jwt.ValidationError:
			//nos guardamos el tipo de error de validacion
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			//comprobamos token expirado
			case jwt.ValidationErrorExpired:
				descError = "Expired Token"
			case jwt.ValidationErrorSignatureInvalid:
				descError = "Wrong Signature"
			case jwt.ValidationErrorClaimsInvalid:
				descError = "Wrong Claims"
			default:
				//descError = fmt.Sprintln("Wrong token %v", err.Error())
				descError = "Wrong Token"
			}
		//si no es un error de validacion
		default:
			descError = "Wrong Token"

		}
	}
	//Si el error esta vacio
	if descError == "" {
		//aplantillamos el token
		contenidoJSON := token.Claims.(*structs.Claims)
		//informamos el user con la decodificación del token
		contenido.User = contenidoJSON.User
		//informamos el rol con la decodificación del token
		contenido.Rol = contenidoJSON.Rol
	}
	return contenido, descError, nil
}
