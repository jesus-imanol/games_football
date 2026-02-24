package main

import (
	dependenciesretas "games-football-api/src/retas/infraestructure/dependencies_retas"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.Default()

	// Configuración de CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Authorization"},
		MaxAge:           12 * time.Hour,
	}))

	// Ruta raíz - Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "online",
			"message": "API Games Football está en línea ✓",
			"version": "1.0.0",
			"endpoints": gin.H{
				"websocket": "/ws/retas",
			},
		})
	})

	dependenciesretas.InitRetas(r)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
