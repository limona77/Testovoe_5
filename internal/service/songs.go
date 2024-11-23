package service

import (
	custom_errors "Testovoe_5/internal/custom-errors"
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/repository"
	"Testovoe_5/internal/service/api"
	"context"
	"fmt"
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
	slog.Debug("Retrieving songs", map[string]any{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
	})
	songs, err := s.songRepository.Songs(c, filter, limit, offset)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return nil, err
	}
	slog.Info("Song successfully created", "songs", songs)
	return songs, nil
}

func (s *SongService) Text(c context.Context, song model.Song, limit, offset int) (text string, err error) {
	op := "service.Text"

	slog.Debug("Retrieving text", map[string]any{
		"song":   song,
		"limit":  limit,
		"offset": offset,
	})
	text, err = s.songRepository.Text(c, song, limit, offset)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return "", err
	}
	slog.Info("Text successfully retrieved", "text", text)
	return text, nil
}

func (s *SongService) Delete(c context.Context, groupName, songName string) error {
	op := "service.Delete"
	slog.Debug("Deleting song", map[string]any{
		"groupName": groupName,
		"songName":  songName,
	})
	err := s.songRepository.Delete(c, groupName, songName)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return err
	}
	slog.Info("Song successfully deleted")
	return nil
}

func (s *SongService) Update(c context.Context, song model.Song) error {
	op := "service.Update"
	slog.Debug("Updating song", "song", song)
	err := s.songRepository.Update(c, song)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return err
	}
	slog.Info("Song successfully updated")
	return nil
}

func (s *SongService) Create(c context.Context, song model.Song) error {
	op := "service.Create"
	slog.Debug("Retrieving info about song", "song", song)
	songInfo, err := s.apiClient.GetSongInfo(song.GroupName, song.SongName)
	if err != nil {
		return fmt.Errorf("%s:%w", op, custom_errors.ErrNoSongInfo)
	}
	song.ReleaseDate = songInfo.ReleaseDate
	song.Text = songInfo.Text
	song.Link = songInfo.Link
	slog.Debug("Creating song", "song", song)
	err = s.songRepository.Create(c, song)
	if err != nil {
		slog.Errorf("%s: %s", op, err)
		return err
	}
	slog.Info("Song successfully created")
	return nil
}
