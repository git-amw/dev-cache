package clustorNode

import (
	"fmt"
	"github/git-amw/devcache/store"
	"log"
	"net/http"
)

type ClustorNode struct {
	address string
	store   *store.Store
}

func NewNode(address string, capacity int) *ClustorNode {
	return &ClustorNode{
		address: address,
		store:   store.NewStore(capacity),
	}
}

func (n *ClustorNode) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/getdata", n.handleGetData)
	mux.HandleFunc("/setdata", n.handleSetData)
	mux.HandleFunc("/health", n.handleHealth)
	log.Printf("Node listening on %s", n.address)
	log.Fatal(http.ListenAndServe(n.address, mux))
}

func (n *ClustorNode) handleGetData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value, ok := n.store.Get(key)
	log.Println("Key found", value)
	fmt.Fprint(w, value, ok)
}

func (n *ClustorNode) handleSetData(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		http.Error(w, "key or value not provided", http.StatusBadRequest)
		return
	}
	ok := n.store.Set(key, value)
	log.Println("Key and value :", key, value)
	fmt.Fprint(w, ok)
}

func (n *ClustorNode) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
