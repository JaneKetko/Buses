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
	// listener, err := net.Listen("tcp", ":0")
	// if err != nil {
	// 	log.Fatal(err)
	// }listener.Addr().(*net.TCPAddr).Port
	log.Printf("Client has started working...\nAddress: http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, s.router))
}
