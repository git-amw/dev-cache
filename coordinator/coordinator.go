package coordinator

import (
	"github/git-amw/devcache/clustorNode"
	"github/git-amw/devcache/hashing"
	"log"
	"net/http"
	"time"
)

type Coordinator struct {
	ring       *hashing.HashRing
	servernode map[string]*clustorNode.ClustorNode
}

func NewCoordinator() *Coordinator {
	return &Coordinator{
		ring:       hashing.NewHashRing(),
		servernode: make(map[string]*clustorNode.ClustorNode),
	}
}

func (c *Coordinator) StartFailureDetection() {
	log.Println("Starting failure detection")
	go func() {
		for {
			for address := range c.servernode {
				resp, err := http.Get("http://" + address + "/health")
				if err != nil || resp.StatusCode != http.StatusOK {
					log.Println("Node " + address + " is down, removing from ring")
					c.RemoveNode(address)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func (c *Coordinator) RemoveNode(address string) {
	delete(c.servernode, address)
	c.ring.RemoveNode(address)
}

func (c *Coordinator) AddNode(address string, capacity int) {
	node := clustorNode.NewNode(address, capacity)
	c.ring.AddNode(address)
	c.servernode[address] = node
	go node.Start()
}

func (c *Coordinator) RouteRequest(key string) string {
	nodeAddress := c.ring.GetNode(key)
	return nodeAddress
}
