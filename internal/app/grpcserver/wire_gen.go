// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package grpcserver

import (
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/app"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation/handler"
	"github.com/terui-ryota/offer-item/internal/application/usecase"
	"github.com/terui-ryota/offer-item/internal/common"
	config2 "github.com/terui-ryota/offer-item/internal/common/config"
)

// Injectors from wire.go:

// gRPCサーバー初期化
func InitializeApp() (common.App, error) {
	grpcConfig := config.LoadConfig()
	database := grpcConfig.Database
	db := config2.LoadDB(database)
	offerItemUsecase := usecase.NewOfferItemUsecase(db)
	offerItemHandler := handler.NewOfferItemHandler(offerItemUsecase)
	commonApp := app.NewApp(offerItemHandler, grpcConfig)
	return commonApp, nil
}
