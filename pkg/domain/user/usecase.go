package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/your-org/jvairv2/pkg/domain/ability"
	"github.com/your-org/jvairv2/pkg/domain/assigned_role"
	"github.com/your-org/jvairv2/pkg/domain/role"
)

var (
	ErrInvalidCredentials   = errors.New("credenciales inválidas")
	ErrUserInactive         = errors.New("usuario inactivo")
	ErrDuplicateEmail       = errors.New("email ya está en uso")
	ErrUserNotFound         = errors.New("usuario no encontrado")
	ErrAssignedRoleNotFound = errors.New("asignación de rol no encontrada")
)

// UseCase define los casos de uso para la gestión de usuarios
type UseCase struct {
	repo             Repository
	assignedRoleRepo assigned_role.Repository
	roleRepo         role.Repository
}

// NewUseCase crea una nueva instancia del caso de uso de usuarios
func NewUseCase(repo Repository, assignedRoleRepo assigned_role.Repository, roleRepo role.Repository) *UseCase {
	return &UseCase{
		repo:             repo,
		assignedRoleRepo: assignedRoleRepo,
		roleRepo:         roleRepo,
	}
}

// GetByID obtiene un usuario por su ID
func (uc *UseCase) GetByID(ctx context.Context, id string) (*User, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetByEmail obtiene un usuario por su email
func (uc *UseCase) GetByEmail(ctx context.Context, email string) (*User, error) {
	return uc.repo.GetByEmail(ctx, email)
}

// Create crea un nuevo usuario
func (uc *UseCase) Create(ctx context.Context, user *User) error {
	// Verificar si el repositorio está inicializado
	if uc.repo == nil {
		log.Printf("ERROR: Repositorio de usuarios no inicializado")
		return errors.New("repositorio de usuarios no inicializado")
	}

	// Verificar si ya existe un usuario con el mismo email
	_, err := uc.repo.GetByEmail(ctx, user.Email)
	if err == nil {
		log.Printf("ERROR: Email ya en uso: %s", user.Email)
		return ErrDuplicateEmail
	} else if !errors.Is(err, ErrUserNotFound) {
		log.Printf("ERROR al verificar email existente: %v", err)
		return fmt.Errorf("error al verificar email existente: %w", err)
	}

	log.Printf("Email disponible para nuevo usuario: %s", user.Email)

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR al generar hash de contraseña: %v", err)
		return fmt.Errorf("error al generar hash de contraseña: %w", err)
	}
	user.Password = string(hashedPassword)

	// Establecer valores predeterminados
	now := time.Now()
	user.CreatedAt = &now
	user.UpdatedAt = &now
	user.IsActive = true

	// Crear el usuario
	log.Printf("Intentando crear usuario: %s <%s>", user.Name, user.Email)
	err = uc.repo.Create(ctx, user)
	if err != nil {
		log.Printf("ERROR al crear usuario en la base de datos: %v", err)
		return fmt.Errorf("error al crear usuario en la base de datos: %w", err)
	}
	log.Printf("Usuario creado exitosamente con ID: %d", user.ID)

	// Si se especificó un rol, asignar el rol al usuario
	if user.RoleID != nil {
		roleID, err := uc.getRoleID(*user.RoleID)
		if err != nil {
			// Si falla la obtención del rol, registrar el error pero continuar
			// ya que el usuario ya fue creado
			log.Printf("Error al obtener rol para el usuario %d: %v", user.ID, err)
		} else {
			assignedRole := &assigned_role.AssignedRole{
				RoleID:     roleID,
				EntityID:   user.ID,
				EntityType: "App\\Models\\User",
				Restricted: false,
			}

			err = uc.assignedRoleRepo.Assign(ctx, assignedRole)
			if err != nil {
				// Si falla la asignación del rol, registrar el error pero continuar
				// ya que el usuario ya fue creado
				log.Printf("Error al asignar rol al usuario %d: %v", user.ID, err)
			}
		}
	}

	return nil
}

// Update actualiza un usuario existente
func (uc *UseCase) Update(ctx context.Context, user *User) error {
	// Verificar si el usuario existe
	userIDStr := strconv.FormatInt(user.ID, 10)
	existingUser, err := uc.repo.GetByID(ctx, userIDStr)
	if err != nil {
		return err
	}

	// Verificar si el email ya está en uso por otro usuario
	if user.Email != existingUser.Email {
		otherUser, err := uc.repo.GetByEmail(ctx, user.Email)
		if err == nil && otherUser != nil && otherUser.ID != user.ID {
			return ErrDuplicateEmail
		} else if err != nil && !errors.Is(err, ErrUserNotFound) {
			return err
		}
	}

	// Si se proporciona una nueva contraseña, hashearla
	if user.Password != "" && user.Password != existingUser.Password {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		// Mantener la contraseña existente
		user.Password = existingUser.Password
	}

	// Actualizar timestamp
	now := time.Now()
	user.UpdatedAt = &now

	// Actualizar el usuario
	err = uc.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Si se especificó un nuevo rol, actualizar la asignación de rol
	if user.RoleID != nil && (existingUser.RoleID == nil || *user.RoleID != *existingUser.RoleID) {
		// Obtener el ID numérico del rol
		roleID, err := uc.getRoleID(*user.RoleID)
		if err != nil {
			return err
		}

		// Verificar si ya tiene algún rol asignado
		userIDStr := strconv.FormatInt(user.ID, 10)
		roles, err := uc.repo.GetUserRoles(ctx, userIDStr)
		if err != nil && !errors.Is(err, ErrUserNotFound) {
			return err
		}

		if len(roles) > 0 {
			// Revocar roles existentes
			for _, r := range roles {
				err = uc.assignedRoleRepo.Revoke(ctx, r.ID, user.ID, "App\\Models\\User")
				if err != nil && !errors.Is(err, ErrAssignedRoleNotFound) {
					return err
				}
			}
		}

		// Asignar el nuevo rol
		assignedRole := &assigned_role.AssignedRole{
			RoleID:     roleID,
			EntityID:   user.ID,
			EntityType: "App\\Models\\User",
			Restricted: false,
		}

		err = uc.assignedRoleRepo.Assign(ctx, assignedRole)
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete elimina un usuario (soft delete)
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// List obtiene una lista paginada de usuarios con filtros opcionales
func (uc *UseCase) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*User, int, error) {
	return uc.repo.List(ctx, filters, page, pageSize)
}

// VerifyCredentials verifica las credenciales de un usuario
func (uc *UseCase) VerifyCredentials(ctx context.Context, email, password string) (*User, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GetUserRoles obtiene los roles de un usuario
func (uc *UseCase) GetUserRoles(ctx context.Context, userID string) ([]*role.Role, error) {
	return uc.repo.GetUserRoles(ctx, userID)
}

// GetUserAbilities obtiene las habilidades de un usuario
func (uc *UseCase) GetUserAbilities(ctx context.Context, userID string) ([]*ability.Ability, error) {
	return uc.repo.GetUserAbilities(ctx, userID)
}

// HasAbility verifica si un usuario tiene una habilidad específica
func (uc *UseCase) HasAbility(ctx context.Context, userID string, abilityName string) (bool, error) {
	abilities, err := uc.repo.GetUserAbilities(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, a := range abilities {
		if a.Name == abilityName {
			return true, nil
		}
	}

	return false, nil
}

// Helper para obtener el ID numérico de un rol a partir de su ID de string
func (uc *UseCase) getRoleID(roleIDStr string) (int64, error) {
	// Verificar si el repositorio está inicializado
	if uc.roleRepo == nil {
		log.Printf("ERROR CRÍTICO: Repositorio de roles no inicializado")
		return 0, errors.New("repositorio de roles no inicializado")
	}

	log.Printf("Intentando obtener rol con ID/nombre: '%s'", roleIDStr)

	// Intentar obtener el rol por ID
	roleID, err := parseRoleID(roleIDStr)
	if err == nil {
		log.Printf("Interpretado como ID numérico: %d", roleID)
		role, err := uc.roleRepo.GetByID(context.Background(), roleID)
		if err == nil && role != nil {
			log.Printf("Rol encontrado por ID: %d, nombre: %s", role.ID, role.Name)
			return roleID, nil
		}
		// Si hay un error específico, registrarlo
		if err != nil {
			log.Printf("ERROR al obtener rol por ID %d: %v", roleID, err)
		} else {
			log.Printf("ERROR: Rol con ID %d no encontrado o es nil", roleID)
		}
	} else {
		log.Printf("No se pudo interpretar '%s' como ID numérico: %v", roleIDStr, err)
	}

	// Si no se encontró por ID, intentar por nombre
	log.Printf("Intentando obtener rol por nombre: '%s'", roleIDStr)
	r, err := uc.roleRepo.GetByName(context.Background(), roleIDStr)
	if err != nil {
		// Registrar el error específico
		log.Printf("ERROR al obtener rol por nombre '%s': %v", roleIDStr, err)
		return 0, fmt.Errorf("error al obtener rol por nombre '%s': %w", roleIDStr, err)
	}

	log.Printf("Rol encontrado por nombre: '%s', ID: %d", roleIDStr, r.ID)
	return r.ID, nil
}

// Helper para convertir un string a int64
func parseRoleID(roleID string) (int64, error) {
	return strconv.ParseInt(roleID, 10, 64)
}
