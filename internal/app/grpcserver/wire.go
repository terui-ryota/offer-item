//go:generate go run github.com/google/wire/cmd/wire
//go:build wireinject

package grpcserver

import (
	"github.com/google/wire"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/app"
	grpcConf "github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation"
	"github.com/terui-ryota/offer-item/internal/application"
	"github.com/terui-ryota/offer-item/internal/common"
	"github.com/terui-ryota/offer-item/internal/common/config"
	"github.com/terui-ryota/offer-item/internal/infrastructure"
)

// gRPCサーバー初期化
//func InitializeApp() (common.App, error) {
//	wire.Build(
//		app.NewApp,
//		grpcConf.LoadConfig,
//		wire.FieldsOf(new(*grpcConf.GRPCConfig), "Database", "Rakuten", "Validation"),
//		config.LoadDB,
//		infrastructure.WireSet,
//		application.WireSet,
//		presentation.WireSet,
//		grpcConf.LoadHttpClient,
//	)
//	return nil, nil
//}

func InitializeApp() (common.App, error) {
	wire.Build(
		app.NewApp,
		grpcConf.LoadConfig,
		wire.FieldsOf(new(*grpcConf.GRPCConfig), "Database", "Rakuten", "Validation"),
		config.LoadDB,
		infrastructure.WireSet,
		application.WireSet,
		presentation.WireSet,
		grpcConf.LoadHttpClient,
	)
	return nil, nil
}
