package grpcserver

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/domain"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/routemanager/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetRoutes(t *testing.T) {
	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewServer(routeman, cfg)

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
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

	routestrg.On("GetAllData", context.Background()).Return(testCases[0].expectedRoutes, testCases[0].expectedError)
	t.Run(testCases[0].name, func(t *testing.T) {
		_, err := s.GetRoutes(context.Background(), &proto.Nothing{})
		log.Println(err)
		require.NoError(t, err)
	})

	var rtstrg mocks.RouteStorage
	routeman = routemanager.NewRouteManager(&rtstrg)
	s = NewServer(routeman, cfg)

	rtstrg.On("GetAllData", context.Background()).Return(testCases[1].expectedRoutes, testCases[1].expectedError)
	t.Run(testCases[1].name, func(t *testing.T) {
		_, err := s.GetRoutes(context.Background(), &proto.Nothing{})
		log.Println(err)
		require.Equal(t, err, testCases[1].expectedError)
	})
	routestrg.AssertExpectations(t)
	rtstrg.AssertExpectations(t)
}

func TestGetRoute(t *testing.T) {
	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewServer(routeman, cfg)

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	ctx := context.Background()

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
			expectedError: errors.New(domain.ErrNoRoutes),
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
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewServer(routeman, cfg)
	ctx := context.Background()

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
			expectedError: errors.New(domain.ErrNoRoutes),
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
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewServer(routeman, cfg)

	ticket := &domain.Ticket{
		Points: domain.Points{
			StartPoint: "Minsk",
			EndPoint:   "Vitebsk",
		},
		StartTime: time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
		Cost:      1000,
		Place:     10,
	}

	ctx := context.Background()

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
			expectedError:  errors.New(domain.ErrNoRoutes),
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
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	s := NewServer(routeman, cfg)

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
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
			Start:     time.Date(2019, 04, 12, 10, 0, 0, 0, time.UTC),
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
			Start:     time.Date(2019, 04, 10, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	ctx := context.Background()

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
			date:           "2019-04-12",
			endPoint:       "Minsk",
			expectedRoutes: routes[:2],
			expectedError:  nil,
			expTotalError:  nil,
		},
		{
			name:           "no routes by endpoint",
			date:           "2019-04-12",
			endPoint:       "Grodno",
			expectedRoutes: nil,
			expectedError:  errors.New(domain.ErrNoRoutesByEndPoint),
			expTotalError:  errors.New(domain.ErrNoRoutesByEndPoint),
		},
		{
			name:           "no routes by date",
			date:           "2022-04-12",
			endPoint:       "Mir",
			expectedRoutes: routes[2:],
			expectedError:  nil,
			expTotalError:  errors.New(domain.ErrNoRoutes),
		},
		{
			name:           "invalid date",
			date:           "2022-04-",
			endPoint:       "Mir",
			expectedRoutes: routes[2:],
			expectedError:  nil,
			expTotalError:  errors.New(domain.ErrInvalidDate),
		},
	}
	for _, tc := range testCases {
		routestrg.On("RoutesByEndPoint", ctx, tc.endPoint).Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.SearchRoutes(ctx, &proto.Search{Start: tc.date, EndPoint: tc.endPoint})
			require.Equal(t, tc.expTotalError, err)
		})
	}
	routestrg.AssertExpectations(t)
}
