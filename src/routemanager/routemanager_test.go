package routemanager

import (
	"errors"
	"testing"
	"time"

	"github.com/JaneKetko/Buses/src/domain"
	"github.com/JaneKetko/Buses/src/routemanager/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutesByEndPoint(t *testing.T) {
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

	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)

	testCases := []struct {
		name           string
		date           time.Time
		endPoint       string
		expectedRoutes []domain.Route
		expectedError  error
		expTotalRoutes []domain.Route
		expTotalError  error
	}{
		{
			name:           "successful test",
			date:           time.Date(2019, 04, 12, 10, 0, 0, 0, time.UTC),
			endPoint:       "Minsk",
			expectedRoutes: routes[:2],
			expectedError:  nil,
			expTotalRoutes: routes[1:2],
			expTotalError:  nil,
		},
		{
			name:           "no routes by endpoint",
			date:           time.Date(2019, 04, 10, 10, 0, 0, 0, time.UTC),
			endPoint:       "Grodno",
			expectedRoutes: nil,
			expectedError:  errors.New("no such routes by this endpoint"),
			expTotalRoutes: nil,
			expTotalError:  errors.New("no such routes by this endpoint"),
		},
		{
			name:           "no routes by date",
			date:           time.Date(2022, 04, 10, 10, 0, 0, 0, time.UTC),
			endPoint:       "Mir",
			expectedRoutes: routes[2:],
			expectedError:  nil,
			expTotalRoutes: nil,
			expTotalError:  errors.New("no such routes"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("RoutesByEndPoint", tc.endPoint).Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rt, err := routeman.ChooseRoutesByDateAndPoint(tc.date, tc.endPoint)
			require.Equal(t, tc.expTotalError, err)
			assert.Equal(t, tc.expTotalRoutes, rt)
		})
	}
}

func TestCreateNewRoute(t *testing.T) {
	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)

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
			expTotalError: errors.New("date is invalid"),
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

	for _, tc := range testCases {
		routestrg.On("AddRoute",
			tc.route).
			Return(tc.expectedID, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := routeman.CreateNewRoute(tc.route)
			require.Equal(t, tc.expTotalError, err)
		})
	}
}

func TestGetAllRoutes(t *testing.T) {
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
	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)

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
	}

	for _, tc := range testCases {
		routestrg.On("GetAllData").Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rt, err := routeman.GetAllRoutes()
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedRoutes, rt)
		})
	}
}

func TestGetRouteByID(t *testing.T) {
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
	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)

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
			expectedError: errors.New("no such route"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("RouteByID", tc.routeID).Return(tc.expectedRoute, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rt, err := routeman.GetRouteByID(tc.routeID)
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedRoute, rt)
		})
	}
}

func TestDeleteRouteByID(t *testing.T) {
	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)

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
			expectedError: errors.New("no such route"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("DeleteRow", tc.routeID).Return(tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := routeman.DeleteRouteByID(tc.routeID)
			require.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTakePlaceInBus(t *testing.T) {
	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)
	ticket := &domain.Ticket{
		Points: domain.Points{
			StartPoint: "Minsk",
			EndPoint:   "Vitebsk",
		},
		StartTime: time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
		Cost:      1000,
		Place:     10,
	}

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
			expectedError:  errors.New("no such route"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("TakePlace", tc.routeID).Return(tc.expectedTicket, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := routeman.TakePlaceInBus(tc.routeID)
			require.Equal(t, tc.expectedError, err)
		})
	}
}
