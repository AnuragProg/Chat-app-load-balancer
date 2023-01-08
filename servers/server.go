package servers

import (
	"fmt"
	"sync"
	"net/url"
	"net/http"
	"container/heap"
	"net/http/httputil"
)

// For creating locks for the heap
var serverMutex = sync.RWMutex{}

type Server struct{
	Servers *ServerStatuses
}

func createNewServer(host string)*ServerStatus{
	target, err := url.Parse(host)
	if err != nil{
		fmt.Println(err.Error())
		return nil
	}
	return &ServerStatus{
		Traffic: 0,
		Proxy: &httputil.ReverseProxy{
			Director: func (req *http.Request)  {
				req.URL.Path = "/"	
				req.URL.Host = target.Host
				req.Host = target.Host
				req.URL.Scheme = target.Scheme
			},
		},
		Index: -1, // will be handled by the heap
		Host: host,
	}
}

func (server *Server)findServer(host string) *ServerStatus{
	for _, s := range *server.Servers{
		if s.Host == host{
			return s
		}
	}
	return nil
}

func (server *Server) GetLeastTrafficServer() *ServerStatus{
	serverMutex.RLock()
	defer serverMutex.RUnlock()

	return server.Servers.Seek().(*ServerStatus)
}

// host format => http://something.com | https://something.com
func (server *Server) AddServer(host string){
	serverMutex.Lock()
	defer serverMutex.Unlock()

	fmt.Println("Requested: Add server: ", host)

	newServer := createNewServer(host)
	heap.Push(server.Servers, newServer)
	fmt.Println("Added server: ", host)
}

func (server *Server) RemoveServer(host string){
	serverMutex.Lock()
	defer serverMutex.Unlock()

	fmt.Println("Requested: Remove server: ", host)

	s := server.findServer(host)
	if s==nil{return}

	heap.Remove(server.Servers, s.Index)	
	fmt.Println("Removed server: ", host)
}

// Used when server status is retrieved with least traffiic
// So that the pointer can be passed directly for updating the traffic
func (server *Server) EnqueueRequest(s *ServerStatus){
	serverMutex.Lock()
	defer serverMutex.Unlock()

	server.Servers.Update(s, s.Traffic+1)
	
	fmt.Println("Enqueued request")
	fmt.Println("Current Servers Status: ")
	server.Servers.printStatuses()
}

func (server *Server) EnqueueRequestToHost(host string){
	serverMutex.Lock()
	defer serverMutex.Unlock()

	s := server.findServer(host)
	if s == nil{
		return
	}

	server.Servers.Update(s, s.Traffic+1)
	fmt.Println("Enqueued request")
	fmt.Println("Current Servers Status: ")
	server.Servers.printStatuses()
}

func (server *Server) DequeueRequestFromHost(host string){
	serverMutex.Lock()
	defer serverMutex.Unlock()

	s := server.findServer(host)
	if s == nil{
		return
	}

	server.Servers.Update(s, s.Traffic-1)
	fmt.Println("Dequeued request")
	fmt.Println("Current Servers Status: ")
	server.Servers.printStatuses()
}

func (servers *ServerStatuses)printStatuses(){
	for _, s:= range *servers{
		fmt.Println(s.Host, " has traffic = ", s.Traffic)
	}
}

func SetupServers() *Server{
	serverHeap := &ServerStatuses{}

	// Initializing the heap
	heap.Init(serverHeap)

	return &Server{Servers: serverHeap}
}