package commodity

import (
	"context"
	"errors"
	"fmt"

	"github.com/FairleyC/space-sim-service/internal/data"
)

var (
	ErrFetchingCommodity = errors.New("failed to fetch commodity by id")
	ErrCommodityNotFound = errors.New("commodity not found")
	ErrNotImplemented    = errors.New("not implemented")
)

type Commodity struct {
	ID    string
	Name  string
	Price float64
}

// Store - this interface defines all methods
// our service needs to operate.
type Store interface {
	GetCommodityById(context.Context, string) (Commodity, error)
	GetCommoditiesByPagination(context.Context, data.Pagination) ([]Commodity, error)
	CreateCommodity(context.Context, Commodity) (Commodity, error)
	RemoveCommodity(context.Context, string) error
	UpdateCommodityPrice(context.Context, string, float64) (Commodity, error)
}

type CommodityService interface {
	GetAllCommodity(ctx context.Context, pagination data.Pagination) ([]Commodity, error)
	GetCommodity(ctx context.Context, id string) (Commodity, error)
	CreateCommodity(ctx context.Context, commodity Commodity) (Commodity, error)
	UpdateCommodityPrice(ctx context.Context, id string, price float64) (Commodity, error)
	RemoveCommodity(ctx context.Context, id string) error
}

// Service - is the struct on which all our
// logic will be built on top of
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetCommodity(ctx context.Context, id string) (Commodity, error) {
	commodity, err := s.Store.GetCommodityById(ctx, id)
	if err != nil {
		return Commodity{}, err
	}

	return commodity, nil
}

func (s *Service) GetAllCommodity(ctx context.Context, pagination data.Pagination) ([]Commodity, error) {
	commodities, err := s.Store.GetCommoditiesByPagination(ctx, pagination)
	if err != nil {
		return nil, fmt.Errorf("error getting commodities by pagination: %w", err)
	}

	return commodities, nil
}

func (s *Service) UpdateCommodityPrice(ctx context.Context, id string, price float64) (Commodity, error) {
	updatedCommodity, err := s.Store.UpdateCommodityPrice(ctx, id, price)
	if err != nil {
		return Commodity{}, fmt.Errorf("error updating commodity price: %w", err)
	}

	return updatedCommodity, nil
}

func (s *Service) CreateCommodity(ctx context.Context, commodity Commodity) (Commodity, error) {
	createdCommodity, err := s.Store.CreateCommodity(ctx, commodity)
	if err != nil {
		return Commodity{}, fmt.Errorf("error creating commodity: %w", err)
	}

	return createdCommodity, nil
}

func (s *Service) RemoveCommodity(ctx context.Context, id string) error {
	err := s.Store.RemoveCommodity(ctx, id)
	if err != nil {
		return fmt.Errorf("error removing commodity: %w", err)
	}

	return nil
}
