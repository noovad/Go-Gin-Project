//go:build wireinject
// +build wireinject

// wire_gen.go

package api

import (
	"go-gin-project/api/controller"
	"go-gin-project/api/repository"
	"go-gin-project/api/service"
	"go-gin-project/config"

	"github.com/google/wire"
)

func InitializeTagsController() *controller.TagsController {
	wire.Build(
		controller.NewTagsController,
		service.NewTagsServiceImpl,
		repository.NewTagsREpositoryImpl,
		config.DatabaseConnection,
		config.NewValidator,
	)
	return &controller.TagsController{}
}
