package servers

import(
	"net/http/httputil"
	"container/heap"
)

type ServerStatus struct{
	Traffic uint64
	Proxy *httputil.ReverseProxy
	Index int
	Host string
}

type ServerStatuses []*ServerStatus

func (servers ServerStatuses) Len() int {
	return len(servers)
}

// Return true if element i should have a higher priority than element j.
func (servers ServerStatuses) Less(i int, j int) bool {
	return servers[i].Traffic < servers[j].Traffic
}

func (servers ServerStatuses) Swap(i int, j int) {
	servers[i], servers[j] = servers[j], servers[i]
	servers[i].Index = i
	servers[j].Index = j
}

/**
The Push method adds an element to the heap by appending it to the end of the slice
and then calling the up function to fix the heap invariant.
*/
func (servers *ServerStatuses) Push(x any) {
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
func (servers *ServerStatuses) Pop() any {
	old := *servers
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	// item.Index = -1 // for safety
	*servers = old[0:n-1]
	return item
}

func (servers *ServerStatuses) Update(item *ServerStatus, traffic uint64){
	item.Traffic = traffic
	heap.Fix(servers, item.Index)
}

func (servers *ServerStatuses) Seek() any{
	if len(*servers) == 0{
		return nil
	}
	return (*servers)[0]
}
