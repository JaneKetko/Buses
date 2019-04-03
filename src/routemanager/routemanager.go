package routemanager

import (
	"context"
	"errors"
	"time"

	"github.com/JaneKetko/Buses/src/domain"
)

//RouteStorage - interface for database methods.
type RouteStorage interface {
	GetAllData(ctx context.Context) ([]domain.Route, error)
	RouteByID(ctx context.Context, id int) (*domain.Route, error)
	DeleteRow(ctx context.Context, id int) error
	RoutesByEndPoint(ctx context.Context, point string) ([]domain.Route, error)
	AddRoute(ctx context.Context, route *domain.Route) (int, error)
	TakePlace(ctx context.Context, id int) (*domain.Ticket, error)
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
func (r RouteManager) GetRoutes(ctx context.Context) ([]domain.Route, error) {
	return r.storage.GetAllData(ctx)
}

//GetRouteByID gets route by id.
func (r RouteManager) GetRouteByID(ctx context.Context, id int) (*domain.Route, error) {
	return r.storage.RouteByID(ctx, id)
}

//CreateNewRoute creates new route in database.
func (r *RouteManager) CreateRoute(ctx context.Context, route *domain.Route) error {
	if route.Start.Before(time.Now()) {
		return errors.New(domain.ErrInvalidDate)
	}
	id, err := r.storage.AddRoute(ctx, route)

	if err != nil {
		return err
	}
	route.ID = id
	return nil
}

//DeleteRouteByID deletes route from all routes by id.
func (r *RouteManager) DeleteRoute(ctx context.Context, id int) error {
	return r.storage.DeleteRow(ctx, id)
}

//ChooseRoutesByDateAndPoint chooses routes by date and point.
func (r RouteManager) SearchRoutes(ctx context.Context,
	date time.Time, endpoint string) ([]domain.Route, error) {

	routes, err := r.storage.RoutesByEndPoint(ctx, endpoint)
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
		return nil, errors.New(domain.ErrNoRoutes)
	}
	return routesDate, nil
}

//TakePlaceInBus takes one place in bus by client
func (r RouteManager) BuyTicket(ctx context.Context, id int) (*domain.Ticket, error) {
	return r.storage.TakePlace(ctx, id)
}
