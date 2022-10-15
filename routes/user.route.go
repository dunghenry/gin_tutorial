package routes

import (
	"trandung/gin_tutorial/controllers"
	"trandung/gin_tutorial/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/auth")
	users.GET("/users", controllers.GetUsers)
	users.POST("/register", controllers.Register)
	users.POST("/login", controllers.Login)
	users.GET("/users/:id", middleware.VerifyToken(), controllers.GetUserById)
}
