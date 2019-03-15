package dbmanager

import (
	"database/sql"
	"testing"

	"github.com/JaneKetko/Buses/src/config"

	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

func DBOpen() (*sql.DB, error) {
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

	db, err := DBOpen()
	assert.NoError(t, err)

	id1, err := insertRoute(db, 7, 32, 44, "2019-02-24 08:30:00", 15.2)
	assert.NoError(t, err)
	id2, err := insertRoute(db, 7, 32, 44, "02-24 08:30:00", 15.2)
	assert.NoError(t, err)

	dbmanager := NewDBManager(db)
	_, err = dbmanager.RouteByID(int(id1))
	assert.NoError(t, err)
	_, err = dbmanager.RouteByID(int(id2))
	assert.EqualError(t, err, "errors with types")

	_, err = db.Exec("DELETE FROM route where id_route=?", id1)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM route where id_route=?", id2)
	assert.NoError(t, err)
}

func TestAddRoute(t *testing.T) {

	db, err := DBOpen()
	assert.NoError(t, err)
	dbmanager := NewDBManager(db)

	id, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	assert.NoError(t, err)

	id, err = dbmanager.AddRoute("Minsk", "Lida", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM points where startpoint=? && endpoint=?", "Minsk", "Lida")
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	assert.NoError(t, err)
}

func TestGetAllData(t *testing.T) {
	db, err := DBOpen()
	assert.NoError(t, err)
	dbmanager := NewDBManager(db)

	id1, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)
	id2, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)
	id3, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)

	routes, err := dbmanager.GetAllData()
	assert.NoError(t, err)

	assert.Equal(t, 3, len(routes))

	_, err = db.Exec("DELETE FROM route where id_route=?", id1)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM route where id_route=?", id2)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM route where id_route=?", id3)
	assert.NoError(t, err)
}

func TestDeleteRoute(t *testing.T) {
	db, err := DBOpen()
	assert.NoError(t, err)
	dbmanager := NewDBManager(db)

	id, err := dbmanager.AddRoute("Minsk", "Vitebsk", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)
	err = dbmanager.DeleteRow(id)
	assert.NoError(t, err)

	_, err = dbmanager.RouteByID(id)

	assert.EqualError(t, err, "no such route")
}

func TestFindRoute(t *testing.T) {
	db, err := DBOpen()
	assert.NoError(t, err)
	dbmanager := NewDBManager(db)

	id1, err := dbmanager.AddRoute("Minsk", "Gomel", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)
	id2, err := dbmanager.AddRoute("Minsk", "Gomel", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)
	id3, err := dbmanager.AddRoute("Minsk", "Lida", "2019-02-24 08:30:00", 15.2, 32, 44)
	assert.NoError(t, err)

	routes, err := dbmanager.FindRoute("Gomel")
	assert.NoError(t, err)

	assert.Equal(t, 2, len(routes))
	_, err = db.Exec("DELETE FROM route where id_route=?", id1)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM route where id_route=?", id2)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM route where id_route=?", id3)
	assert.NoError(t, err)
}
