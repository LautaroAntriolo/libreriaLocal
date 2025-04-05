# Mi Librería

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23.2-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-AGPL-green)](LICENSE)

Una librería local escrita en Go que permite cargar libros por ISBN o URL de Mercado Libre, almacenar sus metadatos (portada, título, autor, etc.) y guardar notas asociadas a cada libro.

> **Nota:** Este proyecto está en desarrollo activo. Aunque es funcional, aún hay características pendientes como la implementación completa de las notas.

## Tabla de Contenidos

- [Características](#características)
- [Requisitos Previos](#requisitos-previos)
- [Instalación](#instalación)
- [Uso](#uso)
- [Ejemplos](#ejemplos)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Contribuciones](#contribuciones)
- [Licencia](#licencia)
- [Contacto](#contacto)

## Características

- Cargar libros por ISBN o URL de Mercado Libre.
- Almacenar metadatos (portada, título, autor, etc.) localmente.
- Guardar notas asociadas a cada libro (en desarrollo).
- Creación automática de tablas en la base de datos al iniciar la aplicación.

## Requisitos Previos

- [XAMPP](https://www.apachefriends.org/index.html) instalado con MySQL activado. 
- Go >= 1.23.2.
- Dependencias externas: `github.com/go-sql-driver/mysql`.

## Instalación

1. Clona el repositorio:

   ```bash
   git clone https://github.com/tu-usuario/mi-libreria.git
   cd mi-libreria
   ```
2. Instala las dependencias:

   ```
   `go mod tidy
   ```
3. Configura la conexión a la base de datos:
   crea un archivo llamado .env con las credenciales de MySQL:

   ```
   `{
       "db_user": "root",
       "db_password": "",
       "db_host": "127.0.0.1:3306"
   }
   ```
4. Inicializa la base de datos:
   Al ejecutar la aplicaci}on, las tablas se crearán automaticamente si no existen.
5. Ejecuta la aplicación:

   ```
   `go run main.go
   ```

## Uso


## Contribuciones

¡Las contribuciones son bienvenidas! Por favor, sigue estos pasos:

1. Haz un fork del repositorio.
2. Crea una nueva rama (`git checkout -b feature/nueva-funcionalidad`).
3. Haz commit de tus cambios (`git commit -m 'Añadir nueva funcionalidad'`).
4. Sube tus cambios (`git push origin feature/nueva-funcionalidad`).
5. Abre un pull request.

**Nota: **Todas las contribuciones deben mantener el proyecto local y seguro.

## Licencia

Este proyecto está bajo la licencia [AGPL ](https://chat.qwen.ai/c/LICENSE), que permite el uso libre del software pero prohíbe su uso comercial sin permiso explícito.

## Contacto

Si tienes preguntas o sugerencias, no dudes en contactarme:

* Email: [antriololautaro@gmail.com](mailto:antriololautaro@gmail.com)
* LinkedIn: [Lautaro Antriolo](https://www.linkedin.com/in/lautaro-antriolo/)
