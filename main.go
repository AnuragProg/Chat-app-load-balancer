package main

import (
	"fmt"
	"os"
	"load-balancer/controllers"
	"load-balancer/servers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine{
	router := gin.Default()

	// Setting up servers for use
	servers := servers.SetupServers()

	// Adding routes
	v1 := router.Group("/v1")	
	{
		v1.GET("/", controllers.EnqueueRequest(servers))
		v1.POST("/completed", controllers.DequeueRequest(servers))
	}
	return router
}

func main(){
	godotenv.Load()
	router := setupRouter()
	router.Run(fmt.Sprintf(":%v", os.Getenv("PORT")))
}

