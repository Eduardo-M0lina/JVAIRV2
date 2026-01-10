package scripts

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/your-org/jvairv2/configs"
	"github.com/your-org/jvairv2/pkg/repository/mysql"
	"golang.org/x/crypto/bcrypt"
)

// CreateTestUser crea un usuario de prueba
func CreateTestUser() {
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

	// Verificar si el usuario de prueba ya existe
	var exists int
	err = dbConn.GetDB().QueryRow("SELECT COUNT(*) FROM users WHERE email = 'test@example.com'").Scan(&exists)
	if err != nil {
		log.Fatalf("Error al verificar si el usuario de prueba existe: %v", err)
	}

	if exists > 0 {
		fmt.Println("El usuario de prueba ya existe")
		return
	}

	// Obtener el ID del rol de usuario normal
	var roleID int64
	err = dbConn.GetDB().QueryRow("SELECT id FROM roles WHERE label = 'User' LIMIT 1").Scan(&roleID)
	if err != nil {
		// Si el rol no existe, crearlo
		now := time.Now()
		result, err := dbConn.GetDB().Exec(
			"INSERT INTO roles (label, created_at, updated_at) VALUES (?, ?, ?)",
			"User", now, now,
		)
		if err != nil {
			log.Fatalf("Error al crear el rol de usuario: %v", err)
		}
		roleID, err = result.LastInsertId()
		if err != nil {
			log.Fatalf("Error al obtener el ID del rol de usuario: %v", err)
		}
		fmt.Printf("Rol de usuario creado con ID: %d\n", roleID)
	} else {
		fmt.Printf("Rol de usuario encontrado con ID: %d\n", roleID)
	}

	// Crear el usuario de prueba
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error al generar el hash de la contraseña: %v", err)
	}

	now := time.Now()
	result, err := dbConn.GetDB().Exec(
		"INSERT INTO users (name, email, password, role_id, email_verified_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"Test User", "test@example.com", string(hashedPassword), roleID, now, now, now,
	)
	if err != nil {
		log.Fatalf("Error al crear el usuario de prueba: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error al obtener el ID del usuario de prueba: %v", err)
	}

	fmt.Printf("Usuario de prueba creado con ID: %d\n", userID)
}
