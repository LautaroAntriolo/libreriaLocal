package rutas

import (
	"ProyectoWEB/conectar"
	"ProyectoWEB/utilidades"
	"html/template"
	"net/http"
)

func LibroImagen(w http.ResponseWriter, r *http.Request) {

	conectar.Conectar()
	defer conectar.CerrarConexion()

	// Cargar la plantilla HTML
	tmpl := template.Must(template.ParseFiles("templates/libreria/cargar_IMG.html",
		"templates/principales/navbarLogin.html",
		"templates/principales/navbarNotLogin.html",
		utilidades.Frontend))

	// Pasar datos a la plantilla
	data := DatosDeSeguridad(w, r)
	
	tmpl.ExecuteTemplate(w, "cargar_IMG.html", data)
}