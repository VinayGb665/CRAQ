package storage

import (
	"sync"
	"time"
	"log"
)

type LogEntry struct {
	Key, Value string
}
type Storage struct {
	Mu sync.Mutex
	Logs []LogEntry
	Maps map[string]string
}


func(s *Storage) Add(entry LogEntry){
	start := time.Now()
	s.Mu.Lock()
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