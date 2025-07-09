package api

import (
	"net/http"
	"strconv"

	"microx/internal/middleware"
	"microx/internal/service"

	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	timelineService service.TimelineService
}

// NewTimelineHandler crea una nueva instancia del handler de timeline
func NewTimelineHandler(timelineService service.TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
	}
}

// GetTimeline maneja la obtención del timeline personal del usuario
func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userID := middleware.GetUserID(c)

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

	tweets, err := h.timelineService.GetTimeline(c.Request.Context(), userID, limit, offset)
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

// RefreshTimeline maneja la actualización del timeline desde la base de datos
func (h *TimelineHandler) RefreshTimeline(c *gin.Context) {
	userID := middleware.GetUserID(c)

	err := h.timelineService.RefreshTimeline(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Timeline refreshed successfully",
	})
}
