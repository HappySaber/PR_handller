package routes

import (
	"PR/internal/controllers"
	"PR/internal/services"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, log *slog.Logger) {
	// создаём сервисы и контроллеры прямо здесь
	testS := services.NewTestUserService()
	testC := controllers.NewTestUserController(testS, log)

	teamS := services.NewTeamService()
	teamC := controllers.NewTeamController(teamS, log)

	userS := services.NewUserService()
	userC := controllers.NewUserController(userS, log)

	prS := services.NewPullRequestService()
	prC := controllers.NewPullRequestController(prS, log)

	// test
	r.POST("/test/add", testC.CreateTestUsers)
	r.DELETE("/test/add", testC.DeleteTestUsers)

	// team
	teams := r.Group("/team")
	{
		teams.POST("/add", teamC.Create)
		teams.GET("/get", teamC.GetTeamMembers)
	}

	// users
	users := r.Group("/users")
	{
		users.GET("/getReview", userC.GetReviews)
		users.PUT("/setisactive", userC.SetIsActive)
	}

	// pull requests
	prs := r.Group("/pull")
	{
		prs.POST("/create", prC.Create)
		prs.PUT("/merge", prC.Merge)
		prs.PUT("/reassign", prC.Reassign)
	}
}
