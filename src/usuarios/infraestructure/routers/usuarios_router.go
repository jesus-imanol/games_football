package routers

import (
	"games-football-api/src/usuarios/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func UsuariosRouter(r *gin.Engine, loginController *controllers.LoginController, registerController *controllers.RegisterController) {
	usuariosGroup := r.Group("/api/usuarios")
	{
		usuariosGroup.POST("/login", loginController.HandleLogin)
		usuariosGroup.POST("/register", registerController.HandleRegister)
	}
}
