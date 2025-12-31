package scripts

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/your-org/jvairv2/configs"
	"github.com/your-org/jvairv2/pkg/repository/mysql"
)

// TestDB prueba la conexión a la base de datos
func TestDB() {
	// Cargar configuración
	configPath := filepath.Join(".", "configs")
	config, err := configs.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Conectar a la base de datos
	dbConn, err := mysql.NewConnection(&config.DB)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	// Verificar la conexión a la base de datos
	err = dbConn.GetDB().Ping()
	if err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}

	fmt.Println("Conexión a la base de datos establecida correctamente")

	// Listar tablas en la base de datos
	rows, err := dbConn.GetDB().Query("SHOW TABLES")
	if err != nil {
		log.Fatalf("Error al listar tablas: %v", err)
	}
	defer func() { _ = rows.Close() }()

	fmt.Println("Tablas en la base de datos:")
	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatalf("Error al escanear nombre de tabla: %v", err)
		}
		fmt.Println("-", tableName)
	}
}
