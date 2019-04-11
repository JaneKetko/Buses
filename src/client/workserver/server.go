package workserver

import (
	"log"
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
func (s *Server) RunServer(addr string) {
	log.Println("Client has started working...")
	log.Fatal(http.ListenAndServe(addr, s.router))
}
