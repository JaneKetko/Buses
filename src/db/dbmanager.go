package dbmanager

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/structs"
)

//RouteDB - struct for describing route from db.
type RouteDB struct {
	IDRoute    int
	Starttime  string
	Cost       float32
	Freeseats  int
	Allseats   int
	IDPoint    int
	Startpoint string
	Endpoint   string
}

//DBManager - struct for storing database.
type DBManager struct {
	db *sql.DB
}

//NewDBManager - constructor for DBManager.
func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{db}
}

//ConvertTypes - convert RouteDB to Route.
func convertTypes(routeDB RouteDB) (structs.Route, error) {
	var route structs.Route
	date, err := time.Parse("2006-01-02 15:04:05", routeDB.Starttime)
	if err != nil {
		return route, err
	}
	route = structs.Route{ID: routeDB.IDRoute,
		Points: structs.Points{StartPoint: routeDB.Startpoint,
			EndPoint: routeDB.Endpoint},
		Start:     date,
		Cost:      routeDB.Cost,
		FreeSeats: routeDB.Freeseats,
		AllSeats:  routeDB.Allseats}
	return route, nil
}

//Open - method that opens connection with database.
func Open(config *config.Config) (*sql.DB, error) {

	db, err := sql.Open("mysql",
		config.Login+":"+config.Passwd+"@tcp("+config.Hostname+":"+strconv.Itoa(config.Port)+")"+"/"+config.DBName)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.New("database hasn't connected")
	}
	return db, nil
}

//GetAllData - get full data from db.
func (db *DBManager) GetAllData() ([]structs.Route, error) {
	rows, err := db.db.Query(`SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats,
		p.id_points, p.startpoint, p.endpoint
		FROM route r JOIN points p ON r.id_points = p.id_points`)
	if err != nil {
		return nil, errors.New("data hasn't read")
	}
	defer func() {
		err = rows.Close()
	}()

	var dbr RouteDB
	var routes []structs.Route
	for rows.Next() {
		err = rows.Scan(&dbr.IDRoute, &dbr.Starttime, &dbr.Cost, &dbr.Freeseats,
			&dbr.Allseats, &dbr.IDPoint, &dbr.Startpoint, &dbr.Endpoint)
		if err != nil {
			return nil, errors.New("no data")
		}
		route, err := convertTypes(dbr)
		if err != nil {
			return nil, errors.New("errors with types")
		}

		routes = append(routes, route)
	}
	return routes, nil
}

//RouteByID - find route by id in database.
func (db *DBManager) RouteByID(id int) (structs.Route, error) {
	rows, err := db.db.Query("SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats, "+
		"p.id_points, p.startpoint, p.endpoint "+
		"FROM route r JOIN points p on r.id_points = p.id_points WHERE r.id_route=?", id)
	var routeDB RouteDB
	var route structs.Route
	if err != nil {
		return route, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
	}()

	if !rows.Next() {
		return route, errors.New("no such route")
	}

	err = rows.Scan(&routeDB.IDRoute, &routeDB.Starttime, &routeDB.Cost, &routeDB.Freeseats,
		&routeDB.Allseats, &routeDB.IDPoint, &routeDB.Startpoint, &routeDB.Endpoint)
	if err != nil {
		return route, errors.New("something is wrong")
	}

	route, err = convertTypes(routeDB)
	if err != nil {
		return route, errors.New("errors with types")
	}
	return route, nil
}

//DeleteRow - delete row from database by id.
func (db *DBManager) DeleteRow(id int) error {
	stmtIns, err := db.db.Prepare("DELETE FROM route where id_route=?")
	if err != nil {
		return errors.New(err.Error())
	}
	rows, err := stmtIns.Exec(id)
	if err != nil {
		return errors.New(err.Error())
	}
	if n, _ := rows.RowsAffected(); n == 0 {
		return errors.New("no such route")
	}
	return nil
}

//FindRoute - find row in database by date and endpoint.
func (db *DBManager) FindRoute(point string) ([]structs.Route, error) {
	rows, err := db.db.Query("SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats, "+
		"p.id_points, p.startpoint, p.endpoint "+
		"FROM route r JOIN points p on r.id_points = p.id_points WHERE p.endpoint=?", point)
	if err != nil {
		return nil, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
	}()

	var dbr RouteDB
	var routes []structs.Route
	for rows.Next() {
		err = rows.Scan(&dbr.IDRoute, &dbr.Starttime, &dbr.Cost, &dbr.Freeseats,
			&dbr.Allseats, &dbr.IDPoint, &dbr.Startpoint, &dbr.Endpoint)
		if err != nil {
			return nil, errors.New("no data")
		}
		route, err := convertTypes(dbr)
		if err != nil {
			return nil, errors.New("errors with types")
		}
		routes = append(routes, route)
	}

	return routes, nil
}

func insertPoint(db *sql.DB, startpoint, endpoint string) (int64, error) {
	stmtIn, err := db.Prepare("INSERT INTO points (startpoint, endpoint) VALUES( ?, ? )")
	if err != nil {
		return 0, errors.New(err.Error())
	}
	defer func() {
		err = stmtIn.Close()
	}()

	row, err := stmtIn.Exec(startpoint, endpoint)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	pointID, err := row.LastInsertId()
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return pointID, nil
}

func insertRoute(db *sql.DB, id, freeseats, allseats int, datetime string, cost float32) (int64, error) {
	stmtIns, err := db.Prepare(`INSERT INTO route (id_points, starttime, cost, freeseats, allseats)
			VALUES( ?, ?, ?, ?, ? )`)
	if err != nil {
		return 0, errors.New(err.Error())
	}
	defer func() {
		err = stmtIns.Close()
	}()

	rowRoute, err := stmtIns.Exec(id, datetime, cost, freeseats, allseats)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	idRoute, err := rowRoute.LastInsertId()
	if err != nil {
		return 0, errors.New(err.Error())
	}
	return idRoute, nil
}

//AddRoute - method of adding route to database.
func (db *DBManager) AddRoute(startpoint, endpoint, datetime string,
	cost float32, freeseats, allseats int) (int, error) {
	rows, err := db.db.Query("SELECT * FROM points WHERE startpoint=? AND endpoint=?",
		startpoint, endpoint)
	if err != nil {
		return 0, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
	}()

	var idRoute int64
	if !rows.Next() {
		pointID, err := insertPoint(db.db, startpoint, endpoint)
		if err != nil {
			return 0, errors.New(err.Error())
		}
		idRoute, err = insertRoute(db.db, int(pointID), freeseats, allseats, datetime, cost)
		if err != nil {
			return 0, errors.New(err.Error())
		}
	} else {
		var id int
		var start, end string
		err = rows.Scan(&id, &start, &end)
		if err != nil {
			return 0, errors.New(err.Error())
		}

		idRoute, err = insertRoute(db.db, id, freeseats, allseats, datetime, cost)
		if err != nil {
			return 0, errors.New(err.Error())
		}
	}

	return int(idRoute), nil
}
