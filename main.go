package main

import (
	"log"
	"net/http"
	"sync/core"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func otjs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "ot.js")
}

func appjs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "app.js")
}

func utiljs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "util.js")
}

func main() {
	t := core.NewTranx()
	http.HandleFunc("/", home)
	http.HandleFunc("/app.js", appjs)
	http.HandleFunc("/ot.js", otjs)
	http.HandleFunc("/util.js", utiljs)

	http.HandleFunc("/sync", t.Fuck)
	log.Fatal(http.ListenAndServe(addr, nil))
}

var addr = "0.0.0.0:7777"
