package routers

import (
	"games-football-api/src/retas/infraestructure/controllers"

	"github.com/gin-gonic/gin"
)

func RetasRouter(r *gin.Engine, wsController *controllers.WebSocketController) {
	retasGroup := r.Group("/ws")
	{
		retasGroup.GET("/retas", wsController.HandleWebSocket)
	}
}
