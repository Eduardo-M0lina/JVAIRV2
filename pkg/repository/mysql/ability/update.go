package ability

import (
	"context"
	"time"

	domainAbility "github.com/your-org/jvairv2/pkg/domain/ability"
)

// Update actualiza una ability existente
func (r *Repository) Update(ctx context.Context, ability *domainAbility.Ability) error {
	// Verificar si la ability existe
	existingAbility, err := r.GetByID(ctx, ability.ID)
	if err != nil {
		return err
	}

	// Verificar si el nombre ya está en uso por otra ability
	if ability.Name != existingAbility.Name {
		otherAbility, err := r.GetByName(ctx, ability.Name)
		if err == nil && otherAbility != nil && otherAbility.ID != ability.ID {
			return ErrDuplicateName
		} else if err != ErrAbilityNotFound && err != nil {
			return err
		}
	}

	// Actualizar timestamp
	now := time.Now()
	ability.UpdatedAt = &now

	// Preparar la consulta
	query := `
		UPDATE abilities
		SET name = ?, title = ?, entity_id = ?, entity_type = ?,
		    only_owned = ?, options = ?, scope = ?, updated_at = ?
		WHERE id = ?
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

	args = append(args, now, ability.ID)

	// Ejecutar la consulta
	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}
