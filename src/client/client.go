package client

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/JaneKetko/Buses/api/proto"
)

//Client - info struct for client.
type Client struct {
	Username string
	Password string
	grpccl   proto.BusesManagerClient
}

//NewClient - constructor for Client.
func NewClient(username, passwd string, client proto.BusesManagerClient) *Client {
	return &Client{
		Username: username,
		Password: passwd,
		grpccl:   client,
	}
}

//JSON marshal content to webpage.
func JSON(w http.ResponseWriter, code int, reply interface{}) {
	response, err := json.Marshal(reply)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}

//BuyTicket takes one place in bus by client.
func (c *Client) BuyTicket(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reply, err := c.grpccl.BuyTicket(r.Context(), &proto.IDRequest{ID: int64(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, reply)
}

//Handlers initializes main handle functions.
func (c *Client) Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/routes/buy/{id}", c.BuyTicket).Methods(http.MethodPost)
	return router
}

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
