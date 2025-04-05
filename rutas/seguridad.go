package rutas

import (
	"ProyectoWEB/conectar"
	"ProyectoWEB/modelos"
	"ProyectoWEB/utilidades"
	"ProyectoWEB/validaciones"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Seguridad_registro(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/seguridad/registro.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))
	data := DatosDeSeguridad(response, request)
	template.ExecuteTemplate(response, "registro.html", data)
}

func Seguridad_registro_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""
	if len(request.FormValue("nombre")) == 0 {
		mensaje = mensaje + "El campo nombre está vacío"
	}
	if len(request.FormValue("correo")) == 0 {
		mensaje = mensaje + "El campo correo está vacío"
	}
	if validaciones.Regex_correo.FindStringSubmatch(request.FormValue("correo")) == nil {
		mensaje = mensaje + "El mail ingresado no es válido"
	}
	if !validaciones.ValidarPassword(request.FormValue("password")) {
		mensaje = mensaje + "La contraseña debe tener un largo entre 6 y 12 caracteres "
	}

	if mensaje != "" {
		utilidades.CrearMensajesFlash(response, request, "danger", mensaje)
		http.Redirect(response, request, "/seguridad/registro", http.StatusSeeOther)
		return
	}

	conectar.Conectar()
	defer conectar.CerrarConexion()

	sql := "INSERT INTO usuarios_library (nombre, correo, telefono, password, fecha_registro, imagen_perfil) VALUES (?, ?, ?, ?, ?, ?);"
	costo := 8
	bytes, err := bcrypt.GenerateFromPassword([]byte(request.FormValue("password")), costo)
	if err != nil {
		utilidades.CrearMensajesFlash(response, request, "danger", "Error al generar hash de contraseña")
		http.Redirect(response, request, "/seguridad/registro", http.StatusSeeOther)
		return
	}
	var count int
	err = conectar.Db.QueryRow("SELECT COUNT(*) FROM usuarios_library WHERE correo = ?", request.FormValue("correo")).Scan(&count)
	if err != nil || count > 0 {
		utilidades.CrearMensajesFlash(response, request, "danger", "El correo ya está registrado")
		http.Redirect(response, request, "/seguridad/registro", http.StatusSeeOther)
		return
	}

	// Ejecutar la consulta y capturar el resultado en 'result'
	result, err := conectar.Db.Exec(sql,
		request.FormValue("nombre"),
		request.FormValue("correo"),
		request.FormValue("telefono"),
		string(bytes),
		time.Now(),
		"sin imagen",
	)
	if err != nil {
		utilidades.CrearMensajesFlash(response, request, "danger", "Error al crear el registro: "+err.Error())
		http.Redirect(response, request, "/seguridad/registro", http.StatusSeeOther)
		return
	}

	// Obtener el ID del usuario recién insertado
	idUsuario, err := result.LastInsertId()
	if err != nil {
		utilidades.CrearMensajesFlash(response, request, "danger", "Error al obtener ID del usuario")
		http.Redirect(response, request, "/seguridad/registro", http.StatusSeeOther)
		return
	}

	// Crear sesión automáticamente
	session, _ := utilidades.Store.Get(request, "session-name")
	session.Options.MaxAge = 0 // Sesión persistente hasta cerrar navegador
	session.Values["usuario_id"] = strconv.FormatInt(idUsuario, 10)
	session.Values["usuario_name"] = request.FormValue("nombre")
	err = session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	utilidades.CrearMensajesFlash(response, request, "success", "¡Registro exitoso! Sesión iniciada automáticamente")
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

func Seguridad_login(response http.ResponseWriter, request *http.Request) {
	// Definir un mapa de funciones personalizadas
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s) // Esto permite renderizar HTML sin escaparlo
		},
	}
	template := template.Must(template.New("login.html").Funcs(funcMap).
		ParseFiles("templates/seguridad/login.html",
			"templates/principales/navbarLogin.html",
			"templates/principales/navbarNotLogin.html",
			utilidades.Frontend))
	css_sesion, css_mensaje := utilidades.RetornarMensajesFlash(response, request)

	data := map[string]string{

		"css":     css_sesion,
		"mensaje": css_mensaje,
	}
	template.ExecuteTemplate(response, "login.html", data)
}
func Seguridad_login_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""
	if len(request.FormValue("nombre")) == 0 {
		mensaje = mensaje + "El campo nombre está vacío"
	}
	if len(request.FormValue("correo")) == 0 {
		mensaje = mensaje + "El campo correo está vacío"
	}
	if validaciones.Regex_correo.FindStringSubmatch(request.FormValue("correo")) == nil {
		mensaje = mensaje + "El mail ingresado no es válido"
	}
	if !validaciones.ValidarPassword(request.FormValue("password")) {
		// if validaciones.ValidarPassword(password)==false{
		mensaje = mensaje + "La contraseña debe tener un largo entre 6 y 12 caracteres "
	}
	if mensaje != "" {
		// fmt.Fprintln(response, mensaje)
		// return
		utilidades.CrearMensajesFlash(response, request, "danger", mensaje)
		http.Redirect(response, request, "/seguridad/login", http.StatusSeeOther)
	}

	// Nos conectamos a la db y comparamos los password
	conectar.Conectar()
	sql := "SELECT id, nombre, correo, telefono, password FROM usuarios_library where correo=?"
	datos, err := conectar.Db.Query(sql, request.FormValue("correo"))
	if err != nil {
		fmt.Println(err)
	}
	defer conectar.CerrarConexion()
	var dato modelos.Usuario
	for datos.Next() {
		errNext := datos.Scan(&dato.Id, &dato.Nombre, &dato.Correo, &dato.Telefono, &dato.Password)
		if errNext != nil {
			utilidades.CrearMensajesFlash(response, request, "danger", "Las credenciales son inválidas")
			http.Redirect(response, request, "/seguridad/login", http.StatusSeeOther)
		}
	}
	// Ahora comparamos los hash! EL hash que genero de l contraseña del formulario y el de labase de datos

	passwordBytes := []byte(request.FormValue("password"))
	passwordDB := []byte(dato.Password)
	errPassword := bcrypt.CompareHashAndPassword(passwordDB, passwordBytes)
	if errPassword != nil {
		utilidades.CrearMensajesFlash(response, request, "danger", "Las credenciales son inválidas")
		http.Redirect(response, request, "/seguridad/login", http.StatusSeeOther)
	} else {
		session, _ := utilidades.Store.Get(request, "session-name")

		session.Options.MaxAge = 0 // 0 significa que la cookie expira cuando se cierra el navegador

		session.Values["usuario_id"] = strconv.Itoa(dato.Id)
		session.Values["usuario_name"] = dato.Nombre

		err2 := session.Save(request, response)
		if err2 != nil {
			http.Error(response, err2.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(response, request, "/", http.StatusSeeOther)
	}
}

func Seguridad_protegida(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	response.Header().Set("Pragma", "no-cache")
	response.Header().Set("Expires", "0")
	template := template.Must(template.ParseFiles("templates/seguridad/registro.html", utilidades.Frontend))
	css_sesion, css_mensaje := utilidades.RetornarMensajesFlash(response, request)
	usuario_id, usuario_name := utilidades.RetornarLogin(request)
	data := map[string]string{

		"css":          css_sesion,
		"mensaje":      css_mensaje,
		"usuario_id":   usuario_id,
		"usuario_name": usuario_name,
	}
	template.Execute(response, data)
}
func DatosDeSeguridad(response http.ResponseWriter, request *http.Request) map[string]interface{} {
	usuario_id, usuario_name := utilidades.RetornarLogin(request)
	css_sesion, css_mensaje := utilidades.RetornarMensajesFlash(response, request)
	data := map[string]interface{}{

		"css":          css_sesion,
		"mensaje":      css_mensaje,
		"usuario_id":   usuario_id,
		"usuario_name": usuario_name,
	}

	return data
}
func Seguridad_logout(response http.ResponseWriter, request *http.Request) {
	// 1. Invalidar la sesión en el servidor
	session, _ := utilidades.Store.Get(request, "session-name")
	session.Options.MaxAge = -1 // Marca la sesión para eliminación inmediata
	err := session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Eliminar la cookie del cliente
	cookie := &http.Cookie{
		Name:     "session-name",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,                               // Elimina inmediatamente la cookie
		Expires:  time.Now().Add(-100 * time.Hour), // Fecha en pasado
		HttpOnly: true,
	}
	http.SetCookie(response, cookie)
	// Agregar un script para limpiar datos locales (opcional)

	utilidades.CrearMensajesFlash(response, request, "primary", "Se ha cerrado tu sesión exitosamente")
	http.Redirect(response, request, "/seguridad/login", http.StatusSeeOther)
}

func Pagina404(w http.ResponseWriter, r *http.Request, mensaje string) {
    data := map[string]interface{}{
        "MensajeError": mensaje,
    }

    
    tmpl := template.Must(template.ParseFiles("templates/principales/error404.html", "templates/principales/navbarLogin.html", "templates/principales/navbarNotLogin.html", utilidades.Frontend))
    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Error al renderizar la página de error", http.StatusInternalServerError)
    }
}
