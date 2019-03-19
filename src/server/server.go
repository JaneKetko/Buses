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

	"github.com/gorilla/mux"
)

//BusStation - struct for describing bus station: manager for work with route info and configuration for server.
type BusStation struct {
	routes *routemanager.RouteManager
	config *config.Config
}

//NewBusStation - constructor for BusStation.
func NewBusStation(r *routemanager.RouteManager, c *config.Config) *BusStation {
	return &BusStation{routes: r,
		config: c,
	}
}

func (b *BusStation) getRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rts, err := b.routes.GetAllRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rserver := make([]routeServer, 0, 0)
	for _, rt := range rts {
		route := routeToRouteServer(rt)
		rserver = append(rserver, route)
	}

	err = json.NewEncoder(w).Encode(rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *BusStation) getRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idparam := params["id"]

	if idparam == "" {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idparam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	route, err := b.routes.GetRouteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rserver := routeToRouteServer(*route)
	err = json.NewEncoder(w).Encode(&rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *BusStation) createRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rserver routeServer
	err := json.NewDecoder(r.Body).Decode(&rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	route := routeServerToRoute(rserver)
	err = b.routes.CreateNewRoute(&route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rsencode := routeToRouteServer(route)
	err = json.NewEncoder(w).Encode(&rsencode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (b *BusStation) deleteRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = b.routes.DeleteRouteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rts, err := b.routes.GetAllRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rserver := make([]routeServer, 0, 0)
	for _, rt := range rts {
		route := routeToRouteServer(rt)
		rserver = append(rserver, route)
	}

	err = json.NewEncoder(w).Encode(rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *BusStation) searchRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	searchDate := params["date"]
	endpoint := params["point"]
	date, err := time.Parse("2006-01-02", searchDate)
	if err != nil {
		http.Error(w, "Invalid date argument!", http.StatusBadRequest)
		return
	}

	routesDate, err := b.routes.ChooseRoutesByDateAndPoint(date, endpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rserver := make([]routeServer, 0, 0)
	for _, rt := range routesDate {
		route := routeToRouteServer(rt)
		rserver = append(rserver, route)
	}
	err = json.NewEncoder(w).Encode(rserver)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

//StartServer - Start work with server
func (b *BusStation) StartServer() {
	r := mux.NewRouter()

	fmt.Printf("Started server at http://localhost%v.\n", ":"+strconv.Itoa(b.config.PortServer))
	r.HandleFunc("/route_search", b.searchRoutes).Queries("date", "{date}", "point", "{point}").Methods(http.MethodGet)
	r.HandleFunc("/routes", b.getRoutes).Methods(http.MethodGet)
	r.HandleFunc("/routes", b.createRoute).Methods(http.MethodPost)
	r.HandleFunc("/routes/{id}", b.getRoute).Methods(http.MethodGet)
	r.HandleFunc("/routes/{id}", b.deleteRoute).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(b.config.PortServer), r))
}
