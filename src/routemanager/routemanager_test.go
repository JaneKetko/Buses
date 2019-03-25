package routemanager

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/JaneKetko/Buses/src/domain"

	"github.com/JaneKetko/Buses/src/routemanager/mocks"
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
				EndPoint:   "Minsk",
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
		date           time.Time
		endPoint       string
		expectedRoutes []domain.Route
		expectedError  error
		expTotalRoutes []domain.Route
		expTotalError  error
	}{
		{
			date:           time.Date(2019, 04, 10, 10, 0, 0, 0, time.UTC),
			endPoint:       "Minsk",
			expectedRoutes: routes,
			expectedError:  nil,
			expTotalRoutes: routes[2:],
			expTotalError:  nil,
		},
		{
			date:           time.Date(2019, 04, 10, 10, 0, 0, 0, time.UTC),
			endPoint:       "Grodno",
			expectedRoutes: nil,
			expectedError:  errors.New("no such routes by this endpoint"),
			expTotalRoutes: nil,
			expTotalError:  errors.New("no such routes by this endpoint"),
		},
		{
			date:           time.Date(2022, 04, 10, 10, 0, 0, 0, time.UTC),
			endPoint:       "Minsk",
			expectedRoutes: routes,
			expectedError:  nil,
			expTotalRoutes: nil,
			expTotalError:  errors.New("no such routes"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("RoutesByEndPoint", tc.endPoint).Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		rt, err := routeman.ChooseRoutesByDateAndPoint(tc.date, tc.endPoint)
		assert.Equal(t, tc.expTotalError, err)
		assert.Equal(t, tc.expTotalRoutes, rt)
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
				EndPoint:   "Minsk",
			},
			Start:     time.Date(2021, 04, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	testCases := []struct {
		route         domain.Route
		expectedID    int
		expectedError error
		expTotalError error
	}{
		{
			route:         routes[0],
			expectedID:    1,
			expectedError: nil,
			expTotalError: errors.New("date is invalid"),
		},
		{
			route:         routes[1],
			expectedID:    2,
			expectedError: errors.New("smth bad"),
			expTotalError: errors.New("smth bad"),
		},
		{
			route:         routes[2],
			expectedID:    3,
			expectedError: nil,
			expTotalError: nil,
		},
	}

	for _, tc := range testCases {
		routestrg.On("AddRoute",
			tc.route.Points.StartPoint,
			tc.route.Points.EndPoint,
			tc.route.Start.Format("2006-01-02 15:04:05"),
			tc.route.Cost,
			tc.route.FreeSeats,
			tc.route.AllSeats).
			Return(tc.expectedID, tc.expectedError)
	}

	for _, tc := range testCases {
		err := routeman.CreateNewRoute(&tc.route)
		assert.Equal(t, tc.expTotalError, err)
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
		expectedRoutes []domain.Route
		expectedError  error
	}{
		{
			expectedRoutes: routes,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		routestrg.On("GetAllData").Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		rt, err := routeman.GetAllRoutes()
		assert.Equal(t, tc.expectedError, err)
		assert.Equal(t, tc.expectedRoutes, rt)
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
		routeID       int
		expectedRoute *domain.Route
		expectedError error
	}{
		{
			routeID:       1,
			expectedRoute: &routes[0],
			expectedError: nil,
		},
		{
			routeID:       2,
			expectedRoute: nil,
			expectedError: errors.New("no such route"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("RouteByID", tc.routeID).Return(tc.expectedRoute, tc.expectedError)
	}

	for _, tc := range testCases {
		rt, err := routeman.GetRouteByID(tc.routeID)
		assert.Equal(t, tc.expectedError, err)
		assert.Equal(t, tc.expectedRoute, rt)
	}
}

func TestDeleteRouteByID(t *testing.T) {
	var routestrg mocks.RouteStorage
	routeman := NewRouteManager(&routestrg)

	testCases := []struct {
		routeID       int
		expectedError error
	}{
		{
			routeID:       1,
			expectedError: nil,
		},
		{
			routeID:       2,
			expectedError: errors.New("no such route"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("DeleteRow", tc.routeID).Return(tc.expectedError)
	}

	for _, tc := range testCases {
		err := routeman.DeleteRouteByID(tc.routeID)
		assert.Equal(t, tc.expectedError, err)
	}
}
