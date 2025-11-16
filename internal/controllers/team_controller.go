package controllers

import (
	"PR/internal/models"
	"PR/internal/services"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	service *services.TeamService
	log     *slog.Logger
}

func NewTeamController(service *services.TeamService, log *slog.Logger) *TeamController {
	return &TeamController{
		service: service,
		log:     log,
	}
}

func (tc *TeamController) Create(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindBodyWithJSON(&team); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		tc.log.Error("invalid json", "error", message)
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}
	if err := tc.service.Create(&team); err != nil {
		c.JSON(400, models.NewErrorResponce(models.ErrTeamExists, "team_name already exists"))
		return
	}
	tc.log.Info("team created", "team", team.Name)
	c.JSON(201, team)
}

func (tc *TeamController) GetTeamMembers(c *gin.Context) {
	var team models.Team

	if err := c.ShouldBindJSON(&team); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		tc.log.Error("invalid json", "error", message)
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}

	if err := tc.service.GetTeamMembers(&team); err != nil {
		c.JSON(404, models.NewErrorResponce(models.ErrNotFound, "team did not found"))
		return
	}
	tc.log.Info("got team members", "team", team.Name, "team members", team.Members)
	c.JSON(200, team)
}
