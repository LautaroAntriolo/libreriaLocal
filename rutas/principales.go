package rutas

import (
	"ProyectoWEB/conectar"
	// "ProyectoWEB/modelos"
	"ProyectoWEB/utilidades"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type AlertData struct {
	CSS     string
	Mensaje string
}

func Home(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/principales/home.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))
	data := DatosDeSeguridad(response, request)
	template.ExecuteTemplate(response, "home.html", data)
}

func Nosotros(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.New("").Funcs(template.FuncMap{
		"formatDate": func(dateString string) string {
			if dateString == "" {
				return ""
			}
			t, err := time.Parse("2006-01-02 15:04:05", dateString)
			if err != nil {
				return dateString
			}
			return t.Format("02 Jan 2006")
		},
	}).ParseFiles(
		"templates/principales/Nosotros.html",
		"templates/principales/navbarLogin.html",
		"templates/principales/navbarNotLogin.html",
		utilidades.Frontend,
	))

	data := DatosDeSeguridad(response, request)

	conectar.Conectar()
	defer conectar.CerrarConexion()

	sql := `
      SELECT 
    tl.nombre, 
    tl.portada, 
    tl.autor, 
    tl.puntaje, 
    ul.nombre AS usuario_nombre, 
    tl.fecha_creacion
FROM 
    todosloslibros tl
LEFT JOIN 
    usuarios_library ul ON tl.UsuarioId = ul.id
ORDER BY 
    tl.fecha_creacion DESC
LIMIT 4;    

    `

	var libros []struct {
		Nombre        string
		Url_imagen    string
		Autor         string
		Puntaje       int
		UsuarioNombre string
		FechaCreacion string
	}

	datos, err := conectar.Db.Query(sql)
	if err != nil {
		http.Error(response, "Error en la consulta: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer datos.Close()

	for datos.Next() {
		var libro struct {
			Nombre        string
			Url_imagen    string
			Autor         string
			Puntaje       int
			UsuarioNombre string
			FechaCreacion string
		}

		var fechaBytes []byte
		err := datos.Scan(
			&libro.Nombre,
			&libro.Url_imagen,
			&libro.Autor,
			&libro.Puntaje,
			&libro.UsuarioNombre,
			&fechaBytes,
		)

		if err != nil {
			fmt.Println("Error al escanear:", err)
			continue
		}

		libro.FechaCreacion = string(fechaBytes)
		libros = append(libros, libro)
	}

	data["libros"] = libros
	template.ExecuteTemplate(response, "Nosotros.html", data)
}
