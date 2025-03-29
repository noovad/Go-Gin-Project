package service

import (
	"errors"
	"go-gin-project/api/repository"
	"go-gin-project/data"
	"go-gin-project/helper"
	"go-gin-project/model"

	"github.com/go-playground/validator/v10"
)

type TagsService interface {
	Create(tag data.TagRequest) error
	FindAll() ([]data.TagResponse, error)
	FindById(tagId string) (data.TagResponse, error)
	Update(tagId string, tag data.TagRequest) error
	Delete(tagId int) error
}

func NewTagsServiceImpl(tagsRepository repository.TagsRepository, validate *validator.Validate) TagsService {
	return &TagsServiceImpl{
		TagsRepository: tagsRepository,
		Validate:       validate,
	}
}

type TagsServiceImpl struct {
	TagsRepository repository.TagsRepository
	Validate       *validator.Validate
}

func (t *TagsServiceImpl) Create(tag data.TagRequest) error {
	err := t.Validate.Struct(tag)
	if err != nil {
		return helper.ErrFailedValidationWrap(err)
	}
	tagModel := model.Tags{
		Name: tag.Name,
	}
	return t.TagsRepository.Save(tagModel)
}

func (t *TagsServiceImpl) FindAll() ([]data.TagResponse, error) {
	result, err := t.TagsRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var tags []data.TagResponse
	for _, value := range result {
		tag := data.TagResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (t *TagsServiceImpl) FindById(tagId string) (data.TagResponse, error) {
	tagData, err := t.TagsRepository.FindById(tagId)
	if err != nil {
		return data.TagResponse{}, err
	}

	tagResponse := data.TagResponse{
		Id:   tagData.Id,
		Name: tagData.Name,
	}

	return tagResponse, nil
}

func (t *TagsServiceImpl) Update(tagId string, tag data.TagRequest) error {
	err := t.Validate.Struct(tag)
	if err != nil {
		return helper.ErrFailedValidationWrap(err)
	}

	tagData, err := t.TagsRepository.FindById(tagId)
	if err != nil {
		return err
	}

	if tagData.Id == 0 {
		return helper.ErrNotFound
	}

	tagData.Name = tag.Name
	return t.TagsRepository.Update(tagData)
}

func (t *TagsServiceImpl) Delete(tagId int) error {
	err := t.TagsRepository.Delete(tagId)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return nil
}
