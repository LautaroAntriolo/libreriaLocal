package proteccion

import (
	"ProyectoWEB/utilidades"
	"net/http"
)

/*
// para que sea interpretada como un mildewhert necesita el parametro nect de tipo handled
func Proteccion(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(response, request) // => Esto es equivalente a decirle que podemos continuar!
	}
}
*/
// Actualiza tu función Proteccion existente en proteccion/proteccion.go
func Proteccion(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Establecer encabezados de control de caché
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Obtener sesión
		session, err := utilidades.Store.Get(r, "session-name")

		// Verificar si el usuario está autenticado o si la sesión tiene error
		if err != nil || session.Values["usuario_id"] == nil {
			utilidades.CrearMensajesFlash(w, r, "danger", "Debes iniciar sesión para acceder a esta página")
			http.Redirect(w, r, "/seguridad/login", http.StatusSeeOther)
			return
		}

		// Actualizar la sesión para que expire cuando se cierre el navegador
		session.Options.MaxAge = 0
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next(w, r)
	}
}

func NoCacheMiddleware(next http.Handler) http.Handler {
	// Obtener sesión

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Establecer encabezados para prevenir el almacenamiento en caché
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Llamar al siguiente manejador
		next.ServeHTTP(w, r)
	})
}
