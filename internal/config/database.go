package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

// DatabaseConfig contiene la configuración de las bases de datos
type DatabaseConfig struct {
	MySQL *sql.DB
	Redis *redis.Client
}

// NewDatabaseConfig crea las conexiones a MySQL y Redis
func NewDatabaseConfig() (*DatabaseConfig, error) {
	// Configurar MySQL
	mysqlDB, err := connectMySQL()
	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL: %w", err)
	}

	// Configurar Redis
	redisClient, err := connectRedis()
	if err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	return &DatabaseConfig{
		MySQL: mysqlDB,
		Redis: redisClient,
	}, nil
}

// connectMySQL establece la conexión con MySQL
func connectMySQL() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	database := getEnv("DB_NAME", "microx")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, password, host, port, database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to MySQL")
	return db, nil
}

// connectRedis establece la conexión con Redis
func connectRedis() (*redis.Client, error) {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbStr := getEnv("REDIS_DB", "0")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
		PoolSize: 10,
	})

	// Verificar conexión
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to Redis")
	return client, nil
}

// getEnv obtiene una variable de entorno con valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Close cierra todas las conexiones
func (c *DatabaseConfig) Close() error {
	if err := c.MySQL.Close(); err != nil {
		return fmt.Errorf("error closing MySQL: %w", err)
	}

	if err := c.Redis.Close(); err != nil {
		return fmt.Errorf("error closing Redis: %w", err)
	}

	return nil
}
