package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FairleyC/space-sim-service/internal/data"
	"github.com/FairleyC/space-sim-service/internal/services/solarSystem"
	"github.com/google/uuid"
)

type SolarSystemRow struct {
	ID   string
	Name sql.NullString
}

func convertSolarSystemRowToSolarSystem(row SolarSystemRow) solarSystem.SolarSystem {
	return solarSystem.SolarSystem{
		ID:   row.ID,
		Name: row.Name.String,
	}
}

func (d *Database) GetSolarSystemById(ctx context.Context, id string) (solarSystem.SolarSystem, error) {

	var solarSystemRow SolarSystemRow
	row := d.Pool.QueryRow(ctx, `
		SELECT id, name
		FROM solar_systems
		WHERE id = $1
	`, id)

	err := row.Scan(&solarSystemRow.ID, &solarSystemRow.Name)
	if err != nil {
		return solarSystem.SolarSystem{}, solarSystem.ErrSolarSystemNotFound
	}

	return convertSolarSystemRowToSolarSystem(solarSystemRow), nil
}

func (d *Database) GetSolarSystemsByPagination(ctx context.Context, pagination data.Pagination) ([]solarSystem.SolarSystem, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	orderBy := pagination.GetOrderByField([]string{"name"}, "createdat")
	direction := pagination.GetOrderByDirection()

	rows, err := d.Pool.Query(ctx, `
		SELECT id, name
		FROM solar_systems
		ORDER BY `+orderBy+` `+direction+`
		LIMIT $1
		OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("error getting solar systems by pagination: %w", err)
	}

	defer rows.Close()

	solarSystems := []solarSystem.SolarSystem{}
	for rows.Next() {
		var solarSystemRow SolarSystemRow
		err := rows.Scan(&solarSystemRow.ID, &solarSystemRow.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning solar system row: %w", err)
		}

		solarSystems = append(solarSystems, convertSolarSystemRowToSolarSystem(solarSystemRow))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return solarSystems, nil
}

func (d *Database) CreateSolarSystem(ctx context.Context, newSolarSystem solarSystem.SolarSystem) (solarSystem.SolarSystem, error) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return solarSystem.SolarSystem{}, fmt.Errorf("error generating uuid: %w", err)
	}

	newSolarSystem.ID = newUuid.String()
	newRow := SolarSystemRow{
		ID:   newSolarSystem.ID,
		Name: sql.NullString{String: newSolarSystem.Name, Valid: true},
	}

	rows, err := d.Pool.Query(ctx, `
		INSERT INTO solar_systems (id, name)
		VALUES ($1, $2)
		RETURNING id, name
	`, newRow.ID, newRow.Name)

	if err != nil {
		return solarSystem.SolarSystem{}, fmt.Errorf("error creating solar system: %w", err)
	}

	defer rows.Close()

	return newSolarSystem, nil
}

func (d *Database) RemoveSolarSystem(ctx context.Context, id string) error {
	_, err := d.Pool.Exec(ctx, `
		DELETE FROM solar_systems
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error deleting solar system: %w", err)
	}

	return nil
}
