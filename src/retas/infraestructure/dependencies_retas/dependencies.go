package dependenciesretas

import (
	"games-football-api/src/core"
	"games-football-api/src/retas/application"
	"games-football-api/src/retas/infraestructure/adapters"
	"games-football-api/src/retas/infraestructure/controllers"
	"games-football-api/src/retas/infraestructure/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func InitRetas(r *gin.Engine) {
	// Inicializar la conexión a la base de datos
	db, err := core.NewMySQL()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Crear el Hub de WebSocket
	hub := adapters.NewHub()
	go hub.Run() // Ejecutar el hub en un goroutine

	// Crear el repositorio
	retaRepo := adapters.NewMySQLRetaRepository(db)

	// Crear los casos de uso
	unirseUseCase := application.NewUnirseRetaUseCase(retaRepo)
	crearRetaUseCase := application.NewCrearRetaUseCase(retaRepo)

	// Crear el controller
	wsController := controllers.NewWebSocketController(hub, unirseUseCase, crearRetaUseCase)

	// Registrar las rutas
	routers.RetasRouter(r, wsController)

	log.Println("Módulo de Retas inicializado correctamente")
}
