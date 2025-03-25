package routes

import (
	"example.com/go-mongo-auth/controllers"
	"example.com/go-mongo-auth/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes atașează rutele la router-ul Gin.
func SetupRoutes(router *gin.Engine) {
	// Rute publice
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	// Rute protejate
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", controllers.Protected)
	}
}