package controllers

import (
	"PR/internal/models"
	"PR/internal/services"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type TestUserController struct {
	service *services.TestUserService
	log     *slog.Logger
}

func NewTestUserController(service *services.TestUserService, log *slog.Logger) *TestUserController {
	return &TestUserController{
		service: service,
		log:     log,
	}
}

func (tuc *TestUserController) CreateTestUsers(c *gin.Context) {
	if os.Getenv("ENV") != "local" {
		message := "Can be done only in local"
		tuc.log.Error("CreateTestUsers: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}

	users, err := tuc.service.CreateTestUsers(8)
	if err != nil {
		tuc.log.Error("failed to create test users", "error", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("INTERNAL_ERROR", "failed to create test users"))
		return
	}

	tuc.log.Info("created test users", "count", len(users))
	c.JSON(http.StatusCreated, gin.H{"users": users})
}

func (tuc *TestUserController) DeleteTestUsers(c *gin.Context) {
	if os.Getenv("ENV") != "local" {
		message := "Can be done only in local"
		tuc.log.Error("DeleteTestUsers: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}
	if err := tuc.service.DeleteTestUsers(); err != nil {
		tuc.log.Error("failed to delete test users", "error", err)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("INTERNAL_ERROR", "failed to delete test users"))
		return
	}

	tuc.log.Info("deleted test users")
	c.JSON(http.StatusOK, gin.H{"deleted_ids": ""})
}
