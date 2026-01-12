package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// Create crea un nuevo usuario
func (r *Repository) Create(ctx context.Context, u *user.User) error {
	log.Printf("[REPO] Creando usuario: %s <%s>", u.Name, u.Email)

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[REPO] ERROR al generar hash de contraseña: %v", err)
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
	}

	query := `
		INSERT INTO users (name, email, password, role_id,
		                  email_verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var emailVerifiedAtValue interface{}
	if u.EmailVerifiedAt != nil {
		emailVerifiedAtValue = *u.EmailVerifiedAt
	} else {
		emailVerifiedAtValue = nil
	}

	// Manejar el role_id que puede ser nulo
	var roleIDValue interface{}
	if u.RoleID != nil {
		roleIDValue = *u.RoleID
		log.Printf("[REPO] Asignando role_id: %s al usuario", *u.RoleID)
	} else {
		roleIDValue = nil
		log.Printf("[REPO] No se asignó role_id al usuario")
	}

	now := time.Now()
	nowPtr := &now

	log.Printf("[REPO] Ejecutando query de inserción para usuario: %s", u.Email)
	result, err := r.db.ExecContext(ctx, query,
		u.Name, u.Email, string(hashedPassword), roleIDValue,
		emailVerifiedAtValue, nowPtr, nowPtr,
	)

	if err != nil {
		log.Printf("[REPO] ERROR en la inserción de usuario: %v", err)
		return fmt.Errorf("error en la inserción de usuario: %w", err)
	}

	// Obtener el ID generado
	lastID, err := result.LastInsertId()
	if err != nil {
		log.Printf("[REPO] ERROR al obtener el ID generado: %v", err)
		return fmt.Errorf("error al obtener el ID generado: %w", err)
	}

	u.ID = lastID
	u.CreatedAt = nowPtr
	u.UpdatedAt = nowPtr

	log.Printf("[REPO] Usuario creado exitosamente con ID: %d", u.ID)
	return nil
}
