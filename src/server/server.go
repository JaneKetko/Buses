package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JaneKetko/Buses/src/config"
	routemanager "github.com/JaneKetko/Buses/src/controller"
	"github.com/JaneKetko/Buses/src/structs"

	"github.com/gorilla/mux"
)

//BusStation - struct for describing bus station
type BusStation struct {
	Routes *routemanager.RouteManager
	Config *config.Config
}

//NewBusStation - constructor for BusStation
func NewBusStation(routes *routemanager.RouteManager, config *config.Config) *BusStation {
	return &BusStation{Routes: routes,
		Config: config,
	}
}

// //WorkDB - initialize database
// func WorkDB(db *dbmanager.DBManager) *BusStation {
// 	routes := routemanager.RouteManager{
// 		Db: db,
// 	}

// 	return &BusStation{&routes}
// }

func (b *BusStation) getRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rts, err := b.Routes.GetAllRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(rts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *BusStation) getRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	route, err := b.Routes.GetRouteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *BusStation) createRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var route structs.Route
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = b.Routes.CreateNewRoute(&route)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(route)
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

	err = b.Routes.DeleteRouteByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rts, err := b.Routes.GetAllRoutes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(rts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (b *BusStation) searchRoutes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	searchDate := params["date"]
	point := params["point"]
	date, err := time.Parse("2006-01-02", searchDate)
	if err != nil {
		http.Error(w, "Invalid date argument!", http.StatusBadRequest)
		return
	}

	routesDate, err := b.Routes.ChooseRoutesByDateAndPoint(date, point)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(routesDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

//StartServer - Start work with server
func (b *BusStation) StartServer() {
	r := mux.NewRouter()

	fmt.Printf("Started server at http://localhost%v.\n", ":"+strconv.Itoa(b.Config.PortServer))
	r.HandleFunc("/route_search", b.searchRoutes).Queries("date", "{date}", "point", "{point}").Methods(http.MethodGet)
	r.HandleFunc("/routes", b.getRoutes).Methods(http.MethodGet)
	r.HandleFunc("/routes", b.createRoute).Methods(http.MethodPost)
	r.HandleFunc("/routes/{id}", b.getRoute).Methods(http.MethodGet)
	r.HandleFunc("/routes/{id}", b.deleteRoute).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(b.Config.PortServer), r))
}
