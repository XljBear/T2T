package storages

import (
	"log"
	"strings"
	"sync"
	"time"
)

type StorageData struct {
	Value  interface{}
	Expire *time.Time
}

type Storage struct {
	Lock          sync.Mutex
	Data          map[string]*StorageData
	ReleaseSignal chan bool
}

func (s *Storage) Set(key string, value interface{}, expire *time.Time) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Data[key] = &StorageData{
		Value:  value,
		Expire: expire,
	}
}
func (s *Storage) Get(key string) (interface{}, bool) {
	v, ok := s.Data[key]
	if !ok {
		return nil, ok
	}
	return v.Value, ok
}
func (s *Storage) Delete(key string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Data, key)
}
func (s *Storage) Exists(key string) bool {
	_, ok := s.Data[key]
	return ok
}
func (s *Storage) DeleteWithPrefix(prefix string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for key := range s.Data {
		if strings.HasPrefix(key, prefix) {
			delete(s.Data, key)
		}
	}
}
func (s *Storage) Clean() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Data = make(map[string]*StorageData)
}
func (s *Storage) Init() {
	s.Clean()
	s.ReleaseSignal = make(chan bool)
	go func() {
		for {
			for k, v := range s.Data {
				if v.Expire != nil && time.Now().After(*v.Expire) {
					s.Delete(k)
				}
			}
			select {
			case <-s.ReleaseSignal:
				s.Clean()
				log.Println("Storage released.")
				break
			case <-time.After(time.Second):
				continue
			}
		}
	}()
}
func (s *Storage) Release() {
	s.ReleaseSignal <- true
}

var StorageInstance *Storage

func Init() {
	if StorageInstance != nil {
		StorageInstance.Release()
	}
	StorageInstance = &Storage{}
	StorageInstance.Init()
	log.Println("Storage initialized.")
}
func Release() {
	if StorageInstance != nil {
		StorageInstance.Release()
	}
	StorageInstance = nil
}
