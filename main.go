package main

import (
	"github/git-amw/devcache/coordinator"
	"io"
	"log"
	"net/http"
)

func main() {
	coordinator := coordinator.NewCoordinator()
	coordinator.AddNode("localhost:3000")
	coordinator.AddNode("localhost:3001")
	coordinator.AddNode("localhost:3002")

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		nodeAddress := coordinator.RouteRequest(key)
		resp, err := http.Get("http://" + nodeAddress + "/getdata?key=" + key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		nodeAddress := coordinator.RouteRequest(key)
		resp, err := http.Get("http://" + nodeAddress + "/setdata?key=" + key + "&value=" + value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	log.Println("Coordinator listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
