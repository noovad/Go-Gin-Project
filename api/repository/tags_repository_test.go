package repository_test

import (
	"go-gin-project/api/repository"
	"go-gin-project/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Tags{})
	return db
}

func createMockData(db *gorm.DB) {
	mockTags := []model.Tags{
		{Id: 1, Name: "Tag1"},
		{Id: 2, Name: "Tag2"},
	}
	db.Create(&mockTags)
}

func TestSaveTag(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTagsRepositoryImpl(db)
	t.Run("should save new tag", func(t *testing.T) {
		tag := model.Tags{Id: 3, Name: "TestTag"}
		err := repo.Save(tag)
		assert.Nil(t, err)

		var count int64
		db.Model(&model.Tags{}).Where("id = ?", tag.Id).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("should return error when db is closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		tag := model.Tags{Id: 4, Name: "ErrorTag"}
		err := repo.Save(tag)
		assert.NotNil(t, err)
	})
}

func TestFindAllTags(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTagsRepositoryImpl(db)
	t.Run("should return all tags", func(t *testing.T) {
		createMockData(db)
		tags, err := repo.FindAll()
		assert.Nil(t, err)
		assert.Len(t, tags, 2)
	})

	t.Run("should return error when db is closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		_, err := repo.FindAll()
		assert.NotNil(t, err)
	})
}

func TestFindTagById(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTagsRepositoryImpl(db)
	createMockData(db)
	t.Run("should find existing tag", func(t *testing.T) {
		tag, err := repo.FindById("1")
		assert.Nil(t, err)
		assert.Equal(t, "Tag1", tag.Name)
	})

	t.Run("should return not found for non-existent tag", func(t *testing.T) {
		_, err := repo.FindById("999")
		assert.NotNil(t, err)
		assert.Equal(t, "resource not found", err.Error())
	})

	t.Run("should return error when db is closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		_, err := repo.FindById("1")
		assert.NotNil(t, err)
	})
}

func TestUpdateTag(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTagsRepositoryImpl(db)
	t.Run("should update existing tag", func(t *testing.T) {
		createMockData(db)
		updatedTag := model.Tags{Id: 1, Name: "UpdatedTag"}
		err := repo.Update(updatedTag)
		assert.Nil(t, err)

		var tag model.Tags
		db.First(&tag, "1")
		assert.Equal(t, "UpdatedTag", tag.Name)
	})

	t.Run("should return error when db is closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		updatedTag := model.Tags{Id: 1, Name: "ErrorTag"}
		err := repo.Update(updatedTag)
		assert.NotNil(t, err)
	})
}

func TestDeleteTag(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTagsRepositoryImpl(db)
	t.Run("should delete existing tag", func(t *testing.T) {
		createMockData(db)
		err := repo.Delete(1)
		assert.Nil(t, err)

		var count int64
		db.Model(&model.Tags{}).Where("id = ?", 1).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("should return not found error for non-existent tag", func(t *testing.T) {
		err := repo.Delete(999)
		assert.NotNil(t, err)
		assert.Equal(t, "resource not found", err.Error())
	})

	t.Run("should return error when db is closed", func(t *testing.T) {
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.Delete(1)
		assert.NotNil(t, err)
	})
}
