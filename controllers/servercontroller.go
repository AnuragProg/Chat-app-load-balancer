package controllers

import (
	"load-balancer/servers"

	"github.com/gin-gonic/gin"
)

// Route the request to minimum traffic server
func EnqueueRequest(server *servers.Servers) gin.HandlerFunc{
	return func(c *gin.Context) {

		s := server.Seek().(*servers.ServerStatus)

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
func DequeueRequest(server *servers.Servers) gin.HandlerFunc{
	return func (c *gin.Context)  {

		var host struct{Host string `json:"host"`}

		if err := c.BindJSON(&host); err != nil{
			c.JSON(404, gin.H{"message":err.Error()})
			return
		}

		go server.DequeueRequestForHost(host.Host)
	}
}