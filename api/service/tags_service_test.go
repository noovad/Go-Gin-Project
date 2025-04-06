package service_test

import (
	"errors"
	"go-gin-project/api/service"
	"go-gin-project/data"
	"go-gin-project/model"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTagsRepository struct {
	mock.Mock
}

func (m *MockTagsRepository) Save(tag model.Tags) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *MockTagsRepository) FindAll() ([]model.Tags, error) {
	args := m.Called()
	tags, ok := args.Get(0).([]model.Tags)
	if !ok {
		return nil, errors.New("invalid type assertion for FindAll")
	}
	return tags, args.Error(1)
}

func (m *MockTagsRepository) FindById(tagId string) (model.Tags, error) {
	args := m.Called(tagId)
	tag, ok := args.Get(0).(model.Tags)
	if !ok {
		return model.Tags{}, errors.New("invalid type assertion for FindById")
	}
	return tag, args.Error(1)
}

func (m *MockTagsRepository) Update(tag model.Tags) error {
	args := m.Called(tag)
	return args.Error(0)
}

func (m *MockTagsRepository) Delete(tagId int) error {
	args := m.Called(tagId)
	return args.Error(0)
}

func setupTest() (*MockTagsRepository, service.TagsService) {
	mockRepo := new(MockTagsRepository)
	validate := validator.New()
	tagsService := service.NewTagsServiceImpl(mockRepo, validate)
	return mockRepo, tagsService
}

func TestCreateTag(t *testing.T) {
	mockRepo, tagsService := setupTest()
	t.Run("should create a tag successfully", func(t *testing.T) {
		mockRepo.On("Save", mock.Anything).Return(nil).Once()

		tagRequest := data.TagRequest{Name: "NewTag"}
		err := tagsService.Create(tagRequest)
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail to create a tag with invalid data", func(t *testing.T) {
		tagRequest := data.TagRequest{Name: ""}
		err := tagsService.Create(tagRequest)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})
}

func TestFindAllTags(t *testing.T) {
	mockRepo, tagsService := setupTest()
	t.Run("should find all tags successfully", func(t *testing.T) {
		tags := []model.Tags{
			{Id: 1, Name: "Tag1"},
			{Id: 2, Name: "Tag2"},
		}
		mockRepo.On("FindAll").Return(tags, nil).Once()

		result, err := tagsService.FindAll()
		assert.Nil(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when finding all tags fails", func(t *testing.T) {
		mockRepo.On("FindAll").Return([]model.Tags{}, errors.New("database error")).Once()

		_, err := tagsService.FindAll()
		assert.NotNil(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestFindTagById(t *testing.T) {
	mockRepo, tagsService := setupTest()
	t.Run("should find a tag by ID successfully", func(t *testing.T) {
		mockRepo.On("FindById", "1").Return(model.Tags{Id: 1, Name: "Tag1"}, nil).Once()

		result, err := tagsService.FindById("1")
		assert.Nil(t, err)
		assert.Equal(t, "Tag1", result.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when tag ID is not found", func(t *testing.T) {
		mockRepo.On("FindById", "999").Return(model.Tags{}, errors.New("resource not found")).Once()

		_, err := tagsService.FindById("999")
		assert.NotNil(t, err)
		assert.Equal(t, "resource not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateTag(t *testing.T) {
	mockRepo, tagsService := setupTest()
	t.Run("should update a tag successfully", func(t *testing.T) {
		mockRepo.On("FindById", "1").Return(model.Tags{Id: 1, Name: "OldTag"}, nil).Once()
		mockRepo.On("Update", mock.Anything).Return(nil).Once()

		err := tagsService.Update("1", data.TagRequest{Name: "UpdatedTag"})
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail to update a tag with invalid data", func(t *testing.T) {
		err := tagsService.Update("1", data.TagRequest{Name: ""})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("should return error when tag ID is not found for update", func(t *testing.T) {
		mockRepo.On("FindById", "0").Return(model.Tags{}, errors.New("resource not found")).Once()

		err := tagsService.Update("0", data.TagRequest{Name: "test"})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "resource not found")
	})

	t.Run("should return error when updating a non-existent tag", func(t *testing.T) {
		mockRepo.On("FindById", "999").Return(model.Tags{}, nil).Once()

		err := tagsService.Update("999", data.TagRequest{Name: "NonExistentTag"})
		assert.NotNil(t, err)
		assert.Equal(t, "resource not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteTag(t *testing.T) {
	mockRepo, tagsService := setupTest()
	t.Run("should delete a tag successfully", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()

		err := tagsService.Delete(1)
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when deleting a tag fails", func(t *testing.T) {
		mockRepo.On("Delete", 999).Return(errors.New("delete failed")).Once()

		err := tagsService.Delete(999)
		assert.NotNil(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
