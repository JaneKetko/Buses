package grpcserver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/routemanager/mocks"
	"github.com/JaneKetko/Buses/src/stores/domain"
	"github.com/JaneKetko/Buses/src/stores/serverstore"

	"github.com/stretchr/testify/require"
)

func TestGetRoutes(t *testing.T) {
	cfg := &config.Config{
		PortGRPCServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewGRPSServer(routeman, cfg.PortGRPCServer)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	testCases := []struct {
		name           string
		expectedRoutes []domain.Route
		expectedError  error
	}{
		{
			name:           "successful test",
			expectedRoutes: routes,
			expectedError:  nil,
		},
		{
			name:           "errors",
			expectedRoutes: nil,
			expectedError:  errors.New("smth bad"),
		},
	}

	routestrg.On("GetCurrentData", ctx).Return(testCases[0].expectedRoutes, testCases[0].expectedError)
	t.Run(testCases[0].name, func(t *testing.T) {
		_, err := s.GetRoutes(ctx, &proto.Nothing{})
		require.NoError(t, err)
	})

	var rtstrg mocks.RouteStorage
	routeman = routemanager.NewRouteManager(&rtstrg)
	s = NewGRPSServer(routeman, cfg.PortGRPCServer)

	rtstrg.On("GetCurrentData", ctx).Return(testCases[1].expectedRoutes, testCases[1].expectedError)
	t.Run(testCases[1].name, func(t *testing.T) {
		_, err := s.GetRoutes(ctx, &proto.Nothing{})
		require.Equal(t, err, testCases[1].expectedError)
	})
	routestrg.AssertExpectations(t)
	rtstrg.AssertExpectations(t)
}

func TestGetRoute(t *testing.T) {
	cfg := &config.Config{
		PortGRPCServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewGRPSServer(routeman, cfg.PortGRPCServer)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	testCases := []struct {
		name          string
		routeID       int
		expectedRoute *domain.Route
		expectedError error
	}{
		{
			name:          "successful test",
			routeID:       1,
			expectedRoute: &routes[0],
			expectedError: nil,
		},
		{
			name:          "no route",
			routeID:       2,
			expectedRoute: nil,
			expectedError: domain.ErrNoRoutes,
		},
	}

	for _, tc := range testCases {
		routestrg.On("RouteByID", ctx, tc.routeID).Return(tc.expectedRoute, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.GetRoute(ctx, &proto.IDRequest{ID: int64(tc.routeID)})
			require.Equal(t, tc.expectedError, err)
		})
	}
	routestrg.AssertExpectations(t)
}

func TestDeleteRoute(t *testing.T) {

	cfg := &config.Config{
		PortGRPCServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewGRPSServer(routeman, cfg.PortGRPCServer)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCases := []struct {
		name          string
		routeID       int
		expectedError error
	}{
		{
			name:          "successful test",
			routeID:       1,
			expectedError: nil,
		},
		{
			name:          "no route",
			routeID:       2,
			expectedError: domain.ErrNoRoutes,
		},
	}

	for _, tc := range testCases {
		routestrg.On("DeleteRow", ctx, tc.routeID).Return(tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.DeleteRoute(ctx, &proto.IDRequest{ID: int64(tc.routeID)})
			require.Equal(t, tc.expectedError, err)
		})
	}
	routestrg.AssertExpectations(t)
}

func TestBuyTicket(t *testing.T) {
	cfg := &config.Config{
		PortGRPCServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewGRPSServer(routeman, cfg.PortGRPCServer)

	ticket := &domain.Ticket{
		Points: domain.Points{
			StartPoint: "Minsk",
			EndPoint:   "Vitebsk",
		},
		StartTime: time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC),
		Cost:      1000,
		Place:     10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCases := []struct {
		name           string
		routeID        int
		expectedTicket *domain.Ticket
		expectedError  error
	}{
		{
			name:           "successful test",
			routeID:        1,
			expectedTicket: ticket,
			expectedError:  nil,
		},
		{
			name:           "errors",
			routeID:        2,
			expectedTicket: nil,
			expectedError:  domain.ErrNoRoutes,
		},
	}

	for _, tc := range testCases {
		routestrg.On("TakePlace", ctx, tc.routeID).Return(tc.expectedTicket, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.BuyTicket(ctx, &proto.IDRequest{ID: int64(tc.routeID)})
			require.Equal(t, tc.expectedError, err)
		})
	}
	routestrg.AssertExpectations(t)
}

func TestSearchRoutes(t *testing.T) {
	cfg := &config.Config{
		PortGRPCServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewGRPSServer(routeman, cfg.PortGRPCServer)

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			ID: 2,
			Points: domain.Points{
				StartPoint: "Grodno",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2020, 04, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			ID: 3,
			Points: domain.Points{
				StartPoint: "Pinsk",
				EndPoint:   "Mir",
			},
			Start:     time.Date(2020, 04, 10, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCases := []struct {
		name           string
		date           string
		endPoint       string
		expectedRoutes []domain.Route
		expectedError  error
		expTotalError  error
	}{
		{
			name:           "successful test",
			date:           "2020-04-12",
			endPoint:       "Minsk",
			expectedRoutes: routes[:2],
			expectedError:  nil,
			expTotalError:  nil,
		},
		{
			name:           "no routes by endpoint",
			date:           "2020-04-12",
			endPoint:       "Grodno",
			expectedRoutes: nil,
			expectedError:  domain.ErrNoRoutesByEndPoint,
			expTotalError:  domain.ErrNoRoutesByEndPoint,
		},
		{
			name:           "no routes by date",
			date:           "2022-04-12",
			endPoint:       "Mir",
			expectedRoutes: routes[2:],
			expectedError:  nil,
			expTotalError:  domain.ErrNoRoutes,
		},
		{
			name:           "invalid date",
			date:           "2022-04-",
			endPoint:       "Mir",
			expectedRoutes: routes[2:],
			expectedError:  nil,
			expTotalError:  domain.ErrInvalidDate,
		},
	}
	for _, tc := range testCases {
		routestrg.On("RoutesByEndPoint", ctx, tc.endPoint).Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.SearchRoutes(ctx, &proto.Search{StartTime: tc.date, EndPoint: tc.endPoint})
			require.Equal(t, tc.expTotalError, err)
		})
	}
	routestrg.AssertExpectations(t)
}

func TestCreateRoute(t *testing.T) {
	cfg := &config.Config{
		PortGRPCServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewGRPSServer(routeman, cfg.PortGRPCServer)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	routes := []domain.Route{
		{
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2002, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			Points: domain.Points{
				StartPoint: "Grodno",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2020, 04, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			Points: domain.Points{
				StartPoint: "Grodno",
				EndPoint:   "Mir",
			},
			Start:     time.Date(2021, 04, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	testCases := []struct {
		name          string
		route         *domain.Route
		expectedID    int
		expectedError error
		expTotalError error
	}{
		{
			name:          "invalid date",
			route:         &routes[0],
			expectedID:    1,
			expectedError: nil,
			expTotalError: domain.ErrInvalidDate,
		},
		{
			name:          "errors",
			route:         &routes[1],
			expectedID:    2,
			expectedError: errors.New("smth bad"),
			expTotalError: errors.New("smth bad"),
		},
		{
			name:          "successful test",
			route:         &routes[2],
			expectedID:    3,
			expectedError: nil,
			expTotalError: nil,
		},
	}

	for _, tc := range testCases[1:] {
		routestrg.On("AddRoute", ctx,
			tc.route).
			Return(tc.expectedID, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, err := serverstore.RouteToPType(*tc.route)
			require.NoError(t, err)
			_, err = s.CreateRoute(ctx, &proto.SingleRoute{Route: r})
			require.Equal(t, tc.expTotalError, err)
		})
	}
	routestrg.AssertExpectations(t)
}
