package workserver

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

//Server is struct for client service.
type Server struct {
	router *mux.Router
}

//NewServer - constructor for Server.
func NewServer(c *Client) *Server {
	return &Server{c.Handlers()}
}

//RunServer starts working with client.
func (s *Server) RunServer() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Client has started working...\nAddress: http://localhost:%d", listener.Addr().(*net.TCPAddr).Port)
	log.Fatal(http.Serve(listener, s.router))
}
