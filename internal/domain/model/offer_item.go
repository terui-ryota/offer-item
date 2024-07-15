//go:generate go run github.com/terui-ryota/gen-getter -type=OfferItem

package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/terui-ryota/offer-item/pkg/apperr"
)

// オファー案件
//
//go:generate go run github.com/terui-ryota/gen-getter -type=OfferItem
type OfferItem struct {
	// オファー案件名
	name string
	// 依頼案件ID
	id OfferItemID
	// 案件
	item *Item
	// DF案件
	dfItem *DFItem
	// クーポンPickバナーID
	couponBannerID *BannerID
	// 特別報酬料率
	specialRate float64
	// 特別報酬金額
	specialAmount int
	// サンプルの有無
	hasSample bool
	// 事前審査の有無
	needsPreliminaryReview bool
	// 事後審査の有無
	needsAfterReview bool
	// PRマークやハッシュタグをつけるか(広告主により自動で設置させたくないケースがある)
	needsPRMark bool
	// 投稿必須フラグ
	postRequired bool
	// 抽選の有無
	hasLottery bool
	// 投稿先サービス
	postTarget PostTarget
	// クーポンPickの有無
	hasCoupon bool
	// 特単の有無
	hasSpecialCommission bool
	// 商品特徴
	productFeatures string
	// ブログ投稿時の注意点・懸念点
	cautionaryPoints string
	// ブログ投稿の際の参考情報
	referenceInfo string
	// その他特記事項
	otherInfo string
	// メール設定: 案件参加依頼通知
	isInvitationMailSent bool
	// メール設定: 参加者へ案件詳細
	isOfferDetailMailSent bool
	// メール設定: 事前審査合格通知
	isPassedPreliminaryReviewMailSent bool
	// メール設定: 事前審査不合格通知
	isFailedPreliminaryReviewMailSent bool
	// メール設定: 記事投稿通知
	isArticlePostMailSent bool
	// メール設定: 事後審査合格通知
	isPassedAfterReviewMailSent bool
	// メール設定: 事後審査不合格通知
	isFailedAfterReviewMailSent bool
	// 案件が終了しているかどうか
	isClosed bool
	// 作成日時
	createdAt time.Time
	// スケジュールリスト
	schedules ScheduleList
	// ItemInfoは楽天などの商品情報が消されても管理面への影響を与えないために、バックエンドのDBにキャッシュするために使用します
	draftedItemInfo *ItemInfo
	// PickInfoはDBに保存してあるitem_id, df_item_id, banner_ids(クーポン)の値を格納します
	pickInfo *PickInfo
}

// PickInfoはDBに保存してあるitem_id, df_item_id, banner_ids(クーポン)の値を格納します
//
//go:generate go run github.com/terui-ryota/gen-getter -type=PickInfo
type PickInfo struct {
	// アイテムID
	itemID ItemID
	// DFアイテムID
	dfItemID *DFItemID
	// バナーID
	bannerIDs []BannerID
}

func NewPickInfo(itemID ItemID, optDfItemID *DFItemID, bannerIDs []BannerID) (*PickInfo, error) {
	if len(itemID) == 0 {
		return nil, errors.New("ID should not be empty.")
	}

	var dfItemID *DFItemID
	if optDfItemID != nil && *optDfItemID != "" {
		dfItemID = optDfItemID
	}

	return &PickInfo{
		itemID:    itemID,
		dfItemID:  dfItemID,
		bannerIDs: bannerIDs,
	}, nil
}

func NewOfferItem(
	offerItemID OfferItemID,
	name string,
	item *Item,
	dfItem *DFItem,
	couponBannerID *string,
	specialRate float64,
	specialAmount int,
	hasSample bool,
	needsPreliminaryReview bool,
	needsAfterReview bool,
	needsPRMark bool,
	postRequired bool,
	postTarget PostTarget,
	hasCoupon bool,
	hasSpecialCommission bool,
	hasLottery bool,
	productFeatures string,
	cautionaryPoints string,
	referenceInfo string,
	otherInfo string,
	isInvitationMailSent bool,
	isOfferDetailMailSent bool,
	isPassedPreliminaryReviewMailSent bool,
	isFailedPreliminaryReviewMailSent bool,
	isArticlePostMailSent bool,
	isPassedAfterReviewMailSent bool,
	isFailedAfterReviewMailSent bool,
	isClosed bool,
	scheduleList ScheduleList,
	draftedItemInfo *ItemInfo,
) (*OfferItem, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if productFeatures == "" {
		return nil, errors.New("productFeatures is required")
	}
	if cautionaryPoints == "" {
		return nil, errors.New("cautionaryPoints is required")
	}
	if referenceInfo == "" {
		return nil, errors.New("referenceInfo is required")
	}
	if otherInfo == "" {
		return nil, errors.New("otherInfo is required")
	}
	if item == nil {
		return nil, errors.New("item is required")
	}
	if specialRate < 0 {
		return nil, errors.New("specialRate must be greater than 0")
	}
	if specialAmount < 0 {
		return nil, errors.New("specialAmount must be greater than 0")
	}
	var bannerID *BannerID
	if couponBannerID != nil {
		b, err := NewBannerID(*couponBannerID)
		if err != nil {
			return nil, fmt.Errorf("NewBannerID: %w", err)
		}
		bannerID = b
	}

	if hasSpecialCommission {
		// 特単の有無がtrueの場合、特単料率または特単金額のどちらかが必須
		if specialRate <= 0 && specialAmount <= 0 {
			return nil, errors.New("specialRate or specialAmount is required")
		}
		// 特単料率と特単金額は片方のみ設定可能
		if specialRate > 0 && specialAmount > 0 {
			return nil, errors.New("specialRate and specialAmount cannot be set at the same time")
		}

		// 特単料率が設定されている場合、報酬タイプが定率であることを確認
		if item.MinCommissionRate().CommissionType() == CommissionTypeFixedRate && specialAmount > 0 {
			return nil, errors.New("specialAmount cannot be set when commissionType is CommissionTypeFixedRate")
		}

		// 特単金額が設定されている場合、報酬タイプが定額であることを確認
		if item.MinCommissionRate().CommissionType() == CommissionTypeFixedAmount && specialRate > 0 {
			return nil, errors.New("specialRate cannot be set when commissionType is CommissionTypeFixedAmount")
		}
	} else {
		// 特単の有無がfalseの場合、特単料率と特単金額は設定不可
		if specialRate > 0 || specialAmount > 0 {
			return nil, errors.New("specialRate and specialAmount cannot be set when hasSpecialCommission is false")
		}
	}

	scheduleMap := map[ScheduleType]*Schedule{}
	for _, schedule := range scheduleList {
		schedule.offerItemID = offerItemID
		scheduleMap[schedule.scheduleType] = schedule
	}

	// 「参加募集」「記事投稿」「支払い」は必須の為設定されているか確認する
	for _, mustScheduleType := range MustScheduleTypeValues() {
		if scheduleMap[mustScheduleType] == nil {
			return nil, apperr.OfferItemValidationError.Wrap(errors.New("mustScheduleType is required"))
		}
	}

	// スケジュールタイプが不正な値の場合はエラー
	if len(scheduleMap) > len(ScheduleTypeValues()) || len(scheduleMap) < len(MustScheduleTypeValues()) {
		return nil, apperr.OfferItemValidationError.Wrap(errors.New("scheduleType is invalid"))
	}

	var DFItemID *DFItemID
	if dfItem.Exists() {
		tmpDFItemID := dfItem.ID()
		DFItemID = &tmpDFItemID
	}

	// TODO: bannerIDは今後複数対応を行う
	var bannerIDs []BannerID
	if couponBannerID != nil && *couponBannerID != "" {
		bannerIDs = append(bannerIDs, BannerID(*couponBannerID))
	}

	PickInfo, err := NewPickInfo(item.ID(), DFItemID, bannerIDs)
	if err != nil {
		return nil, fmt.Errorf("model.NewPickInfo: %w", err)
	}

	return &OfferItem{
		id:                                offerItemID,
		name:                              name,
		item:                              item,
		dfItem:                            dfItem,
		couponBannerID:                    bannerID,
		specialAmount:                     specialAmount,
		specialRate:                       specialRate,
		hasSample:                         hasSample,
		needsPreliminaryReview:            needsPreliminaryReview,
		needsAfterReview:                  needsAfterReview,
		needsPRMark:                       needsPRMark,
		postRequired:                      postRequired,
		postTarget:                        postTarget,
		hasCoupon:                         hasCoupon,
		hasSpecialCommission:              hasSpecialCommission,
		hasLottery:                        hasLottery,
		productFeatures:                   productFeatures,
		cautionaryPoints:                  cautionaryPoints,
		referenceInfo:                     referenceInfo,
		otherInfo:                         otherInfo,
		isInvitationMailSent:              isInvitationMailSent,
		isOfferDetailMailSent:             isOfferDetailMailSent,
		isPassedPreliminaryReviewMailSent: isPassedPreliminaryReviewMailSent,
		isFailedPreliminaryReviewMailSent: isFailedPreliminaryReviewMailSent,
		isArticlePostMailSent:             isArticlePostMailSent,
		isPassedAfterReviewMailSent:       isPassedAfterReviewMailSent,
		isFailedAfterReviewMailSent:       isFailedAfterReviewMailSent,
		isClosed:                          isClosed,
		schedules:                         scheduleList,
		draftedItemInfo:                   draftedItemInfo,
		pickInfo:                          PickInfo,
	}, nil
}

func NewOfferItemFromRepository(
	id OfferItemID,
	name string,
	item *Item,
	dfItem *DFItem,
	couponBannerID *BannerID,
	specialRate float64,
	specialAmount int,
	hasSample,
	needsPreliminaryReview,
	needsAfterReview,
	needsPRMark,
	postRequired,
	hasCoupon,
	hasSpecialCommission,
	hasLottery bool,
	postTarget PostTarget,
	productFeatures,
	cautionaryPoints,
	referenceInfo,
	otherInfo string,
	isInvitationMailSent,
	isOfferDetailMailSent,
	isPassedPreliminaryReviewMailSent,
	isFailedPreliminaryReviewMailSent,
	isArticlePostMailSent,
	isPassedAfterReviewMailSent bool,
	isFailedAfterReviewMailSent bool,
	isClosed bool,
	createdAt time.Time,
	schedules ScheduleList,
	draftedItemInfo *ItemInfo,
	pickInfo *PickInfo,
) *OfferItem {
	return &OfferItem{
		id:                                id,
		name:                              name,
		item:                              item,
		dfItem:                            dfItem,
		couponBannerID:                    couponBannerID,
		specialAmount:                     specialAmount,
		specialRate:                       specialRate,
		hasSample:                         hasSample,
		hasLottery:                        hasLottery,
		needsPreliminaryReview:            needsPreliminaryReview,
		needsAfterReview:                  needsAfterReview,
		needsPRMark:                       needsPRMark,
		postRequired:                      postRequired,
		postTarget:                        postTarget,
		hasCoupon:                         hasCoupon,
		hasSpecialCommission:              hasSpecialCommission,
		productFeatures:                   productFeatures,
		cautionaryPoints:                  cautionaryPoints,
		referenceInfo:                     referenceInfo,
		otherInfo:                         otherInfo,
		isInvitationMailSent:              isInvitationMailSent,
		isOfferDetailMailSent:             isOfferDetailMailSent,
		isPassedPreliminaryReviewMailSent: isPassedPreliminaryReviewMailSent,
		isFailedPreliminaryReviewMailSent: isFailedPreliminaryReviewMailSent,
		isArticlePostMailSent:             isArticlePostMailSent,
		isPassedAfterReviewMailSent:       isPassedAfterReviewMailSent,
		isFailedAfterReviewMailSent:       isFailedAfterReviewMailSent,
		isClosed:                          isClosed,
		createdAt:                         createdAt,
		schedules:                         schedules,
		draftedItemInfo:                   draftedItemInfo,
		pickInfo:                          pickInfo,
	}
}

func (o *OfferItem) SetName(v string) error {
	if v == "" {
		return errors.New("name is required")
	}
	o.name = v
	return nil
}

// NOTE: QueryService へ移行検討
func (o *OfferItem) SetItem(item *Item) error {
	if item == nil {
		return errors.New("item is required")
	}
	o.item = item
	return nil
}

// NOTE: QueryService へ移行検討
func (o *OfferItem) SetDFItem(dfItem *DFItem) {
	o.dfItem = dfItem
}

func (o *OfferItem) SetHasSample(v bool) {
	o.hasSample = v
}

func (o *OfferItem) SetSpecialCommission(hasSpecialCommission bool, specialRate float64, specialAmount int, commissionType CommissionType) error {
	if specialRate < 0 {
		return errors.New("specialRate must be greater than 0")
	}
	if specialAmount < 0 {
		return errors.New("specialAmount must be greater than 0")
	}

	if hasSpecialCommission {
		// 特単の有無がtrueの場合、特単料率または特単金額のどちらかが必須
		if specialRate <= 0 && specialAmount <= 0 {
			return errors.New("specialRate or specialAmount is required")
		}
		// 特単料率と特単金額は片方のみ設定可能
		if specialRate > 0 && specialAmount > 0 {
			return errors.New("specialRate and specialAmount cannot be set at the same time")
		}

		// 特単料率が設定されている場合、報酬タイプが定率であることを確認
		if commissionType == CommissionTypeFixedRate && specialAmount > 0 {
			return errors.New("specialAmount cannot be set when commissionType is CommissionTypeFixedRate")
		}

		// 特単金額が設定されている場合、報酬タイプが定額であることを確認
		if commissionType == CommissionTypeFixedAmount && specialRate > 0 {
			return errors.New("specialRate cannot be set when commissionType is CommissionTypeFixedAmount")
		}
	} else {
		// 特単の有無がfalseの場合、特単料率と特単金額は設定不可
		if specialRate > 0 || specialAmount > 0 {
			return errors.New("specialRate and specialAmount cannot be set when hasSpecialCommission is false")
		}
	}

	o.hasSpecialCommission = hasSpecialCommission
	o.specialRate = specialRate
	o.specialAmount = specialAmount
	return nil
}

func (o *OfferItem) SetNeedsPreliminaryReview(v bool) {
	o.needsPreliminaryReview = v
}

func (o *OfferItem) SetNeedsAfterReview(v bool) {
	o.needsAfterReview = v
}

func (o *OfferItem) SetNeedsPRMark(v bool) {
	o.needsPRMark = v
}

func (o *OfferItem) SetPostRequired(v bool) {
	o.postRequired = v
}

func (o *OfferItem) SetPostTarget(v PostTarget) {
	o.postTarget = v
}

func (o *OfferItem) SetCoupon(hasCoupon bool, couponBannerID *string) error {
	if !hasCoupon {
		if couponBannerID != nil {
			return errors.New("coupon banner id must be not empty")
		}
	}
	var b *BannerID
	if couponBannerID != nil {
		var err error
		b, err = NewBannerID(*couponBannerID)
		if err != nil {
			return fmt.Errorf("NewBannerID: %w", err)
		}
	}
	o.hasCoupon = hasCoupon
	o.couponBannerID = b
	return nil
}

func (o *OfferItem) SetHasLottery(v bool) {
	o.hasLottery = v
}

func (o *OfferItem) SetProductFeatures(v string) error {
	if v == "" {
		return errors.New("productFeatures is required")
	}
	o.productFeatures = v
	return nil
}

func (o *OfferItem) SetCautionaryPoints(v string) error {
	if v == "" {
		return errors.New("cautionaryPoints is required")
	}
	o.cautionaryPoints = v
	return nil
}

func (o *OfferItem) SetReferenceInfo(v string) error {
	if v == "" {
		return errors.New("referenceInfo is required")
	}
	o.referenceInfo = v
	return nil
}

func (o *OfferItem) SetOtherInfo(v string) error {
	if v == "" {
		return errors.New("otherInfo is required")
	}
	o.otherInfo = v
	return nil
}

func (o *OfferItem) SetIsInvitationMailSent(v bool) {
	o.isInvitationMailSent = v
}

func (o *OfferItem) SetIsOfferDetailMailSent(v bool) {
	o.isOfferDetailMailSent = v
}

func (o *OfferItem) SetIsArticlePostMailSent(v bool) {
	o.isArticlePostMailSent = v
}

func (o *OfferItem) SetIsPassedPreliminaryReviewMailSent(v bool) {
	o.isPassedPreliminaryReviewMailSent = v
}

func (o *OfferItem) SetIsFailedPreliminaryReviewMailSent(v bool) {
	o.isFailedPreliminaryReviewMailSent = v
}

func (o *OfferItem) SetIsPassedAfterReviewMailSent(v bool) {
	o.isPassedAfterReviewMailSent = v
}

func (o *OfferItem) SetIsFailedAfterReviewMailSent(v bool) {
	o.isFailedAfterReviewMailSent = v
}

func (o *OfferItem) SetIsClosed(v bool) {
	o.isClosed = v
}

func (o *OfferItem) SetDraftedItemInfo(name, contentName, imageURL, url string, minCommission, maxCommission *Commission, offerItemID OfferItemID) error {

	fmt.Println("=======================-")
	fmt.Println("name: ", name)
	fmt.Println("contentName: ", contentName)
	fmt.Println("imageURL: ", imageURL)
	fmt.Println("url: ", url)
	fmt.Println("=======================-")

	if err := validateItemInfo(name, contentName, imageURL, url, minCommission, maxCommission); err != nil {
		return fmt.Errorf("validateItemInfo: %w", err)
	}

	o.draftedItemInfo = &ItemInfo{
		offerItemID:   offerItemID,
		name:          name,
		contentName:   contentName,
		imageURL:      imageURL,
		url:           url,
		minCommission: minCommission,
		maxCommission: maxCommission,
	}
	return nil
}

func (o *OfferItem) SetPickInfo(itemID ItemID, dfItemID *DFItemID, bannerIDs []BannerID) {
	o.pickInfo = &PickInfo{
		itemID:    itemID,
		dfItemID:  dfItemID,
		bannerIDs: bannerIDs,
	}
}

// オファー案件リスト
type OfferItemList []*OfferItem

// ItemIdentifiers はOfferItemListから案件ID、DF案件IDを取得して、ItemIdentifiersを返します
func (oil OfferItemList) ItemIdentifiers() ItemIdentifiers {
	itemIdentifierMap := make(map[ItemIdentifier]struct{})
	for _, offerItem := range oil {
		var dfItemID DFItemID
		if offerItem.dfItem != nil {
			dfItemID = offerItem.dfItem.ID()
		}

		itemIdentifier := NewItemIdentifier(offerItem.item.ID(), dfItemID)
		itemIdentifierMap[*itemIdentifier] = struct{}{}
	}

	identifiers := make(ItemIdentifiers, 0, len(itemIdentifierMap))
	for itemIdentifier := range itemIdentifierMap {
		identifiers = append(identifiers, itemIdentifier)
	}

	return identifiers
}

// オファー案件ID
type OfferItemID string

func (oi OfferItemID) String() string {
	return string(oi)
}

// オファー案件検索結果
type ListOfferItemResult struct {
	// オファー案件リスト
	offerItems OfferItemList
	// リスト取得結果
	listResult *ListResult
}

func NewListOfferItemResult(offerItems OfferItemList, totalCount int) (*ListOfferItemResult, error) {
	listResult, err := NewListResult(len(offerItems), totalCount)
	if err != nil {
		return nil, err
	}

	return &ListOfferItemResult{
		offerItems: offerItems,
		listResult: listResult,
	}, nil
}

func (o *ListOfferItemResult) OfferItems() OfferItemList {
	return o.offerItems
}

func (o *ListOfferItemResult) ListResult() *ListResult {
	return o.listResult
}

// GetIDs オファー案件の ID の一覧を取得する
func (oil OfferItemList) GetIDs() OfferItemIDList {
	list := make(OfferItemIDList, 0, len(oil))
	for _, offerItem := range oil {
		list = append(list, offerItem.ID())
	}
	return list
}

type OfferItemIDList []OfferItemID

func (oil OfferItemIDList) String() []string {
	result := make([]string, 0, len(oil))
	for _, v := range oil {
		result = append(result, v.String())
	}
	return result
}

func (oil OfferItemIDList) InterfaceSlice() []interface{} {
	result := make([]interface{}, 0, len(oil))
	for _, v := range oil {
		result = append(result, v)
	}
	return result
}

//go:generate go run github.com/terui-ryota/gen-getter -type=AssigneeOfferItemPair
type AssigneeOfferItemPair struct {
	// アサイニー
	assignee *Assignee
	// オファー案件
	offerItem *OfferItem
}

func NewAssigneeOfferItemPair(assignee *Assignee, offerItem *OfferItem) *AssigneeOfferItemPair {
	return &AssigneeOfferItemPair{
		assignee:  assignee,
		offerItem: offerItem,
	}
}
