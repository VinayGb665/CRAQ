package storage

import (
	"sync"
	"time"
	"log"
	"strconv"
	"os"
)

type LogEntry struct {
	Id int
	Key, Value string
}
type Storage struct {
	Mu sync.Mutex
	Logs []LogEntry
	Maps map[string]string
	Topics map[string]Topic
	DataDirectory string
	// Consumers 

}

type Consumer struct {
	Topic string
	LastReadEventId int
}

type Topic struct {
	Mu sync.Mutex
	DataDirectory string
	Version int
	// Events []string
	LastWriteId int
	Name string
	Dirty bool
}




func(t *Topic) Push(event *[]byte) {
	t.IncrementVersion()
	t.WriteToFile(event)
	event = nil
	
}

func(t *Topic) Read(offset int, size int) []byte {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	log.Println("Reading from file: ", t.Name)
	file, err := os.Open(t.DataDirectory + "/" + t.Name + "." + strconv.Itoa(t.Version))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()
	file.Seek(int64(offset), 0)
	data := make([]byte, size)
	_, err = file.Read(data)
	if err != nil {
		log.Println(err)
		return nil
	}
	return data
}

func (t *Topic) WriteToFile(event *[]byte) (bool) {
	t.Mu.Lock()
	// file, err := os.OpenFile(t.Name, os.O_APPEND|os.O_WRONLY, 0600)
	file , err := os.OpenFile(t.DataDirectory + "/" + t.Name + "." + strconv.Itoa(t.Version), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return false
	}
	defer file.Close()

	if _, err = file.Write(*event); err != nil {
		log.Println(err)
		return false
	}
	
	t.Dirty = true
	t.Mu.Unlock()
	event = nil
	return true
	
}

func (t *Topic) MarkClean() {
	// log.Println("Marking topic as clean, before lock")
	t.Mu.Lock()

	// log.Println("Marking topic as clean, after lock")
	t.Dirty = false
	t.Mu.Unlock()
	return 
}

func (t *Topic) IncrementVersion(){
	t.Mu.Lock()
	t.Version = t.Version + 1
	t.Mu.Unlock()
}

func(s *Storage) Add(entry LogEntry){
	start := time.Now()
	s.Mu.Lock()
	entry.Id = len(s.Logs)
	s.Logs = append(s.Logs, entry)
	s.Maps[entry.Key] = entry.Value

	s.Mu.Unlock()
	log.Println("Add took: ", time.Since(start))
}



func(s *Storage) Get(key string) (string, bool){
	start := time.Now()
	s.Mu.Lock()
	defer s.Mu.Unlock()
	value, ok := s.Maps[key]
	log.Println("Get took: ", time.Since(start))
	return value, ok
}

func(s *Storage) GetLogs() []LogEntry{
	s.Mu.Lock()
	logs := s.Logs
	s.Mu.Unlock()
	return logs
}



