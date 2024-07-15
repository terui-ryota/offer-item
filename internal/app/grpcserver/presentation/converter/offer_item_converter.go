package converter

import (
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/protofiles/go/offer_item"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OfferItemModelToPB(m *model.OfferItem) (*offer_item.OfferItem, error) {
	var (
		item   *offer_item.OfferItem_Item
		dfItem *offer_item.OfferItem_DfItem
	)

	if m.Item().Exists() {
		item = &offer_item.OfferItem_Item{Item: ItemModelToPB(m.Item())}
	}

	if m.DfItem().Exists() {
		dfItem = &offer_item.OfferItem_DfItem{DfItem: dfItemModelToPB(m.DfItem())}
	}

	schedules := make([]*offer_item.Schedule, 0, len(m.Schedules()))
	for _, schedule := range m.Schedules() {
		schedules = append(schedules, ScheduleModelToPB(schedule))
	}

	if m.DraftedItemInfo() == nil {
		return nil, apperr.OfferItemInternalError.Wrap(errors.New(fmt.Sprintf("DraftedItemInfo is nil for OfferItem with ID: %s", m.ID().String())))
	}
	draftedItemInfo := ItemInfoModelToPB(m.DraftedItemInfo())

	if m.PickInfo() == nil {
		return nil, apperr.OfferItemInternalError.Wrap(errors.New(fmt.Sprintf("PickInfo is nil for OfferItem with ID: %s", m.ID().String())))
	}
	pickInfo := PickInfoModelToPB(m.PickInfo())

	offerItemPB := &offer_item.OfferItem{
		Id:             m.ID().String(),
		Name:           m.Name(),
		OptionalItem:   item,
		OptionalDfItem: dfItem,
		OptionalCouponBannerId: func() *offer_item.OfferItem_CouponBannerId {
			id := m.CouponBannerID()
			if id == nil {
				return nil
			}
			return &offer_item.OfferItem_CouponBannerId{
				CouponBannerId: id.String(),
			}
		}(),
		SpecialRate:                       m.SpecialRate(),
		SpecialAmount:                     int64(m.SpecialAmount()),
		HasSample:                         m.HasSample(),
		NeedsPreliminaryReview:            m.NeedsPreliminaryReview(),
		NeedsAfterReview:                  m.NeedsAfterReview(),
		NeedsPrMark:                       m.NeedsPRMark(),
		PostRequired:                      m.PostRequired(),
		PostTarget:                        PostTargetModelToPB(m.PostTarget()),
		HasCoupon:                         m.HasCoupon(),
		HasSpecialCommission:              m.HasSpecialCommission(),
		HasLottery:                        m.HasLottery(),
		ProductFeatures:                   m.ProductFeatures(),
		CautionaryPoints:                  m.CautionaryPoints(),
		ReferenceInfo:                     m.ReferenceInfo(),
		OtherInfo:                         m.OtherInfo(),
		IsInvitationMailSent:              m.IsInvitationMailSent(),
		IsOfferDetailMailSent:             m.IsOfferDetailMailSent(),
		IsPassedPreliminaryReviewMailSent: m.IsPassedPreliminaryReviewMailSent(),
		IsFailedPreliminaryReviewMailSent: m.IsFailedPreliminaryReviewMailSent(),
		IsArticlePostMailSent:             m.IsArticlePostMailSent(),
		IsPassedAfterReviewMailSent:       m.IsPassedAfterReviewMailSent(),
		IsFailedAfterReviewMailSent:       m.IsFailedAfterReviewMailSent(),
		CreatedAt:                         timestamppb.New(m.CreatedAt()),
		Schedules:                         schedules,
		DraftedItemInfo:                   draftedItemInfo,
		PickInfo:                          pickInfo,
	}

	return offerItemPB, nil
}

func ItemModelToPB(m *model.Item) *offer_item.Item {
	return &offer_item.Item{
		Id:                m.ID().String(),
		Img:               m.Img(),
		Name:              m.Name(),
		IsDf:              m.IsDF(),
		MinCommissionRate: CommissionModelToPB(m.MinCommissionRate()),
		MaxCommissionRate: CommissionModelToPB(m.MaxCommissionRate()),
		Urls:              PlatformURLsModelToPB(m.Urls()),
		Tieup:             m.HasTieup(),
		ContentName:       m.ContentName(),
		EnabledSelfBack:   m.EnabledSelfBack(),
	}
}

func dfItemModelToPB(m *model.DFItem) *offer_item.DFItem {
	return &offer_item.DFItem{
		Id:                m.ID().String(),
		Img:               m.Img(),
		Name:              m.Name(),
		MinCommissionRate: CommissionModelToPB(m.MinCommissionRate()),
		MaxCommissionRate: CommissionModelToPB(m.MaxCommissionRate()),
		Urls:              PlatformURLsModelToPB(m.Urls()),
	}
}

func PostTargetModelToPB(m model.PostTarget) offer_item.PostTarget {
	switch m {
	case model.PostTargetAmeba:
		return offer_item.PostTarget_AMEBA
	case model.PostTargetX:
		return offer_item.PostTarget_X
	case model.PostTargetInstagram:
		return offer_item.PostTarget_INSTAGRAM
	default:
		return offer_item.PostTarget_INSTAGRAM
	}
}

func PostTargetDTOToModel(dtoPostTarget dto.PostTarget) model.PostTarget {
	switch dtoPostTarget {
	case dto.PostTarget_AMEBA:
		return model.PostTargetAmeba
	case dto.PostTarget_X:
		return model.PostTargetX
	case dto.PostTarget_INSTAGRAM:
		return model.PostTargetInstagram
	default:
		return model.PostTargetTypeUnknown
	}
}

func QuestionnaireModelToPB(m *model.Questionnaire) *offer_item.Questionnaire {
	return &offer_item.Questionnaire{
		Description: m.Description(),
		Questions: func() []*offer_item.Questionnaire_Question {
			res := make([]*offer_item.Questionnaire_Question, 0, len(m.Questions()))
			for _, q := range m.Questions() {
				res = append(res, &offer_item.Questionnaire_Question{
					Id:           q.ID().String(),
					QuestionType: offer_item.Questionnaire_QuestionType(q.QuestionType()),
					Title:        q.Title(),
					ImageUrl:     q.ImageURL(),
					Options:      q.Options(),
				})
			}
			return res
		}(),
	}
}

func QuestionAnswerModelToPB(m *model.QuestionAnswer) *offer_item.QuestionAnswer {
	return &offer_item.QuestionAnswer{
		QuestionId: m.QuestionID().String(),
		Content:    m.Content(),
	}
}

func ItemInfoModelToPB(m *model.ItemInfo) *offer_item.ItemInfo {
	var minCommission *offer_item.Commission
	if m.MinCommission() != nil {
		minCommission = CommissionModelToPB(m.MinCommission())
	}

	var maxCommission *offer_item.Commission
	if m.MaxCommission() != nil {
		maxCommission = CommissionModelToPB(m.MaxCommission())
	}

	return &offer_item.ItemInfo{
		Name:          m.Name(),
		ContentName:   m.ContentName(),
		ImageUrl:      m.ImageURL(),
		Url:           m.URL(),
		MinCommission: minCommission,
		MaxCommission: maxCommission,
	}
}

func PickInfoModelToPB(m *model.PickInfo) *offer_item.PickInfo {
	var optionalDfItemID *offer_item.PickInfo_DfItemId
	if m.DfItemID() != nil && m.DfItemID().String() != "" {
		optionalDfItemID = &offer_item.PickInfo_DfItemId{
			DfItemId: m.DfItemID().String(),
		}
	}

	bannerIDs := make([]string, 0, len(m.BannerIDs()))
	for _, bannerID := range m.BannerIDs() {
		bannerIDs = append(bannerIDs, bannerID.String())
	}

	return &offer_item.PickInfo{
		ItemId:           m.ItemID().String(),
		OptionalDfItemId: optionalDfItemID,
		BannerIds:        bannerIDs,
	}
}
