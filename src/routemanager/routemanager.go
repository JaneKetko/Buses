package routemanager

import (
	"errors"
	"time"

	"github.com/JaneKetko/Buses/src/domain"
)

//WorkDB - interface for database methods.
type RouteStorage interface {
	GetAllData() ([]domain.Route, error)
	RouteByID(id int) (*domain.Route, error)
	DeleteRow(id int) error
	FindRoute(point string) ([]domain.Route, error)
	AddRoute(startpoint, endpoint, datetime string,
		cost float32, freeseats, allseats int) (int, error)
}

//RouteManager - struct for slice of routes.
type RouteManager struct {
	//Db *dbmanager.DBManager
	storage RouteStorage
}

//NewRouteManager creates new object of RouteManager struct.
func NewRouteManager(storage RouteStorage) *RouteManager {
	return &RouteManager{storage: storage}
}

//GetAllRoutes gets all routes.
func (r RouteManager) GetAllRoutes() ([]domain.Route, error) {

	routes, err := r.storage.GetAllData()
	if err != nil {
		return nil, err
	}
	return routes, nil
}

//GetRouteByID gets route by id.
func (r RouteManager) GetRouteByID(id int) (*domain.Route, error) {

	route, err := r.storage.RouteByID(id)
	if err != nil {
		return nil, err
	}

	return route, nil
}

//CreateNewRoute creates new route in database.
func (r *RouteManager) CreateNewRoute(route *domain.Route) error {
	if route.Start.Before(time.Now()) {
		return errors.New("date is invalid")
	}
	id, err := r.storage.AddRoute(route.Points.StartPoint,
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

//DeleteRouteByID deletes route from all routes by id.
func (r *RouteManager) DeleteRouteByID(id int) error {

	err := r.storage.DeleteRow(id)
	if err != nil {
		return err
	}
	return nil
}

//ChooseRoutesByDateAndPoint chooses routes by date and point.
func (r RouteManager) ChooseRoutesByDateAndPoint(date time.Time, point string) ([]domain.Route, error) {

	routes, err := r.storage.FindRoute(point)
	if err != nil {
		return nil, err
	}

	var routesDate []domain.Route
	for _, route := range routes {
		diff := route.Start.Sub(date).Hours()
		if diff >= 0 && diff < 24 {
			routesDate = append(routesDate, route)
		}
	}
	if len(routesDate) == 0 {
		return nil, errors.New("no such routes")
	}
	return routesDate, nil
}
