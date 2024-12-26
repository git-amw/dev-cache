package hashing

import (
	"hash/fnv"
	"sort"
)

type HashRing struct {
	nodeAddressList []string
	hashmap         map[uint32]string
}

func NewHashRing() *HashRing {
	return &HashRing{
		nodeAddressList: []string{},
		hashmap:         make(map[uint32]string),
	}
}

func (h *HashRing) AddNode(address string) {
	hash := hashKey(address)
	h.nodeAddressList = append(h.nodeAddressList, address)
	h.hashmap[hash] = address
	sort.Slice(h.nodeAddressList, func(i, j int) bool {
		return hashKey(h.nodeAddressList[i]) < hashKey(h.nodeAddressList[j])
	})
}

func (h *HashRing) GetNode(key string) string {
	if len(h.nodeAddressList) == 0 {
		return "No node in the list is available"
	}
	hash := hashKey(key)
	for _, address := range h.nodeAddressList {
		if hash <= hashKey(address) {
			return address
		}
	}
	return h.nodeAddressList[0]
}

func (h *HashRing) RemoveNode(address string) {
	hash := hashKey(address)
	delete(h.hashmap, hash)
	for i, n := range h.nodeAddressList {
		if n == address {
			h.nodeAddressList = append(h.nodeAddressList[:i], h.nodeAddressList[i+1:]...)
			break
		}
	}
}

func hashKey(key string) uint32 {
	hash := fnv.New32a()
	hash.Write([]byte(key))
	return hash.Sum32()
}
