package controllers

import (
	"fmt"
	"load-balancer/servers"

	"github.com/gin-gonic/gin"
)

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
		host := c.Request.RemoteAddr
		fmt.Println(host, " => Completed a request")
		go server.DequeueRequestFromHost(host)
	}
}

// TODO Authenticate whether request coming from valid server
func AddNewServerController(server *servers.Server) gin.HandlerFunc{
	return func(c *gin.Context){
		host := c.Request.RemoteAddr

		fmt.Println("RemoteAddr => ", host)
		fmt.Println("RequestURI => ", c.Request.RequestURI)
		fmt.Println("RemoteIP   => ", c.RemoteIP())
		go server.AddServer(host)
	}
}

// TODO Authenticate whether request coming from valid server
func RemoveServerController(server *servers.Server) gin.HandlerFunc{
	return func(c *gin.Context){
		host := c.Request.RemoteAddr
		go server.RemoveServer(host)
	}
}