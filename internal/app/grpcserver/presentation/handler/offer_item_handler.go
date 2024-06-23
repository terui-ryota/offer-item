package handler

import (
	"context"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/application/usecase"
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	pb "github.com/terui-ryota/protofiles/go/offer_item"
)

type offerItemHandler struct {
	offerItemUsecase usecase.OfferItemUsecase
	pb.UnimplementedOfferItemHandlerServer
}

func NewOfferItemHandler(offerItemUsecase usecase.OfferItemUsecase) *offerItemHandler {
	return &offerItemHandler{
		offerItemUsecase: offerItemUsecase,
	}
}

func (h *offerItemHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	fmt.Print("============HealthCheck===============")
	// TODO: HealthCheckの実装を追加
	return &pb.HealthCheckRes{
		Request: req,
		Num:     1,
	}, nil
}

// オファー案件を作成する
func (h *offerItemHandler) SaveOfferItem(ctx context.Context, req *pb.SaveOfferItemRequest) (*pb.SaveOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// DTOに変換する
	offerItemDTO := dto.SaveOfferItemPBToDTO(req.GetOfferItem())

	// 作成する
	if err := h.offerItemUsecase.SaveOfferItem(ctx, offerItemDTO); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.CreateOfferItem: %w", err)
	}

	return &pb.SaveOfferItemResponse{
		Request: req,
	}, nil
}
