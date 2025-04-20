package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FairleyC/space-sim-service/internal/data"
	"github.com/FairleyC/space-sim-service/internal/services/commodity"
	"github.com/google/uuid"
)

type CommodityRow struct {
	ID         string
	Name       sql.NullString
	UnitMass   sql.NullFloat64
	UnitVolume sql.NullFloat64
}

func convertCommodityRowToCommodity(row CommodityRow) commodity.Commodity {
	return commodity.Commodity{
		ID:         row.ID,
		Name:       row.Name.String,
		UnitMass:   row.UnitMass.Float64,
		UnitVolume: row.UnitVolume.Float64,
	}
}

func (d *Database) GetCommodityById(ctx context.Context, id string) (commodity.Commodity, error) {

	var commodityRow CommodityRow
	row := d.Pool.QueryRow(ctx, `
		SELECT id, name, unitmass, unitvolume
		FROM commodities
		WHERE id = $1
	`, id)

	err := row.Scan(&commodityRow.ID, &commodityRow.Name, &commodityRow.UnitMass, &commodityRow.UnitVolume)
	if err != nil {
		return commodity.Commodity{}, commodity.ErrCommodityNotFound
	}

	return convertCommodityRowToCommodity(commodityRow), nil
}

func (d *Database) GetCommoditiesByPagination(ctx context.Context, pagination data.Pagination) ([]commodity.Commodity, error) {
	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	orderBy := pagination.GetOrderByField([]data.AllowedField{
		{
			FieldName:          "unitmass",
			FormattedFieldName: "unit_mass",
		},
		{
			FieldName:          "unitvolume",
			FormattedFieldName: "unit_volume",
		},
		{
			FieldName:          "name",
			FormattedFieldName: "name",
		},
	}, "created_at")
	direction := pagination.GetOrderByDirection()

	rows, err := d.Pool.Query(ctx, `
		SELECT id, name, unit_mass, unit_volume
		FROM commodities
		ORDER BY `+orderBy+` `+direction+`
		LIMIT $1
		OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("error getting commodities by pagination: %w", err)
	}

	defer rows.Close()

	commodities := []commodity.Commodity{}
	for rows.Next() {
		var commodityRow CommodityRow
		err := rows.Scan(&commodityRow.ID, &commodityRow.Name, &commodityRow.UnitMass, &commodityRow.UnitVolume)
		if err != nil {
			return nil, fmt.Errorf("error scanning commodity row: %w", err)
		}

		commodities = append(commodities, convertCommodityRowToCommodity(commodityRow))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return commodities, nil
}

func (d *Database) CreateCommodity(ctx context.Context, newCommodity commodity.Commodity) (commodity.Commodity, error) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return commodity.Commodity{}, fmt.Errorf("error generating uuid: %w", err)
	}

	newCommodity.ID = newUuid.String()
	newRow := CommodityRow{
		ID:         newCommodity.ID,
		Name:       sql.NullString{String: newCommodity.Name, Valid: true},
		UnitMass:   sql.NullFloat64{Float64: newCommodity.UnitMass, Valid: true},
		UnitVolume: sql.NullFloat64{Float64: newCommodity.UnitVolume, Valid: true},
	}

	rows, err := d.Pool.Query(ctx, `
		INSERT INTO commodities (id, name, unit_mass, unit_volume)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, unit_mass, unit_volume
	`, newRow.ID, newRow.Name, newRow.UnitMass, newRow.UnitVolume)

	if err != nil {
		return commodity.Commodity{}, fmt.Errorf("error creating commodity: %w", err)
	}

	defer rows.Close()

	return newCommodity, nil
}

func (d *Database) RemoveCommodity(ctx context.Context, id string) error {
	_, err := d.Pool.Exec(ctx, `
		DELETE FROM commodities
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("error deleting commodity: %w", err)
	}

	return nil
}
