package cachenode

import (
	"fmt"
	"github/git-amw/devcache/store"
	"log"
	"net/http"
)

type Node struct {
	address string
	store   *store.Store
}

func NewNode(address string) *Node {
	return &Node{
		address: address,
		store:   store.NewStore(),
	}
}

func (n *Node) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getdata", n.handleGetData)
	mux.HandleFunc("/setdata", n.handleSetData)
	mux.HandleFunc("/health", n.handleHealth)
	log.Printf("Node listening on %s", n.address)
	log.Fatal(http.ListenAndServe(n.address, mux))
}

func (n *Node) handleGetData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, ok := n.store.Get(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	log.Println("Key found", value)
	fmt.Fprint(w, value)
}

func (n *Node) handleSetData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		http.Error(w, "key or value not provided", http.StatusBadRequest)
		return
	}
	n.store.Set(key, value)
	log.Println("Key and value :", key, value)
	fmt.Fprint(w, value)
}

func (n *Node) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
