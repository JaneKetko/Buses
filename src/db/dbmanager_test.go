// +build integration

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

func TestInsertDeletePoint(t *testing.T) {

	db, err := DBOpen()
	assert.NoError(t, err)

	id, err := insertPoint(db, "Minsk", "Vitebsk")
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM points where id_points=?", id)
	assert.NoError(t, err)
}

func TestInsertRoute(t *testing.T) {

	db, err := DBOpen()
	assert.NoError(t, err)

	id, err := insertRoute(db, 7, 32, 44, "2019-02-24 08:30:00", 15.2)
	assert.NoError(t, err)

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	assert.NoError(t, err)
}

func TestRouteID(t *testing.T) {

	db, err := DBOpen()
	if err != nil {
		t.Errorf("error with opening testdb: %s", err)
	}

	id, err := insertRoute(db, 7, 32, 44, "2019-02-24 08:30:00", 15.2)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	dbmanager := NewDBManager(db)
	_, err = dbmanager.RouteByID(int(id))
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	_, err = db.Exec("DELETE FROM route where id_route=?", id)
	if err != nil {
		t.Errorf("errors with deleting row: %s", err)
	}
}
