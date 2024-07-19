package handler

import (
	"context"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation/converter"
	"github.com/terui-ryota/offer-item/internal/application/usecase"
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	offer_item "github.com/terui-ryota/protofiles/go/offer_item"
)

func NewOfferItemHandler(offerItemUsecase usecase.OfferItemUsecase) offer_item.OfferItemHandlerServer {
	return &offerItemHandler{
		offerItemUsecase: offerItemUsecase,
	}
}

type offerItemHandler struct {
	offerItemUsecase usecase.OfferItemUsecase
	offer_item.UnimplementedOfferItemHandlerServer
}

func (h *offerItemHandler) HealthCheck(ctx context.Context, req *offer_item.HealthCheckReq) (*offer_item.HealthCheckRes, error) {
	fmt.Print("============HealthCheck===============")
	// TODO: HealthCheckの実装を追加
	return &offer_item.HealthCheckRes{
		Request: req,
		Num:     1,
	}, nil
}

// オファー案件を作成する
func (h *offerItemHandler) SaveOfferItem(ctx context.Context, req *offer_item.SaveOfferItemRequest) (*offer_item.SaveOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// DTOに変換する
	offerItemDTO := dto.SaveOfferItemPBToDTO(req.GetOfferItem())

	// 作成する
	if err := h.offerItemUsecase.SaveOfferItem(ctx, offerItemDTO); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.CreateOfferItem: %w", err)
	}

	return &offer_item.SaveOfferItemResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) GetOfferItem(ctx context.Context, req *offer_item.GetOfferItemRequest) (*offer_item.GetOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	offerItem, err := h.offerItemUsecase.GetOfferItem(ctx, offerItemID)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.GetOfferItem: %w", err)
	}
	offerItemPB, err := converter.OfferItemModelToPB(offerItem)
	if err != nil {
		return nil, fmt.Errorf("converter.OfferItemModelToPB: %w", err)
	}

	return &offer_item.GetOfferItemResponse{
		Request:   req,
		OfferItem: offerItemPB,
	}, nil
}

// オファー案件一覧取得
func (h *offerItemHandler) ListOfferItem(ctx context.Context, req *offer_item.ListOfferItemRequest) (*offer_item.ListOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	condition, err := converter.ListConditionPBToModel(req.GetCondition())
	if err != nil {
		return nil, fmt.Errorf("converter.ListConditionPBToModel: %w", err)
	}

	// オファー案件一覧を取得
	result, err := h.offerItemUsecase.ListOfferItem(ctx, condition)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.ListOfferItem: %w", err)
	}

	// protoに変換する
	offerItemPBs := make([]*offer_item.OfferItem, 0, len(result.OfferItems()))
	for _, offerItem := range result.OfferItems() {
		offerItemPB, err := converter.OfferItemModelToPB(offerItem)
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemModelToPB: %w", err)
		}
		offerItemPBs = append(offerItemPBs, offerItemPB)
	}

	return &offer_item.ListOfferItemResponse{
		Request:    req,
		OfferItems: offerItemPBs,
		Result:     converter.ListResultModelToPB(result.ListResult()),
	}, nil
}
