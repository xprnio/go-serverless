package database

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Function struct {
	Id          string              `json:"id"`
	Image       string              `json:"image"`
	Environment FunctionEnvironment `json:"environment"`
}

func (d *Database) GetFunction(id uuid.UUID) (*Function, error) {
	r := d.connection.QueryRow(
		"select id, image, environment from functions where id = ? limit 1",
		id.String(),
	)

	return d.scanFunction(r)
}

func (d *Database) scanFunction(row Scannable) (*Function, error) {
	var id, image string
	var env string

	err := row.Scan(&id, &image, &env)
	if err != nil {
		return nil, err
	}

	environment, err := ParseEnvironment(env)
	if err != nil {
		return nil, err
	}

	f := &Function{
		Id:          id,
		Image:       image,
		Environment: environment,
	}
	return f, nil
}

func (d *Database) GetFunctions() ([]*Function, error) {
	r, err := d.connection.Query(
		"select id, image, environment from functions",
	)
	if err != nil {
		return nil, err
	}

	result := []*Function{}
	for r.Next() {
		f, err := d.scanFunction(r)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}

	return result, nil
}

func (d *Database) SaveFunction(f *Function) (*Function, error) {
	if f.Id == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}

		f.Id = id.String()
		_, err = d.connection.Exec(
			"insert into functions (id, image, environment) values (?, ?, ?)",
			f.Id, f.Image,
			f.Environment.String(),
		)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	_, err := d.connection.Exec(
		"update functions set image = ?, environment = ? where id = ?",
		f.Image, f.Environment.String(),
		f.Id,
	)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type FunctionEnvironment map[string]string

func ParseEnvironment(raw string) (FunctionEnvironment, error) {
	env := make(FunctionEnvironment)

	r := strings.NewReader(raw)
	s := bufio.NewScanner(r)

	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, "=", 2)
		switch len(parts) {
		case 0:
		case 1:
			return env, fmt.Errorf("invalid environment")
		}

		key, value := parts[0], parts[1]
		env[key] = value
	}

	return env, nil
}

func (env FunctionEnvironment) String() string {
	var result string

	for key, value := range env {
		result += fmt.Sprintf("%s=%s\n", key, value)
	}

	return result
}
