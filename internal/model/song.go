package model

import "time"

type Song struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"groupName"`
	SongName    string    `json:"songName"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
