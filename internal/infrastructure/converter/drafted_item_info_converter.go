package converter

import (
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
)

func ConvertDraftedItemModelToEntity(draftedItemInfo *model.ItemInfo) entity.DraftedItemInfo {
	return entity.DraftedItemInfo{
		OfferItemID:       draftedItemInfo.OfferItemID().String(),
		Name:              draftedItemInfo.Name(),
		ContentName:       draftedItemInfo.ContentName(),
		ImageURL:          draftedItemInfo.ImageURL(),
		URL:               draftedItemInfo.URL(),
		MinCommission:     float64(draftedItemInfo.MinCommission().CalculatedRate()),
		MinCommissionType: int(draftedItemInfo.MinCommission().CommissionType()),
		MaxCommission:     float64(draftedItemInfo.MaxCommission().CalculatedRate()),
		MaxCommissionType: int(draftedItemInfo.MaxCommission().CommissionType()),
	}
}

func ConvertDraftedItemToModel(e *entity.DraftedItemInfo) (*model.ItemInfo, error) {
	minCommission, err := model.NewCommission(
		model.CommissionType(e.MinCommissionType),
		float32(e.MinCommission),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create minCommission: %w", err)
	}

	maxCommission, err := model.NewCommission(
		model.CommissionType(e.MaxCommissionType),
		float32(e.MaxCommission),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create maxCommission: %w", err)
	}

	return model.NewDraftedItemInfoFromRepository(
		model.OfferItemID(e.OfferItemID),
		e.Name,
		e.ContentName,
		e.ImageURL,
		e.URL,
		minCommission,
		maxCommission,
	), nil
}
