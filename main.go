package main

import (
	"swapnilbarai/go-crud-auth/config"
	"swapnilbarai/go-crud-auth/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Initialise()
	router := gin.Default()
	routes.RegisterAuthRoutes(router)
	routes.RegisterUserRoutes(router)
	router.Run()

}
