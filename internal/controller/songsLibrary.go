package controller

import (
	custom_errors "Testovoe_5/internal/custom-errors"
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/service"
	"errors"

	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type songsLibrary struct {
	songsService service.ISongService
}

func newSongsRoutes(g *gin.Engine, songsService service.ISongService) {
	sL := &songsLibrary{songsService: songsService}
	g.GET("/songs", sL.Songs)
	g.GET("/text/:group_name/:song_name/", sL.Text)
}

// Songs — метод контроллера для получения списка песен с фильтрацией и пагинацией.
func (sL *songsLibrary) Songs(c *gin.Context) {
	song := model.Song{
		GroupName: c.Query("groupName"),
		SongName:  c.Query("songName"),
		Text:      c.Query("text"),
		Link:      c.Query("link"),
	}

	if releaseDate := c.Query("releaseDate"); releaseDate != "" {
		if date, err := time.Parse("2006-01-02", releaseDate); err == nil {
			song.ReleaseDate = date
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid releaseDate format. Use YYYY-MM-DD."})
			return
		}
	}

	if createdAt := c.Query("createdAt"); createdAt != "" {
		if date, err := time.Parse("2006-01-02", createdAt); err == nil {
			song.CreatedAt = date
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid createdAt format. Use YYYY-MM-DD."})
			return
		}
	}

	if updatedAt := c.Query("updatedAt"); updatedAt != "" {
		if date, err := time.Parse("2006-01-02", updatedAt); err == nil {
			song.UpdatedAt = date
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid updatedAt format. Use YYYY-MM-DD."})
			return
		}
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit. Must be a positive integer."})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset. Must be a non-negative integer."})
		return
	}

	songs, err := sL.songsService.Songs(c.Request.Context(), song, limit, offset)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No songs found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// Text — метод контроллера для получения текста песни.
func (sL *songsLibrary) Text(c *gin.Context) {
	// Извлекаем параметры из URL
	groupName := c.Param("group_name")
	songName := c.Param("song_name")

	// Извлекаем параметры пагинации из query string
	limit := c.DefaultQuery("limit", "10")  // если limit не передан, по умолчанию будет 10
	offset := c.DefaultQuery("offset", "0") // если offset не передан, по умолчанию будет 0

	// Преобразуем limit и offset в int
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	// Создаем фильтр на основе группы и песни
	song := model.Song{
		GroupName: groupName,
		SongName:  songName,
	}

	// Получаем текст песни с пагинацией
	text, err := sL.songsService.Text(c, song, limitInt, offsetInt)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		} else if errors.Is(err, custom_errors.ErrOffsetOutOfRange) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "offset exceeds the available number of verses"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Возвращаем текст песни
	c.JSON(http.StatusOK, gin.H{
		"group_name": groupName,
		"song_name":  songName,
		"text":       text,
	})
}
