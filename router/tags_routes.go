package router

import (
	"go-gin-project/api"

	"github.com/gin-gonic/gin"
)

func TagsRouter(router *gin.Engine) {
	controller := api.InitializeTagsController()

	tagsRouter := router.Group("/tags")
	{
		tagsRouter.GET("", controller.FindAll)
		tagsRouter.GET("/:tagId", controller.FindById)
		tagsRouter.POST("", controller.Create)
		tagsRouter.PATCH("/:tagId", controller.Update)
		tagsRouter.DELETE("/:tagId", controller.Delete)
	}
}
