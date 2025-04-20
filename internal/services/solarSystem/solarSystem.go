package solarSystem

import (
	"context"
	"errors"

	"github.com/FairleyC/space-sim-service/internal/data"
)

var (
	ErrFindingSolarSystem  = errors.New("failed to find solar system by id")
	ErrSolarSystemNotFound = errors.New("solar system not found")
	ErrNotImplemented      = errors.New("not implemented")
)

type SolarSystem struct {
	ID   string
	Name string
}

type SolarSystemWithCommodityMarkets struct {
	ID               string
	Name             string
	CommodityMarkets []CommodityMarket
}

type CommodityMarket struct {
	ID             string
	BasePrice      float64
	DemandQuantity int
	CommodityName  string
}

type CommodityMarketUpdate struct {
	BasePrice      float64
	DemandQuantity int
}

type Store interface {
	GetSolarSystemById(context.Context, string) (SolarSystemWithCommodityMarkets, error)
	GetSolarSystemsByPagination(context.Context, data.Pagination) ([]SolarSystem, error)
	CreateSolarSystem(context.Context, SolarSystem) (SolarSystem, error)
	RemoveSolarSystem(context.Context, string) error
	GetCommodityMarketsBySolarSystemId(context.Context, string) ([]CommodityMarket, error)
	CreateCommodityMarket(context.Context, string, float64, int, string) (CommodityMarket, error)
	RemoveCommodityMarket(context.Context, string) error
	UpdateCommodityMarket(context.Context, string, CommodityMarketUpdate) (CommodityMarket, error)
	RemoveAllCommodityMarketsBySolarSystemId(context.Context, string) error
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) FindSolarSystem(ctx context.Context, id string) (SolarSystemWithCommodityMarkets, error) {
	solarSystem, err := s.Store.GetSolarSystemById(ctx, id)
	if err != nil {
		return SolarSystemWithCommodityMarkets{}, err
	}

	return solarSystem, nil
}

func (s *Service) FindAllSolarSystems(ctx context.Context, pagination data.Pagination) ([]SolarSystem, error) {
	solarSystems, err := s.Store.GetSolarSystemsByPagination(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return solarSystems, nil
}

func (s *Service) CreateSolarSystem(ctx context.Context, solarSystem SolarSystem) (SolarSystem, error) {
	newSolarSystem, err := s.Store.CreateSolarSystem(ctx, solarSystem)
	if err != nil {
		return SolarSystem{}, err
	}

	return newSolarSystem, nil
}

func (s *Service) RemoveSolarSystem(ctx context.Context, id string) error {
	err := s.Store.RemoveSolarSystem(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateCommodityMarket(ctx context.Context, solarSystemId string, basePrice float64, demandQuantity int, commodityId string) (CommodityMarket, error) {
	newCommodityMarket, err := s.Store.CreateCommodityMarket(ctx, solarSystemId, basePrice, demandQuantity, commodityId)
	if err != nil {
		return CommodityMarket{}, err
	}

	return newCommodityMarket, nil
}

func (s *Service) RemoveCommodityMarket(ctx context.Context, id string) error {
	err := s.Store.RemoveCommodityMarket(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateCommodityMarket(ctx context.Context, commodityMarketId string, commodityMarketUpdate CommodityMarketUpdate) (CommodityMarket, error) {
	updatedCommodityMarket, err := s.Store.UpdateCommodityMarket(ctx, commodityMarketId, commodityMarketUpdate)
	if err != nil {
		return CommodityMarket{}, err
	}

	return updatedCommodityMarket, nil
}
