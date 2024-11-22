package repository

import (
	custom_errors "Testovoe_5/internal/custom-errors"
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/pkg/postgres"
	"context"
	"database/sql"
	"errors"
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

func (sR *SongsRepository) Songs(c context.Context, filter model.Song, limit, offset int) (songs []model.Song, err error) {
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

func (sR *SongsRepository) Text(c context.Context, song model.Song, limit, offset int) (text string, err error) {
	const op = "repository.Text"

	query := "SELECT text FROM public.songs WHERE group_name = $1 AND song_name = $2"
	var songText string

	err = sR.DB.Pool.QueryRow(c, query, song.GroupName, song.SongName).Scan(&songText)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, custom_errors.ErrNoRows)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}
	couplets := strings.Split(songText, "\n\n")

	start := offset
	end := limit

	if start >= len(couplets) || start < 0 {
		return "", custom_errors.ErrOffsetOutOfRange
	}

	if end > len(couplets) {
		end = len(couplets)
	}
	var selectedVerses []string
	if start == end {
		selectedVerses = couplets[start : end+1]
	} else {
		selectedVerses = couplets[start:end]
	}

	text = strings.Join(selectedVerses, "\n")
	return text, nil

}

func (sR *SongsRepository) Delete(c context.Context, groupName, songName string) error {
	const op = "repository.Delete"

	query := "DELETE FROM public.songs WHERE group_name = $1 AND song_name = $2"

	result, err := sR.DB.Pool.Exec(c, query, groupName, songName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, custom_errors.ErrNoRows)
	}

	return nil
}

func (sR *SongsRepository) Update(c context.Context, song model.Song) error {
	const op = "repository.Update"

	query := "UPDATE public.songs SET updated_at = NOW()"
	args := []interface{}{}
	paramIndex := 1

	if song.GroupName != "" {
		query += fmt.Sprintf(", group_name = $%v", paramIndex)
		args = append(args, song.GroupName)
		paramIndex++
	}
	if song.SongName != "" {
		query += fmt.Sprintf(", song_name = $%v", paramIndex)
		args = append(args, song.SongName)
		paramIndex++
	}
	if !song.ReleaseDate.IsZero() {
		query += fmt.Sprintf(", release_date = $%v", paramIndex)
		args = append(args, song.ReleaseDate)
		paramIndex++
	}
	if song.Text != "" {
		query += fmt.Sprintf(", text = $%v", paramIndex)
		args = append(args, song.Text)
		paramIndex++
	}
	if song.Link != "" {
		query += fmt.Sprintf(", link = $%v", paramIndex)
		args = append(args, song.Link)
		paramIndex++
	}

	query += " WHERE id = $" + fmt.Sprint(len(args)+1)
	slog.Info("query", query, "len", len(args), "args", args)
	args = append(args, song.ID)

	result, err := sR.DB.Pool.Exec(c, query, args...)
	if err != nil {
		return fmt.Errorf("%s: failed to update song with ID %d: %w", op, song.ID, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, custom_errors.ErrNoRows)
	}

	return nil
}

func (sR *SongsRepository) Create(c context.Context, song model.Song) error {
	const op = "repository.Create"

	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	slog.Info("song", song.Text)
	_, err := sR.DB.Pool.Query(c, query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)

	}
	return nil
}
