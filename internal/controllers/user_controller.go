package controllers

import (
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
	type SetIsActiveRequest struct {
		UserID   string `json:"user_id" binding:"required"`
		IsActive bool   `json:"is_active" binding:"required"`
	}

	var req SetIsActiveRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		uc.log.Error("invalid json", "error", message)
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}

	err := uc.service.SetIsActive()
	if err != nil {
		uc.log.Error("failed to set is_active", "error", err)
		c.JSON(500, models.NewErrorResponce("INTERNAL_ERROR", "failed to update"))
		return
	}

	c.JSON(200, gin.H{"message": "updated"})
}
