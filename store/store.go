package store

import (
	"log"
	"sync"
)

type Store struct {
	cacheMap    map[string]*Node
	capacity    int
	currentSize int
	cacheList   *DoublyLinkedList
	mu          sync.RWMutex
}

func NewStore(capacity int) *Store {
	return &Store{
		cacheMap:    make(map[string]*Node),
		capacity:    capacity,
		currentSize: 0,
		cacheList:   &DoublyLinkedList{},
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.cacheMap[key] == nil {
		node := s.cacheMap[key]
		s.cacheList.moveToFront(node)
		return node.value, true
	}
	return "", false
}

func (s *Store) Set(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cacheMap[key] != nil {
		node := s.cacheMap[key]
		node.value = value
		s.cacheList.moveToFront(node)
	} else {
		newnode := &Node{key: key, value: value, prev: nil, next: nil}
		s.cacheMap[key] = newnode
		s.cacheList.addToFront(newnode)
		s.currentSize++
		if s.currentSize > s.capacity {
			deletnode := s.cacheList.removeTail()
			delete(s.cacheMap, deletnode.key)
			s.currentSize--
		}
	}
	log.Print(s.currentSize)
	return true
}
