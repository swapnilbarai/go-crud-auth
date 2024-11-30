package routes

import (
	"swapnilbarai/go-crud-auth/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/signup/", controllers.SignUpUser)
		authRoute.POST("/signin", controllers.SignInUser)
		authRoute.GET("/refresh", controllers.RefreshToken)
	}

}
