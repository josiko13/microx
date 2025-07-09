-- Migración inicial para MicroX
-- Crea todas las tablas necesarias para la plataforma de microblogging

-- Crear base de datos si no existe
CREATE DATABASE IF NOT EXISTS microx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE microx;

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de tweets
CREATE TABLE IF NOT EXISTS tweets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_user_created (user_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla de follows (relaciones de seguimiento)
CREATE TABLE IF NOT EXISTS follows (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_follow (follower_id, following_id),
    INDEX idx_follower_id (follower_id),
    INDEX idx_following_id (following_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insertar algunos usuarios de prueba
INSERT INTO users (username, email) VALUES
('jose', 'jose@example.com'),
('rocio', 'rocio@example.com'),
('yanina', 'yanina@example.com'),
('axel', 'axel@example.com'),
('memo', 'memo@example.com')
ON DUPLICATE KEY UPDATE username = username;

-- Insertar algunos tweets de prueba
INSERT INTO tweets (user_id, content) VALUES
(1, '¡Hola gente! Este es mi primer tweet en MicroX'),
(2, 'Hola aqui probando esta nueva plataforma!!!'),
(3, 'Hola buen dia!! Aqui cuidando la planta de mandarina. Abrazo'),
(1, 'Luego de mi primer tweet aqui va otro sabiendo que va creciendo.'),
(2, 'Chicos! recuerden en llevar sus tareas a clase mañana miercoles.'),
(4, 'Me encanta andar en moto por el barrio. Si me ven me saludan.'),
(5, 'Trabajar, trabajar... no me queda otra. Abrazo')
ON DUPLICATE KEY UPDATE content = content;

-- Insertar algunas relaciones de follow de prueba
INSERT INTO follows (follower_id, following_id) VALUES
(1, 2), -- Jose sigue a Rocio
(1, 3), -- Jose sigue a Yanina
(2, 1), -- Rocio sigue a Jose
(2, 4), -- Rocio sigue a Axel
(3, 1), -- Yanina sigue a Jose
(3, 2), -- Yanina sigue a Rocio
(4, 1), -- Axel sigue a Jose
(4, 3), -- Axel sigue a Yanina
(5, 1), -- Memo sigue a Jose
(5, 2)  -- Memo sigue a Rocio
ON DUPLICATE KEY UPDATE follower_id = follower_id; 