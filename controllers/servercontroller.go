package controllers

import (
	"fmt"
	"load-balancer/servers"

	"github.com/gin-gonic/gin"
)

type HostRequestBody struct{
	Host string `json:"host"`
}

// Route the request to minimum traffic server
func EnqueueRequest(server *servers.Server) gin.HandlerFunc{
	return func(c *gin.Context) {

		s := server.GetLeastTrafficServer()

		if s == nil{
			c.JSON(404, "no servers available")
			return
		}

		fmt.Println("Forwarding request to `", s.Host, "` with traffic:`", s.Traffic,"`")

		// Increment traffic counter		
		go server.EnqueueRequest(s)

		// Changing request to route it to server
		s.Proxy.Director(c.Request)

		// Deleting the Host header
		c.Request.Header.Del("Host")

		// Directing the request to server
		s.Proxy.ServeHTTP(c.Writer, c.Request)
	}
}


// request from server that a request has been completed
func DequeueRequest(server *servers.Server) gin.HandlerFunc{
	return func (c *gin.Context)  {
		var reqBody HostRequestBody

		if err := c.BindJSON(&reqBody); err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return 
		}

		go server.DequeueRequestFromHost(reqBody.Host)
	}
}

// TODO Authenticate whether request coming from valid server
func AddNewServerController(server *servers.Server) gin.HandlerFunc{
	return func(c *gin.Context){
		var reqBody HostRequestBody

		if err := c.BindJSON(&reqBody); err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return 
		}
		go server.AddServer(reqBody.Host)
	}
}

// TODO Authenticate whether request coming from valid server
func RemoveServerController(server *servers.Server) gin.HandlerFunc{
	return func(c *gin.Context){
		var reqBody HostRequestBody

		if err := c.BindJSON(&reqBody); err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return 
		}
		go server.RemoveServer(reqBody.Host)
	}
}