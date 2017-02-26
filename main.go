package main

import (
	"log"
	"net/http"
	"sync/core"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	t := core.NewTranx()
	http.HandleFunc("/", home)
	http.HandleFunc("/sync", t.Fuck)
	log.Fatal(http.ListenAndServe(addr, nil))
}

var addr = "0.0.0.0:7777"
