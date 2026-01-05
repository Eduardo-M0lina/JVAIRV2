package ability

import (
	"context"
	"time"

	"github.com/your-org/jvairv2/pkg/domain/ability"
)

// Create crea una nueva ability
func (r *Repository) Create(ctx context.Context, ability *ability.Ability) error {
	// Verificar si ya existe una ability con el mismo nombre
	existingAbility, err := r.GetByName(ctx, ability.Name)
	if err == nil && existingAbility != nil {
		return ErrDuplicateName
	} else if err != ErrAbilityNotFound {
		return err
	}

	// Establecer timestamps
	now := time.Now()
	ability.CreatedAt = &now
	ability.UpdatedAt = &now

	// Preparar la consulta
	query := `
		INSERT INTO abilities (name, title, entity_id, entity_type, only_owned, options, scope, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Preparar los argumentos
	var args []interface{}
	args = append(args, ability.Name)

	// Manejar campos opcionales
	if ability.Title != nil {
		args = append(args, *ability.Title)
	} else {
		args = append(args, nil)
	}

	if ability.EntityID != nil {
		args = append(args, *ability.EntityID)
	} else {
		args = append(args, nil)
	}

	if ability.EntityType != nil {
		args = append(args, *ability.EntityType)
	} else {
		args = append(args, nil)
	}

	args = append(args, ability.OnlyOwned)

	if ability.Options != nil {
		args = append(args, *ability.Options)
	} else {
		args = append(args, nil)
	}

	if ability.Scope != nil {
		args = append(args, *ability.Scope)
	} else {
		args = append(args, nil)
	}

	args = append(args, now, now)

	// Ejecutar la consulta
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// Obtener el ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	ability.ID = id
	return nil
}
