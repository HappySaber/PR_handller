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
	service *services.PullRequestService
	log     *slog.Logger
}

func NewPullRequestController(service *services.PullRequestService, log *slog.Logger) *PullRequestController {
	return &PullRequestController{
		service: service,
		log:     log,
	}
}

func (prc *PullRequestController) Create(c *gin.Context) {
	var req dto.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := fmt.Sprintf("invalid JSON: %v", err)
		prc.log.Warn("CreatePR: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}

	pr := models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
	}

	if err := prc.service.Create(&pr); err != nil {
		if err.Error() == "author not found" {
			prc.log.Warn("CreatePR: author not found", slog.String("author_id", req.AuthorID))
			c.JSON(404, models.NewErrorResponse(models.ErrNotFound, "author or team not found"))
			return
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			prc.log.Warn("CreatePR: PR already exists", slog.String("pr_id", req.PullRequestID))
			c.JSON(409, models.NewErrorResponse(models.ErrPRExistst, "PR id already exists"))
			return
		}
		prc.log.Error("CreatePR: failed to create PR", slog.String("error", err.Error()))
		c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to create PR"))
		return
	}

	c.JSON(201, gin.H{"pr": pr})
}

func (prc *PullRequestController) Merge(c *gin.Context) {
	var req dto.MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := fmt.Sprintf("invalid JSON: %v", err)
		prc.log.Warn("MergePR: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}

	pr, err := prc.service.Merge(req.PullRequestID)
	if err != nil {
		prc.log.Error("MergePR: failed to merge PR", slog.String("error", err.Error()))
		c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to merge PR"))
		return
	}
	if pr == nil {
		prc.log.Warn("MergePR: PR not found or already merged", slog.String("pr_id", req.PullRequestID))
		c.JSON(404, models.NewErrorResponse(models.ErrNotFound, "PR not found or already merged"))
		return
	}

	resp := dto.PullRequestResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		MergedAt:          pr.MergedAt,
	}

	prc.log.Info("MergePR: PR merged successfully", slog.String("pr_id", pr.PullRequestID))
	c.JSON(200, gin.H{"pr": resp})
}

func (prc *PullRequestController) Reassign(c *gin.Context) {
	var req dto.ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := fmt.Sprintf("invalid JSON: %v", err)
		prc.log.Warn("ReassignPR: invalid request", slog.String("error", message))
		c.JSON(400, models.NewErrorResponse(models.ErrInvalidReq, message))
		return
	}

	pr, replacedBy, err := prc.service.Reassign(req.PullRequestID, req.OldUserID)
	if err != nil {
		switch err.Error() {
		case "not found":
			prc.log.Warn("ReassignPR: PR not found", slog.String("pr_id", req.PullRequestID))
			c.JSON(404, models.NewErrorResponse(models.ErrNotFound, "PR not found"))
		case "PR_MERGED":
			prc.log.Warn("ReassignPR: PR already merged", slog.String("pr_id", req.PullRequestID))
			c.JSON(409, models.NewErrorResponse(models.ErrPRMerged, "cannot reassign on merged PR"))
		case "NOT_ASSIGNED":
			prc.log.Warn("ReassignPR: reviewer not assigned", slog.String("pr_id", req.PullRequestID), slog.String("user_id", req.OldUserID))
			c.JSON(409, models.NewErrorResponse(models.ErrNotAssigned, "reviewer is not assigned to this PR"))
		case "NO_CANDIDATE":
			prc.log.Warn("ReassignPR: no candidate found", slog.String("pr_id", req.PullRequestID))
			c.JSON(409, models.NewErrorResponse(models.ErrNoCandidate, "no active replacement candidate in team"))
		default:
			prc.log.Error("ReassignPR: failed to reassign", slog.String("error", err.Error()))
			c.JSON(500, models.NewErrorResponse(models.ErrInternal, "failed to reassign PR"))
		}
		return
	}

	resp := dto.PullRequestResponse{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}

	prc.log.Info("ReassignPR: reviewer reassigned successfully", slog.String("pr_id", pr.PullRequestID), slog.String("replaced_by", replacedBy))
	c.JSON(200, gin.H{
		"pr":          resp,
		"replaced_by": replacedBy,
	})
}
