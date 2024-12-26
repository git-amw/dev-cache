package coordinator

import (
	"github/git-amw/devcache/cacheNode"
	"github/git-amw/devcache/hashing"
	"log"
	"net/http"
)

type Coordinator struct {
	ring       *hashing.HashRing
	servernode map[string]*cachenode.Node
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		ring:       hashing.NewHashRing(),
		servernode: make(map[string]*cachenode.Node),
	}
}

func (c *Coordinator) StartFailureDetection() {
	go func() {
		for {
			for address := range c.servernode {
				resp, err := http.Get("http://" + address + "/health")
				if err != nil || resp.StatusCode != http.StatusOK {
					log.Println("Node " + address + " is down, removing from ring")
					c.RemoveNode(address)
				}
			}
		}
	}()
}

func (c *Coordinator) RemoveNode(address string) {
	delete(c.servernode, address)
	c.ring.RemoveNode(address)
}

func (c *Coordinator) AddNode(address string) {
	node := cachenode.NewNode(address)
	c.ring.AddNode(address)
	c.servernode[address] = node
	go node.Start()
}

func (c *Coordinator) RouteRequest(key string) string {
	nodeAddress := c.ring.GetNode(key)
	return nodeAddress
}
