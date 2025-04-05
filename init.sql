CREATE TABLE todosloslibros (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    isbn TEXT DEFAULT NULL,
    nombre TEXT NOT NULL,
    autor TEXT NULL,
    editorial TEXT NULL,
    pagina TEXT NULL,
    propiedad TEXT NULL,
    regalo_para TEXT DEFAULT NULL,
    portada TEXT DEFAULT NULL,
    Comentarios TEXT DEFAULT "",
    Descripcion TEXT DEFAULT "",
    fecha_creacion TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    leido INTEGER DEFAULT 0,
    porcentaje_de_lectura INTEGER DEFAULT 0,
    fecha_lectura DATE DEFAULT NULL,
    puntaje INTEGER DEFAULT 0,
    UsuarioId INTEGER DEFAULT 0
);

CREATE TABLE usuarios_library (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    correo TEXT NOT NULL,
    telefono TEXT NULL,
    password TEXT NOT NULL,
    fecha_registro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    imagen_perfil TEXT DEFAULT NULL
);