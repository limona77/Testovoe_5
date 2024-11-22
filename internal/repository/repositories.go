package repository

import (
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/pkg/postgres"
	"context"
)

type ISongRepository interface {
	Songs(c context.Context, filter model.Song, limit, offset int) ([]model.Song, error)
	Text(c context.Context, song model.Song, limit, offset int) (text string, err error)
	Delete(c context.Context, groupName, songName string) error
	Update(c context.Context, song model.Song) error
	Create(c context.Context, song model.Song) error
}

type Repositories struct {
	ISongRepository
}

func NewRepositories(db *postgres.DB) *Repositories {
	return &Repositories{
		ISongRepository: NewSongsRepository(db),
	}
}
