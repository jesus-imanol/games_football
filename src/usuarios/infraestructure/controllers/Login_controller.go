package controllers

import (
	"games-football-api/src/usuarios/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	loginUseCase *application.LoginUseCase
}

func NewLoginController(loginUseCase *application.LoginUseCase) *LoginController {
	return &LoginController{
		loginUseCase: loginUseCase,
	}
}

// LoginRequest representa el cuerpo de la petición de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HandleLogin maneja la petición POST de login
func (lc *LoginController) HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"mensaje": "Campos requeridos: username, password",
		})
		return
	}

	usuario, err := lc.loginUseCase.Execute(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"mensaje": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"mensaje": "Login exitoso",
		"usuario": usuario,
	})
}
