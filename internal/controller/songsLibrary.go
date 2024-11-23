package controller

import (
	custom_errors "Testovoe_5/internal/custom-errors"
	"Testovoe_5/internal/model"
	"Testovoe_5/internal/service"
	"errors"
	"github.com/gookit/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type songsLibrary struct {
	songsService service.ISongService
}
type SongInput struct {
	GroupName string `json:"groupName" binding:"required"`
	SongName  string `json:"songName" binding:"required"`
}

func newSongsRoutes(g *gin.RouterGroup, songsService service.ISongService) {
	sL := &songsLibrary{songsService: songsService}
	g.GET("/", sL.Songs)
	g.GET("/text/:groupName/:songName/", sL.Text)
	g.DELETE("/:groupName/:songName/", sL.Delete)
	g.PATCH("/:id/", sL.Update)
	g.POST("/:releaseDate/:text/:link/", sL.Create)
}

// Songs retrieves a list of songs with optional filtering and pagination.
// @Summary Get a list of songs
// @Tags songs
// @Description Retrieves a list of songs based on query parameters like group name, song name, and release date.
// @Accept json
// @Produce json
// @Param groupName query string false "Filter songs by group name"
// @Param songName query string false "Filter songs by song name"
// @Param text query string false "Filter songs by text content"
// @Param link query string false "Filter songs by link"
// @Param releaseDate query string false "Filter songs by release date (format: YYYY-MM-DD)"
// @Param createdAt query string false "Filter songs by creation date (format: YYYY-MM-DD)"
// @Param updatedAt query string false "Filter songs by last updated date (format: YYYY-MM-DD)"
// @Param limit query int false "Number of songs to return per page (pagination)" default(10)
// @Param offset query int false "Number of songs to skip before starting to return results (pagination)" default(0)
// @Success 200 {array} model.Song "List o songs"
// @Failure 400 {object} map[string]string "Invalid input or parameters"
// @Failure 404 {object} map[string]string "No songs found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /songs [get]
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
	slog.Debug("Received request to retrieve songs", map[string]any{
		"method":  c.Request.Method,
		"url":     c.Request.URL,
		"filters": song,
		"offset":  offset,
		"limit":   limit,
	})
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
	slog.Info("Successfully retrieved songs", "songs", songs)
}

// Text — retrieves the text of a song by group name and song name with optional pagination
// @Summary Get song text
// @Tags songs
// @Description Retrieve the text of a song by group name and song name with optional pagination
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the group"
// @Param songName path string true "Name of the song"
// @Param limit query int false "Maximum number of verses to return (default: 10)"
// @Param offset query int false "Starting verse index (default: 0)"
// @Success 200 {object} map[string]interface{} "Returns group name, song name, and text"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters or offset out of range"
// @Failure 404 {object} map[string]interface{} "Song not found "
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /songs/text/{groupName}/{songName}/ [get]
func (sL *songsLibrary) Text(c *gin.Context) {
	groupName := c.Param("groupName")
	songName := c.Param("songName")

	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")

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
	song := model.Song{
		GroupName: groupName,
		SongName:  songName,
	}
	slog.Debug("Received request to retrieve song text", map[string]any{
		"method": c.Request.Method,
		"url":    c.Request.URL,
		"song":   song,
		"offset": offsetInt,
		"limit":  limitInt,
	})
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

	c.JSON(http.StatusOK, gin.H{
		"groupName": groupName,
		"songName":  songName,
		"text":      text,
	})
	slog.Info("Successfully retrieved song text", "song", song, "text", text)
}

// Delete — delete song
// @Summary Delete a song
// @Tags songs
// @Description Delete a song by group name and song name
// @Accept json
// @Produce json
// @Param groupName path string true "Name of the group"
// @Param songName path string true "Name of the song"
// @Success 200 {string} string "Song was deleted successfully"
// @Failure 404 {object} map[string]interface{} "Song not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /songs/{groupName}/{songName} [delete]
func (sL *songsLibrary) Delete(c *gin.Context) {
	groupName := c.Param("groupName")
	songName := c.Param("songName")

	slog.Debug("Received request to delete song", map[string]any{
		"method":    c.Request.Method,
		"url":       c.Request.URL,
		"groupName": groupName,
		"songName":  songName,
	})
	err := sL.songsService.Delete(c, groupName, songName)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, "song was deleted")
	slog.Info("Successfully deleted song", "groupName", groupName, "songName", songName)
}

// Update — update song
// @Summary Partially update a song
// @Tags songs
// @Description Update specific fields of a song by ID
// @Accept json
// @Produce json
// @Param id path int true "ID of the song"
// @Param groupName query string false "Updated name of the group"
// @Param songName query string false "Updated name of the song"
// @Param text query string false "Updated text of the song"
// @Param link query string false "Updated link of the song"
// @Param release_date query string false "Updated release date of the song in YYYY-MM-DD format"
// @Success 200 {string} string "Song was updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or parameters"
// @Failure 404 {object} map[string]interface{} "Song not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /songs/{id} [patch]
func (sL *songsLibrary) Update(c *gin.Context) {
	song := model.Song{
		GroupName: c.Query("groupName"),
		SongName:  c.Query("songName"),
		Text:      c.Query("text"),
		Link:      c.Query("link"),
	}
	id := c.Param("id")

	if idInt, err := strconv.Atoi(id); err == nil {
		song.ID = idInt
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}
	if releaseDate := c.Query("release_date"); releaseDate != "" {
		if date, err := time.Parse("2006-01-02", releaseDate); err == nil {
			song.ReleaseDate = date
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid releaseDate format. Use YYYY-MM-DD."})
			return
		}
	}

	// do not update createdAt и updatedAt
	song.CreatedAt = time.Time{}
	song.UpdatedAt = time.Time{}

	slog.Debug("Received request to update song", map[string]any{
		"method":      c.Request.Method,
		"url":         c.Request.URL,
		"song":        song,
		"releaseDate": song.ReleaseDate,
	})
	err := sL.songsService.Update(c, song)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, "song was updated")
	slog.Info("Successfully updated song", "song", song)
}

// Create — create song
// @Summary Create a new song
// @Tags songs
// @Description Add a new song to the library using path parameters and request body
// @Accept json
// @Produce json
// @Param releaseDate path string true "Release date of the song in YYYY-MM-DD format"
// @Param text path string true "Lyrics of the song"
// @Param link path string true "URL to the song"
// // @Param input body SongInput true "Additional song details (groupName, songName)"
// @Success 200 {object} map[string]interface{} "Song successfully added"
// @Failure 400 {object} map[string]interface{} "Invalid request or parameters"
// @Failure 409 {object} map[string]interface{} "Song already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /songs/{releaseDate}/{text}/{link} [post]
func (sL *songsLibrary) Create(c *gin.Context) {

	releaseDateStr := c.Param("releaseDate")
	text := c.Param("text")
	link := c.Param("link")

	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text is required"})
		return
	}
	if link == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Link is required"})
		return
	}
	if releaseDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Release date is required"})
		return
	}
	var songInput SongInput
	if err := c.ShouldBindJSON(&songInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	var song model.Song
	song.Text = text
	song.Link = link
	song.SongName = songInput.SongName
	song.GroupName = songInput.GroupName
	if date, err := time.Parse("2006-01-02", releaseDateStr); err == nil {
		song.ReleaseDate = date
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid releaseDate format. Use YYYY-MM-DD."})
		return
	}
	slog.Debug("Received request to create song", map[string]any{
		"method": c.Request.Method,
		"url":    c.Request.URL,
		"song":   song,
	})
	err := sL.songsService.Create(c, song)
	if err != nil {
		if errors.Is(err, custom_errors.ErrNoSongInfo) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song information from external API"})
		} else if errors.Is(err, custom_errors.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Song already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song successfully added"})
	slog.Info("Successfully added song", "song", song)
}
