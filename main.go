package main

import (
	"fmt"
	"net"
	"net/rpc"
	"chainr/master" 
	// "chainr/utils"

)

func monitorReplicas(){
	for {
		fmt.Println("Monitoring")
		// check if replica is alive
		// if not, then remove it from the list
		// if it is alive, then do nothing
		// if it is dead, then add it to the list
	}
}


func main(){

	m := new(master.Master)
	rpc.Register(m)
	rpc.HandleHTTP()
	
	l, e := net.Listen("tcp", ":1234")
	if e != nil{
		fmt.Println(e)
	}
	// go utils.RunMonitor(&master.ReplicaClients)
	go master.RunMonitor()
	rpc.Accept(l)

	
	// fmt.Println(master.Add(1,2))
}