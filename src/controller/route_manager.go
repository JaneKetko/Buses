package routemanager

import (
	"errors"
	"time"

	"github.com/JaneKetko/Buses/src/structs"
)

//WorkDB - interface for database methods
type WorkDB interface {
	GetAllData() ([]structs.Route, error)
	RouteByID(id int) (structs.Route, error)
	DeleteRow(id int) error
	FindRoute(point string) ([]structs.Route, error)
	AddRoute(startpoint, endpoint, datetime string,
		cost float32, freeseats, allseats int) (int, error)
}

//RouteManager - struct for slice of routes
type RouteManager struct {
	//Db *dbmanager.DBManager
	Work WorkDB
}

//NewRouteManager - create new object of RouteManager struct
func NewRouteManager(work WorkDB) *RouteManager {
	return &RouteManager{Work: work}
}

//GetAllRoutes - method that gets all routes
func (r RouteManager) GetAllRoutes() ([]structs.Route, error) {

	routes, err := r.Work.GetAllData()
	if err != nil {
		return nil, err
	}
	return routes, nil
}

//GetRouteByID - Method that gets route by id
func (r RouteManager) GetRouteByID(id int) (structs.Route, error) {

	route, err := r.Work.RouteByID(id)
	if err != nil {
		return route, err
	}

	return route, nil
}

//CreateNewRoute - Method of creating new route
func (r *RouteManager) CreateNewRoute(route *structs.Route) error {
	if route.Start.Before(time.Now()) {
		return errors.New("Date is invalid")
	}
	id, err := r.Work.AddRoute(route.Points.StartPoint,
		route.Points.EndPoint,
		route.Start.Format("2006-01-02 15:04:05"),
		route.Cost,
		route.FreeSeats,
		route.AllSeats)

	if err != nil {
		return err
	}
	route.ID = id
	return nil
}

//DeleteRouteByID - Method of deleting route from all routes by id
func (r *RouteManager) DeleteRouteByID(id int) error {

	err := r.Work.DeleteRow(id)
	if err != nil {
		return err
	}
	return nil
}

//ChooseRoutesByDateAndPoint - Method that choose routes by date and point
func (r RouteManager) ChooseRoutesByDateAndPoint(date time.Time, point string) ([]structs.Route, error) {

	routes, err := r.Work.FindRoute(point)
	if err != nil {
		return nil, err
	}

	var routesDate []structs.Route
	for _, route := range routes {
		diff := route.Start.Sub(date).Hours()
		if diff >= 0 && diff < 24 {
			routesDate = append(routesDate, route)
		}
	}
	if len(routesDate) == 0 {
		return routesDate, errors.New("No such routes")
	}
	return routesDate, nil
}
