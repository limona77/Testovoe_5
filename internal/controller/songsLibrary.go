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
}

// Songs — метод контроллера для получения списка песен с фильтрацией и пагинацией.
func (sL *songsLibrary) Songs(c *gin.Context) {
	filter := model.Song{
		GroupName: c.Query("groupName"),
		SongName:  c.Query("songName"),
		Text:      c.Query("text"),
		Link:      c.Query("link"),
	}

	if releaseDate := c.Query("releaseDate"); releaseDate != "" {
		if date, err := time.Parse("2006-01-02", releaseDate); err == nil {
			filter.ReleaseDate = date
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid releaseDate format. Use YYYY-MM-DD."})
			return
		}
	}

	if createdAt := c.Query("createdAt"); createdAt != "" {
		if date, err := time.Parse("2006-01-02", createdAt); err == nil {
			filter.CreatedAt = date
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid createdAt format. Use YYYY-MM-DD."})
			return
		}
	}

	if updatedAt := c.Query("updatedAt"); updatedAt != "" {
		if date, err := time.Parse("2006-01-02", updatedAt); err == nil {
			filter.UpdatedAt = date
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

	songs, err := sL.songsService.Songs(c.Request.Context(), filter, limit, offset)
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
