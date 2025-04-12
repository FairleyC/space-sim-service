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
	Price sql.NullFloat64
}

func convertCommodityRowToCommodity(row CommodityRow) commodity.Commodity {
	return commodity.Commodity{
		ID:    row.ID,
		Name:  row.Name.String,
		Price: row.Price.Float64,
	}
}

func (d *Database) GetCommodityById(ctx context.Context, id string) (commodity.Commodity, error) {

	var commodityRow CommodityRow
	row := d.Pool.QueryRow(ctx, `
		SELECT id, name, price
		FROM commodities
		WHERE id = $1
	`, id)

	err := row.Scan(&commodityRow.ID, &commodityRow.Name, &commodityRow.Price)
	if err != nil {
		return commodity.Commodity{}, commodity.ErrCommodityNotFound
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
		Price: sql.NullFloat64{Float64: newCommodity.Price, Valid: true},
	}

	rows, err := d.Pool.Query(ctx, `
		INSERT INTO commodities (id, name, price)
		VALUES ($1, $2, $3)
		RETURNING id, name, price
	`, newRow.ID, newRow.Name, newRow.Price)

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
		SET price = $1
		WHERE id = $2
		RETURNING id, name, price
	`, price, id)

	err := row.Scan(&commodityRow.ID, &commodityRow.Name, &commodityRow.Price)
	if err != nil {
		return commodity.Commodity{}, fmt.Errorf("error updating commodity price: %w", err)
	}

	return convertCommodityRowToCommodity(commodityRow), nil
}
