package master

import (
	"log"
	"time"
	"sync"
	"chainr/utils"
	"chainr/types"
)





var ReplicasMapping = types.Clients{
	Clients : []types.ReplicaClient{},
}



var wg sync.WaitGroup


type Master struct {
	Mu sync.Mutex
	Clients types.Clients

}



func (t *Master) AddClient(client types.AddClientRequest, reply *types.GetReplicasReply) error {
	log.Println("Adding client: ", client.Ip, client.Port)
	// log.Println("Client added: ", t.Clients)
	t.Mu.Lock()
	defer t.Mu.Unlock()

	
	t.Clients.Add(client)	
	reply.Replicas = t.Clients.Clients
	// t.Mu.Unlock()

	return nil
}

func (t *Master) WriteData(data types.WriteRequest, reply *types.WriteResponse) error {
	for i:=0; i<len(t.Clients.Clients); i++{
		log.Println("Writing to client", i)

	}
	reply.Success = true
	return nil
}

func (t *Master) GetReadReplica(request types.GetReplicasRequest, reply *types.ReplicaClient) error {
	ReplicasMapping.Mu.Lock()
	defer ReplicasMapping.Mu.Unlock()
	lastIndex := len(t.Clients.Clients) - 1	
	reply.Ip = t.Clients.Clients[lastIndex].Ip
	reply.Port = t.Clients.Clients[lastIndex].Port
	return nil
}



func RunMonitor(t *Master) {
	log.Println("Starting monitor")
	
	for {
		t.Mu.Lock()

		MonitorList(t)
		t.Mu.Unlock()
		time.Sleep(30 * time.Second)
	}

	
}


/* func UpdateClients(clients []ReplicaClient) {
	log.Println("Updating clients" , clients)

	for i:=0; i<len(clients); i++{
		
		rc, err :=rpc.Dial("tcp", clients[i].Ip + ":" + c.Clients[i].Port)
		if err != nil{
			log.Println("Error in dialing client: ", err)
			return
		}
		log.Println("Dialed client: ", c.Clients[i].Ip + ":" + c.Clients[i].Port)
		resp := UpdateClientsResponse{}
		err = rc.Call("Replica.UpdateClients", &UpdateClientsRequest{Index: i, Clients: c.Clients}, &resp)
		log.Println("Updated client: ", resp)
		if err != nil{
			log.Println("Error in updating clients", err)

		} else{
			log.Println(resp)
		}

	}
	return 
}
 */

func MonitorList(t *Master) {
	newClients := []types.ReplicaClient{}
	count := 0
	t.Clients.Mu.Lock()
	clientsToCheck := t.Clients.Clients
	
	t.Clients.Mu.Unlock()
	log.Println("Replica List", clientsToCheck, len(clientsToCheck))

	dropIndices := []int{}
	for i:=0; i<len(clientsToCheck); i++{

		_, err := utils.GetRPCConnection(clientsToCheck[i].Ip, clientsToCheck[i].Port)
		if err != nil {
			log.Println("Error dialing: ", err)
			count++
			dropIndices = append(dropIndices, i)
			log.Println("Updated")
			
			continue
		} else{
			newClients = append(newClients, clientsToCheck[i])
		}
	}
	if count >0{
		t.Clients.Mu.Lock()
		t.Clients.Clients = newClients
		t.Clients.UpdateClients(newClients)
		t.Clients.Mu.Unlock()
	}
	



}

func(t *Master) GetWriteReplica(request types.GetReplicasRequest, reply *types.ReplicaClient) error {
	t.Clients.Mu.Lock()
	defer t.Clients.Mu.Unlock()
	reply.Ip = t.Clients.Clients[0].Ip
	reply.Port = t.Clients.Clients[0].Port
	return nil
}