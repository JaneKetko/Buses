package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/stores/domain"
	sst "github.com/JaneKetko/Buses/src/stores/serverstore"

	"github.com/gorilla/mux"
)

//RESTServer - struct for describing bus station: manager for work with route info and configuration for server.
type RESTServer struct {
	routes *routemanager.RouteManager
	config *config.Config
}

//NewRESTServer - constructor for BusStation.
func NewRESTServer(r *routemanager.RouteManager, c *config.Config) *RESTServer {
	return &RESTServer{
		routes: r,
		config: c,
	}
}

func (b *RESTServer) getRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rts, err := b.routes.GetRoutes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rserver := make([]sst.RouteForServer, 0)
	for _, rt := range rts {
		route := sst.RouteToRouteServer(rt)
		rserver = append(rserver, route)
	}

	err = json.NewEncoder(w).Encode(rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *RESTServer) getCurrentRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rts, err := b.routes.GetCurrentRoutes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rserver := make([]sst.RouteForServer, 0)
	for _, rt := range rts {
		route := sst.RouteToRouteServer(rt)
		rserver = append(rserver, route)
	}

	err = json.NewEncoder(w).Encode(rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *RESTServer) getRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idparam := params["id"]
	id, err := strconv.Atoi(idparam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	route, err := b.routes.GetRouteByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rserver := sst.RouteToRouteServer(*route)
	err = json.NewEncoder(w).Encode(&rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *RESTServer) createRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rserver sst.RouteForServer
	err := json.NewDecoder(r.Body).Decode(&rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	route := sst.RouteServerToRoute(rserver)
	err = b.routes.CreateRoute(r.Context(), &route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rsencode := sst.RouteToRouteServer(route)
	err = json.NewEncoder(w).Encode(&rsencode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (b *RESTServer) deleteRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = b.routes.DeleteRoute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("the route was deleted successfully"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *RESTServer) searchRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	searchDate := params["date"]
	endpoint := params["point"]
	if endpoint == "" {
		http.Error(w, domain.ErrInvalidArg.Error(), http.StatusBadRequest)
		return
	}
	date, err := time.Parse("2006-01-02", searchDate)
	if err != nil {
		http.Error(w, domain.ErrInvalidDate.Error(), http.StatusBadRequest)
		return
	}

	routesDate, err := b.routes.SearchRoutes(r.Context(), date, endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rserver := make([]sst.RouteForServer, 0)
	for _, rt := range routesDate {
		route := sst.RouteToRouteServer(rt)
		rserver = append(rserver, route)
	}
	err = json.NewEncoder(w).Encode(rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *RESTServer) managerHandlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/route_search", b.searchRoutes).Queries("date", "{date}", "point", "{point}").
		Methods(http.MethodGet)
	router.HandleFunc("/routes", b.getRoutes).Methods(http.MethodGet)
	router.HandleFunc("/buses", b.getCurrentRoutes).Methods(http.MethodGet)
	router.HandleFunc("/routes/add", b.createRoute).Methods(http.MethodPost)
	router.HandleFunc("/routes/{id}", b.getRoute).Methods(http.MethodGet)
	router.HandleFunc("/routes/{id}", b.deleteRoute).Methods(http.MethodDelete)
	return router
}

//RunServer - Start work with server.
func (b *RESTServer) RunServer() {
	fmt.Printf("Started server at http://localhost%v.\n", b.config.PortRESTServer)
	router := b.managerHandlers()
	log.Fatal(http.ListenAndServe(b.config.PortRESTServer, router))
}
