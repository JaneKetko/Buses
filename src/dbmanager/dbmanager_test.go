// +build integration

package dbmanager

import (
	"database/sql"
	"testing"

	"github.com/JaneKetko/Buses/src/config"
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
	id1, err := dbmanager.insertRoute(7, 32, 44, 1500, "2019-02-24 08:30:00")
	require.NoError(t, err)
	_, err = dbmanager.insertRoute(7, 32, 44, 1520, "02-24 08:30:00")
	require.Error(t, err, "invalid format of date")

	_, err = dbmanager.RouteByID(int(id1))
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id1)
	assert.NoError(t, err)
}

func TestAddRoute(t *testing.T) {

	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	_, err = dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02- 08:30:00", 1520, 32, 44)
	require.Error(t, err)

	id, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	assert.NoError(t, err)

	id, err = dbmanager.AddRoute("Minsk", "Lida", "2019-02-24 08:30:00", 1520, 32, 44)
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

	id1, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)
	id2, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)
	id3, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)

	routes, err := dbmanager.GetAllData()
	assert.NoError(t, err)

	assert.Equal(t, 3, len(routes))

	_, err = db.Exec("DELETE FROM route where id_route in (?, ?, ?)", id1, id2, id3)
	assert.NoError(t, err)
}

func TestDeleteRoute(t *testing.T) {
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	id, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)
	err = dbmanager.DeleteRow(id)
	require.NoError(t, err)

	_, err = dbmanager.RouteByID(id)

	assert.EqualError(t, err, "no such route")
}

func TestFindRoute(t *testing.T) {
	db, err := dbOpen()
	require.NoError(t, err)
	dbmanager := NewDBManager(db)

	id1, err := dbmanager.AddRoute("Minsk", "Gomel", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)
	id2, err := dbmanager.AddRoute("Minsk", "Gomel", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)
	id3, err := dbmanager.AddRoute("Minsk", "Lida", "2019-02-24 08:30:00", 1520, 32, 44)
	require.NoError(t, err)

	routes, err := dbmanager.RoutesByEndPoint("Gomel")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(routes))
	_, err = db.Exec("DELETE FROM route where id_route in (?, ?, ?)", id1, id2, id3)
	assert.NoError(t, err)
}
