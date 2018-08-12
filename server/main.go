package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/gorilla/mux"
)

var (
	connections = NewConnections()
	port        = flag.String("port", ":80", "Server port")
)

func main() {
	flag.Parse()
	go connections.Run()
	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(connections, w, r)
	}).Methods("GET")

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte("PONG"))
	}).Methods("GET")

	router.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		data, err := json.Marshal(connections.Channels)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(data)
	}).Methods("GET")

	router.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

	router.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		var files []os.FileInfo
		filepath.Walk("logs/", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			files = append(files, info)
			return nil
		})
		sort.Slice(files[:], func(i, j int) bool {
			return files[i].Size() > files[j].Size()
		})

		var output []*Channel
		var max = 5
		if len(files) < 5 {
			max = len(files)
		}
		for i := 0; i < max; i++ {
			f, _ := ioutil.ReadFile("logs/" + files[i].Name())

			channel := &Channel{}
			json.Unmarshal(f, channel)
			channel.Logfile = files[i].Name()
			output = append(output, channel)
		}

		resp, err := json.Marshal(output)
		if err != nil {
			fmt.Println(err)
		}

		w.Write(resp)
	}).Methods("GET")

	router.HandleFunc("/logs/{filename:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		filename := vars["filename"]

		f, err := ioutil.ReadFile("logs/" + filename)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(f)
	}).Methods("GET")

	fmt.Println("Starting " + *port)
	panic(http.ListenAndServe(*port, router))
}
