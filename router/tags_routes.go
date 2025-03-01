package router

import (
	"todo-project/api/controller"
	"todo-project/api/repository"
	"todo-project/api/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func TagsRouter(db *gorm.DB, router *gin.Engine, validate *validator.Validate) {
	tagRepository := repository.NewTagsREpositoryImpl(db)
	tagService := service.NewTagsServiceImpl(tagRepository, validate)
	tagsController := controller.NewTagsController(tagService)

	tagsRouter := router.Group("/tags")
	{
		tagsRouter.GET("", tagsController.FindAll)
		tagsRouter.GET("/:tagId", tagsController.FindById)
		tagsRouter.POST("", tagsController.Create)
		tagsRouter.PATCH("/:tagId", tagsController.Update)
		tagsRouter.PATCH("/:tagId", tagsController.Update)
		tagsRouter.DELETE("/:tagId", tagsController.Delete)
	}
}
