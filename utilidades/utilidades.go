package utilidades

import (
	"os"
	"net/http"
	"github.com/gorilla/sessions"
	gomail "gopkg.in/gomail.v2"
	"fmt"
)
var Frontend string = "templates/layout/frontend.html"

// genero una variable llamada Store, que generará una instancia del nomrbe de la cookie con que se genera la cookie.
// El mensaje flash es una cookie que se genera por un segundo
var Store = sessions.NewCookieStore([]byte("session-name"))
func init() {
    Store.Options = &sessions.Options{
        Path:     "/",
        // MaxAge:   86400 * 30,  // 30 días
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    }
}

func RetornarMensajesFlash(response http.ResponseWriter, request *http.Request)(string, string){
	session, _ := Store.Get(request, "flash-session")
	fm := session.Flashes("css")
	session.Save(request,response)
	css_sesion := ""
	if len(fm)== 0{
		css_sesion = ""
	}else{
		css_sesion = fm[0].(string)
	}
	fm2 := session.Flashes("mensaje")
	session.Save(request, response)
	css_mensaje := ""
	if len(fm2)== 0{
		css_mensaje = ""
	}else{
		css_mensaje = fm2[0].(string)
	}
	return css_sesion, css_mensaje

}

func CrearMensajesFlash(response http.ResponseWriter, request *http.Request,css string, mensaje string){
	session, err := Store.Get(request, "flash-session")
	if err != nil{
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return // siempre un return para que no se caiga todo cuando termine.
	}
	session.AddFlash(css, "css") // una cookie con el estilo css de las alertas de bootstrap.
	session.AddFlash(mensaje, "mensaje")
	session.Save(request, response)
}

func EnviarCorreo(To string, subject string, bodyHTML string, attach string) error { // Modifica la firma de la función para retornar un error
    msg := gomail.NewMessage()
    msg.SetHeader("From", os.Getenv("MAIL_SALIDA"))
    msg.SetHeader("To", To)
    msg.SetHeader("Subject", subject)
    msg.SetBody("text/html", bodyHTML)

	if attach != "" { // Verifica si la ruta del archivo adjunto no está vacía
        msg.Attach(attach)
    }

    n := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("MAIL_SALIDA"), os.Getenv("CONTRA_SEGURA"))

    if err := n.DialAndSend(msg); err != nil {
        fmt.Println("Error al enviar el correo:", err)
        return err // Retorna el error
    }
    return nil // Retorna nil si no hay error
}



func RetornarLogin(request *http.Request)(string, string){
	session, _ := Store.Get(request, "session-name")
	usuario_id:=""
	usuario_name:=""
	if session.Values["usuario_id"]!=nil{
		usuario_id_t, _ :=session.Values["usuario_id"].(string)
		usuario_id = usuario_id_t
	}
	if session.Values["usuario_name"]!=nil{
		usuario_name_t, _ :=session.Values["usuario_name"].(string)
		usuario_name = usuario_name_t
	}
	return usuario_id, usuario_name
}
