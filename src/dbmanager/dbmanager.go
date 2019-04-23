package dbmanager

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/JaneKetko/Buses/src/stores/domain"
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

//DBConfig - struct for database config info.
type DBConfig struct {
	Login    string
	Passwd   string
	Hostname string
	Port     int
	DBName   string
}

//Open opens connection with database.
func Open(cfg *DBConfig) (*sql.DB, error) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Login, cfg.Passwd, cfg.Hostname, strconv.Itoa(cfg.Port), cfg.DBName))

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database hasn't connected: %s", err)
	}
	return db, nil
}

func getData(ctx context.Context, dbmngr *DBManager, request string) ([]domain.Route, error) {
	var dbr RouteDB
	var routes []domain.Route

	rows, err := dbmngr.db.QueryContext(ctx, request)

	if err != nil {
		return nil, fmt.Errorf("data hasn't read: %s", err)
	}
	for rows.Next() {
		err = rows.Scan(&dbr.idRoute, &dbr.startTime, &dbr.cost, &dbr.freeSeats,
			&dbr.allSeats, &dbr.idPoint, &dbr.startPoint, &dbr.endPoint)
		if err != nil {
			return nil, err
		}
		route, err := convertTypes(dbr)
		if err != nil {
			return nil, domain.ErrTypes
		}

		routes = append(routes, route)
	}

	if len(routes) == 0 {
		return nil, domain.ErrNoRoutes
	}

	return routes, nil
}

//GetAllData gets full data from db.
func (dbmanager *DBManager) GetAllData(ctx context.Context) ([]domain.Route, error) {

	req := `SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats,
	p.id_points, p.startpoint, p.endpoint
	FROM route r 
	JOIN points p 
	ON r.id_points = p.id_points`

	return getData(ctx, dbmanager, req)
}

//GetCurrentData gets current buses. Their start time equal or greater than today.
func (dbmanager *DBManager) GetCurrentData(ctx context.Context) ([]domain.Route, error) {

	req := `SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats,
	p.id_points, p.startpoint, p.endpoint
	FROM route r 
	JOIN points p 
	ON r.id_points = p.id_points
	WHERE r.starttime >= DATE(NOW())`

	return getData(ctx, dbmanager, req)
}

//RouteByID finds route by id in database.
func (dbmanager *DBManager) RouteByID(ctx context.Context, id int) (*domain.Route, error) {

	var routeDB RouteDB
	err := dbmanager.db.QueryRowContext(ctx, `SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats,
	 p.id_points, p.startpoint, p.endpoint 
	 FROM route r JOIN points p on r.id_points = p.id_points WHERE r.id_route=?`, id).
		Scan(&routeDB.idRoute, &routeDB.startTime, &routeDB.cost, &routeDB.freeSeats,
			&routeDB.allSeats, &routeDB.idPoint, &routeDB.startPoint, &routeDB.endPoint)

	if err != nil {
		return nil, domain.ErrNoRoutes
	}

	route, err := convertTypes(routeDB)
	if err != nil {
		return nil, domain.ErrTypes
	}
	return &route, nil
}

//DeleteRow deletes row from database by id.
func (dbmanager *DBManager) DeleteRow(ctx context.Context, id int) error {

	res, err := dbmanager.db.ExecContext(ctx, "DELETE FROM route where id_route=?", id)
	if err != nil {
		return fmt.Errorf("cannot delete route: %s", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("something wrong with deleting: %s", err)
	}

	if n == 0 {
		return domain.ErrNoRoutes
	}
	return nil
}

//RoutesByEndPoint finds row in database by date and endpoint.
func (dbmanager *DBManager) RoutesByEndPoint(ctx context.Context, endpoint string) ([]domain.Route, error) {

	rows, err := dbmanager.db.QueryContext(ctx, `SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats, 
	p.id_points, p.startpoint, p.endpoint 
	FROM route r JOIN points p on r.id_points = p.id_points WHERE p.endpoint=?`, endpoint)

	if err != nil {
		return nil, fmt.Errorf("data hasn't read: %s", err)
	}

	var dbr RouteDB
	var routes []domain.Route
	for rows.Next() {
		err = rows.Scan(&dbr.idRoute, &dbr.startTime, &dbr.cost, &dbr.freeSeats,
			&dbr.allSeats, &dbr.idPoint, &dbr.startPoint, &dbr.endPoint)
		if err != nil {
			return nil, err
		}
		route, err := convertTypes(dbr)
		if err != nil {
			return nil, domain.ErrTypes
		}
		routes = append(routes, route)
	}

	if len(routes) == 0 {
		return nil, domain.ErrNoRoutesByEndPoint
	}
	return routes, nil
}

func (dbmanager *DBManager) insertPoint(ctx context.Context, startpoint, endpoint string) (int64, error) {

	res, err := dbmanager.db.ExecContext(ctx, "INSERT INTO points (startpoint, endpoint) VALUES( ?, ? )",
		startpoint, endpoint)

	if err != nil {
		return 0, fmt.Errorf("cannot insert points: %s", err)
	}

	pointID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return pointID, nil
}

type insertRouteStorage struct {
	IDPoints  int
	FreeSeats int
	AllSeats  int
	Cost      int
	DateTime  string
}

func (dbmanager *DBManager) insertRoute(ctx context.Context, rtinfo *insertRouteStorage) (int64, error) {

	date, err := time.Parse("2006-01-02 15:04:05", rtinfo.DateTime)
	if err != nil {
		return 0, err
	}

	res, err := dbmanager.db.ExecContext(ctx, `INSERT INTO route (id_points, starttime, cost, freeseats, allseats)
	VALUES( ?, ?, ?, ?, ? )`, rtinfo.IDPoints, date, rtinfo.Cost, rtinfo.FreeSeats, rtinfo.AllSeats)

	if err != nil {
		return 0, fmt.Errorf("cannot insert route: %s", err)
	}

	idRoute, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return idRoute, nil
}

//AddRoute adds route to database.
func (dbmanager *DBManager) AddRoute(ctx context.Context, r *domain.Route) (int, error) {

	row, err := dbmanager.db.QueryContext(ctx, "SELECT id_points FROM points WHERE startpoint=? AND endpoint=?",
		r.Points.StartPoint, r.Points.EndPoint)

	if err != nil {
		return 0, fmt.Errorf("data hasn't read: %s", err)
	}

	var pointID int64
	if !row.Next() {
		pointID, err = dbmanager.insertPoint(ctx, r.Points.StartPoint, r.Points.EndPoint)
		if err != nil {
			return 0, err
		}
	} else {
		err = row.Scan(&pointID)
		if err != nil {
			return 0, err
		}
	}
	idRoute, err := dbmanager.insertRoute(ctx, &insertRouteStorage{
		IDPoints:  int(pointID),
		FreeSeats: r.FreeSeats,
		AllSeats:  r.AllSeats,
		Cost:      r.Cost,
		DateTime:  r.Start.Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		return 0, err
	}
	return int(idRoute), nil
}

//TakePlace takes one place in bus.
func (dbmanager *DBManager) TakePlace(ctx context.Context, id int) (*domain.Ticket, error) {

	var routeDB RouteDB
	err := dbmanager.db.QueryRowContext(ctx, `SELECT r.id_route, r.starttime, r.cost, r.freeseats, r.allseats,
	 p.id_points, p.startpoint, p.endpoint 
	 FROM route r JOIN points p on r.id_points = p.id_points WHERE r.id_route=?`, id).
		Scan(&routeDB.idRoute, &routeDB.startTime, &routeDB.cost, &routeDB.freeSeats,
			&routeDB.allSeats, &routeDB.idPoint, &routeDB.startPoint, &routeDB.endPoint)

	if err != nil {
		return nil, domain.ErrNoRoutes
	}

	route, err := convertTypes(routeDB)
	if err != nil {
		return nil, domain.ErrTypes
	}

	if route.FreeSeats == 0 {
		return nil, domain.ErrNoFreeSeats
	}
	_, err = dbmanager.db.Query("UPDATE route SET freeseats = freeseats - 1 where id_route=?", id)
	if err != nil {
		return nil, err
	}

	ticket := &domain.Ticket{
		Points: domain.Points{
			StartPoint: route.Points.StartPoint,
			EndPoint:   route.Points.EndPoint},
		StartTime: route.Start,
		Cost:      route.Cost,
		Place:     route.AllSeats - route.FreeSeats + 1,
	}

	return ticket, nil
}
