package presentation

import (
	"github.com/google/wire"
	pb "github.com/terui-ryota/protofiles/go/offer_item"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation/handler"
)

var WireSet = wire.NewSet(
	handler.NewOfferItemHandler,
	wire.Bind(new(pb.OfferItemHandlerServer), new(*handler.OfferItemHandler)),
)
