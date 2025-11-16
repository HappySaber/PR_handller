package controllers

import (
	"PR/internal/controllers/dto"
	"PR/internal/models"
	"PR/internal/services"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
	log     *slog.Logger
}

func NewUserController(service *services.UserService, log *slog.Logger) *UserController {
	return &UserController{
		service: service,
		log:     log,
	}
}

func (uc *UserController) SetIsActive(c *gin.Context) {
	var req dto.SetIsActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		message := fmt.Sprintf("invalid json: %v", err)
		uc.log.Error("SetIsActive: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}
	userResp, err := uc.service.SetIsActive(req.UserID, *req.IsActive)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(404, models.NewErrorResponse("NOT_FOUND", "user not found"))
			return
		}

		uc.log.Error("SetIsActive: failed to update user", slog.String("error", err.Error()), slog.String("user_id", req.UserID))
		c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to update user"))
		return
	}

	uc.log.Info("SetIsActive: user updated successfully",
		slog.String("user_id", req.UserID),
		slog.Bool("is_active", *req.IsActive),
	)

	c.JSON(200, gin.H{"user": userResp})
}

func (uc *UserController) GetReviews(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		uc.log.Warn("GetReviews: invalid JSON", slog.String("error", err.Error()))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, "missing or invalid user_id"))
		return
	}

	reviews, err := uc.service.GetReviews(req.UserID)
	if err != nil {
		uc.log.Error("GetReviews: failed to fetch pull requests", slog.String("error", err.Error()), slog.String("user_id", req.UserID))
		c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to fetch pull requests"))
		return
	}

	uc.log.Info("GetReviews: fetched reviews successfully",
		slog.String("user_id", req.UserID),
		slog.Int("reviews_count", len(reviews)),
	)

	c.JSON(200, gin.H{
		"user_id":       req.UserID,
		"pull_requests": reviews,
	})
}

func (uc *UserController) GetReviewStats(c *gin.Context) {
	stats, err := uc.service.GetReviewStats()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch review stats"})
		return
	}

	c.JSON(200, gin.H{"stats": stats})
}
