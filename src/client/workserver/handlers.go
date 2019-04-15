package workserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/stores/domain"
	sst "github.com/JaneKetko/Buses/src/stores/serverstore"

	"github.com/gorilla/mux"
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
		return
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

	ticket, err := sst.TicketPTypeToJSON(reply.Ticket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusCreated, ticket)
}

//ViewBuses gets all current buses for user.
func (c *Client) ViewBuses(w http.ResponseWriter, r *http.Request) {

	reply, err := c.grpccl.GetRoutes(r.Context(), &proto.Nothing{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes := make([]*sst.RouteForServer, 0)
	for _, rt := range reply.BusRoutes {
		route, err := sst.PtypeBusToJSON(rt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		routes = append(routes, route)
	}

	JSON(w, http.StatusOK, routes)
}

//FindBusByID finds bus by id.
func (c *Client) FindBusByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idparam := params["id"]
	id, err := strconv.Atoi(idparam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reply, err := c.grpccl.GetRoute(r.Context(), &proto.IDRequest{ID: int64(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	route, err := sst.PtypeBusToJSON(reply.Route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JSON(w, http.StatusOK, route)
}

//SearchBuses finds buses by date and endpoint.
func (c *Client) SearchBuses(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchDate := params["date"]
	endpoint := params["point"]
	if endpoint == "" {
		http.Error(w, domain.ErrInvalidArg.Error(), http.StatusBadRequest)
		return
	}
	reply, err := c.grpccl.SearchRoutes(r.Context(), &proto.Search{StartTime: searchDate, EndPoint: endpoint})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes := make([]*sst.RouteForServer, 0)
	for _, rt := range reply.BusRoutes {
		route, err := sst.PtypeBusToJSON(rt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		routes = append(routes, route)
	}

	JSON(w, http.StatusOK, routes)
}

//Handlers initializes main handle functions.
func (c *Client) Handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/%s/routes/buy/{id}", c.Username), c.BuyTicket).Methods(http.MethodPost)
	router.HandleFunc(fmt.Sprintf("/%s/buses", c.Username), c.ViewBuses).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/%s/route_search", c.Username), c.SearchBuses).
		Queries("date", "{date}", "point", "{point}").
		Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/%s/routes/{id}", c.Username), c.FindBusByID).Methods(http.MethodGet)
	return router
}
