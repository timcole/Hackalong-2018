package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	connections = NewConnections()
)

func main() {
	go connections.Run()
	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(connections, w, r)
	}).Methods("GET")

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	}).Methods("GET")

	router.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(connections.Channels)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(data)
	}).Methods("GET")

	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, err := json.Marshal(struct {
			Players  int `json:"players"`
			Channels int `json:"channels"`
		}{
			Players:  len(connections.Clients),
			Channels: len(connections.Channels),
		})
		if err != nil {
			fmt.Println(err)
		}
		w.Write(data)
	}).Methods("GET")

	fmt.Println("Starting :80")
	panic(http.ListenAndServe(":80", router))
}
