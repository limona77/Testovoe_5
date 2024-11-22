package api

import (
	"Testovoe_5/internal/config"
	"Testovoe_5/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest = errors.New("incorrect request")
	ErrNoResponse = errors.New("no response from API")
)

type ApiClient struct {
	config *config.Config
}

func NewApiClient(config *config.Config) *ApiClient {
	return &ApiClient{config: config}
}

func (ac *ApiClient) GetSongInfo(group, song string) (*model.Song, error) {
	resp, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", ac.config.ExternalApi, group, song))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, ErrBadRequest
		} else {
			return nil, ErrNoResponse
		}
	}

	var songInfo model.Song
	if err := json.NewDecoder(resp.Body).Decode(&songInfo); err != nil {
		return nil, err
	}

	return &songInfo, nil
}
