package replica

import (
	"net"
	"net/rpc"
	"fmt"
	"chainr/utils"
	"chainr/types"
	"chainr/storage"
)




type ReplicaClient struct {
	Ip, ClientId string
}

type GetReplicasReply struct {
	Replicas []ReplicaClient
}

var replicaStorage = storage.Storage{
	Logs: []storage.LogEntry{},
	Maps: make(map[string]string),
}

type Replica int

func(t *Replica) StoreValue(request *types.WriteRequest, resp *types.WriteResponse ) error {
	fmt.Println("Storing value, size of value: ", len(request.Value))
	replicaStorage.Add(storage.LogEntry{Key: request.Key, Value: request.Value})
	*resp = types.WriteResponse{true}
	return nil
}

func(t *Replica) ReadValue(request *types.ClientReadRequest, resp *types.ClientReadResponse) error {
	fmt.Println("Trying to read value: ", request.Key)
	value, ok := replicaStorage.Get(request.Key)
	resp.Value = value
	resp.Success = ok
	return nil
}



func StartReplicaServer() {

	replicaServerPort := utils.GetPortForReplica()
	fmt.Println(" Starting replica Server  at Port: ", replicaServerPort)
	replica := new(Replica)
	rpc.Register(replica)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":" + replicaServerPort)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("Adding replica to monitor")
	utils.AddReplicaToMonitor("localhost", replicaServerPort)
	rpc.Accept(l)
	

}

