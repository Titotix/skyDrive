package main

import (
	"fmt"
	//	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
)

func httpServer() {
	//fooHandler := new(Handler)
	//http.Handle("/foo", fooHandler)

	http.HandleFunc("/storage", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		if r.Method == "POST" || r.Method == "PUT" {
			err := r.ParseForm()
			if err != nil {
				log.Print("ParseForm failed :", err)
			}

		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

	//	rtr := mux.NewRouter()
	//	rtr.HandleFunc("/user/{name:[a-z]+}/profile", profile).Methods("GET")
	//
	//	http.Handle("/", rtr)
	//
	//	log.Println("Listening...")
	//	http.ListenAndServe(":3000", nil)

}

//func profile(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//	name := params["name"]
//	w.Write([]byte("Hello " + name))
//}
