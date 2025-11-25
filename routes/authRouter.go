package routes

import (
	controllers "github.com/Aaryansingh20/jwt/controllers"
	"github.com/gin-gonic/gin"
)

// this is when the user has not signed up. userRouter is when the user has logged in
// and has the token.
func AuthRoutes(incomingRoutes *gin.Engine) {
    incomingRoutes.POST("users/signup", controllers.SignUp())
    incomingRoutes.POST("user/login", controllers.Login())
}