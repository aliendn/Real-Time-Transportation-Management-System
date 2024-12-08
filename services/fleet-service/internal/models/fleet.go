package models

import (
	"database/sql"
	"fmt"
	"services/fleet-service/internal/db"
	"sync"
)

type Vehicle struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}

var VehiclePool = sync.Pool{
	New: func() interface{} {
		return &Vehicle{}
	},
}

func FetchVehicleByID(id int) (*Vehicle, error) {
	vehicle := VehiclePool.Get().(*Vehicle)
	defer VehiclePool.Put(vehicle)

	query := "SELECT id, name, capacity, status FROM vehicles WHERE id = $1"
	err := db.DB.QueryRow(query, id).Scan(&vehicle.ID, &vehicle.Name, &vehicle.Capacity, &vehicle.Status)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no vehicle found with id %d", id)
	} else if err != nil {
		return nil, err
	}

	return vehicle, nil
}
