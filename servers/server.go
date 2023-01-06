package servers

import(
	"fmt"
	"sync"
	"net/url"
	"net/http/httputil"
	"container/heap"
)

type ServerStatus struct{
	Traffic uint64
	Proxy *httputil.ReverseProxy
	Index int
	Mutex sync.Mutex
	Host string
}

type Servers []*ServerStatus

var serverMutex sync.Mutex = sync.Mutex{}

func (servers Servers) Len() int {
	return len(servers)
}

// Return true if element i should have a higher priority than element j.
func (servers Servers) Less(i int, j int) bool {
	return servers[i].Traffic < servers[j].Traffic
}

func (servers Servers) Swap(i int, j int) {
	servers[i].Mutex.Lock()
	servers[j].Mutex.Lock()
	defer func(){
		servers[i].Mutex.Unlock()
		servers[j].Mutex.Unlock()
	}()

	servers[i], servers[j] = servers[j], servers[i]
	servers[i].Index = i
	servers[j].Index = j

}

/**
The Push method adds an element to the heap by appending it to the end of the slice
and then calling the up function to fix the heap invariant.
*/
func (servers *Servers) Push(x any) {
	n := len(*servers)
	item := x.(*ServerStatus)
	item.Index = n
	*servers = append(*servers, item)
}

/**
The Pop method removes and returns the element 
with the highest priority by swapping it 
with the last element in the slice and then calling the down function to fix the heap invariant.
*/
func (servers *Servers) Pop() any {
	old := *servers
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	// item.Index = -1 // for safety
	*servers = old[0:n-1]
	return item
}

func (servers *Servers) update(item *ServerStatus, traffic uint64){
	item.Traffic = traffic
	heap.Fix(servers, item.Index)
}

func (servers *Servers) Seek() any{
	if len(*servers) == 0{
		return nil
	}
	return (*servers)[0]
}


func CreateNewServer(host string) *ServerStatus{
	return &ServerStatus{
			Traffic: 0,
			Proxy: httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",                            
				Host: host,         
			}),                                         
			Index: 0, // will be updated while pushing
			Mutex: sync.Mutex{},
			Host: host,
		}
}

// pushing new server to the server
func (servers *Servers)PushNewServer(host string){
	serverMutex.Lock()
	defer serverMutex.Unlock()

	heap.Push(servers, CreateNewServer(host))
}

// poping the server out of the server
// server with least traffic
func (servers *Servers)PopServer() *ServerStatus{
	serverMutex.Lock()
	defer serverMutex.Unlock()

	return heap.Pop(servers).(*ServerStatus)
}

// find server in the heap
func (servers *Servers) FindServer(host string) *ServerStatus{
	for _, s := range *servers{
		if s.Host == host{
			return s
		}
	}
	return nil
}

// enqueue request for given server
func (servers *Servers) EnqueueRequest(server *ServerStatus){
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	servers.update(server, server.Traffic+1)
}

// request enqueued to server for given host
func (servers *Servers) EnqueueRequestForHost(host string){

	server := servers.FindServer(host)
	if server == nil{
		fmt.Println("Unable to find server with hostname ", host)
		return
	}

	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	servers.update(server, server.Traffic+1)
}


// request dequeued from the server for given host
func (servers *Servers) DequeueRequestForHost(host string){

	server := servers.FindServer(host)
	if server == nil{
		fmt.Println("Unable to find server with hostname ", host)
		return
	}
	server.Mutex.Lock()
	defer server.Mutex.Unlock()

	servers.update(server, server.Traffic-1)
}


func SetupServers() *Servers{
	serverHeap := Servers{
		&ServerStatus{
			Traffic : 10,
			Proxy : httputil.NewSingleHostReverseProxy(&url.URL{
					Scheme: "http",
					Host: "localhost:5000",                    
				}),
			Index: 0, 
			Mutex: sync.Mutex{},
			Host: "localhost:5000",
		},
		&ServerStatus{
			Traffic: 0,
			Proxy: httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",                            
				Host: "localhost:4000",         
			}),                                         
			Index: 1,
			Mutex: sync.Mutex{},
			Host: "localhost:4000",
		},
	}
	heap.Init(&serverHeap)

	return &serverHeap
}