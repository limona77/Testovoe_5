package service

import (
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/repository"
	"Testovoe_5/internal/service/api"
	"context"
)

type ISongService interface {
	Songs(c context.Context, filter model.Song, limit, offset int) ([]model.Song, error)
	Text(c context.Context, song model.Song, limit, offset int) (text string, err error)
	Delete(c context.Context, groupName, songName string) error
	Update(c context.Context, song model.Song) error
	Create(c context.Context, song model.Song) error
}

type Services struct {
	ISongService
}
type ServicesDeps struct {
	Repository *repository.Repositories
	ClientApi  *api.ApiClient
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		ISongService: NewSongService(deps.Repository, deps.ClientApi),
	}
}
