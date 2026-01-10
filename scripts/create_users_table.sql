-- Crear tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    is_change_password BOOLEAN DEFAULT FALSE,
    role_id VARCHAR(36),
    email_verified_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Crear tabla de roles
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Crear tabla de asignación de roles
CREATE TABLE IF NOT EXISTS assigned_roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role_id VARCHAR(36) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    INDEX (entity_id, entity_type)
);

-- Crear tabla de habilidades/permisos
CREATE TABLE IF NOT EXISTS abilities (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Crear tabla de permisos
CREATE TABLE IF NOT EXISTS permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ability_id VARCHAR(36) NOT NULL,
    entity_id VARCHAR(36) NOT NULL,
    entity_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (ability_id) REFERENCES abilities(id),
    INDEX (entity_id, entity_type)
);

-- Insertar rol de administrador
INSERT INTO roles (id, name, title, description)
VALUES ('1', 'administrator', 'Administrador', 'Rol con acceso completo al sistema');

-- Insertar usuario administrador (contraseña: admin123)
-- La contraseña está hasheada con bcrypt
INSERT INTO users (id, name, email, password, is_change_password, role_id, is_active)
VALUES ('1', 'Admin', 'admin@example.com', '$2a$10$JmQwJyFKSf4G9HXgYF.qAOLVH5vKhM5uWGZYmRrA2UD8JDzX1Lp.e', FALSE, '1', TRUE);

-- Asignar rol de administrador al usuario
INSERT INTO assigned_roles (role_id, entity_id, entity_type)
VALUES ('1', '1', 'App\\Models\\User');
