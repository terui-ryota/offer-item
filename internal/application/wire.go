package application

import (
	"github.com/google/wire"

	"github.com/terui-ryota/offer-item/internal/application/usecase"
)

var WireSet = wire.NewSet(
	usecase.NewOfferItemUsecase,
)
