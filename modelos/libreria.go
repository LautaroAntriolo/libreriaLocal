package modelos

import (
	"database/sql"
	"strconv"
	"time"
)

type Libro struct {
	UsuarioID    string
	Nombre       string
	Url_imagen   string
	Autor        string
	Id           int
	Isbn         string
	Editorial    string
	Descripcion  string
	Comentarios  string
	Critica      string
	Leido        string
	Puntaje      sql.NullInt64
	FechaLectura *time.Time
}

// Método para obtener el puntaje como string
func (l Libro) GetPuntaje() string {
	if l.Puntaje.Valid {
		return strconv.FormatInt(l.Puntaje.Int64, 10)
	}
	return "Sin puntaje"
}

type Libros []Libro

type ClienteHttp struct {
	Css     string
	Mensaje string
	Datos   Libros
}
