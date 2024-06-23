//go:generate go run github.com/terui-ryota/gen-getter -type=OfferItem

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
)

func TestOfferItem_SetCoupon(t *testing.T) {
	type fields struct {
		name                              string
		id                                OfferItemID
		item                              *Item
		dfItem                            *DFItem
		couponBannerID                    *BannerID
		specialRate                       float64
		specialAmount                     int
		hasSample                         bool
		needsPreliminaryReview            bool
		needsAfterReview                  bool
		postRequired                      bool
		hasLottery                        bool
		postTarget                        PostTarget
		hasCoupon                         bool
		hasSpecialCommission              bool
		productFeatures                   string
		cautionaryPoints                  string
		referenceInfo                     string
		otherInfo                         string
		isInvitationMailSent              bool
		isOfferDetailMailSent             bool
		isPassedPreliminaryReviewMailSent bool
		isFailedPreliminaryReviewMailSent bool
		isArticlePostMailSent             bool
		isPassedAfterReviewMailSent       bool
		isFailedAfterReviewMailSent       bool
		isClosed                          bool
		createdAt                         time.Time
		schedules                         ScheduleList
	}
	type args struct {
		hasCoupon      bool
		couponBannerID *string
	}
	bannerIDPtr := func(s BannerID) *BannerID {
		return &s
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantHasCoupon      bool
		wantCouponBannerID *BannerID
		wantErr            bool
	}{
		{
			name:   "正常系。 hasCoupon=true x couponBannerID 設定なし。運用上、 couponBannerID は後から発行される場合があるため",
			fields: fields{},
			args: args{
				hasCoupon:      true,
				couponBannerID: nil,
			},
			wantHasCoupon:      true,
			wantCouponBannerID: nil,
			wantErr:            false,
		},
		{
			name:   "正常系。 hasCoupon=true x couponBannerID 設定あり",
			fields: fields{},
			args: args{
				hasCoupon:      true,
				couponBannerID: null.StringFrom("xxx").Ptr(),
			},
			wantHasCoupon:      true,
			wantCouponBannerID: bannerIDPtr("xxx"),
			wantErr:            false,
		},
		{
			name:   "正常系。 hasCoupon=false x couponBannerID 設定なし",
			fields: fields{},
			args: args{
				hasCoupon:      false,
				couponBannerID: nil,
			},
			wantHasCoupon:      false,
			wantCouponBannerID: nil,
			wantErr:            false,
		},
		{
			name: "正常系。 hasCoupon=false x couponBannerID 設定なし。更新前にcouponBannerIDの登録あり",
			fields: fields{
				couponBannerID: bannerIDPtr("xxx"),
				hasCoupon:      true,
			},
			args: args{
				hasCoupon:      false,
				couponBannerID: nil,
			},
			wantHasCoupon:      false,
			wantCouponBannerID: nil,
			wantErr:            false,
		},
		{
			name:   "異常系。 hasCoupon=false x couponBannerID 設定あり。 couponBannerID を設定する場合、 hasCoupon=true にする必要がある",
			fields: fields{},
			args: args{
				hasCoupon:      false,
				couponBannerID: null.StringFrom("xxx").Ptr(),
			},
			wantHasCoupon:      false,
			wantCouponBannerID: nil,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OfferItem{
				name:                              tt.fields.name,
				id:                                tt.fields.id,
				item:                              tt.fields.item,
				dfItem:                            tt.fields.dfItem,
				couponBannerID:                    tt.fields.couponBannerID,
				specialRate:                       tt.fields.specialRate,
				specialAmount:                     tt.fields.specialAmount,
				hasSample:                         tt.fields.hasSample,
				needsPreliminaryReview:            tt.fields.needsPreliminaryReview,
				needsAfterReview:                  tt.fields.needsAfterReview,
				postRequired:                      tt.fields.postRequired,
				hasLottery:                        tt.fields.hasLottery,
				postTarget:                        tt.fields.postTarget,
				hasCoupon:                         tt.fields.hasCoupon,
				hasSpecialCommission:              tt.fields.hasSpecialCommission,
				productFeatures:                   tt.fields.productFeatures,
				cautionaryPoints:                  tt.fields.cautionaryPoints,
				referenceInfo:                     tt.fields.referenceInfo,
				otherInfo:                         tt.fields.otherInfo,
				isInvitationMailSent:              tt.fields.isInvitationMailSent,
				isOfferDetailMailSent:             tt.fields.isOfferDetailMailSent,
				isPassedPreliminaryReviewMailSent: tt.fields.isPassedPreliminaryReviewMailSent,
				isFailedPreliminaryReviewMailSent: tt.fields.isFailedPreliminaryReviewMailSent,
				isArticlePostMailSent:             tt.fields.isArticlePostMailSent,
				isPassedAfterReviewMailSent:       tt.fields.isPassedAfterReviewMailSent,
				isFailedAfterReviewMailSent:       tt.fields.isFailedAfterReviewMailSent,
				isClosed:                          tt.fields.isClosed,
				createdAt:                         tt.fields.createdAt,
				schedules:                         tt.fields.schedules,
			}
			if err := o.SetCoupon(tt.args.hasCoupon, tt.args.couponBannerID); (err != nil) != tt.wantErr {
				t.Errorf("OfferItem.SetCoupon() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantCouponBannerID, o.couponBannerID)
			assert.Equal(t, tt.wantHasCoupon, o.hasCoupon)
		})
	}
}
