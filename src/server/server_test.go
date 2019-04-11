package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/mock"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/routemanager"
	"github.com/JaneKetko/Buses/src/routemanager/mocks"
	"github.com/JaneKetko/Buses/src/stores/domain"
	sst "github.com/JaneKetko/Buses/src/stores/serverstore"
)

func forGetRoutes(t *testing.T, method, path string) {
	cfg := &config.Config{
		PortRESTServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	serv := NewRESTServer(routeman, cfg)

	s := serv.managerHandlers()
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
			Start:     time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC),
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
		{
			name:           "errors",
			expectedStatus: http.StatusInternalServerError,
			expectedRoutes: nil,
			expectedError:  errors.New("smth bad"),
		},
	}

	routestrg.On(method, mock.Anything).Return(testCases[0].expectedRoutes, testCases[0].expectedError)

	t.Run(testCases[0].name, func(t *testing.T) {
		res := e.Request(http.MethodGet, path).Expect()
		res.Status(testCases[0].expectedStatus)
	})

	var rtstrg mocks.RouteStorage
	routeman = routemanager.NewRouteManager(&rtstrg)
	busstation := NewRESTServer(routeman, cfg)

	s = busstation.managerHandlers()
	server = httptest.NewServer(s)
	defer server.Close()
	e = httpexpect.New(t, server.URL)

	rtstrg.On(method, mock.Anything).Return(testCases[1].expectedRoutes, testCases[1].expectedError)
	t.Run(testCases[1].name, func(t *testing.T) {
		res := e.Request(http.MethodGet, path).Expect()
		res.Status(testCases[1].expectedStatus)
	})

	routestrg.AssertExpectations(t)
	rtstrg.AssertExpectations(t)
}

func TestGetRoutes(t *testing.T) {
	forGetRoutes(t, "GetAllData", "/routes")
}

func TestGetCurrentRoutes(t *testing.T) {
	forGetRoutes(t, "GetCurrentData", "/buses")
}

func TestGetRoute(t *testing.T) {

	cfg := &config.Config{
		PortRESTServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewRESTServer(routeman, cfg)

	s := busstation.managerHandlers()
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
			Start:     time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC),
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
			expectedError:  domain.ErrNoRoutes,
		},
		{
			name:           "invalid id",
			paramID:        "df2",
			expectedStatus: http.StatusBadRequest,
			expectedRoute:  nil,
			expectedError:  domain.ErrNoRoutes,
		},
	}
	for _, tc := range testCases[:2] {
		routestrg.On("RouteByID", mock.Anything, tc.routeID).Return(tc.expectedRoute, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodGet, "/routes/"+tc.paramID).Expect()
			res.Status(tc.expectedStatus)
		})
	}

	routestrg.AssertExpectations(t)
}

func TestCreateRoute(t *testing.T) {
	cfg := &config.Config{
		PortRESTServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewRESTServer(routeman, cfg)

	s := busstation.managerHandlers()
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
		route          *domain.Route
		expectedStatus int
		expectedID     int
		expectedError  error
	}{
		{
			name:           "errors",
			route:          &routes[0],
			expectedStatus: http.StatusBadRequest,
			expectedID:     1,
			expectedError:  domain.ErrInvalidDate,
		},
		{
			name:           "successful test",
			route:          &routes[1],
			expectedStatus: http.StatusOK,
			expectedID:     1,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases[1:] {
		routestrg.On("AddRoute", mock.Anything,
			tc.route).
			Return(tc.expectedID, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodPost, "/routes/add").WithHeader("Content-Type", "application/json").
				WithJSON(sst.RouteToRouteServer(*tc.route)).Expect()
			res.Status(tc.expectedStatus)
		})
	}
	routestrg.AssertExpectations(t)
}

func TestDeleteRoute(t *testing.T) {
	cfg := &config.Config{
		PortRESTServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewRESTServer(routeman, cfg)

	s := busstation.managerHandlers()
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
			expectedError:  domain.ErrNoRoutes,
		},
		{
			name:           "invalid id",
			paramID:        "df2",
			expectedStatus: http.StatusBadRequest,
			expectedError:  domain.ErrNoRoutes,
		},
	}

	for _, tc := range testCases[:2] {
		routestrg.On("DeleteRow", mock.Anything, tc.routeID).Return(tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodDelete, "/routes/"+tc.paramID).Expect()
			res.Status(tc.expectedStatus)
		})
	}
	routestrg.AssertExpectations(t)
}
func TestSearchRoutes(t *testing.T) {
	cfg := &config.Config{
		PortRESTServer: ":8000",
	}
	var routestrg mocks.RouteStorage
	routeman := routemanager.NewRouteManager(&routestrg)
	busstation := NewRESTServer(routeman, cfg)

	s := busstation.managerHandlers()
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
			date:           "2020-04-12",
			endPoint:       "Minsk",
			expectedStatus: http.StatusOK,
			expectedRoutes: routes[:2],
			expectedError:  nil,
		},
		{
			name:           "no routes by endpoint",
			date:           "2020-04-10",
			endPoint:       "Grodno",
			expectedStatus: http.StatusInternalServerError,
			expectedRoutes: nil,
			expectedError:  domain.ErrNoRoutesByEndPoint,
		},
		{
			name:           "invalid date argument",
			date:           "2019-04",
			endPoint:       "Grodno",
			expectedStatus: http.StatusBadRequest,
			expectedRoutes: nil,
			expectedError:  domain.ErrNoRoutesByEndPoint,
		},
	}

	for _, tc := range testCases {
		routestrg.On("RoutesByEndPoint", mock.Anything, tc.endPoint).Return(tc.expectedRoutes, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodGet, "/route_search").
				WithQueryString(fmt.Sprintf("date=%s&point=%s", tc.date, tc.endPoint)).Expect()
			res.Status(tc.expectedStatus)
		})
	}
	routestrg.AssertExpectations(t)
}
