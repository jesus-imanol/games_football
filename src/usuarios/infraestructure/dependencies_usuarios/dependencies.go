package dependenciesusuarios

import (
	"games-football-api/src/core"
	"games-football-api/src/usuarios/application"
	"games-football-api/src/usuarios/infraestructure/adapters"
	"games-football-api/src/usuarios/infraestructure/controllers"
	"games-football-api/src/usuarios/infraestructure/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func InitUsuarios(r *gin.Engine) {
	// Inicializar la conexión a la base de datos
	db, err := core.NewMySQL()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Crear el repositorio
	usuarioRepo := adapters.NewMySQLUsuarioRepository(db)

	// Crear los casos de uso
	loginUseCase := application.NewLoginUseCase(usuarioRepo)
	registerUseCase := application.NewRegisterUseCase(usuarioRepo)

	// Crear los controladores
	loginController := controllers.NewLoginController(loginUseCase)
	registerController := controllers.NewRegisterController(registerUseCase)

	// Registrar las rutas
	routers.UsuariosRouter(r, loginController, registerController)

	log.Println("Módulo de Usuarios inicializado correctamente")
}
