package replica

import (
	"net"
	"net/rpc"
	"fmt"
	"sync"
	"runtime/debug"
	"chainr/utils"
	"chainr/types"
	"log"
	"runtime"
	"chainr/storage"
	"os"
)




type ReplicaClient struct {
	Ip, ClientId string
}

type GetReplicasReply struct {
	Replicas []ReplicaClient
}

// var replicaStorage = storage.Storage{
// 	Logs: []storage.LogEntry{},
// 	Maps: make(map[string]string),
// }

type Replica struct {
	Mu sync.Mutex
	ReplicaID string
	index int
	Clients []types.ReplicaClient
	DataDirectory string
	// storage *storage.Storage
	Topics map[string]storage.Topic

}

func(r *Replica) UpdateClients(request *types.UpdateClientsRequest, resp *types.UpdateClientsResponse) error {
	log.Println("Updating clients", request.Clients)
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Clients = request.Clients
	resp.Success = true
	return nil
}
/* 
func(r *Replica) StoreValue(request *types.WriteRequest, resp *types.WriteResponse ) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.storage.Add(storage.LogEntry{Key: request.Key, Value: request.Value})
	fmt.Println("Storing value, size of value: ", len(request.Value))
	// replicaStorage.Add(storage.LogEntry{Key: request.Key, Value: request.Value})
	*resp = types.WriteResponse{true}
	return nil
}

func(t *Replica) ReadValue(request *types.ClientReadRequest, resp *types.ClientReadResponse) error {
	fmt.Println("Trying to read value: ", request.Key)
	// value, ok := replicaStorage.Get(request.Key)
	value, ok := t.storage.Get(request.Key)
	resp.Value = value
	resp.Success = ok
	return nil
}
 */
func(t *Replica) PushToTopic( request *types.PushToTopicRequest, resp *types.PushToTopicResponse) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	log.Println("Recieved a push request for topic: ", request.TopicName)
	topic, ok := t.Topics[request.TopicName]
	
	if !ok {
		log.Println("Topic does not exist, creating new topic")
		t.Topics[request.TopicName] = storage.Topic{
			Name: request.TopicName,
			DataDirectory: t.DataDirectory,
			LastWriteId: 0,
			Version: 0,
			Dirty: false,
		}
		topic = t.Topics[request.TopicName]
	}

	
	
	// topic.Push(request)
	// topic.IncrementVersion()
	// topic.WriteToFile(&request.Value)
	debug.FreeOSMemory()
	
	log.Println("Finished writing in my replica", t.Clients)
	value := request.Value
	i :=0
	for i=0; i<len(t.Clients); i++ {
		if t.Clients[i].ReplicaID == t.ReplicaID {
			if i+1 < len(t.Clients) {
				log.Println("Pushing to client: ", t.Clients[i+1].Port, t.Clients[i].Port, )
				success := utils.PushToReplica(t.Clients[i+1], request.TopicName, value)
				topic.MarkClean()
				t.Topics[request.TopicName] = topic
				// log.Println("Success: ", success, t.Topics)
				resp.Success = success
				value = nil
				// topic = nil
				request = nil
				runtime.GC()
				// request.Value = nil

				log.Println("Finished running GC")
				return nil
			} else {
				topic.MarkClean()
				t.Topics[request.TopicName] = topic
				request.Value = nil

				log.Println("Reached tail of replicas, backpropagating")
				resp.Success = true
				request = nil
				return nil
			}
			
		}
	}	
	// resp.Success = true
	return nil
}

func(t *Replica) ReadObject( request *types.ReadObjectRequest, resp *types.ReadObjectResponse) error {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	fmt.Println("Trying to read value: ", request.TopicName)
	// Check if topic exists
	topic, ok := t.Topics[request.TopicName]


	if !ok {
		resp.Success = false
		return nil
	}
	topic.Mu.Lock()
	log.Println("Version: ", topic.Version)
	topic.Mu.Unlock()
	readResult := topic.Read(request.Offset, request.Size)
	resp.Success = true
	resp.Value = readResult
	topic.Mu.Lock()
	resp.Dirty = topic.Dirty
	topic.Mu.Unlock()
	return nil
	// value, ok := replicaStorage.Get(request.Key)
	// value, ok := t.Topics[request.TopicName].Get(request.Key)
	// resp.Value = value
	// resp.Success = ok
	return nil
}


func StartReplicaServer(dataDirectory string) {

	if err := os.MkdirAll(dataDirectory, 0777); err != nil { 
		log.Println(err)
	}
	log.Println("Created data directory: ", dataDirectory)
	replicaServerPort := utils.GetPortForReplica()
	replicaId := utils.RandomGen(20)
	fmt.Println(" Starting replica Server  at Port: ", replicaServerPort, " with replica id: ", replicaId)
	replica := Replica{
		ReplicaID: replicaId,
		index: 0,
		Clients: []types.ReplicaClient{},
		Topics: make(map[string]storage.Topic),
		DataDirectory: dataDirectory,

	}
	rpc.Register(&replica)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":" + replicaServerPort)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("Adding replica to monitor")
	clients, err:= utils.AddReplicaToMonitor("localhost", replicaServerPort, replicaId)
	replica.Clients = clients.Replicas
	// replica.Mu.Lock()
	
	log.Println("Error adding replica to monitor: ", err)
	rpc.Accept(l)




	

}

