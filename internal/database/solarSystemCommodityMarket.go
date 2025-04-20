package database

import (
	"context"
	"fmt"

	"github.com/FairleyC/space-sim-service/internal/services/solarSystem"
	"github.com/google/uuid"
)

type SolarSystemCommodityMarketRow struct {
	ID             string
	BasePrice      float64
	DemandQuantity int
	CommodityID    string
	SolarSystemID  string
}

type SolarSystemCommodityMarketRowWithCommodityName struct {
	SolarSystemCommodityMarketRow
	CommodityName string
}

func convertSolarSystemCommodityMarketRowWithCommodityNameToSolarSystemCommodityMarket(row SolarSystemCommodityMarketRowWithCommodityName) solarSystem.CommodityMarket {
	return solarSystem.CommodityMarket{
		ID:             row.ID,
		BasePrice:      row.BasePrice,
		DemandQuantity: row.DemandQuantity,
		CommodityName:  row.CommodityName,
	}
}

func convertSolarSystemCommodityMarketRowToSolarSystemCommodityMarket(row SolarSystemCommodityMarketRow, commodityName string) solarSystem.CommodityMarket {
	return solarSystem.CommodityMarket{
		ID:             row.ID,
		BasePrice:      row.BasePrice,
		DemandQuantity: row.DemandQuantity,
		CommodityName:  commodityName,
	}
}

func (d *Database) GetCommodityMarketsBySolarSystemId(ctx context.Context, solarSystemId string) ([]solarSystem.CommodityMarket, error) {
	rows, err := d.Pool.Query(ctx, `
		SELECT market.id, market.base_price, market.demand_quantity, market.commodity_id, market.solar_system_id, commodity.name
		FROM solar_system_commodity_markets market
		JOIN commodities commodity ON market.commodity_id = commodity.id
		WHERE market.solar_system_id = $1
	`, solarSystemId)

	if err != nil {
		return []solarSystem.CommodityMarket{}, err
	}

	defer rows.Close()

	commodityMarkets := []solarSystem.CommodityMarket{}
	for rows.Next() {
		var row SolarSystemCommodityMarketRowWithCommodityName
		err := rows.Scan(&row.ID, &row.BasePrice, &row.DemandQuantity, &row.CommodityID, &row.SolarSystemID, &row.CommodityName)
		if err != nil {
			return []solarSystem.CommodityMarket{}, err
		}

		commodityMarkets = append(commodityMarkets, convertSolarSystemCommodityMarketRowWithCommodityNameToSolarSystemCommodityMarket(row))
	}

	return commodityMarkets, nil
}

func (d *Database) GetCommodityMarketById(ctx context.Context, id string) (solarSystem.CommodityMarket, error) {
	var marketRow SolarSystemCommodityMarketRowWithCommodityName
	row := d.Pool.QueryRow(ctx, `
		SELECT market.id, market.base_price, market.demand_quantity, commodity.name
		FROM solar_system_commodity_markets market
		JOIN commodities commodity ON market.commodity_id = commodity.id
		WHERE market.id = $1
	`, id)

	err := row.Scan(&marketRow.ID, &marketRow.BasePrice, &marketRow.DemandQuantity, &marketRow.CommodityName)
	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error scanning commodity market: %w", err)
	}

	return convertSolarSystemCommodityMarketRowWithCommodityNameToSolarSystemCommodityMarket(marketRow), nil
}

func (d *Database) CreateCommodityMarket(ctx context.Context, solarSystemId string, basePrice float64, demandQuantity int, commodityId string) (solarSystem.CommodityMarket, error) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error generating uuid: %w", err)
	}

	_, err = d.Pool.Exec(ctx, `
		INSERT INTO solar_system_commodity_markets (id, base_price, demand_quantity, commodity_id, solar_system_id)
		VALUES ($1, $2, $3, $4, $5)
	`, newUuid.String(), basePrice, demandQuantity, commodityId, solarSystemId)

	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error creating commodity market: %w", err)
	}

	commodityMarket, err := d.GetCommodityMarketById(ctx, newUuid.String())
	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error getting commodity market by id: %w", err)
	}

	return commodityMarket, nil
}

func (d *Database) UpdateCommodityMarket(ctx context.Context, commodityMarketId string, updatedCommodityMarket solarSystem.CommodityMarketUpdate) (solarSystem.CommodityMarket, error) {
	rows, err := d.Pool.Query(ctx, `
		UPDATE solar_system_commodity_markets
		SET base_price = $1, demand_quantity = $2
		WHERE id = $3
		RETURNING id, base_price, demand_quantity, commodity_id, solar_system_id
	`, updatedCommodityMarket.BasePrice, updatedCommodityMarket.DemandQuantity, commodityMarketId)

	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error updating commodity market: %w", err)
	}

	defer rows.Close()

	var row SolarSystemCommodityMarketRow
	err = rows.Scan(&row.ID, &row.BasePrice, &row.DemandQuantity, &row.CommodityID, &row.SolarSystemID)
	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error scanning commodity market: %w", err)
	}

	commodity, err := d.GetCommodityById(ctx, row.CommodityID)
	if err != nil {
		return solarSystem.CommodityMarket{}, fmt.Errorf("error getting commodity by id: %w", err)
	}

	commodityMarket := convertSolarSystemCommodityMarketRowToSolarSystemCommodityMarket(row, commodity.Name)

	return commodityMarket, nil
}

func (d *Database) RemoveCommodityMarket(ctx context.Context, id string) error {
	_, err := d.Pool.Exec(ctx, `
		DELETE FROM solar_system_commodity_markets
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("error deleting commodity market: %w", err)
	}

	return nil
}

func (d *Database) RemoveAllCommodityMarketsBySolarSystemId(ctx context.Context, solarSystemId string) error {
	_, err := d.Pool.Exec(ctx, `
		DELETE FROM solar_system_commodity_markets
		WHERE solar_system_id = $1
	`, solarSystemId)

	if err != nil {
		return fmt.Errorf("error deleting all commodity markets by solar system id: %w", err)
	}

	return nil
}

func (d *Database) RemoveAllCommodityMarketsByCommodityId(ctx context.Context, commodityId string) error {
	_, err := d.Pool.Exec(ctx, `
		DELETE FROM solar_system_commodity_markets
		WHERE commodity_id = $1
	`, commodityId)

	if err != nil {
		return fmt.Errorf("error deleting all commodity markets by commodity id: %w", err)
	}

	return nil
}
