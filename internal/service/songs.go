package service

import (
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/repository"
	"context"
	"github.com/gookit/slog"
)

type SongService struct {
	songRepository repository.ISongRepository
}

func NewSongService(songsRepository repository.ISongRepository) *SongService {
	return &SongService{songRepository: songsRepository}
}

func (s *SongService) Songs(c context.Context, filter model.Song, limit, offset int) ([]model.Song, error) {
	songs, err := s.songRepository.Songs(c, filter, limit, offset)
	if err != nil {
		slog.Error("error get songs", err)
		return nil, err
	}
	return songs, nil
}

func (s *SongService) Text(c context.Context, song model.Song, limit, offset int) (text string, err error) {
	text, err = s.songRepository.Text(c, song, limit, offset)
	if err != nil {
		slog.Error("error get text", err)
		return "", err
	}
	return text, nil
}
