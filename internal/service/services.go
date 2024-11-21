package service

import (
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/repository"
	"context"
)

type ISongService interface {
	Songs(c context.Context, filter model.Song, limit, offset int) ([]model.Song, error)
}

type Services struct {
	ISongService
}
type ServicesDeps struct {
	Repository *repository.Repositories
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		ISongService: NewSongService(deps.Repository),
	}
}
