package controllers

import (
	"PR/internal/controllers/dto"
	"PR/internal/models"
	"PR/internal/services"
	"database/sql"
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
		message := fmt.Sprintf("invalid JSON: %v", err)
		tc.log.Warn("CreateTeam: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
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
		if err.Error() == "team exists" {
			tc.log.Warn("CreateTeam: team already exists", slog.String("team_name", req.TeamName))
			c.JSON(400, models.NewErrorResponse(models.ErrTeamExists, "team_name already exists"))
			return
		}
		tc.log.Error("CreateTeam: failed to create team", slog.String("error", err.Error()))
		c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to create team"))
		return
	}

	tc.log.Info("CreateTeam: team created successfully", slog.String("team_name", req.TeamName))
	c.JSON(201, gin.H{"team": team})
}

func (tc *TeamController) GetTeamMembers(c *gin.Context) {
	type GetTeamMembersRequest struct {
		TeamName string `json:"team_name" binding:"required"`
	}

	var req GetTeamMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := "team_name is required in JSON body"
		tc.log.Warn("GetTeamMembers: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}

	team, err := tc.service.GetTeamMembers(req.TeamName)
	tc.log.Info("params from JSON", "team_name", req.TeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			tc.log.Warn("GetTeamMembers: team not found", slog.String("team_name", req.TeamName))
			c.JSON(404, models.NewErrorResponse(models.ErrNotFound, "team not found"))
			return
		}
		tc.log.Error("GetTeamMembers: failed to fetch team members", slog.String("error", err.Error()))
		c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to fetch team members"))
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

	tc.log.Info("GetTeamMembers: fetched team members successfully", slog.String("team_name", team.Name))
	c.JSON(200, gin.H{"team": resp})
}
