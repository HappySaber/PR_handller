package controllers

import (
	"PR/internal/models"
	"PR/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TeamController struct {
	service *services.TeamService
}

func NewTeamController(service *services.TeamService) *TeamController {
	return &TeamController{
		service: service,
	}
}

func (tc *TeamController) Create(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindBodyWithJSON(&team); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}
	if err := tc.service.Create(&team); err != nil {
		c.JSON(400, models.NewErrorResponce(models.ErrTeamExists, "team_name already exists"))
		return
	}
	c.JSON(201, team)
}

func (tc *TeamController) GetTeamMembers(c *gin.Context) {
	var team models.Team

	if err := c.ShouldBindJSON(&team); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}

	if err := tc.service.GetTeamMembers(&team); err != nil {
		c.JSON(404, models.NewErrorResponce(models.ErrNotFound, "team did not found"))
		return
	}

	c.JSON(200, team)
}
