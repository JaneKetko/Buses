package workserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JaneKetko/Buses/api/proto"
	"github.com/JaneKetko/Buses/src/client/workserver/mocks"
	"github.com/JaneKetko/Buses/src/stores/domain"

	"github.com/gavv/httpexpect"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBuyTicket(t *testing.T) {
	cl := &mocks.BusesManagerClient{}
	client := NewClient("jane", "jane", cl)
	router := client.Handlers()
	serv := httptest.NewServer(router)
	defer serv.Close()
	e := httpexpect.New(t, serv.URL)

	date, err := ptypes.TimestampProto(time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	ticket := &proto.Ticket{
		Points: &proto.RoutePoints{
			StartPoint: "Minsk",
			EndPoint:   "Vitebsk",
		},
		Start: date,
		Cost:  1000,
		Place: 10,
	}

	testCases := []struct {
		name           string
		routeID        int
		paramID        string
		expectedTicket *proto.Ticket
		expectedError  error
		expectedStatus int
	}{
		{
			name:           "successful test",
			routeID:        1,
			paramID:        "1",
			expectedTicket: ticket,
			expectedError:  nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid id",
			routeID:        1,
			paramID:        "sdvsd",
			expectedTicket: ticket,
			expectedError:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "errors",
			routeID:        2,
			paramID:        "2",
			expectedTicket: nil,
			expectedError:  domain.ErrNoRoutes,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		cl.On("BuyTicket", mock.Anything, &proto.IDRequest{ID: int64(tc.routeID)}).
			Return(&proto.TicketResponse{Ticket: tc.expectedTicket}, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodPost, "/jane/routes/buy/"+tc.paramID).Expect()
			res.Status(tc.expectedStatus)
		})
	}

	cl.AssertExpectations(t)
}

func TestViewBuses(t *testing.T) {
	cl := &mocks.BusesManagerClient{}
	client := NewClient("jane", "jane", cl)
	router := client.Handlers()
	serv := httptest.NewServer(router)
	defer serv.Close()
	e := httpexpect.New(t, serv.URL)

	date, err := ptypes.TimestampProto(time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	routes := []*proto.BusRoute{
		{
			ID: 1,
			Points: &proto.RoutePoints{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     date,
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}

	testCases := []struct {
		name           string
		expectedStatus int
		expectedRoutes []*proto.BusRoute
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

	cl.On("GetRoutes", mock.Anything, &proto.Nothing{}).
		Return(&proto.ListRoutes{BusRoutes: testCases[0].expectedRoutes}, testCases[0].expectedError)

	t.Run(testCases[0].name, func(t *testing.T) {
		res := e.Request(http.MethodGet, "/jane/buses").Expect()
		res.Status(testCases[0].expectedStatus)
	})

	cl = &mocks.BusesManagerClient{}
	client = NewClient("jane", "jane", cl)
	router = client.Handlers()
	serv = httptest.NewServer(router)
	defer serv.Close()
	e = httpexpect.New(t, serv.URL)

	cl.On("GetRoutes", mock.Anything, &proto.Nothing{}).
		Return(&proto.ListRoutes{BusRoutes: testCases[1].expectedRoutes}, testCases[1].expectedError)

	t.Run(testCases[1].name, func(t *testing.T) {
		res := e.Request(http.MethodGet, "/jane/buses").Expect()
		res.Status(testCases[1].expectedStatus)
	})
	cl.AssertExpectations(t)
}

func TestFindBusByID(t *testing.T) {
	cl := &mocks.BusesManagerClient{}
	client := NewClient("jane", "jane", cl)
	router := client.Handlers()
	serv := httptest.NewServer(router)
	defer serv.Close()
	e := httpexpect.New(t, serv.URL)

	date, err := ptypes.TimestampProto(time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	routes := []*proto.BusRoute{
		{
			ID: 1,
			Points: &proto.RoutePoints{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     date,
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
		expectedRoute  *proto.BusRoute
		expectedError  error
	}{
		{
			name:           "successful test",
			routeID:        1,
			paramID:        "1",
			expectedStatus: http.StatusOK,
			expectedRoute:  routes[0],
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
		cl.On("GetRoute", mock.Anything, &proto.IDRequest{ID: int64(tc.routeID)}).
			Return(&proto.SingleRoute{Route: tc.expectedRoute}, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodGet, "/jane/routes/"+tc.paramID).Expect()
			res.Status(tc.expectedStatus)
		})
	}
	cl.AssertExpectations(t)
}

func TestSearchBuses(t *testing.T) {
	cl := &mocks.BusesManagerClient{}
	client := NewClient("jane", "jane", cl)
	router := client.Handlers()
	serv := httptest.NewServer(router)
	defer serv.Close()
	e := httpexpect.New(t, serv.URL)

	date1, err := ptypes.TimestampProto(time.Date(2020, 04, 23, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	date2, err := ptypes.TimestampProto(time.Date(2020, 04, 12, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	date3, err := ptypes.TimestampProto(time.Date(2020, 04, 10, 10, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	routes := []*proto.BusRoute{
		{
			ID: 1,
			Points: &proto.RoutePoints{
				StartPoint: "Vitebsk",
				EndPoint:   "Minsk",
			},
			Start:     date1,
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			ID: 2,
			Points: &proto.RoutePoints{
				StartPoint: "Grodno",
				EndPoint:   "Minsk",
			},
			Start:     date2,
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			ID: 3,
			Points: &proto.RoutePoints{
				StartPoint: "Pinsk",
				EndPoint:   "Mir",
			},
			Start:     date3,
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
		expectedRoutes []*proto.BusRoute
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
			expectedStatus: http.StatusInternalServerError,
			expectedRoutes: nil,
			expectedError:  domain.ErrNoRoutesByEndPoint,
		},
		{
			name:           "invalid arguments",
			date:           "2020-04-10",
			endPoint:       "",
			expectedStatus: http.StatusBadRequest,
			expectedRoutes: nil,
			expectedError:  domain.ErrNoRoutesByEndPoint,
		},
	}

	for _, tc := range testCases[:3] {
		cl.On("SearchRoutes", mock.Anything, &proto.Search{StartTime: tc.date, EndPoint: tc.endPoint}).
			Return(&proto.ListRoutes{BusRoutes: tc.expectedRoutes}, tc.expectedError)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res := e.Request(http.MethodGet, "/jane/route_search").
				WithQueryString(fmt.Sprintf("date=%s&point=%s", tc.date, tc.endPoint)).Expect()
			res.Status(tc.expectedStatus)
		})
	}
	cl.AssertExpectations(t)
}
