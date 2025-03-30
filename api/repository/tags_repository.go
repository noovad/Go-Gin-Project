package repository

import (
	"errors"
	"go-gin-project/helper"
	"go-gin-project/model"

	"gorm.io/gorm"
)

type TagsRepository interface {
	Save(tag model.Tags) error
	FindAll() ([]model.Tags, error)
	FindById(tagId string) (tag model.Tags, err error)
	Update(tag model.Tags) error
	Delete(tagId int) error
}

func NewTagsRepositoryImpl(Db *gorm.DB) TagsRepository {
	return &TagsRepositoryImpl{Db: Db}
}

type TagsRepositoryImpl struct {
	Db *gorm.DB
}

func (t *TagsRepositoryImpl) Save(tag model.Tags) error {
	result := t.Db.Create(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *TagsRepositoryImpl) FindAll() ([]model.Tags, error) {
	var tags []model.Tags
	result := t.Db.Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

func (t *TagsRepositoryImpl) FindById(tagId string) (tagModel model.Tags, err error) {
	var tag model.Tags
	result := t.Db.First(&tag, tagId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.Tags{}, helper.ErrNotFound
	} else if result.Error != nil {
		return model.Tags{}, result.Error
	}

	return tag, nil
}

func (t *TagsRepositoryImpl) Update(tags model.Tags) error {
	result := t.Db.Model(&tags).Updates(tags)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *TagsRepositoryImpl) Delete(tagsId int) error {
	deleteResult := t.Db.Delete(&model.Tags{}, tagsId)
	if deleteResult.Error != nil {
		return deleteResult.Error
	}
	if deleteResult.RowsAffected == 0 {
		return helper.ErrNotFound
	}
	return nil
}
