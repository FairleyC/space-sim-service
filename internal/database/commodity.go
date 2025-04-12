package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FairleyC/space-sim-service/internal/commodity"
	"github.com/google/uuid"
)

type CommodityRow struct {
	ID    string
	Name  sql.NullString
	Value sql.NullFloat64
}

func convertCommodityRowToCommodity(row CommodityRow) commodity.Commodity {
	return commodity.Commodity{
		ID:    row.ID,
		Name:  row.Name.String,
		Value: row.Value.Float64,
	}
}

func (d *Database) GetCommodityById(ctx context.Context, id string) (commodity.Commodity, error) {

	var commodityRow CommodityRow
	row := d.Pool.QueryRow(ctx, `
		SELECT id, name, value
		FROM commodities
		WHERE id = $1
	`, id)

	err := row.Scan(&commodityRow.ID, &commodityRow.Name, &commodityRow.Value)
	if err != nil {
		return commodity.Commodity{}, fmt.Errorf("error fetching commodity by id: %w", err)
	}

	return convertCommodityRowToCommodity(commodityRow), nil
}

func (d *Database) CreateCommodity(ctx context.Context, newCommodity commodity.Commodity) (commodity.Commodity, error) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return commodity.Commodity{}, fmt.Errorf("error generating uuid: %w", err)
	}

	newCommodity.ID = newUuid.String()
	newRow := CommodityRow{
		ID:    newCommodity.ID,
		Name:  sql.NullString{String: newCommodity.Name, Valid: true},
		Value: sql.NullFloat64{Float64: newCommodity.Value, Valid: true},
	}

	rows, err := d.Pool.Query(ctx, `
		INSERT INTO commodities (id, name, value)
		VALUES ($1, $2, $3)
		RETURNING id, name, value
	`, newRow.ID, newRow.Name, newRow.Value)

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

func (d *Database) UpdateCommodityPrice(ctx context.Context, id string, price float64) (commodity.Commodity, error) {
	var commodityRow CommodityRow
	row := d.Pool.QueryRow(ctx, `
		UPDATE commodities
		SET value = $1
		WHERE id = $2
		RETURNING id, name, value
	`, price, id)

	err := row.Scan(&commodityRow.ID, &commodityRow.Name, &commodityRow.Value)
	if err != nil {
		return commodity.Commodity{}, fmt.Errorf("error updating commodity price: %w", err)
	}

	return convertCommodityRowToCommodity(commodityRow), nil
}
