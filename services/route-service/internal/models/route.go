package models

import (
	"database/sql"
	"fmt"
	"route-service/internal/db"
)

type Route struct {
	ID          string  `json:"id"`
	Start       string  `json:"start"`
	Destination string  `json:"destination"`
	Distance    float64 `json:"distance"`
}

func FetchRouteByID(id string) (*Route, error) {
	var route Route
	query := "SELECT id, start, destination, distance FROM routes WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	err := row.Scan(&route.ID, &route.Start, &route.Destination, &route.Distance)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no route found with id %s", id)
	} else if err != nil {
		return nil, err
	}

	return &route, nil
}

func SaveRoute(route *Route) error {
	query := "INSERT INTO routes (id, start, destination, distance) VALUES ($1, $2, $3, $4)"
	_, err := db.DB.Exec(query, route.ID, route.Start, route.Destination, route.Distance)
	return err
}
