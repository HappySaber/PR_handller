package controllers

import (
	"PR/internal/controllers/dto"
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
	var req dto.CreateTeamRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		message := fmt.Sprintf("invalid json: %v", err.Error())
		tc.log.Error("invalid json", "error", message)
		c.JSON(400, models.NewErrorResponce("INVALID_REQUEST", message))
		return
	}

	team := models.Team{
		Name:    req.TeamName,
		Members: make([]models.TeamMember, len(req.Members)),
	}

	for i, m := range req.Members {
		team.Members[i] = models.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	if err := tc.service.Create(&team); err != nil {
		c.JSON(400, models.NewErrorResponce(models.ErrTeamExists, "team_name already exists"))
		return
	}

	tc.log.Info("team created", "team", team.Name)
	c.JSON(201, req)
}

func (tc *TeamController) GetTeamMembers(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(400, gin.H{"error": gin.H{
			"code":    "INVALID_REQUEST",
			"message": "team_name query parameter required",
		}})
		return
	}
	team, err := tc.service.GetTeamMembers(teamName)
	if err != nil {
		c.JSON(404, gin.H{"error": gin.H{
			"code":    "NOT_FOUND",
			"message": "team not found",
		}})
		return
	}
	resp := dto.TeamResponse{
		TeamName: team.Name,
		Members:  make([]dto.TeamMemberDTO, len(team.Members)),
	}

	for i, m := range team.Members {
		resp.Members[i] = dto.TeamMemberDTO{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	tc.log.Info("got team members", "team", team.Name)
	c.JSON(200, resp)
}
