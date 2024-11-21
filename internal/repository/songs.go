package repository

import (
	custom_errors "Testovoe_5/internal/custom-errors"
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/pkg/postgres"
	"context"
	"fmt"
	"github.com/gookit/slog"

	"strings"
)

type SongsRepository struct {
	*postgres.DB
}

func NewSongsRepository(db *postgres.DB) *SongsRepository {
	return &SongsRepository{db}
}

func (sR *SongsRepository) Songs(c context.Context, filter model.Song, limit, offset int) ([]model.Song, error) {
	const op = "repository.GetSongs"

	query := "SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at FROM public.songs"

	var args []interface{}
	conditions := []string{}
	paramIndex := 1

	if filter.GroupName != "" {
		conditions = append(conditions, fmt.Sprintf("group_name = $%v", paramIndex))
		args = append(args, filter.GroupName)
		paramIndex++
	}
	if filter.SongName != "" {
		conditions = append(conditions, fmt.Sprintf("song_name = $%v", paramIndex))
		args = append(args, filter.SongName)
		paramIndex++
	}
	if !filter.ReleaseDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("release_date = $%v", paramIndex))
		args = append(args, filter.ReleaseDate)
		paramIndex++
	}
	if filter.Text != "" {
		conditions = append(conditions, fmt.Sprintf("text ILIKE $%v", paramIndex))
		args = append(args, "%"+filter.Text+"%")
		paramIndex++
	}
	if filter.Link != "" {
		conditions = append(conditions, fmt.Sprintf("link ILIKE $%v", paramIndex))
		args = append(args, "%"+filter.Link+"%")
		paramIndex++
	}
	if !filter.CreatedAt.IsZero() {
		conditions = append(conditions, fmt.Sprintf("created_at = $%v", paramIndex))
		args = append(args, filter.CreatedAt)
		paramIndex++
	}
	if !filter.UpdatedAt.IsZero() {
		conditions = append(conditions, fmt.Sprintf("updated_at = $%v", paramIndex))
		args = append(args, filter.UpdatedAt)
		paramIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY release_date DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	args = append(args, limit, offset)
	rows, err := sR.DB.Pool.Query(c, query, args...)
	if err != nil {
		return []model.Song{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var songs []model.Song
	for rows.Next() {
		var song model.Song
		err = rows.Scan(
			&song.ID,
			&song.SongName,
			&song.GroupName,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			&song.CreatedAt,
			&song.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, song)
	}
	if songs == nil {
		return []model.Song{}, custom_errors.ErrNoRows
	}
	slog.Info(songs)
	return songs, nil
}
