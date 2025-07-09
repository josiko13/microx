package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using system environment variables")
	}

	// Configurar conexión a MySQL
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	database := getEnv("DB_NAME", "microx")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4", user, password, host, port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to MySQL:", err)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging MySQL:", err)
	}

	log.Println("✅ Connected to MySQL")

	// Ejecutar migración por comandos separados
	commands := []string{
		"CREATE DATABASE IF NOT EXISTS microx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		"USE microx",
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_username (username),
			INDEX idx_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS tweets (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id),
			INDEX idx_created_at (created_at),
			INDEX idx_user_created (user_id, created_at DESC)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
		`CREATE TABLE IF NOT EXISTS follows (
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
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for i, command := range commands {
		log.Printf("Executing command %d/%d", i+1, len(commands))
		_, err = db.Exec(command)
		if err != nil {
			log.Fatal("Error executing command:", err)
		}
	}

	// Insertar datos de prueba
	log.Println("Inserting test data...")
	insertCommands := []string{
		"INSERT INTO users (username, email) VALUES ('jose', 'jose@example.com'), ('rocio', 'rocio@example.com'), ('yanina', 'yanina@example.com'), ('axel', 'axel@example.com'), ('memo', 'memo@example.com') ON DUPLICATE KEY UPDATE username = username",
		"INSERT INTO tweets (user_id, content) VALUES (1, '¡Hola gente! Este es mi primer tweet en MicroX'), (2, 'Hola aqui probando esta nueva plataforma!!!'), (3, 'Hola buen dia!! Aqui cuidando la planta de mandarina. Abrazo'), (1, 'Luego de mi primer tweet aqui va otro sabiendo que va creciendo.'), (2, 'Chicos! recuerden en llevar sus tareas a clase mañana miercoles. '), (4, 'Me encanta andar en moto por el barrio. Si me ven me saludan.'), (5, 'Trabajar, trabajar... no me queda otra. Abrazo') ON DUPLICATE KEY UPDATE content = content",
		"INSERT INTO follows (follower_id, following_id) VALUES (1, 2), (1, 3), (2, 1), (2, 4), (3, 1), (3, 2), (4, 1), (4, 3), (5, 1), (5, 2) ON DUPLICATE KEY UPDATE follower_id = follower_id",
	}

	for i, command := range insertCommands {
		log.Printf("Inserting data %d/%d", i+1, len(insertCommands))
		_, err = db.Exec(command)
		if err != nil {
			log.Printf("Warning: Error inserting data: %v", err)
		}
	}

	log.Println("✅ Migration completed successfully")
	log.Printf("✅ Database '%s' is ready", database)
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
