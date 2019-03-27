package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/gorilla/mux"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/domain"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/routemanager/mocks"
)

func TestGetRoutes(t *testing.T) {

	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewBusStation(routeman, cfg)

	s := busstation.managerHandlers(mux.NewRouter())
	server := httptest.NewServer(s)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

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
		expectedStatus int
		expectedRoutes []domain.Route
		expectedError  error
	}{
		{
			name:           "successful test",
			expectedStatus: http.StatusOK,
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
			res := e.Request(http.MethodGet, "/routes").Expect()
			res.Status(tc.expectedStatus)
		})
	}
}

func TestGetRoute(t *testing.T) {

	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewBusStation(routeman, cfg)

	s := busstation.managerHandlers(mux.NewRouter())
	server := httptest.NewServer(s)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

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
		routeID        int
		paramID        string
		expectedStatus int
		expectedRoute  *domain.Route
		expectedError  error
	}{
		{
			name:           "successful test",
			routeID:        1,
			paramID:        "1",
			expectedStatus: http.StatusOK,
			expectedRoute:  &routes[0],
			expectedError:  nil,
		},
		{
			name:           "no route",
			routeID:        2,
			paramID:        "2",
			expectedStatus: http.StatusInternalServerError,
			expectedRoute:  nil,
			expectedError:  errors.New("no such route"),
		},
		{
			name:           "invalid id",
			paramID:        "df2",
			expectedStatus: http.StatusBadRequest,
			expectedRoute:  nil,
			expectedError:  errors.New("no such route"),
		},
	}
	for _, tc := range testCases {
		routestrg.On("RouteByID", tc.routeID).Return(tc.expectedRoute, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodGet, "/routes/"+tc.paramID).Expect()
			res.Status(tc.expectedStatus)
		})
	}
}

func TestCreateRoute(t *testing.T) {
	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewBusStation(routeman, cfg)

	s := busstation.managerHandlers(mux.NewRouter())
	server := httptest.NewServer(s)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

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
			Start:     time.Date(2021, 04, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}
	testCases := []struct {
		name           string
		route          domain.Route
		expectedStatus int
		expectedID     int
		expectedError  error
	}{
		{
			name:           "errors",
			route:          routes[0],
			expectedStatus: http.StatusBadRequest,
			expectedID:     1,
			expectedError:  errors.New("date is invalid"),
		},
		{
			name:           "successful test",
			route:          routes[1],
			expectedStatus: http.StatusOK,
			expectedID:     1,
			expectedError:  nil,
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodPost, "/routes").WithHeader("Content-Type", "application/json").
				WithJSON(routeToRouteServer(tc.route)).Expect()
			res.Status(tc.expectedStatus)
		})
	}
}

func TestDeleteRoute(t *testing.T) {
	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewBusStation(routeman, cfg)

	s := busstation.managerHandlers(mux.NewRouter())
	server := httptest.NewServer(s)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

	testCases := []struct {
		name           string
		routeID        int
		paramID        string
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "successful test",
			routeID:        1,
			paramID:        "1",
			expectedStatus: http.StatusOK,
			expectedError:  nil,
		},
		{
			name:           "no route",
			routeID:        2,
			paramID:        "2",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  errors.New("no such route"),
		},
		{
			name:           "invalid id",
			paramID:        "df2",
			expectedStatus: http.StatusBadRequest,
			expectedError:  errors.New("no such route"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("DeleteRow", tc.routeID).Return(tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodDelete, "/routes/"+tc.paramID).Expect()
			res.Status(tc.expectedStatus)
		})
	}
}

func TestSearchRoutes(t *testing.T) {
	cfg := &config.Config{
		PortServer: 8000,
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewBusStation(routeman, cfg)

	s := busstation.managerHandlers(mux.NewRouter())
	server := httptest.NewServer(s)
	defer server.Close()
	e := httpexpect.New(t, server.URL)

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

	testCases := []struct {
		name           string
		date           string
		endPoint       string
		expectedStatus int
		expectedRoutes []domain.Route
		expectedError  error
	}{
		{
			name:           "successful test",
			date:           "2019-04-12",
			endPoint:       "Minsk",
			expectedStatus: http.StatusOK,
			expectedRoutes: routes[:2],
			expectedError:  nil,
		},
		{
			name:           "no routes by endpoint",
			date:           "2019-04-10",
			endPoint:       "Grodno",
			expectedStatus: http.StatusInternalServerError,
			expectedRoutes: nil,
			expectedError:  errors.New("no such routes by this endpoint"),
		},
		{
			name:           "invalid date argument",
			date:           "2019-04",
			endPoint:       "Grodno",
			expectedStatus: http.StatusBadRequest,
			expectedRoutes: nil,
			expectedError:  errors.New("no such routes by this endpoint"),
		},
	}

	for _, tc := range testCases {
		routestrg.On("RoutesByEndPoint", tc.endPoint).Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodGet, "/route_search").
				WithQueryString(fmt.Sprintf("date=%s&point=%s", tc.date, tc.endPoint)).Expect()
			res.Status(tc.expectedStatus)
		})
	}
}
