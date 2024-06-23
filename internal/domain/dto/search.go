package dto

import "github.com/terui-ryota/offer-item/internal/domain/model"

type SearchOfferItemCriteria struct {
	NameContains  *string
	ItemIDEqual   *model.ItemID
	DfItemIDEqual *model.DFItemID
}
