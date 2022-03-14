package utils

import (
	"fmt"
	"net/rpc"
	"strconv"
	"net"
	"log"
	"time"
	"math/rand"
	"chainr/types"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomGen(n int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[r1.Intn(len(letters))]
    }
    return string(b)
}

func GetRPCConnection(ip string, port string) (client *rpc.Client, err error) {
	client, err = rpc.Dial("tcp", ip+":"+port)
	if err != nil {
		log.Println("Error while talking to replica", ip, err)
		return client, err
	}
	return client, err
}

func GetMonitorClient() (client *rpc.Client, err error) {
	return GetRPCConnection("localhost", "1234")
}

func CheckIfReplicaIsAlive(replicaClient types.ReplicaClient) (bool, error) {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", replicaClient.Ip+":"+replicaClient.Port, timeout)
	if err != nil {
		// fmt.Println("Error while talking to replica", replicaClient.Ip, err)
		return false, err
	}
	if conn != nil {
		defer conn.Close()
		// fmt.Println("Replica is alive")
		return true, nil
	}
	return false, nil
}

func CheckPortAvailable(port int) bool {
	timeout := 10 * time.Second
	ln, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), timeout)
	if err != nil {
		return true
	}
	if ln != nil {
		ln.Close()
		return false
	}
	return true
}

func GetPortForReplica() string {
	var ranges = [100]int{1235, 1236, 1237, 1238, 1239, 1240, 1241, 1242, 1243, 1244, 1245, 1246}
	var replicaServerPort string

	for i := 0; i < len(ranges); i++ {
		fmt.Println("Checking if port is available", ranges[i])
		if CheckPortAvailable(ranges[i]) {
			fmt.Println("Available")
			replicaServerPort = strconv.Itoa(ranges[i])
			break
		}
	}
	return replicaServerPort
}


func RunMonitor(replicas *[]types.ReplicaClient){
	fmt.Println("Starting monitor")
	for {
		replicasLocal := *replicas

		fmt.Println("Replica List", replicasLocal, len(replicasLocal))
		for i := 0; i < len(replicasLocal); i++ {
			fmt.Println("Checking if replica is alive", replicasLocal[i])
			isAlive, err := CheckIfReplicaIsAlive(replicasLocal[i])
			if err != nil {
				fmt.Println("Error while talking to replica", err)
				
			}
			if(isAlive){
				fmt.Println("Replica is alive")
			}else{
				fmt.Println("Replica is dead")
			}
		}
		fmt.Println("Sleeping for 10 seconds")
		time.Sleep(5 * time.Second)
	}
}

func AddReplicaToMonitor(ip string, port string, id string) (clients *types.GetReplicasReply, err error) {
	client, err := GetMonitorClient()
	if err != nil {
		return nil, err	
	}
	var resp = new(types.GetReplicasReply)
	err = client.Call("Master.AddClient", types.AddClientRequest{ip, port, id}, &resp)
	if err != nil {
		fmt.Println("Error while talking to the master monitor:", err)
		return resp, err
	}
	client.Close()
	fmt.Println("Replicas:", resp.Replicas)
	return resp, err
}


func SendWriteToReplica(replicaClient types.ReplicaClient, request types.WriteRequest) (bool, error) {
	client, err := rpc.Dial("tcp", replicaClient.Ip+":"+replicaClient.Port)
	if err != nil {
		fmt.Println("Error while talking to replica", replicaClient.Ip, err)
		return false, err
	}
	var resp = new(types.WriteResponse)
	err = client.Call("Replica.StoreValue", &request, &resp)
	if err != nil {
		fmt.Println("Error while talking to replica", replicaClient.Ip, err)
		return false, err
	}
	client.Close()
	return resp.Success, nil
}


func PushToReplica(writeReplica types.ReplicaClient, topic string, value []byte) bool {
	writeReplicaRPC, err := GetRPCConnection(writeReplica.Ip, writeReplica.Port)
	pushEventRequest := types.PushToTopicRequest{
		TopicName: topic,
		Value: value,
	}
	resp := types.PushToTopicResponse{}
	err = writeReplicaRPC.Call("Replica.PushToTopic", pushEventRequest, &resp)
	if err != nil{
		fmt.Println("Error pushing to topic", err)
		return false
	}
	return resp.Success
}

func ReadFromReplica(readReplica types.ReplicaClient, topic string , offset int , size int) *types.ReadObjectResponse {
	readReplicaRPC, err := GetRPCConnection(readReplica.Ip, readReplica.Port)
	readObjectReq := types.ReadObjectRequest{
		TopicName: topic,
		Offset: offset,
		Size: size,
	}
	resp := types.ReadObjectResponse{}
	err = readReplicaRPC.Call("Replica.ReadObject", readObjectReq, &resp)
	if err != nil{
		fmt.Println("Error pushing to topic", err)
		return nil
	}
	return &resp
}