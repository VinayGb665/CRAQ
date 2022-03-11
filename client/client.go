package client
import (
	// "fmt"
	"log"
	"chainr/utils"	
	"chainr/types"
	"net/rpc"

)

func SendWrite(key string, value string) {
	client, err := utils.GetMonitorClient()
	if err != nil {
		log.Println("error:", err)
	}
	var reply types.WriteResponse
	
	err = client.Call("Master.WriteData", types.WriteRequest{key, value}, &reply)
	if err != nil {
		log.Println("error:", err)
	}
	// fmt.Println("Write response from server:", reply.Success)
}


func SendRead(key string) string {
	
	replicaClient, err := GetReadReplica()
	if err != nil {
		log.Fatal("error:", err)
	}
	response := new(types.ClientReadResponse)
	err = replicaClient.Call("Replica.ReadValue", types.ClientReadRequest{key}, &response)
	if err != nil {
		// log.Println("error:", err)
		return ""
	}
	replicaClient.Close()
	// fmt.Println("Read response from server:", response.Value)
	return response.Value
	// var reply types.ClientReadResponse
	
	// err = client.Call("Master.ReadValue", types.ClientReadRequest{key}, &reply)
	// if err != nil {
	// 	log.Fatal("error:", err)
	// }
	// fmt.Println("Read response from server:", reply.Success, reply.Value)
}


func GetReadReplica() (rClient *rpc.Client, connErr error){
	client, err := utils.GetMonitorClient()
	if err != nil {
		log.Println("error:", err)
	}
	primary := types.ReplicaClient{}
	err = client.Call("Master.GetReadReplica", types.GetReplicasRequest{}, &primary)
	if err != nil {
		log.Fatal("error:", err)
	}
	if primary == (types.ReplicaClient{}) {
		log.Println("error:", "No replica found")
		return nil, err
	}
	client.Close()
	replicaClient, err := utils.GetRPCConnection(primary.Ip, primary.Port)
	if err != nil {
		log.Println("Failed to talk to replica:", err)
		return nil, err
	}
	return replicaClient, nil
}


func ReadFromReplica(replicaClient *rpc.Client, key string, ) string {
	response := new(types.ClientReadResponse)
	err := replicaClient.Call("Replica.ReadValue", types.ClientReadRequest{key}, &response)
	if err != nil {
		// log.Println("error:", err)
		return "-"
	}
	return response.Value
}

func WriteToMaster(client *rpc.Client, key string, value string ) bool {
	var reply types.WriteResponse
	
	err := client.Call("Master.WriteData", types.WriteRequest{key, value}, &reply)
	if err != nil {
		log.Println("error:", err)
		return false
	}
	return reply.Success
	// fmt.Println("Write response from server:", reply.Success)
}