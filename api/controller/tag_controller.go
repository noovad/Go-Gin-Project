package controller

import (
	"errors"
	"go-gin-project/api/service"
	"go-gin-project/data"
	"go-gin-project/helper"
	"go-gin-project/helper/responsejson"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagsController struct {
	tagsService service.TagsService
}

func NewTagsController(service service.TagsService) *TagsController {
	return &TagsController{
		tagsService: service,
	}
}

func (controller *TagsController) Create(ctx *gin.Context) {
	createTagsRequest := data.TagRequest{}
	err := ctx.ShouldBindJSON(&createTagsRequest)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}

	err = controller.tagsService.Create(createTagsRequest)
	if err != nil {
		if errors.Is(err, helper.ErrFailedValidation) {
			responsejson.BadRequest(ctx, err)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}
	responsejson.Success(ctx, "create", nil)
}

func (controller *TagsController) FindAll(ctx *gin.Context) {
	tagResponse, err := controller.tagsService.FindAll()
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}
	responsejson.Success(ctx, "read", tagResponse)
}

func (controller *TagsController) FindById(ctx *gin.Context) {
	tagId := ctx.Param("tagId")

	tagResponse, err := controller.tagsService.FindById(tagId)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			responsejson.NotFound(ctx, err.Error())
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}
	responsejson.Success(ctx, "read", tagResponse)
}

func (controller *TagsController) Update(ctx *gin.Context) {
	tagId := ctx.Param("tagId")
	updateTagsRequest := data.TagRequest{}
	err := ctx.ShouldBindJSON(&updateTagsRequest)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}
	err = controller.tagsService.Update(tagId, updateTagsRequest)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			responsejson.NotFound(ctx, err.Error())
			return
		}
		if errors.Is(err, helper.ErrFailedValidation) {
			responsejson.BadRequest(ctx, err)
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "update", nil)
}

func (controller *TagsController) Delete(ctx *gin.Context) {
	tagId := ctx.Param("tagId")
	id, err := strconv.Atoi(tagId)
	if err != nil {
		responsejson.InternalServerError(ctx, err)
		return
	}
	err = controller.tagsService.Delete(id)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			responsejson.NotFound(ctx, err.Error())
			return
		}
		responsejson.InternalServerError(ctx, err)
		return
	}

	responsejson.Success(ctx, "delete", nil)
}
