package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/your-org/jvairv2/configs"
)

// DBConfig contiene la configuración para la conexión a MySQL
type DBConfig struct {
	Driver          string
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Connection representa una conexión a MySQL
type Connection struct {
	DB *sql.DB
}

// NewConnection crea una nueva conexión a MySQL
func NewConnection(config *configs.DBConfig) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión a la base de datos: %w", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al verificar conexión a la base de datos: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close cierra la conexión a la base de datos
func (c *Connection) Close() error {
	return c.DB.Close()
}

// Ping verifica la conexión a la base de datos
func (c *Connection) Ping() error {
	return c.DB.Ping()
}

// GetDB devuelve la conexión a la base de datos
func (c *Connection) GetDB() *sql.DB {
	return c.DB
}
