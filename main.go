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
	server := servers.SetupServers()

	v1 := router.Group("/v1")

	serverGroup := v1.Group("/server")
	{
		serverGroup.POST("/dequeue", controllers.DequeueRequest(server))
		serverGroup.POST("/add", controllers.AddNewServerController(server))
		serverGroup.POST("/remove", controllers.RemoveServerController(server))
	}

	userGroup := v1.Group("/user")
	{
		userGroup.GET("/", controllers.EnqueueRequest(server))
	}

	return router
}

func main(){
	godotenv.Load()
	router := setupRouter()
	router.Run(fmt.Sprintf(":%v", os.Getenv("PORT")))
}

