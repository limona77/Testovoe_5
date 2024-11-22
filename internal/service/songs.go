package service

import (
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/repository"
	"Testovoe_5/internal/service/api"
	"context"
	"github.com/gookit/slog"
)

type SongService struct {
	songRepository repository.ISongRepository
	apiClient      *api.ApiClient
}

func NewSongService(songsRepository repository.ISongRepository, apiClient *api.ApiClient) *SongService {
	return &SongService{songRepository: songsRepository, apiClient: apiClient}
}

func (s *SongService) Songs(c context.Context, filter model.Song, limit, offset int) ([]model.Song, error) {
	op := "service.Songs"
	songs, err := s.songRepository.Songs(c, filter, limit, offset)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return nil, err
	}
	return songs, nil
}

func (s *SongService) Text(c context.Context, song model.Song, limit, offset int) (text string, err error) {
	op := "service.Text"
	text, err = s.songRepository.Text(c, song, limit, offset)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return "", err
	}
	return text, nil
}

func (s *SongService) Delete(c context.Context, groupName, songName string) error {
	op := "service.Delete"
	err := s.songRepository.Delete(c, groupName, songName)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return err
	}
	return nil
}

func (s *SongService) Update(c context.Context, song model.Song) error {
	op := "service.Update"
	err := s.songRepository.Update(c, song)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return err
	}
	return nil
}

func (s *SongService) Create(c context.Context, song model.Song) error {
	op := "service.Create"
	//songInfo, err := s.apiClient.GetSongInfo(song.GroupName, song.SongName)
	//if err != nil {
	//	return fmt.Errorf("%s:%w", op, custom_errors.ErrNoSongInfo)
	//}
	//song.ReleaseDate = songInfo.ReleaseDate
	//song.Text = songInfo.Text
	//song.Link = songInfo.Link

	err := s.songRepository.Create(c, song)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return err
	}
	return nil
}
