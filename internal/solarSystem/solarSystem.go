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

type Store interface {
	GetSolarSystemById(context.Context, string) (SolarSystem, error)
	GetSolarSystemsByPagination(context.Context, data.Pagination) ([]SolarSystem, error)
	CreateSolarSystem(context.Context, SolarSystem) (SolarSystem, error)
	RemoveSolarSystem(context.Context, string) error
}

type SolarSystemService interface {
	GetAllSolarSystems(ctx context.Context, pagination data.Pagination) ([]SolarSystem, error)
	GetSolarSystem(ctx context.Context, id string) (SolarSystem, error)
	CreateSolarSystem(ctx context.Context, solarSystem SolarSystem) (SolarSystem, error)
	RemoveSolarSystem(ctx context.Context, id string) error
}

type Service struct {
	Store Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) GetSolarSystem(ctx context.Context, id string) (SolarSystem, error) {
	solarSystem, err := s.Store.GetSolarSystemById(ctx, id)
	if err != nil {
		return SolarSystem{}, err
	}

	return solarSystem, nil
}

func (s *Service) GetAllSolarSystems(ctx context.Context, pagination data.Pagination) ([]SolarSystem, error) {
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
