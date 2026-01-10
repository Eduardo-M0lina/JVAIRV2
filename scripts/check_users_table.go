package scripts

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/your-org/jvairv2/configs"
	"github.com/your-org/jvairv2/pkg/repository/mysql"
)

// CheckUsersTable verifica la tabla de usuarios en la base de datos
func CheckUsersTable() {
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

	log.Println("Conexión a la base de datos establecida correctamente")

	// Verificar si la tabla users existe
	var exists int
	err = dbConn.GetDB().QueryRow("SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'users' LIMIT 1").Scan(&exists)
	if err != nil {
		log.Fatalf("Error al verificar la tabla users: %v", err)
	}

	if exists == 1 {
		fmt.Println("La tabla users existe")
	} else {
		fmt.Println("La tabla users no existe")
	}

	// Contar registros en la tabla users
	var count int
	err = dbConn.GetDB().QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Error al contar registros en la tabla users: %v", err)
	}

	fmt.Printf("La tabla users tiene %d registros\n", count)
}
