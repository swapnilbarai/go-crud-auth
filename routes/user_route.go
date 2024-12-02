package routes

import (
	"swapnilbarai/go-crud-auth/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/user")

	userRoutes.Use(controllers.AuthenticationMiddleware)
	{
		userRoutes.GET("/details/:username", controllers.GetUserDetails)
		userRoutes.GET("/revoke/:tokenID", controllers.RevokeUserToken)
		userRoutes.GET("/active/tokens", controllers.ShowActiveTokens)
	}

}
