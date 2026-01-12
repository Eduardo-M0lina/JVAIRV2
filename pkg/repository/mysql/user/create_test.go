package user

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/your-org/jvairv2/pkg/domain/user"
)

func TestCreate_Success(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &user.User{
		Name:     "New Test User",
		Email:    "newuser@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar la expectativa para la consulta INSERT
	// Ya no se espera la consulta GetByEmail porque la validación está en el caso de uso
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO users (name, email, password, role_id,
		                  email_verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)).WithArgs(
		testUser.Name, testUser.Email, sqlmock.AnyArg(), sqlmock.AnyArg(),
		nil, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(123, 1))

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), testUser)

	// Verificar que no haya errores
	assert.NoError(t, err)
	assert.Equal(t, int64(123), testUser.ID)
	assert.NotNil(t, testUser.CreatedAt)
	assert.NotNil(t, testUser.UpdatedAt)

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestCreate_DuplicateEmail y TestCreate_GetByEmailError fueron eliminados
// porque la validación de email duplicado ahora se realiza en el caso de uso,
// no en el repositorio. El repositorio solo se encarga de insertar datos.

func TestCreate_InsertError(t *testing.T) {
	// Configurar el mock de la base de datos
	db, mock, repo := setupMockDB(t)
	defer func() { _ = db.Close() }()

	// Datos de prueba
	roleID := "1"
	testUser := &user.User{
		Name:     "Insert Error User",
		Email:    "inserterror@example.com",
		Password: "password123",
		RoleID:   &roleID,
	}

	// Configurar la expectativa para la consulta INSERT (que debe devolver un error)
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO users (name, email, password, role_id,
		                  email_verified_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)).WithArgs(
		testUser.Name, testUser.Email, sqlmock.AnyArg(), sqlmock.AnyArg(),
		nil, sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnError(sql.ErrConnDone)

	// Ejecutar la función que estamos probando
	err := repo.Create(context.Background(), testUser)

	// Verificar que haya un error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error en la inserción de usuario")

	// Verificar que todas las expectativas se cumplieron
	assert.NoError(t, mock.ExpectationsWereMet())
}
