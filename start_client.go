package main

import (
	"fmt"
	"time"
	// "flag"
	// "sync"
	"chainr/utils"	
	"math/rand"
	"chainr/types"
	"chainr/client"
)


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func writeRandom() string {
	key := "testkey"
	value := randSeq(10)
	// fmt.Println("Writing key:", key, "value:", value)
	client.SendWrite(key, value)
	return value
}

func readRandom() string{
	key := "testkey"
	// fmt.Println("Reading key:", key)
	return client.SendRead(key)
}


func runTest(batchSize int){
	var writeTimes = []int64{}
	var readTimes = []int64{}
	var randomKeys = []string{}
	master, _ := utils.GetMonitorClient()
	primary, _ := client.GetReadReplica()
	client.WriteToMaster(master, "testkey", "testvalue")
	val := client.ReadFromReplica(primary, "testkey")
	fmt.Println("Read value:", val)
	readKey := "testkey"
	
	for i:=0; i<batchSize; i++{
		randVal := randSeq(100)
		randKey := randSeq(10)
		
		start := time.Now()
		randomKeys = append(randomKeys, randKey)
		client.WriteToMaster(master, randKey, randVal)
		writeTimes = append(writeTimes, time.Since(start).Microseconds())
		// fmt.Println("Request number:", i)
		time.Sleep(time.Millisecond * 1)
	}

	for i:=0; i<batchSize; i++{
		start := time.Now()
		readKey = randomKeys[i]
		val = client.ReadFromReplica(primary, readKey)
		readTimes = append(readTimes, time.Since(start).Microseconds())
		
	}
	
	// fmt.Println("Read times:", readTimes)
	// fmt.Println("Write times:", writeTimes)
	var averageReadTime int64 = 0
	for i:=0; i<len(readTimes); i++{
		averageReadTime += readTimes[i]
	}
	// averageReadTime = averageReadTime/float64(len(readTimes))

	var averageWriteTime int64 = 0	
	for i:=0; i<len(writeTimes); i++{
		averageWriteTime += writeTimes[i]
	}
	averageWriteTime = averageWriteTime/int64(len(writeTimes))

	fmt.Println("Average read time:", float64(averageReadTime)/float64(len(readTimes)*1000000))
	fmt.Println("Average write time:", float64(averageWriteTime)/float64(len(writeTimes)*1000000))
	master.Close()
	primary.Close()


}

func main(){	
	// fmt.Println("Starting Client")
	// master, err := utils.GetMonitorClient()
	// if err != nil{
	// 	fmt.Println("Error getting master client")
	// 	return
	// }
	
	// var lastWriteValue string
	// var lastReadValue string
	
/* 
	var wg sync.WaitGroup
	flag.Parse()

	test:= func (batchSize int){
		runTest(batchSize)
		wg.Done()
	}
	wg.Add(1)
	go test(1000)
	// go test(1000)
	wg.Wait()
	 
	return */

	// readPtr := flag.String("read", "nil", "read key")
	// writePtr := flag.Bool("write", false, "write flag")
	// writeKeyPtr := flag.String("key", " ", "write key")
	// writeValuePtr := flag.String("value", " ", "write value")
	

	// if *readPtr != "nil" {
	// 	fmt.Println("Reading key:", *readPtr)
	// 	client.SendRead(*readPtr)
	// } else if *writePtr != false && *writeKeyPtr != " " && *writeValuePtr != " " {
	// 	client.SendWrite(*writeKeyPtr, *writeValuePtr)
	// } else {
	// 	fmt.Println("Please provide a read or write flag")
	// }
	
	// Write part
/* 	
	if err != nil{
		fmt.Println("Error getting monitor client")
		return
	}
	writeReplica := types.ReplicaClient{}
	master.Call("Master.GetWriteReplica", types.GetReplicasRequest{}, &writeReplica)
	fmt.Println("Write replica:", writeReplica)
	fmt.Println("Writing big string")
	bigBytesArray := make([]byte, 1024*1024*100)
	rand.Read(bigBytesArray)
	resp := utils.PushToReplica(writeReplica, "key2", bigBytesArray)
	fmt.Println("Response:", resp)

 */
	// Read part
	// readReplica := types.ReplicaClient{}
	// master.Call("Master.GetReadReplica", types.GetReplicasRequest{}, &readReplica)
	// fmt.Println("Read replica:", readReplica)
	// // valueArray := make([]byte, 1024*1024*100)
	// resp :=  utils.ReadFromReplica(readReplica, "key2", 0, 1000)
	// fmt.Println("Response:", resp) 


	BenchMarkWrites()

}



	// client, err := utils.GetMonitorClient()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// var resp  = new(types.WriteResponse);
	// err = client.Call("Master.WriteData", types.WriteRequest{"a", "b"}, &resp)
	// if err != nil {
	// 	fmt.Println(err)
	// }


	// fmt.Println("Write response from server:", resp.Success)
	

func BenchMarkWrites(){
	fmt.Println("Starting Benchmark")
	master, err := utils.GetMonitorClient()
	if err != nil{
		fmt.Println("Error getting master client")
		return
	}
	writeReplica := types.ReplicaClient{}
	master.Call("Master.GetWriteReplica", types.GetReplicasRequest{}, &writeReplica)
	fmt.Println("Write replica:", writeReplica)
	fmt.Println("Writing big string")

	// Increase size of array to be written exponentially
	// Start at 100KB, Max size is 500MB
	// Array is a byte array
	// For each size take a sample of 10 writes
	size := 1024*100
	runsPerSize := 10
	var bigBytesArray []byte;
	writeTimes := make(map[int][10]int64)
	for {
		fmt.Println("Iteration size:", size)
		// writeTimes[size] = 
		// create a 10 length array of write times
		nestedTimings := [10]int64{}
		for i:=0; i<runsPerSize; i++{
			key:= randSeq(10)
			// Create array of size
			bigBytesArray = make([]byte, size)
			rand.Read(bigBytesArray)
			start := time.Now()

			_ = utils.PushToReplica(writeReplica, key, bigBytesArray)
			nestedTimings[i] = time.Since(start).Microseconds()
			time.Sleep(time.Millisecond * 1) 
		}
		writeTimes[size] = nestedTimings
		// Print average time for current size
		
		var sum int= 0
		for i:=0; i<runsPerSize; i++{
			sum +=int( writeTimes[size][i])
		}
		var avg float64 = float64(sum)/(10*100000)
		fmt.Println("Average time for size:", avg)
		// writeTimes[size] = nestedTimings

		size *= 2
		if size > 1024*1024*10 {
			break
		}
		
	}
	// bigBytesArray := make([]byte, 1024*1024*100)
	
	// rand.Read(bigBytesArray)
	// resp := utils.PushToReplica(writeReplica, "key2", bigBytesArray)
	// fmt.Println("Response:", resp)
}