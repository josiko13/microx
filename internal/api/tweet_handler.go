package api

import (
	"net/http"
	"strconv"

	"microx/internal/middleware"
	"microx/internal/model"
	"microx/internal/service"

	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	tweetService service.TweetService
}

// NewTweetHandler crea una nueva instancia del handler de tweets
func NewTweetHandler(tweetService service.TweetService) *TweetHandler {
	return &TweetHandler{
		tweetService: tweetService,
	}
}

// CreateTweet maneja la creación de un nuevo tweet
func (h *TweetHandler) CreateTweet(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.CreateTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	tweet, err := h.tweetService.CreateTweet(c.Request.Context(), userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tweet created successfully",
		"tweet":   tweet,
	})
}

// GetTweet maneja la obtención de un tweet específico
func (h *TweetHandler) GetTweet(c *gin.Context) {
	tweetIDStr := c.Param("id")
	tweetID, err := strconv.ParseInt(tweetIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid tweet ID format",
		})
		return
	}

	tweet, err := h.tweetService.GetTweet(c.Request.Context(), tweetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweet": tweet,
	})
}

// GetUserTweets maneja la obtención de tweets de un usuario específico
func (h *TweetHandler) GetUserTweets(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Obtener parámetros de paginación
	limit := 20 // Default limit
	offset := 0 // Default offset

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	tweets, err := h.tweetService.GetUserTweets(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tweets": tweets,
		"count":  len(tweets),
		"limit":  limit,
		"offset": offset,
	})
}
