package handler

import (
	"context"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/application/usecase"
	pb "github.com/terui-ryota/protofiles/go/offer_item"
)

type OfferItemHandler struct {
	offerItemUsecase usecase.OfferItemUsecase
	pb.UnimplementedOfferItemHandlerServer
}

func NewOfferItemHandler(offerItemUsecase usecase.OfferItemUsecase) *OfferItemHandler {
	return &OfferItemHandler{
		offerItemUsecase: offerItemUsecase,
	}
}

func (h *OfferItemHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	fmt.Print("============HealthCheck===============")
	// TODO: HealthCheckの実装を追加
	return &pb.HealthCheckRes{
		Request: req,
		Num:     1,
	}, nil
}

func (h *OfferItemHandler) SaveOfferItem(ctx context.Context, req *pb.SaveOfferItemRequest) (*pb.SaveOfferItemResponse, error) {
	// TODO: SaveOfferItemの実装を追加
	return &pb.SaveOfferItemResponse{}, nil
}
