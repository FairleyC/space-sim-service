package commodity

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFetchingCommodity = errors.New("failed to fetch commodity by id")
	ErrNotImplemented = errors.New("not implemented")
)

type Commodity struct {
	ID string
	Name string
	Value float64
}

// Store - this interface defines all methods
// our service needs to operate.
type Store interface {
	GetCommodityById(context.Context, string) (Commodity, error)
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
	fmt.Println("retrieve commodity with id: ", id)

	cmt, err := s.Store.GetCommodityById(ctx, id)
	if err != nil {
		// Print the error for logging and return sanitized error to user.
		fmt.Println(err)
		return Commodity{}, ErrFetchingCommodity
	}

	return cmt, nil
}

func (s *Service) GetAllCommodity(ctx context.Context) ([]Commodity, error) {
	return nil, ErrNotImplemented
}

func (s *Service) UpdateCommodityPrice(ctx context.Context, id string, price float64) (Commodity, error) {
	return Commodity{}, ErrNotImplemented
}

func (s *Service) CreateCommodity(ctx context.Context, commodity Commodity) (Commodity, error) {
	return Commodity{}, ErrNotImplemented
}

func (s *Service) DeleteCommodity(ctx context.Context, id string) error {
	return ErrNotImplemented
}
