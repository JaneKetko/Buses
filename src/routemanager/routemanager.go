package routemanager

import (
	"errors"
	"time"

	"github.com/JaneKetko/Buses/src/domain"
)

//RouteStorage - interface for database methods.
type RouteStorage interface {
	GetAllData() ([]domain.Route, error)
	RouteByID(id int) (*domain.Route, error)
	DeleteRow(id int) error
	RoutesByEndPoint(point string) ([]domain.Route, error)
	AddRoute(*domain.Route) (int, error)
	TakePlace(id int) (*domain.Ticket, error)
}

//RouteManager - struct for slice of routes.
type RouteManager struct {
	storage RouteStorage
}

//NewRouteManager creates new object of RouteManager struct.
func NewRouteManager(storage RouteStorage) *RouteManager {
	return &RouteManager{storage: storage}
}

//GetAllRoutes gets all routes.
func (r RouteManager) GetAllRoutes() ([]domain.Route, error) {
	return r.storage.GetAllData()
}

//GetRouteByID gets route by id.
func (r RouteManager) GetRouteByID(id int) (*domain.Route, error) {
	return r.storage.RouteByID(id)
}

//CreateNewRoute creates new route in database.
func (r *RouteManager) CreateNewRoute(route *domain.Route) error {
	if route.Start.Before(time.Now()) {
		return errors.New("date is invalid")
	}
	id, err := r.storage.AddRoute(route)

	if err != nil {
		return err
	}
	route.ID = id
	return nil
}

//DeleteRouteByID deletes route from all routes by id.
func (r *RouteManager) DeleteRouteByID(id int) error {
	return r.storage.DeleteRow(id)
}

//ChooseRoutesByDateAndPoint chooses routes by date and point.
func (r RouteManager) ChooseRoutesByDateAndPoint(date time.Time, endpoint string) ([]domain.Route, error) {

	routes, err := r.storage.RoutesByEndPoint(endpoint)
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

//TakePlaceInBus takes one place in bus by client
func (r RouteManager) TakePlaceInBus(id int) (*domain.Ticket, error) {
	return r.storage.TakePlace(id)
}
