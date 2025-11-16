package routes

import (
	"PR/internal/controllers"
	"PR/internal/services"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, log *slog.Logger) {
	testS := services.NewTestUserService()
	testC := controllers.NewTestUserController(testS, log)

	teamS := services.NewTeamService()
	teamC := controllers.NewTeamController(teamS, log)

	userS := services.NewUserService()
	userC := controllers.NewUserController(userS, log)

	prS := services.NewPullRequestService()
	prC := controllers.NewPullRequestController(prS, log)

	r.POST("/test/add", testC.CreateTestUsers)
	r.DELETE("/test/add", testC.DeleteTestUsers)

	teams := r.Group("/team")
	{
		teams.POST("/add", teamC.Create)
		teams.GET("/get", teamC.GetTeamMembers)
	}

	users := r.Group("/users")
	{
		users.GET("/getReview", userC.GetReviews)
		users.PUT("/setIsActive", userC.SetIsActive)
	}

	prs := r.Group("/pullRequest")
	{
		prs.POST("/create", prC.Create)
		prs.PUT("/merge", prC.Merge)
		prs.PUT("/reassign", prC.Reassign)
	}

	r.GET("/stats/reviews", userC.GetReviewStats)

}
