package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JaneKetko/Buses/src/config"
	"github.com/JaneKetko/Buses/src/domain"
)

//RouteDB - struct for describing route from db.
type RouteDB struct {
	idRoute    int
	startTime  string
	cost       int
	freeSeats  int
	allSeats   int
	idPoint    int
	startPoint string
	endPoint   string
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
func convertTypes(routeDB RouteDB) (domain.Route, error) {
	var route domain.Route
	date, err := time.Parse("2006-01-02 15:04:05", routeDB.startTime)
	if err != nil {
		return route, err
	}
	route = domain.Route{ID: routeDB.idRoute,
		Points: domain.Points{StartPoint: routeDB.startPoint,
			EndPoint: routeDB.endPoint},
		Start:     date,
		Cost:      routeDB.cost,
		FreeSeats: routeDB.freeSeats,
		AllSeats:  routeDB.allSeats}
	return route, nil
}

//Open opens connection with database.
func Open(cfg *config.Config) (*sql.DB, error) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Login, cfg.Passwd, cfg.Hostname, strconv.Itoa(cfg.Port), cfg.DBName))

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.New("database hasn't connected")
	}
	return db, nil
}

//GetAllData gets full data from db.
func (dbmanager *DBManager) GetAllData() ([]domain.Route, error) {
	rows, err := dbmanager.db.Query(`SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats,
		p.id_points, p.startpoint, p.endpoint
		FROM route r JOIN points p ON r.id_points = p.id_points`)
	if err != nil {
		return nil, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var dbr RouteDB
	var routes []domain.Route
	for rows.Next() {
		err = rows.Scan(&dbr.idRoute, &dbr.startTime, &dbr.cost, &dbr.freeSeats,
			&dbr.allSeats, &dbr.idPoint, &dbr.startPoint, &dbr.endPoint)
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

//RouteByID finds route by id in database.
func (dbmanager *DBManager) RouteByID(id int) (*domain.Route, error) {
	rows, err := dbmanager.db.Query(`SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats, 
	p.id_points, p.startpoint, p.endpoint 
	FROM route r JOIN points p on r.id_points = p.id_points WHERE r.id_route=?`, id)

	if err != nil {
		return nil, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if !rows.Next() {
		return nil, errors.New("no such route")
	}
	var routeDB RouteDB
	err = rows.Scan(&routeDB.idRoute, &routeDB.startTime, &routeDB.cost, &routeDB.freeSeats,
		&routeDB.allSeats, &routeDB.idPoint, &routeDB.startPoint, &routeDB.endPoint)
	if err != nil {
		return nil, errors.New("something is wrong")
	}

	route, err := convertTypes(routeDB)
	if err != nil {
		return nil, errors.New("errors with types")
	}
	return &route, nil
}

//DeleteRow deletes row from database by id.
func (dbmanager *DBManager) DeleteRow(id int) error {
	stmtIns, err := dbmanager.db.Prepare("DELETE FROM route where id_route=?")
	if err != nil {
		return err
	}
	rows, err := stmtIns.Exec(id)
	if err != nil {
		return err
	}
	if n, _ := rows.RowsAffected(); n == 0 {
		return errors.New("no such route")
	}
	return nil
}

//RoutesByEndPoint finds row in database by date and endpoint.
func (dbmanager *DBManager) RoutesByEndPoint(endpoint string) ([]domain.Route, error) {
	rows, err := dbmanager.db.Query(`SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats, 
	p.id_points, p.startpoint, p.endpoint 
	FROM route r JOIN points p on r.id_points = p.id_points WHERE p.endpoint=?`, endpoint)
	if err != nil {
		return nil, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var dbr RouteDB
	var routes []domain.Route
	for rows.Next() {
		err = rows.Scan(&dbr.idRoute, &dbr.startTime, &dbr.cost, &dbr.freeSeats,
			&dbr.allSeats, &dbr.idPoint, &dbr.startPoint, &dbr.endPoint)
		if err != nil {
			return nil, errors.New("no data")
		}
		route, err := convertTypes(dbr)
		if err != nil {
			return nil, errors.New("errors with types")
		}
		routes = append(routes, route)
	}

	if len(routes) == 0 {
		return nil, errors.New("no such routes by this endpoint")
	}
	return routes, nil
}

func (dbmanager *DBManager) insertPoint(startpoint, endpoint string) (int64, error) {
	stmtIn, err := dbmanager.db.Prepare("INSERT INTO points (startpoint, endpoint) VALUES( ?, ? )")
	if err != nil {
		return 0, err
	}
	defer func() {
		err = stmtIn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	row, err := stmtIn.Exec(startpoint, endpoint)
	if err != nil {
		return 0, err
	}

	pointID, err := row.LastInsertId()
	if err != nil {
		return 0, err
	}
	return pointID, nil
}

func (dbmanager *DBManager) insertRoute(id, freeseats, allseats, cost int, datetime string) (int64, error) {

	date, err := time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil {
		return 0, err
	}
	stmtIns, err := dbmanager.db.Prepare(`INSERT INTO route (id_points, starttime, cost, freeseats, allseats)
			VALUES( ?, ?, ?, ?, ? )`)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = stmtIns.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	rowRoute, err := stmtIns.Exec(id, date, cost, freeseats, allseats)
	if err != nil {
		return 0, err
	}

	idRoute, err := rowRoute.LastInsertId()
	if err != nil {
		return 0, err
	}
	return idRoute, nil
}

//AddRoute adds route to database.
func (dbmanager *DBManager) AddRoute(startpoint, endpoint, datetime string,
	cost, freeseats, allseats int) (int, error) {
	rows, err := dbmanager.db.Query("SELECT id_points FROM points WHERE startpoint=? AND endpoint=?",
		startpoint, endpoint)
	if err != nil {
		return 0, errors.New("data hasn't read")
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var pointID int64
	if !rows.Next() {
		pointID, err = dbmanager.insertPoint(startpoint, endpoint)
		if err != nil {
			return 0, err
		}
	} else {
		err = rows.Scan(&pointID)
		if err != nil {
			return 0, err
		}
	}
	idRoute, err := dbmanager.insertRoute(int(pointID), freeseats, allseats, cost, datetime)
	if err != nil {
		return 0, err
	}
	return int(idRoute), nil
}
