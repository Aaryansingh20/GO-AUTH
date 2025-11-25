package routes

import (
	"github.com/Aaryansingh20/jwt/controllers"
	"github.com/Aaryansingh20/jwt/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
    // Create a group for protected routes
    userRoutes := incomingRoutes.Group("")

    // Apply middleware ONLY to this group
    userRoutes.Use(middleware.Authenticate())

    // Protected routes
    userRoutes.GET("/users", controllers.GetUsers())
    userRoutes.GET("/users/:user_id", controllers.GetUserById())
}
