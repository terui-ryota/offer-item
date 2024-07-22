package service

import (
	"context"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/adapter"
	"github.com/terui-ryota/offer-item/internal/domain/model"
)

type OfferItemService interface {
	AddItemInfo(ctx context.Context, offerItems model.OfferItemList) error
}

func NewOfferItemServiceImpl(affiliateItemAdapter adapter.AffiliateItemAdapter) OfferItemService {
	return &offerItemServiceImpl{
		affiliateItemAdapter: affiliateItemAdapter,
	}
}

type offerItemServiceImpl struct {
	affiliateItemAdapter adapter.AffiliateItemAdapter
}

func (o *offerItemServiceImpl) AddItemInfo(ctx context.Context, offerItems model.OfferItemList) error {
	itemIdentifiers := offerItems.ItemIdentifiers()

	// 案件ID、DF案件IDが設定されているオファー案件がない場合、何もしない
	if len(itemIdentifiers) == 0 {
		return nil
	}

	// 案件情報を取得
	itemMap, err := o.affiliateItemAdapter.BulkGetItems(ctx, itemIdentifiers)
	if err != nil {
		return fmt.Errorf("o.affiliateItemAdapter.BulkGetItems: %w", err)
	}

	// 取得した情報を適用する
	for _, offerItem := range offerItems {
		var dfItemID model.DFItemID
		if offerItem.DfItem() != nil {
			dfItemID = offerItem.DfItem().ID()
		}

		// 案件情報を適用
		if items, ok := itemMap[*model.NewItemIdentifier(offerItem.Item().ID(), dfItemID)]; ok {
			if err := offerItem.SetItem(&items.Item); err != nil {
				return fmt.Errorf("offerItem.SetItem: %w", err)
			}
			if items.DFItem.Exists() {
				offerItem.SetDFItem(&items.DFItem)
			}
		}
	}

	return nil
}
