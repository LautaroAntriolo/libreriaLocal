package modelos


import "time"

type Usuario struct {
	Id            int       `json:"id"`
	Nombre        string    `json:"nombre"`
	Correo        string    `json:"correo"`
	Telefono      string    `json:"telefono"`
	Password      string    `json:"password"`
	Imagen_perfil string    `json:"imagen_perfil"`
	FechaRegistro time.Time `json:"fecha_registro"`
}

type Usuarios []Usuario
