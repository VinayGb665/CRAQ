package types
import (
	"sync"
	"log"
)

type WriteRequest struct {
	Key string
	Value string
}

type WriteResponse struct {
	Success bool
}

type ReplicaClient struct {
	Ip, Port string
	Next *ReplicaClient
}

type GetReplicasReply struct {
	Replicas []ReplicaClient
}

type GetReplicasRequest struct {
	ReqType string
}


type ClientReadRequest struct {
	Key string
}
type ClientReadResponse struct {
	Value string
	Success bool
}
type Clients struct { 
	Mu sync.Mutex
	Clients  []ReplicaClient
	HeadReplica *ReplicaClient
	TailReplica *ReplicaClient
}

func(c *Clients) Add(client ReplicaClient){
	c.Mu.Lock()
	
	clientCopy := ReplicaClient{
		Ip: client.Ip,
		Port: client.Port,
		Next: nil,
	}
	
	c.Clients = append(c.Clients, clientCopy)
	if c.HeadReplica == nil { 
		client.Next = nil
		c.HeadReplica = &clientCopy
	}
	if c.TailReplica == nil {
		c.TailReplica = &clientCopy
	} else{
		c.TailReplica.Next = &clientCopy
		c.TailReplica = &clientCopy
	}
	// else{
	// 	c.TailReplica.Next = &clientCopy
	// 	c.TailReplica = &clientCopy
	// }
	c.Mu.Unlock()

}

func(c *Clients) Remove(client ReplicaClient){
	c.Mu.Lock()
	defer c.Mu.Unlock()
	
	var currentHead = c.HeadReplica
	var prevNode *ReplicaClient = nil;
	for {
		log.Println("Current head: ", currentHead.Ip, currentHead.Port)

		if currentHead == nil {
			break
		}
		
		if currentHead.Ip == client.Ip && currentHead.Port == client.Port {
			if prevNode == nil {
				c.HeadReplica = currentHead.Next
			} else {
				prevNode.Next = currentHead.Next
				c.HeadReplica = prevNode
				break
			}
		}
		prevNode = currentHead
		currentHead = currentHead.Next
	}
	return
}
