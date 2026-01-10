package scripts

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-org/jvairv2/configs"
	commonAuth "github.com/your-org/jvairv2/pkg/common/auth"
	domainAuth "github.com/your-org/jvairv2/pkg/domain/auth"
	"github.com/your-org/jvairv2/pkg/repository/mysql"
	mysqlUser "github.com/your-org/jvairv2/pkg/repository/mysql/user"
)

// TestAuth prueba la autenticación de un usuario
func TestAuth() {
	// Cargar configuración
	config, err := configs.LoadConfig("./configs")
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Inicializar conexión a la base de datos
	dbConn, err := mysql.NewConnection(&config.DB)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer func() { _ = dbConn.Close() }()

	// Inicializar repositorio de usuarios
	userRepo := mysqlUser.NewRepository(dbConn.GetDB())

	// Inicializar servicios de autenticación
	tokenStore := commonAuth.NewMemoryTokenStore()
	authService := commonAuth.NewJWTService(
		config.JWT.AccessSecret,
		config.JWT.RefreshSecret,
		config.JWT.AccessExpiration,
		config.JWT.RefreshExpiration,
		tokenStore,
	)

	// Inicializar caso de uso de autenticación
	authUC := domainAuth.NewUseCase(userRepo, authService)

	// Credenciales de prueba (usar las del usuario admin creado anteriormente)
	email := "admin@example.com"
	password := "admin123"

	// Crear solicitud de login
	loginReq := &domainAuth.LoginRequest{
		Email:    email,
		Password: password,
	}

	// Intentar autenticar al usuario
	fmt.Println("Intentando autenticar al usuario:", email)
	resp, err := authUC.Login(context.Background(), loginReq)
	if err != nil {
		log.Fatalf("Error en la autenticación: %v", err)
	}

	// Mostrar información del token
	fmt.Println("Autenticación exitosa!")
	fmt.Println("Usuario:", resp.User.Name)
	fmt.Println("Email:", resp.User.Email)
	fmt.Println("Rol ID:", resp.User.RoleID)
	fmt.Println("Access Token:", resp.AccessToken)
	fmt.Println("Refresh Token:", resp.RefreshToken)
	fmt.Println("Expira en:", resp.ExpiresAt.Format(time.RFC3339))

	// Validar el token generado
	fmt.Println("\nValidando token...")
	valid, err := authUC.ValidateToken(context.Background(), resp.AccessToken)
	if err != nil {
		log.Fatalf("Error al validar el token: %v", err)
	}

	if valid {
		fmt.Println("Token válido!")
	} else {
		fmt.Println("Token inválido!")
	}

	// Obtener usuario a partir del token
	fmt.Println("\nObteniendo usuario a partir del token...")
	user, err := authUC.GetUserFromToken(context.Background(), resp.AccessToken)
	if err != nil {
		log.Fatalf("Error al obtener usuario del token: %v", err)
	}

	fmt.Println("Usuario obtenido del token:")
	fmt.Println("ID:", user.ID)
	fmt.Println("Nombre:", user.Name)
	fmt.Println("Email:", user.Email)
	fmt.Println("Rol ID:", user.RoleID)

	// Refrescar token
	fmt.Println("\nRefrescando token...")
	newTokens, err := authUC.RefreshToken(context.Background(), resp.RefreshToken)
	if err != nil {
		log.Fatalf("Error al refrescar el token: %v", err)
	}

	fmt.Println("Nuevo Access Token:", newTokens.AccessToken)
	fmt.Println("Nuevo Refresh Token:", newTokens.RefreshToken)

	// Cerrar sesión (eliminar token)
	fmt.Println("\nCerrando sesión...")
	err = authUC.Logout(context.Background(), resp.AccessToken)
	if err != nil {
		log.Fatalf("Error al cerrar sesión: %v", err)
	}

	fmt.Println("Sesión cerrada exitosamente!")

	// Intentar validar el token después de cerrar sesión
	fmt.Println("\nIntentando validar token después de cerrar sesión...")
	valid, err = authUC.ValidateToken(context.Background(), resp.AccessToken)
	if err != nil {
		fmt.Printf("Error esperado al validar token después de logout: %v\n", err)
	} else if !valid {
		fmt.Println("Token invalidado correctamente después del logout")
	} else {
		fmt.Println("¡Error! El token sigue siendo válido después del logout")
	}
}
