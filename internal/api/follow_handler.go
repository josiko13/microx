package api

import (
	"net/http"
	"strconv"

	"microx/internal/middleware"
	"microx/internal/service"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService service.FollowService
}

// NewFollowHandler crea una nueva instancia del handler de follows
func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// FollowUser maneja la acción de seguir a un usuario
func (h *FollowHandler) FollowUser(c *gin.Context) {
	followerID := middleware.GetUserID(c)

	followingIDStr := c.Param("user_id")
	followingID, err := strconv.ParseInt(followingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.followService.FollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully followed user",
	})
}

// UnfollowUser maneja la acción de dejar de seguir a un usuario
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followerID := middleware.GetUserID(c)

	followingIDStr := c.Param("user_id")
	followingID, err := strconv.ParseInt(followingIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.followService.UnfollowUser(c.Request.Context(), followerID, followingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unfollowed user",
	})
}

// GetFollowers maneja la obtención de seguidores de un usuario
func (h *FollowHandler) GetFollowers(c *gin.Context) {
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

	followers, err := h.followService.GetFollowers(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"count":     len(followers),
		"limit":     limit,
		"offset":    offset,
	})
}

// GetFollowing maneja la obtención de usuarios seguidos
func (h *FollowHandler) GetFollowing(c *gin.Context) {
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

	following, err := h.followService.GetFollowing(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"following": following,
		"count":     len(following),
		"limit":     limit,
		"offset":    offset,
	})
}
