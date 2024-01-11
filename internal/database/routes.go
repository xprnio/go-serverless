package database

import (
	"github.com/google/uuid"
)

type Route struct {
	Id       string    `json:"id"`
	Path     string    `json:"path"`
	Function *Function `json:"function"`
}

func (d *Database) GetRoute(id uuid.UUID) (*Route, error) {
	r := d.connection.QueryRow(
		"select id, path, function_id from routes where id = ? limit 1",
		id.String(),
	)

	return d.scanRoute(r)
}

func (d *Database) GetRouteByPath(path string) (*Route, error) {
	r := d.connection.QueryRow(
		"select id, path, function_id from routes where path = ? limit 1",
		path,
	)

	return d.scanRoute(r)
}

func (d *Database) GetRoutes() ([]*Route, error) {
	r, err := d.connection.Query(
		"select id, path, function_id from routes",
	)
	if err != nil {
		return nil, err
	}

	result := []*Route{}
	for r.Next() {
		route, err := d.scanRoute(r)
		if err != nil {
			return nil, err
		}
		result = append(result, route)
	}

	return result, nil
}

func (d *Database) CreateRoute(path string, functionId uuid.UUID) (*Route, error) {
	f, err := d.GetFunction(functionId)
	if err != nil {
		return nil, err
	}

	return d.SaveRoute(&Route{
		Path:     path,
		Function: f,
	})
}

func (d *Database) SaveRoute(f *Route) (*Route, error) {
	if f.Function.Id == "" {
		if _, err := d.SaveFunction(f.Function); err != nil {
			return nil, err
		}
	}

	if f.Id == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		f.Id = id.String()

		_, err = d.connection.Exec(
			"insert into routes (id, path, function_id) values (?, ?, ?)",
			f.Id, f.Path,
			f.Function.Id,
		)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	_, err := d.connection.Exec(
		"update routes set path = ?, function_id = ? where id = ?",
		f.Path, f.Function.Id,
		f.Id,
	)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (d *Database) scanRoute(row Scannable) (*Route, error) {
	var id, path string
	var functionId string

	err := row.Scan(&id, &path, &functionId)
	if err != nil {
		return nil, err
	}

	function, err := d.GetFunction(
		uuid.MustParse(functionId),
	)
	if err != nil {
		return nil, err
	}

	f := &Route{
		Id:       id,
		Path:     path,
		Function: function,
	}
	return f, nil
}
