package conectar

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //Este es un "import lateral" (side-effect import).*/
	"os"

	"github.com/joho/godotenv"
)

// función para conectarnos a la BD

var Db *sql.DB 

func Conectar(){
	//Cargo las variables de entorno y me fijo que no tenga algún error: 
	errorVariables := godotenv.Load()
	if errorVariables!=nil{
		panic(errorVariables)
	}
	// Genero la conexión:
	conection, err := sql.Open("mysql",os.Getenv("DB_USER")+":@tcp("+os.Getenv("DB_SERVER")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))	
	if err != nil{
		panic(err)
	}
	Db = conection
}

func CerrarConexion(){
	Db.Close()
}
