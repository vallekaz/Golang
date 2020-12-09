package main

import (
	"testing"

	"github.com/jantome/planificacion/environment"
)

//Test limpieza de la tabla
func Test_limpia(t *testing.T) {
	environment.Loadenvironment("local")
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "Limpia CM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limpia()
		})
	}
}

func Test_planificar(t *testing.T) {
	environment.Loadenvironment("local")
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "Planificacion",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planificar()
		})
	}
}

func Test_creahoras(t *testing.T) {
	environment.Loadenvironment("local")
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "Genera horas",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creahoras()
		})
	}
}
