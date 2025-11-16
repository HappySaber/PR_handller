package controllers

import (
	"PR/internal/controllers/dto"
	"PR/internal/models"
	"PR/internal/services"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type PullRequestController struct {
	service services.PullRequestService
	log     *slog.Logger
}

func NewPullRequestController(service services.PullRequestService, log *slog.Logger) *PullRequestController {
	return &PullRequestController{
		service: service,
		log:     log,
	}
}

func (prc *PullRequestController) Create(c *gin.Context) {
	var req dto.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		prc.log.Error("invalid json", "error", message)
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}

	pr := models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	if err := prc.service.Create(&pr); err != nil {
		if err.Error() == "author not found" {
			c.JSON(404, models.NewErrorResponce(models.ErrNotFound, "author or team not found"))
			return
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			c.JSON(409, models.NewErrorResponce(models.ErrPRExistst, "PR id already exists"))
			return
		}
		prc.log.Error("failed to create PR", "error", err)
		c.JSON(500, models.NewErrorResponce("INTERNAL_ERROR", "failed to create PR"))
		return
	}

	c.JSON(201, gin.H{"pr": pr})
}
