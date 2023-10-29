package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	host := []byte{127, 0, 0, 1}
	port := 8080
	addr := net.TCPAddr{IP: host, Port: port, Zone: ""}

	users := NewUserAPI()
	logger := http.HandlerFunc(logHandler)

	mux.Handle("/users", HandlerChain{logger, users})

	log.Printf("Listening on http://%v\n", addr.String())
	if err := http.ListenAndServe(addr.String(), mux); err != nil {
		log.Fatalln(err)
	}
}
