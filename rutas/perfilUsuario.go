package rutas

import (
	"ProyectoWEB/conectar"
	"ProyectoWEB/modelos"
	"ProyectoWEB/utilidades"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func MiPerfil(response http.ResponseWriter, request *http.Request) {

	conectar.Conectar()
	defer conectar.CerrarConexion()

	data := DatosDeSeguridad(response, request)
	usuarioID := data["usuario_id"]

	sql := "SELECT nombre, correo, telefono, imagen_perfil FROM usuarios_library WHERE id = ?"
	datos, err := conectar.Db.Query(sql, usuarioID)
	if err != nil {
		http.Error(response, "Error en la consulta a la base de datos", http.StatusInternalServerError)
		fmt.Println("Error en la consulta SQL:", err)
		return
	}
	defer datos.Close()

	var usuario []modelos.Usuario
	for datos.Next() {
		var dato modelos.Usuario
		err := datos.Scan(&dato.Nombre, &dato.Correo, &dato.Telefono, &dato.Imagen_perfil)
		if err != nil {
			fmt.Println("Error al escanear los datos:", err)
			continue
		}
		usuario = append(usuario, dato)
	}

	if err = datos.Err(); err != nil {
		fmt.Println("Error al iterar sobre los datos:", err)
	}

	// Agregar los datos del usuario al mapa `data`
	data["Datos"] = usuario

	// Parsear las plantillas
	template := template.Must(template.ParseFiles(
		"templates/perfil/perfil.html",
		"templates/principales/navbarLogin.html",
		utilidades.Frontend,
	))

	// Renderizar la plantilla principal (`perfil.html`) y pasarle los datos
	err = template.ExecuteTemplate(response, "perfil.html", data)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		http.Error(response, "Error al renderizar la página", http.StatusInternalServerError)
	}
}

func EditarPerfilForm(response http.ResponseWriter, request *http.Request) {

	conectar.Conectar()
	defer conectar.CerrarConexion()

	data := DatosDeSeguridad(response, request)
	usuarioID := data["usuario_id"]
	fmt.Println(usuarioID)
	sql := "SELECT id, nombre, correo, telefono, imagen_perfil, password FROM usuarios_library WHERE id = ?"
	datos, err := conectar.Db.Query(sql, usuarioID)
	if err != nil {
		http.Error(response, "Error en la consulta a la base de datos", http.StatusInternalServerError)
		fmt.Println("Error en la consulta SQL:", err)
		return
	}
	defer datos.Close()

	var usuario []modelos.Usuario
	for datos.Next() {
		var dato modelos.Usuario
		err := datos.Scan(&dato.Id, &dato.Nombre, &dato.Correo, &dato.Telefono, &dato.Imagen_perfil, &dato.Password)
		if err != nil {
			fmt.Println("Error al escanear los datos:", err)
			continue
		}
		usuario = append(usuario, dato)
	}

	if err = datos.Err(); err != nil {
		fmt.Println("Error al iterar sobre los datos:", err)
	}

	// Agregar los datos del usuario al mapa `data`
	data["Datos"] = usuario

	// Cargar la plantilla HTML
	tmpl := template.Must(template.ParseFiles(
		"templates/perfil/editarPerfil.html",
		"templates/principales/navbarLogin.html",
		utilidades.Frontend))

	// Corregido: usar el nombre correcto de la plantilla
	err = tmpl.ExecuteTemplate(response, "editarPerfil.html", data)
	if err != nil {
		fmt.Println("Error al ejecutar la plantilla:", err)
		http.Error(response, "Error al renderizar la página", http.StatusInternalServerError)
	}
}

func ActualizarPerfil(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Usuario inválido", http.StatusBadRequest)
		return
	}

	// Recopilar datos del formulario
	nombre := r.FormValue("nombre")
	url_imagen := r.FormValue("url_imagen")
	correo := r.FormValue("correo")
	telefono := r.FormValue("telefono")

	// Actualizar libro en la base de datos
	err = actualizarPerfilEnDB(id, nombre, url_imagen, correo, telefono)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirigir a la página de detalles del libro
	utilidades.CrearMensajesFlash(w, r, "success", "Usuario actualizado")
	http.Redirect(w, r, "/perfil", http.StatusSeeOther)
}

func actualizarPerfilEnDB(id int, nombre, url_imagen string, correo string, telefono string) error {

	conectar.Conectar()
	defer conectar.CerrarConexion()

	// Corregido según la estructura real de la tabla
	query := `UPDATE usuarios_library 
             SET nombre = ?, correo = ?, telefono = ?, imagen_perfil = ? WHERE id = ?`

	_, err := conectar.Db.Exec(query, nombre, correo, telefono, url_imagen, id)

	if err != nil {
		return fmt.Errorf("error al actualizar el usuario: %w", err)
	}

	return nil
}
