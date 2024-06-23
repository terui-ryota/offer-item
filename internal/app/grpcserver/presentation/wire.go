package presentation

import (
	"github.com/google/wire"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation/handler"
)

var WireSet = wire.NewSet(
	handler.NewOfferItemHandler,
)
