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
	service services.UserService
	log     *slog.Logger
}

func NewUserController(service services.UserService, log *slog.Logger) *UserController {
	return &UserController{
		service: service,
		log:     log,
	}
}

func (uc *UserController) SetIsActive(c *gin.Context) {
	var req dto.SetIsActiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		uc.log.Error("invalid json", "error", message)
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}
	err := uc.service.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(404, models.NewErrorResponce("NOT_FOUND", "user not found"))
			return
		}

		uc.log.Error("failed to update user", "error", err)
		c.JSON(500, models.NewErrorResponce("INTERNAL_ERROR", "failed to update"))
		return
	}

	uc.log.Info("user updated", "id", req.UserID, "is_active", req.IsActive)

	c.JSON(200, gin.H{"message": "updated"})
}

func (uc *UserController) GetReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", "missing user_id"))
		return
	}

	reviews, err := uc.service.GetReviews(userID)
	if err != nil {
		c.JSON(500, models.NewErrorResponce("INTERNAL_ERROR", "failed to fetch pull requests"))
		return
	}
	uc.log.Info("fetched user reviews", "user_id", userID, "reviews_count", len(reviews))
	c.JSON(200, gin.H{
		"user_id":       userID,
		"pull_requests": reviews,
	})
}
