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

// CreateAdminRoleAndUser crea un rol de administrador y un usuario administrador
func CreateAdminRoleAndUser() {
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

	// Crear rol de administrador si no existe
	createAdminRole(dbConn)

	// Crear usuario de prueba
	createAdminUser(dbConn)
}

func createAdminRole(dbConn *mysql.Connection) {
	// Verificar si el rol ya existe
	var exists int
	err := dbConn.GetDB().QueryRow("SELECT COUNT(*) FROM roles WHERE label = 'Admin'").Scan(&exists)
	if err != nil {
		log.Fatalf("Error al verificar si el rol de administrador existe: %v", err)
	}

	if exists > 0 {
		fmt.Println("El rol de administrador ya existe")
		return
	}

	// Crear el rol de administrador
	now := time.Now()
	result, err := dbConn.GetDB().Exec(
		"INSERT INTO roles (label, created_at, updated_at) VALUES (?, ?, ?)",
		"Admin", now, now,
	)
	if err != nil {
		log.Fatalf("Error al crear el rol de administrador: %v", err)
	}

	roleID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error al obtener el ID del rol de administrador: %v", err)
	}

	fmt.Printf("Rol de administrador creado con ID: %d\n", roleID)
}

func createAdminUser(dbConn *mysql.Connection) {
	// Verificar si el usuario admin ya existe
	var exists int
	err := dbConn.GetDB().QueryRow("SELECT COUNT(*) FROM users WHERE email = 'admin@example.com'").Scan(&exists)
	if err != nil {
		log.Fatalf("Error al verificar si el usuario admin existe: %v", err)
	}

	if exists > 0 {
		fmt.Println("El usuario admin ya existe")
		return
	}

	// Obtener el ID del rol de administrador
	var roleID int64
	err = dbConn.GetDB().QueryRow("SELECT id FROM roles WHERE label = 'Admin' LIMIT 1").Scan(&roleID)
	if err != nil {
		log.Fatalf("Error al obtener el ID del rol de administrador: %v", err)
	}

	// Crear el usuario admin
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error al generar el hash de la contraseña: %v", err)
	}

	now := time.Now()
	result, err := dbConn.GetDB().Exec(
		"INSERT INTO users (name, email, password, role_id, email_verified_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"Admin", "admin@example.com", string(hashedPassword), roleID, now, now, now,
	)
	if err != nil {
		log.Fatalf("Error al crear el usuario admin: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error al obtener el ID del usuario admin: %v", err)
	}

	fmt.Printf("Usuario admin creado con ID: %d\n", userID)
}
