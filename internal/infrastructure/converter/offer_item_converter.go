package converter

import (
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	null "github.com/volatiletech/null/v8"
)

func ConvertOfferItemToModel(e *entity.OfferItem, schedules model.ScheduleList, draftedItemInfo *model.ItemInfo) (*model.OfferItem, error) {
	var err error
	var item *model.Item

	if len(e.ItemID) > 0 {
		item, err = model.NewItemByItemID(model.ItemID(e.ItemID))
		if err != nil {
			return nil, fmt.Errorf("model.NewItemByItemID: %w", err)
		}
	}

	var DFItemID model.DFItemID
	if e.DFItemID.Valid {
		DFItemID = model.DFItemID(e.DFItemID.String)
	}

	// TODO: bannerIDは今後複数対応を行う
	var bannerIDs []model.BannerID
	if e.CouponBannerID.Valid {
		bannerIDs = append(bannerIDs, model.BannerID(e.CouponBannerID.String))
	}

	pickInfo, err := model.NewPickInfo(model.ItemID(e.ItemID), &DFItemID, bannerIDs)
	if err != nil {
		return nil, fmt.Errorf("model.NewPickInfoByDFItemID: %w", err)
	}

	offerItem := model.NewOfferItemFromRepository(
		model.OfferItemID(e.ID),
		e.Name,
		item,
		model.NewDFItemByID(model.DFItemID(e.DFItemID.String)),
		(*model.BannerID)(e.CouponBannerID.Ptr()),
		e.SpecialRate,
		e.SpecialAmount,
		e.HasSample,
		e.NeedsPreliminaryReview,
		e.NeedsAfterReview,
		e.NeedsPRMark,
		e.PostRequired,
		e.HasCoupon,
		e.HasSpecialCommission,
		e.HasLottery,
		model.PostTarget(e.PostTarget),
		e.ProductFeatures,
		e.CautionaryPoints,
		e.ReferenceInfo,
		e.OtherInfo,
		e.IsInvitationMailSent,
		e.IsOfferDetailMailSent,
		e.IsPassedPreliminaryReviewMailSent,
		e.IsFailedPreliminaryReviewMailSent,
		e.IsArticlePostMailSent,
		e.IsPassedAfterReviewMailSent,
		e.IsFailedAfterReviewMailSent,
		e.IsClosed,
		e.CreatedAt,
		schedules,
		draftedItemInfo,
		pickInfo,
	)
	return offerItem, nil
}

func ConvertOfferItemModelToEntity(offerItem *model.OfferItem) entity.OfferItem {
	var dfItemID string
	if offerItem.DfItem() != nil {
		dfItemID = offerItem.DfItem().ID().String()
	}

	return entity.OfferItem{
		ID:       offerItem.ID().String(),
		Name:     offerItem.Name(),
		ItemID:   offerItem.Item().ID().String(),
		DFItemID: null.StringFrom(dfItemID),
		CouponBannerID: func() null.String {
			if offerItem.CouponBannerID() == nil {
				return null.StringFromPtr(nil)
			}
			return null.StringFrom(offerItem.CouponBannerID().String())
		}(),
		SpecialRate:                       offerItem.SpecialRate(),
		SpecialAmount:                     offerItem.SpecialAmount(),
		HasSample:                         offerItem.HasSample(),
		NeedsPreliminaryReview:            offerItem.NeedsPreliminaryReview(),
		NeedsAfterReview:                  offerItem.NeedsAfterReview(),
		NeedsPRMark:                       offerItem.NeedsPRMark(),
		PostRequired:                      offerItem.PostRequired(),
		PostTarget:                        uint(offerItem.PostTarget()),
		HasCoupon:                         offerItem.HasCoupon(),
		HasSpecialCommission:              offerItem.HasSpecialCommission(),
		HasLottery:                        offerItem.HasLottery(),
		ProductFeatures:                   offerItem.ProductFeatures(),
		CautionaryPoints:                  offerItem.CautionaryPoints(),
		ReferenceInfo:                     offerItem.ReferenceInfo(),
		OtherInfo:                         offerItem.OtherInfo(),
		IsInvitationMailSent:              offerItem.IsInvitationMailSent(),
		IsOfferDetailMailSent:             offerItem.IsOfferDetailMailSent(),
		IsPassedPreliminaryReviewMailSent: offerItem.IsPassedPreliminaryReviewMailSent(),
		IsFailedPreliminaryReviewMailSent: offerItem.IsFailedPreliminaryReviewMailSent(),
		IsArticlePostMailSent:             offerItem.IsArticlePostMailSent(),
		IsPassedAfterReviewMailSent:       offerItem.IsPassedAfterReviewMailSent(),
		IsFailedAfterReviewMailSent:       offerItem.IsFailedAfterReviewMailSent(),
		IsClosed:                          offerItem.IsClosed(),
	}
}
