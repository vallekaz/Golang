package main

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/jantome/apilogin/authentication"
	"github.com/jantome/apilogin/environment"
	"github.com/jantome/apilogin/structs"

	"github.com/jantome/apilogin/db2"

	"github.com/jantome/apilogin/logs"
)

//Test para comprobación de grabación de log's
func TestGrabaLog(t *testing.T) {
	type args struct {
		err2        error
		descripcion string
		tipo        string
	}
	tests := []struct {
		name string
		args args
	}{
		//Grabación Warning
		{
			name: "Warning",
			args: args{
				err2:        nil,
				descripcion: "Prueba Warning",
				tipo:        "w",
			},
		},
		//Grabación informational
		{
			name: "Informational",
			args: args{
				err2:        nil,
				descripcion: "Prueba Informational",
				tipo:        "i",
			},
		},
		//Grabación  Error
		{
			name: "Error",
			args: args{
				err2:        nil,
				descripcion: "Prueba Error",
				tipo:        "e",
			},
		},
		//No graba pero finaliza OK
		{
			name: "Other",
			args: args{
				err2:        nil,
				descripcion: "Prueba other",
				tipo:        "z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs.GrabaLog(tt.args.err2, tt.args.descripcion, tt.args.tipo)
		})
	}
}

//Test para operaciones db2
func TestEjecutaQuery(t *testing.T) {
	//cargamos variables
	environment.Loadenvironment()
	type args struct {
		query string
	}
	tests := []struct {
		name       string
		args       args
		wantResult *sql.Rows
		wantErr    bool
	}{
		{
			name: "Conexión Db2",
			args: args{
				query: "SELECT * FROM batch_usuarios",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := db2.EjecutaQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("EjecutaQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//Si es igual a nil el resultado (ya que no lo hemos informado es cuando debe de fallar)
			if reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("EjecutaQuery() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

//Test carga de variables de entorno
func TestLoadenvironment(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Carga Variables de Sistema",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment.Loadenvironment()
		})
	}
}

//Test para generar token
func TestGenerateJWT(t *testing.T) {
	//cargamos variables de entorno ya que necesitamos el tiempo de vida del token
	environment.Loadenvironment()
	type args struct {
		userLogon string
		rol       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Generación Token",
			args: args{
				userLogon: "Prueba",
				rol:       "admin",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authentication.GenerateJWT(tt.args.userLogon, tt.args.rol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//Si es igual es que no genero el token, por lo que se falla
			if got == tt.want {
				t.Errorf("GenerateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

//Test consulta Token
func TestCompruebaToken(t *testing.T) {
	type args struct {
		reqToken string
	}
	tests := []struct {
		name    string
		args    args
		want    structs.Claims
		want1   string
		wantErr bool
	}{
		// Caso de test Token Expired
		{
			name: "Comprobación de Token",
			args: args{
				reqToken: "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoicHJ1ZWJhIiwicm9sIjoiYWFhIiwiZXhwIjoxNTk3MzU0NDcxLCJpc3MiOiJMb2dpbiBVc2VyIn0.BRZAWUx8wvOV5_p_fFjUAnKRt31DQRbOaMAYoXQmZ06vQYKZNLzcdvteTh-JvK4JsyDG_Jefemf4HkbymoqjvINihKoDLNr9RHOpJ6GPhdq7QaBvz0q7AihbArHrZdTf7fbjyKknUetSQRIaEnwQuiHxfLEftD_rS2dgChay4eOnr0N04e7F2L_2W_NpRRJa91_L_528zyZ8JU32dLRjUjO1K4QSB5-6MOsuTu_gtFAF2ZUazleokNWjTp9Kd566AEp_veGeR333B43tscaG6St0vz1jby7xgiK81QvJ8Gv_4wPcZnwC0Wib8m5cSXvE3MWN_phoudDEC2zeR2vVyQ",
			},
			want1: "Expired Token",
		},
		//Caso de test default wrong token
		{
			name: "Claims Invalid",
			args: args{
				reqToken: "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.e4J1c2VyIjoicHJ1ZWJhIiwicm9sIjoiYWFhIiwiZXhwIjoxNTk3MzU0NDcxLCJpc3MiOiJMb2dpbiBVc2VyIn0.BRZAWUx8wvOV5_p_fFjUAnKRt31DQRbOaMAYoXQmZ06vQYKZNLzcdvteTh-JvK4JsyDG_Jefemf4HkbymoqjvINihKoDLNr9RHOpJ6GPhdq7QaBvz0q7AihbArHrZdTf7fbjyKknUetSQRIaEnwQuiHxfLEftD_rS2dgChay4eOnr0N04e7F2L_2W_NpRRJa91_L_528zyZ8JU32dLRjUjO1K4QSB5-6MOsuTu_gtFAF2ZUazleokNWjTp9Kd566AEp_veGeR333B43tscaG6St0vz1jby7xgiK81QvJ8Gv_4wPcZnwC0Wib8m5cSXvE3MWN_phoudDEC2zeR2vVyQ",
			},
			want1: "Wrong Token",
		},
		//Caso de tes signature
		{
			name: "Claims Invalid",
			args: args{
				reqToken: "eyJhbGciOiJQUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoicHJ1ZWJhIiwicm9sIjoiYWFhIiwiZXhwIjoxNTk3MzU0NDcxLCJpc3MiOiJMb2dpbiBVc2VyIn0.BRZAWUx8wvOV5_p_fFjUAnKRt31DQRbOaMAYoXQmZ06vQYKZNLzcdvteTh-JvK4JsyDG_Jefemf4HkbymoqjvINihKoDLNr9RHOpJ6GPhdq7QaBvz0q7AihbArHrZdTf7fbjyKknUetSQRIaEnwQuiHxfLEftD_rS2dgChay4eOnr0N04e7F2L_2W_NpRRJa91_L_528zyZ8JU32dLRjUjO1K4QSB5-6MOsuTu_gtFAF2ZUazleokNWjTp9Kd566AEp_veGeR333B43tscaG6St0vz1jby7xgiK81QvJ8Gv_4wPcZnwC0Wib8m5cSXvE3MWN_phoudDEC2zeR2vVy3",
			},
			want1: "Wrong Token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := authentication.CompruebaToken(tt.args.reqToken)
			//Si da error es en la apertura
			if (err != nil) != tt.wantErr {
				t.Errorf("CompruebaToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompruebaToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CompruebaToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
