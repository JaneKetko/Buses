//+build testdb

package dbmanager

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

func dbOpen() (*sql.DB, error) {
	config := &config.Config{
		PortServer: 8000,
		Login:      "root",
		Passwd:     "root",
		Hostname:   "172.17.0.2",
		Port:       3306,
		DBName:     "busstationtest",
	}

	db, err := Open(config)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestRouteID(t *testing.T) {

	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)
	id1, err := dbmanager.insertRoute(context.Background(), 7, 32, 44, 1500, "2019-02-24 08:30:00")
	require.NoError(t, err)
	_, err = dbmanager.insertRoute(context.Background(), 7, 32, 44, 1520, "02-24 08:30:00")
	require.Error(t, err, "invalid format of date")

	_, err = dbmanager.RouteByID(context.Background(), int(id1))
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id1)
	assert.NoError(t, err)
}

func TestAddRoute(t *testing.T) {

	routes := []domain.Route{
		{
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Vitebsk",
			},
			Start:     time.Date(2019, 02, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Lida",
			},
			Start:     time.Date(2019, 04, 10, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	id, err := dbmanager.AddRoute(context.Background(), &routes[0])
	require.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	assert.NoError(t, err)

	id, err = dbmanager.AddRoute(context.Background(), &routes[1])
	require.NoError(t, err)
	_, err = db.Exec("DELETE FROM points where startpoint=? && endpoint=?", "Minsk", "Lida")
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	assert.NoError(t, err)
}

func TestGetAllData(t *testing.T) {
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	route := domain.Route{
		Points: domain.Points{
			StartPoint: "Minsk",
			EndPoint:   "Vitebsk",
		},
		Start:     time.Date(2019, 02, 12, 10, 0, 0, 0, time.UTC),
		Cost:      1000,
		FreeSeats: 12,
		AllSeats:  13,
	}
	id1, err := dbmanager.AddRoute(context.Background(), &route)
	require.NoError(t, err)
	id2, err := dbmanager.AddRoute(context.Background(), &route)
	require.NoError(t, err)
	id3, err := dbmanager.AddRoute(context.Background(), &route)
	require.NoError(t, err)

	routes, err := dbmanager.GetAllData(context.Background())
	assert.NoError(t, err)

	assert.Equal(t, 3, len(routes))

	_, err = db.Exec("DELETE FROM route where id_route in (?, ?, ?)", id1, id2, id3)
	assert.NoError(t, err)

	_, err = dbmanager.GetAllData(context.Background())
	assert.EqualError(t, err, domain.ErrNoRoutes)
}

func TestDeleteRoute(t *testing.T) {
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	route := domain.Route{
		Points: domain.Points{
			StartPoint: "Minsk",
			EndPoint:   "Vitebsk",
		},
		Start:     time.Date(2019, 02, 12, 10, 0, 0, 0, time.UTC),
		Cost:      1000,
		FreeSeats: 12,
		AllSeats:  13,
	}
	id, err := dbmanager.AddRoute(context.Background(), &route)
	require.NoError(t, err)
	err = dbmanager.DeleteRow(context.Background(), id)
	require.NoError(t, err)
	err = dbmanager.DeleteRow(context.Background(), id)
	require.EqualError(t, err, domain.ErrNoRoutes)
}

func TestFindRoute(t *testing.T) {
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Vitebsk",
			},
			Start:     time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Vitebsk",
			},
			Start:     time.Date(2019, 02, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Lida",
			},
			Start:     time.Date(2019, 04, 10, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
	}
	id1, err := dbmanager.AddRoute(context.Background(), &routes[0])
	require.NoError(t, err)
	id2, err := dbmanager.AddRoute(context.Background(), &routes[1])
	require.NoError(t, err)
	id3, err := dbmanager.AddRoute(context.Background(), &routes[2])
	require.NoError(t, err)

	rts, err := dbmanager.RoutesByEndPoint(context.Background(), "Vitebsk")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(rts))
	_, err = db.Exec("DELETE FROM route where id_route in (?, ?, ?)", id1, id2, id3)
	assert.NoError(t, err)

	_, err = dbmanager.RoutesByEndPoint(context.Background(), "sdsvds")
	assert.EqualError(t, err, domain.ErrNoRoutesByEndPoint)
}

func TestTakePlace(t *testing.T) {
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)
	routes := []domain.Route{
		{
			ID: 1,
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Vitebsk",
			},
			Start:     time.Date(2019, 04, 23, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 12,
			AllSeats:  13,
		},
		{
			Points: domain.Points{
				StartPoint: "Minsk",
				EndPoint:   "Vitebsk",
			},
			Start:     time.Date(2019, 02, 12, 10, 0, 0, 0, time.UTC),
			Cost:      1000,
			FreeSeats: 0,
			AllSeats:  13,
		},
	}
	id1, err := dbmanager.AddRoute(context.Background(), &routes[0])
	require.NoError(t, err)
	id2, err := dbmanager.AddRoute(context.Background(), &routes[1])
	require.NoError(t, err)

	_, err = dbmanager.TakePlace(context.Background(), id1)
	assert.NoError(t, err)
	_, err = dbmanager.TakePlace(context.Background(), id2)
	assert.EqualError(t, err, domain.ErrNoFreeSeats)

	_, err = dbmanager.TakePlace(context.Background(), 10000000000)
	assert.EqualError(t, err, domain.ErrNoRoutes)

	_, err = db.Exec("DELETE FROM route where id_route=?", id1)
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id2)
	assert.NoError(t, err)
}
