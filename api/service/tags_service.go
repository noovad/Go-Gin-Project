package service

import (
	"go-gin-project/api/repository"
	"go-gin-project/data"
	"go-gin-project/helper"
	"go-gin-project/model"

	"github.com/go-playground/validator/v10"
)

type TagsService interface {
	Create(tags data.CreateTagsRequest)
	Update(tags data.UpdateTagsRequest)
	Delete(tagsId int)
	FindById(tagsId int) data.TagsResponse
	FindAll() []data.TagsResponse
}

func NewTagsServiceImpl(tagRepository repository.TagsRepository, validate *validator.Validate) TagsService {
	return &TagsServiceImpl{
		TagsRepository: tagRepository,
		Validate:       validate,
	}
}

type TagsServiceImpl struct {
	TagsRepository repository.TagsRepository
	Validate       *validator.Validate
}

func (t *TagsServiceImpl) Create(tags data.CreateTagsRequest) {
	err := t.Validate.Struct(tags)
	helper.ErrorPanic(err)
	tagModel := model.Tags{
		Name: tags.Name,
	}
	t.TagsRepository.Save(tagModel)
}

func (t *TagsServiceImpl) Delete(tagsId int) {
	t.TagsRepository.Delete(tagsId)
}

func (t *TagsServiceImpl) FindAll() []data.TagsResponse {
	result := t.TagsRepository.FindAll()

	var tags []data.TagsResponse
	for _, value := range result {
		tag := data.TagsResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *TagsServiceImpl) FindById(tagsId int) data.TagsResponse {
	tagData, err := t.TagsRepository.FindById(tagsId)
	helper.ErrorPanic(err)

	tagResponse := data.TagsResponse{
		Id:   tagData.Id,
		Name: tagData.Name,
	}
	return tagResponse
}

func (t *TagsServiceImpl) Update(tags data.UpdateTagsRequest) {
	tagData, err := t.TagsRepository.FindById(tags.Id)
	helper.ErrorPanic(err)
	tagData.Name = tags.Name
	t.TagsRepository.Update(tagData)
}
