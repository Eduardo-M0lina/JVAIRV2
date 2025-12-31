package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config almacena toda la configuración de la aplicación
type Config struct {
	Server ServerConfig
	DB     DBConfig
}

// ServerConfig almacena la configuración del servidor HTTP
type ServerConfig struct {
	Port         string
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DBConfig almacena la configuración de la base de datos
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

// LoadConfig carga la configuración desde el archivo app.env
func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo de configuración: %w", err)
	}

	var config Config

	// Configuración del servidor
	config.Server.Port = viper.GetString("SERVER_PORT")
	config.Server.Timeout = viper.GetDuration("SERVER_TIMEOUT")
	config.Server.ReadTimeout = viper.GetDuration("SERVER_READ_TIMEOUT")
	config.Server.WriteTimeout = viper.GetDuration("SERVER_WRITE_TIMEOUT")

	// Configuración de la base de datos
	config.DB.Driver = viper.GetString("DB_DRIVER")
	config.DB.Host = viper.GetString("DB_HOST")
	config.DB.Port = viper.GetString("DB_PORT")
	config.DB.User = viper.GetString("DB_USER")
	config.DB.Password = viper.GetString("DB_PASSWORD")
	config.DB.DBName = viper.GetString("DB_NAME")
	config.DB.MaxOpenConns = viper.GetInt("DB_MAX_OPEN_CONNS")
	config.DB.MaxIdleConns = viper.GetInt("DB_MAX_IDLE_CONNS")
	config.DB.ConnMaxLifetime = viper.GetDuration("DB_CONN_MAX_LIFETIME")

	return &config, nil
}
