package types
import (
	"sync"
	"log"
	"net/rpc"
)

type WriteRequest struct {
	Key string
	Value string
}

type WriteResponse struct {
	Success bool
}

type ReplicaClient struct {
	Ip string
	Port string
	ReplicaID  string
	// Conn *rpc.Client
}

type GetReplicasReply struct {
	Replicas []ReplicaClient
}

type GetReplicasRequest struct {
	ReqType string
}

type AddClientRequest struct {
	Ip string
	Port string
	Id string
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
}

type UpdateClientsRequest struct {
	Index int
	Clients []ReplicaClient
}

type UpdateClientsResponse struct {
	Success bool
}

type PushToTopicRequest struct {
	TopicName string
	Value	[]byte
}
type PushToTopicResponse struct {
	Success bool
}

type ReadObjectRequest struct {
	TopicName string
	Offset int
	Size int
}

type ReadObjectResponse struct {
	Value []byte
	Dirty bool
	Success bool
}


func(c *Clients) Add(client AddClientRequest){
	log.Println("Adding client: ", client.Ip + ":" + client.Port)
	c.Mu.Lock()
	clientCopy := ReplicaClient{
		Ip: client.Ip,
		Port: client.Port,
		ReplicaID: client.Id,
		// Conn: conn,
	}
	
	

	tobeUpdatedClients := c.Clients
	newClientList := append(tobeUpdatedClients, clientCopy)
	c.Clients = newClientList
	c.Mu.Unlock()
	c.UpdateClients(tobeUpdatedClients) 
	return 
	
}

func(c *Clients) Remove(client ReplicaClient, currentClients []ReplicaClient ) []ReplicaClient{

	
	for i, replica := range currentClients{

		if replica.Ip == client.Ip && replica.Port == client.Port{
			
			currentClients = append(currentClients[:i], currentClients[i+1:]...)
			break
		}
	}
	return currentClients
}

func(c *Clients) UpdateClients(clients []ReplicaClient) {
	// c.Mu.Lock()
	// defer c.Mu.Unlock()
	log.Println("Updating clients" , clients)

	for i:=0; i<len(clients); i++{
		
		rc, err :=rpc.Dial("tcp", clients[i].Ip + ":" + c.Clients[i].Port, )
		if err != nil{
			log.Println("Error in dialing client: ", err)
			return
		}
		log.Println("Dialed client: ", clients[i].Ip + ":" + clients[i].Port)
		resp := UpdateClientsResponse{}
		err = rc.Call("Replica.UpdateClients", &UpdateClientsRequest{Index: i, Clients: c.Clients}, &resp)
		// log.Println("Updated client: ", c.Clients[i].Ip + ":" + c.Clients[i].Port)
		if err != nil{
			log.Println("Error in updating clients", err)

		} else{
			log.Println(resp)
		}

	}
	return 
}