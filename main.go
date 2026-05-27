package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is home handler!"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("displaying snippet view page"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("displaying form for creating snippet"))
}

func main() {
	// our server mux, which we can use to connect certain end point/url path to certain hendlers
	// for example /home is for Home hanlders
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("server is running on port :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Println(err)
		return
	}
}
