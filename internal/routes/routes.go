package routes

import "github.com/gin-gonic/gin"

func Routes(r *gin.Engine) {
	teams := r.Group("/team").Use()
	{
		teams.POST("/add")
		teams.GET("/get")
	}
}
