CREATE DATABASE IF NOT EXISTS golang_consola;

USE golang_consola;

CREATE TABLE todosloslibros (
    id INT(11) NOT NULL AUTO_INCREMENT,
    isbn VARCHAR(20) DEFAULT NULL,
    nombre VARCHAR(255) NOT NULL,
    autor VARCHAR(255) NOT NULL,
    editorial VARCHAR(100) NOT NULL,
    pagina VARCHAR(255) NOT NULL,
    propiedad VARCHAR(5) NOT NULL,
    regalo_para VARCHAR(100) DEFAULT NULL,
    portada TEXT DEFAULT NULL,
    Comentarios TEXT CHARACTER SET utf8 COLLATE utf8_spanish2_ci NOT NULL,
    Descripcion TEXT CHARACTER SET utf8 COLLATE utf8_spanish_ci NOT NULL,
    fecha_creacion TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    leido TINYINT(1) DEFAULT 0,
    porcentaje_de_lectura INT(11) DEFAULT 0,
    fecha_lectura DATE DEFAULT NULL,
    puntaje INT(11) DEFAULT 0,
    UsuarioId INT(11) DEFAULT 0,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE usuarios_library (
    id INT(11) NOT NULL AUTO_INCREMENT,
    nombre VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    correo VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    telefono VARCHAR(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    password VARCHAR(160) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    fecha_registro TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    imagen_perfil VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
