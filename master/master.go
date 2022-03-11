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
	HeadReplica : nil,
	TailReplica : nil,
}

var wg sync.WaitGroup


type Master int



func (t *Master) AddClient(client types.ReplicaClient, reply *types.GetReplicasReply,) error {
	log.Println("Adding client: ", client.Ip, client.Port)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		ReplicasMapping.Add(client)
		wg.Done()
		
	}()
	wg.Wait()
	reply.Replicas = ReplicasMapping.Clients
	return nil
}

func (t *Master) WriteData(data types.WriteRequest, reply *types.WriteResponse) error {
	
	var currentWriteReplica = ReplicasMapping.HeadReplica
	successfulWriteCount := 0
	for {
		if currentWriteReplica == nil {
			log.Println("No replicas left to write")
			if successfulWriteCount == 0 {
				reply.Success = false
				return nil
			}
		}
		log.Println("Writing data", currentWriteReplica)
		success, err := utils.SendWriteToReplica(*currentWriteReplica, data)
		if err != nil {
			log.Println("Error while writing to replica", err)
			continue
		}
		if success{
			successfulWriteCount++

		}
		currentWriteReplica = currentWriteReplica.Next
		log.Println("Remaining replicas to write: ", len(ReplicasMapping.Clients) - successfulWriteCount)
		if successfulWriteCount == len(ReplicasMapping.Clients) {
			log.Println("Successfully wrote to all replicas")
			break
		}
		// time.Sleep(5 * time.Second)
	}
	
	reply.Success = true
	return nil
}

func (t *Master) GetReadReplica(request types.GetReplicasRequest, reply *types.ReplicaClient) error {
	ReplicasMapping.Mu.Lock()
	defer ReplicasMapping.Mu.Unlock()
	log.Println("Getting read replica", ReplicasMapping.TailReplica)
	if ReplicasMapping.TailReplica == nil {
		log.Println("No replicas available")
		reply = nil
		return nil
	}
	// reply = ReplicasMapping.TailReplica
	reply.Ip = ReplicasMapping.TailReplica.Ip
	reply.Port = ReplicasMapping.TailReplica.Port
	return nil
}


func RunMonitor(){
	log.Println("Starting monitor")
	
	
	for {
		wg.Add(1)

		go func () {
			MonitorList()
			wg.Done()
		}()
		time.Sleep(5 * time.Second)
		wg.Wait()

	}
}

func MonitorList(){
	var PrevNode *types.ReplicaClient
	
	ReplicasMapping.Mu.Lock()
	defer ReplicasMapping.Mu.Unlock()
	var headReplica = ReplicasMapping.HeadReplica

	for {

		if headReplica != nil {

			isAlive, err := utils.CheckIfReplicaIsAlive(*headReplica)
			if err != nil {
				log.Println("Error while talking to replica", err)
				
			}
			if !isAlive {
				log.Println("Replica is dead, removing it from the list")
				if PrevNode == nil {
					ReplicasMapping.HeadReplica = headReplica.Next
				} else {
					PrevNode.Next = headReplica.Next
				}

				if headReplica.Next == nil {
					ReplicasMapping.TailReplica = PrevNode
				}
			} 
			
			PrevNode = headReplica
			headReplica = headReplica.Next


		} else {
			// log.Println("Reached end of chain")
			break
		}
	}

}

