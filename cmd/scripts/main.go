package main

import (
	"fmt"
	"os"

	"github.com/your-org/jvairv2/scripts"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "check-roles-table":
		scripts.CheckRolesTable()
	case "check-users-table":
		scripts.CheckUsersTable()
	case "create-admin-user":
		scripts.CreateAdminRoleAndUser()
	case "create-test-user":
		scripts.CreateTestUser()
	case "test-auth":
		scripts.TestAuth()
	case "test-db":
		scripts.TestDB()
	default:
		fmt.Printf("Comando desconocido: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Uso: go run cmd/scripts/main.go <comando>")
	fmt.Println("Comandos disponibles:")
	fmt.Println("  check-roles-table    - Verificar la tabla de roles")
	fmt.Println("  check-users-table    - Verificar la tabla de usuarios")
	fmt.Println("  create-admin-user    - Crear un usuario administrador")
	fmt.Println("  create-test-user     - Crear un usuario de prueba")
	fmt.Println("  test-auth            - Probar la autenticación")
	fmt.Println("  test-db              - Probar la conexión a la base de datos")
}
