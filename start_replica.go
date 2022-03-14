package main
import (
	"chainr/replica"
	"fmt"
	"runtime/debug"
	"flag"
)

func main(){
	fmt.Println("asdad")
	dirPtr := flag.String("datadir", "", "data directory")
	flag.Parse()
	if *dirPtr == ""{
		fmt.Println("No data directory specified")
		return
	}
	debug.SetGCPercent(10)
	replica.StartReplicaServer(*dirPtr)
	

}