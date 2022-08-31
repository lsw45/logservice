package main

import (
	"log-ext/adapter/controller"
	"log-ext/domain"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	addRouter(engine)

	engine.Run(":8080")
}

func addRouter(r *gin.Engine) {
	search := r.Group("/search")
	{
		searchCtl := controller.NewSearchController(&domain.SearchService{})

		search.GET("/list", searchCtl.List)
	}
	job := r.Group("/job")
	job.GET("/get", controller.JobController.Get)
}
