package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"microx/internal/api"
	"microx/internal/config"
	"microx/internal/middleware"
	"microx/internal/repository"
	"microx/internal/repository/mysql"
	"microx/internal/repository/redis"
	"microx/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using system environment variables")
	}

	// Configurar modo de Gin
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar configuración de base de datos
	dbConfig, err := config.NewDatabaseConfig()
	if err != nil {
		log.Fatal("Failed to connect to databases:", err)
	}
	defer dbConfig.Close()

	// Inicializar repositorios
	userRepo := mysql.NewUserRepository(dbConfig.MySQL)
	tweetRepo := mysql.NewTweetRepository(dbConfig.MySQL)
	followRepo := mysql.NewFollowRepository(dbConfig.MySQL)
	timelineRepo := redis.NewTimelineRepository(dbConfig.Redis)

	// Obtener configuración de la aplicación
	maxTweetLength := getEnvAsInt("MAX_TWEET_LENGTH", 280)

	// Inicializar servicios
	userService := service.NewUserService(userRepo)
	tweetService := service.NewTweetService(tweetRepo, userRepo, timelineRepo, followRepo, maxTweetLength)
	followService := service.NewFollowService(followRepo, userRepo, timelineRepo, tweetRepo)
	timelineService := service.NewTimelineService(timelineRepo, tweetRepo, userRepo, followRepo)

	// Inicializar handlers
	userHandler := api.NewUserHandler(userService)
	tweetHandler := api.NewTweetHandler(tweetService)
	followHandler := api.NewFollowHandler(followService)
	timelineHandler := api.NewTimelineHandler(timelineService)

	// Pre-cargar todos los timelines al iniciar - Esto solo se hace para pruebas.
	// En un entorno de producción, la reconstrucción o precarga del timeline se
	// deberia hacer cuando el usuario se loguea o accede a su timeline, es decir, bajo demanda.
	ctx := context.Background()
	err = timelineService.PreloadAllTimelines(ctx)
	if err != nil {
		log.Printf("Warning: Error preloading timelines: %v", err)
	}

	// Crear router
	r := gin.Default()

	// Configurar rutas
	setupRoutes(r, userHandler, tweetHandler, followHandler, userRepo, timelineHandler, dbConfig)

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Printf("🚀 Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRoutes(r *gin.Engine, userHandler *api.UserHandler, tweetHandler *api.TweetHandler, followHandler *api.FollowHandler, userRepo repository.UserRepository, timelineHandler *api.TimelineHandler, dbConfig *config.DatabaseConfig) {
	// Middleware de autenticación con validación de usuario (para rutas protegidas)
	authWithValidationMiddleware := middleware.AuthWithUserValidationMiddleware(userRepo)

	// Grupo de rutas API
	api := r.Group("/api")
	{
		// Rutas de usuarios (no requieren autenticación para lectura)
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.GET("/:id/stats", userHandler.GetUserStats)
			users.GET("/:id/tweets", tweetHandler.GetUserTweets)
		}

		// Rutas de tweets (requieren autenticación con validación de usuario)
		tweets := api.Group("/tweets")
		tweets.Use(authWithValidationMiddleware)
		{
			tweets.POST("", tweetHandler.CreateTweet)
			tweets.GET("/:id", tweetHandler.GetTweet)
		}

		// Rutas de follow (requieren autenticación con validación de usuario)
		follows := api.Group("/follow")
		follows.Use(authWithValidationMiddleware)
		{
			follows.POST("/:user_id", followHandler.FollowUser)
			follows.DELETE("/:user_id", followHandler.UnfollowUser)
		}

		// Rutas de usuarios para follows (no requieren autenticación para lectura)
		usersFollow := api.Group("/users")
		{
			usersFollow.GET("/:id/followers", followHandler.GetFollowers)
			usersFollow.GET("/:id/following", followHandler.GetFollowing)
		}

		// Rutas de timeline (requieren autenticación)
		timeline := api.Group("/timeline")
		timeline.Use(authWithValidationMiddleware)
		{
			timeline.GET("", timelineHandler.GetTimeline)
			timeline.POST("/refresh", timelineHandler.RefreshTimeline)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "microx",
			"version": "1.0.0",
		})
	})

	// Debug endpoint para inspeccionar Redis
	r.GET("/debug/redis", func(c *gin.Context) {
		if dbConfig.Redis == nil {
			c.JSON(500, gin.H{"error": "Redis not available"})
			return
		}

		ctx := context.Background()

		// Obtener todas las claves
		keys, err := dbConfig.Redis.Keys(ctx, "*").Result()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error getting keys: " + err.Error()})
			return
		}

		result := make(map[string]interface{})

		for _, key := range keys {
			keyType, err := dbConfig.Redis.Type(ctx, key).Result()
			if err != nil {
				continue
			}

			switch keyType {
			case "zset":
				// Para timelines (sorted sets)
				members, err := dbConfig.Redis.ZRange(ctx, key, 0, -1).Result()
				if err == nil {
					result[key] = gin.H{
						"type":    "zset",
						"count":   len(members),
						"members": members,
					}
				}
			case "string":
				// Para valores simples
				value, err := dbConfig.Redis.Get(ctx, key).Result()
				if err == nil {
					result[key] = gin.H{
						"type":  "string",
						"value": value,
					}
				}
			case "hash":
				// Para hashes
				values, err := dbConfig.Redis.HGetAll(ctx, key).Result()
				if err == nil {
					result[key] = gin.H{
						"type":   "hash",
						"values": values,
					}
				}
			default:
				result[key] = gin.H{
					"type": keyType,
					"note": "type not fully inspected",
				}
			}
		}

		c.JSON(200, gin.H{
			"total_keys": len(keys),
			"keys":       result,
		})
	})

	// Debug endpoint para limpiar Redis
	r.DELETE("/debug/redis", func(c *gin.Context) {
		if dbConfig.Redis == nil {
			c.JSON(500, gin.H{"error": "Redis not available"})
			return
		}

		ctx := context.Background()

		// Limpiar todas las claves
		err := dbConfig.Redis.FlushAll(ctx).Err()
		if err != nil {
			c.JSON(500, gin.H{"error": "Error clearing Redis: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message": "Redis cleared successfully",
			"status":  "success",
		})
	})

	// Ruta de información de la API
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "MicroX API - Plataforma de Microblogging",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":   "/health",
				"tweets":   "/api/tweets",
				"follows":  "/api/follow",
				"timeline": "/api/timeline",
			},
		})
	})
}

// getEnvAsInt obtiene una variable de entorno como entero con valor por defecto
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
