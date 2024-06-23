package dto

import (
	"time"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	offer_item "github.com/terui-ryota/protofiles/go/offer_item"
)

type OfferItemDTO struct {
	// オファー案件名
	Name string
	// 依頼案件ID
	ID *string
	// 案件
	ItemID string
	// DF案件
	DfItemID *string
	// クーポンPickのバナー ID
	CouponBannerID *string
	// 特別報酬料率
	SpecialRate float64
	// 特別報酬単価
	SpecialAmount int
	// サンプルの有無
	HasSample bool
	// 事前審査の有無
	NeedsPreliminaryReview bool
	// 事後審査の有無
	NeedsAfterReview bool
	// PRマークやハッシュタグをつけるか(広告主により自動で設置させたくないケースがある)
	NeedsPRMark bool
	// 投稿必須フラグ
	PostRequired bool
	// 投稿先サービス
	PostTarget PostTarget
	// クーポンPickの有無
	HasCoupon bool
	// 特単の有無
	HasSpecialCommission bool
	// 抽選の有無
	HasLottery bool
	// 商品特徴
	ProductFeatures string
	// ブログ投稿時の注意点・懸念点
	CautionaryPoints string
	// ブログ投稿の際の参考情報
	ReferenceInfo string
	// その他特記事項
	OtherInfo string
	// メール設定: 案件参加依頼通知
	IsInvitationMailSent bool
	// メール設定: 当選者へ案件詳細通知
	IsOfferDetailMailSent bool
	// メール設定: 事前審査合格通知
	IsPassedPreliminaryReviewMailSent bool
	// メール設定: 事前審査不合格通知
	IsFailedPreliminaryReviewMailSent bool
	// メール設定: 記事投稿通知
	IsArticlePostMailSent bool
	// メール設定: 事後審査合格通知
	IsPassedAfterReviewMailSent bool
	// メール設定: 事後審査不合格通知
	IsFailedAfterReviewMailSent bool
	// 案件が終了したか
	IsClosed bool
	// スケジュールリスト
	Schedules ScheduleList
	// アサイニーリスト
	Assignees AssigneeList
	// メール設定リスト
	MailSettings MailSettingList
	// アンケート
	Questionnaire *Questionnaire
	// draftedItemInfoは楽天などの商品情報が消されても管理面への影響を与えないために、バックエンドのDBにキャッシュするために使用します
	DraftedItemInfo *ItemInfo
}

type PostTarget int

const (
	PostTarget_UNDEFINED PostTarget = iota
	PostTarget_AMEBA
	PostTarget_X
	PostTarget_INSTAGRAM
)

type AssigneeList []Assignee

func (a AssigneeList) GetAmebaIDs() []string {
	amebaIDs := make([]string, 0, len(a))
	for _, as := range a {
		amebaIDs = append(amebaIDs, as.AmebaID)
	}
	return amebaIDs
}

type Assignee struct {
	AmebaID       string
	Stage         Stage
	WritingFee    int
	DeclineReason *string
	IsDeleted     bool
}

type ScheduleList []Schedule

type Schedule struct {
	ID           *string
	ScheduleType ScheduleType
	StartDate    *time.Time
	EndDate      *time.Time
}

type MailSettingList []MailSetting

type MailSetting struct {
	ID                 *string
	MailTemplateID     string
	MailType           *MailType
	IsAutoDistribution bool
}

type MailTemplate struct {
	ID           *string
	Name         string
	Mail         MailType
	TemplateCode string
}

type MailType struct {
	Stage      Stage
	IsReminder bool
}

// スケジュールタイプ
type ScheduleType int32

type CommissionType int32

const (
	Commission_Type_Unknown CommissionType = iota
	Commission_Type_Fixed_Rate_
	Commission_Type_Fixed_Amount
	Commission_Type_Multi_Fixed_Amount
	Commission_Type_Click
)

func convertCommissionType(itemCommissionType offer_item.ItemCommissionType) CommissionType {
	switch itemCommissionType {
	case offer_item.ItemCommissionType_COMMISSION_TYPE_FIXED_RATE:
		return Commission_Type_Fixed_Rate_
	case offer_item.ItemCommissionType_COMMISSION_TYPE_FIXED_AMOUNT:
		return Commission_Type_Fixed_Amount
	case offer_item.ItemCommissionType_COMMISSION_TYPE_MULTI_FIXED_AMOUNT:
		return Commission_Type_Multi_Fixed_Amount
	case offer_item.ItemCommissionType_COMMISSION_TYPE_CLICK:
		return Commission_Type_Click
	default:
		return Commission_Type_Unknown
	}
}

// ItemInfo ItemInfoは楽天などの商品情報が消されても管理面への影響を与えないために、バックエンドのDBにキャッシュするために使用します
type ItemInfo struct {
	// 商品名
	Name string
	// 会社名
	ContentName string
	// 商品画像URL
	ImageURL string
	// 詳細を見るの遷移URL
	URL string
	// 最低報酬
	MinCommission *Commission
	// 最高報酬単価
	MaxCommission *Commission
}

// 報酬
type Commission struct {
	// 報酬タイプ
	CommissionType CommissionType
	// 計算済みの報酬料率
	CalculatedRate float64
}

const (
	// 不明
	ScheduleType_SCHEDULE_TYPE_UNKNOWN ScheduleType = iota
	// 参加依頼
	ScheduleType_SCHEDULE_TYPE_INVITATION
	// 抽選
	ScheduleType_SCHEDULE_TYPE_LOTTERY
	// 発送
	ScheduleType_SCHEDULE_TYPE_SHIPMENT
	// 下書き提出
	ScheduleType_SCHEDULE_TYPE_DRAFT_SUBMISSION
	// 事前審査
	ScheduleType_SCHEDULE_TYPE_PRE_EXAMINATION
	// 記事投稿
	ScheduleType_SCHEDULE_TYPE_ARTICLE_POSTING
	// 審査
	ScheduleType_SCHEDULE_TYPE_EXAMINATION
	// 支払い
	ScheduleType_SCHEDULE_TYPE_PAYMENT
)

// ステージ
type Stage int32

const (
	// 不明
	Stage_STAGE_UNKNOWN Stage = iota
	// 参加募集前
	Stage_STAGE_BEFORE_INVITATION
	// 参加募集
	Stage_STAGE_INVITATION
	// 抽選
	Stage_STAGE_LOTTERY
	// 抽選落ち
	Stage_STAGE_LOTTERY_LOST
	// 発送
	Stage_STAGE_SHIPMENT
	// 下書き提出
	Stage_STAGE_DRAFT_SUBMISSION
	// 下書き審査
	Stage_STAGE_PRE_EXAMINATION
	// 下書き再審査
	Stage_STAGE_PRE_REEXAMINATION
	// 記事提出
	Stage_STAGE_ARTICLE_POSTING
	// 記事審査
	Stage_STAGE_EXAMINATION
	// 記事再審査
	Stage_STAGE_REEXAMINATION
	// 支払い中
	Stage_STAGE_PAYING
	// 支払い完了
	Stage_STAGE_PAYMENT_COMPLETED
	// 終了(辞退の場合も含む)
	Stage_STAGE_DONE
)

func SaveOfferItemPBToDTO(offerItem *offer_item.SaveOfferItem) *OfferItemDTO {
	var dfItemID *string
	if offerItem.GetDfItemId() != "" {
		s := offerItem.GetDfItemId()
		dfItemID = &s
	}

	var id *string
	if offerItem.GetId() != "" {
		s := offerItem.GetId()
		id = &s
	}

	var couponBannerID *string
	if offerItem.GetOptionalCouponBannerId() != nil {
		i := offerItem.GetCouponBannerId()
		couponBannerID = &i
	}

	return &OfferItemDTO{
		Name:                              offerItem.GetName(),
		ID:                                id,
		ItemID:                            offerItem.GetItemId(),
		DfItemID:                          dfItemID,
		CouponBannerID:                    couponBannerID,
		SpecialRate:                       offerItem.GetSpecialRate(),
		SpecialAmount:                     int(offerItem.GetSpecialAmount()),
		HasSample:                         offerItem.GetHasSample(),
		NeedsPreliminaryReview:            offerItem.GetNeedsPreliminaryReview(),
		NeedsAfterReview:                  offerItem.GetNeedsAfterReview(),
		NeedsPRMark:                       offerItem.GetNeedsPrMark(),
		PostRequired:                      offerItem.GetPostRequired(),
		PostTarget:                        PostTargetPBToDTO(offerItem.GetPostTarget()),
		HasCoupon:                         offerItem.GetHasCoupon(),
		HasSpecialCommission:              offerItem.GetHasSpecialCommission(),
		HasLottery:                        offerItem.GetHasLottery(),
		ProductFeatures:                   offerItem.GetProductFeatures(),
		CautionaryPoints:                  offerItem.GetCautionaryPoints(),
		ReferenceInfo:                     offerItem.GetReferenceInfo(),
		OtherInfo:                         offerItem.GetOtherInfo(),
		IsInvitationMailSent:              offerItem.GetIsInvitationMailSent(),
		IsOfferDetailMailSent:             offerItem.GetIsOfferDetailMailSent(),
		IsPassedPreliminaryReviewMailSent: offerItem.GetIsPassedPreliminaryReviewMailSent(),
		IsFailedPreliminaryReviewMailSent: offerItem.GetIsFailedPreliminaryReviewMailSent(),
		IsArticlePostMailSent:             offerItem.GetIsArticlePostMailSent(),
		IsPassedAfterReviewMailSent:       offerItem.GetIsPassedAfterReviewMailSent(),
		IsFailedAfterReviewMailSent:       offerItem.GetIsFailedAfterReviewMailSent(),
		IsClosed:                          offerItem.GetIsClosed(),
		Schedules:                         SaveScheduleListPBToDTO(offerItem.GetSchedules()),
		Assignees:                         SaveAssigneeListPBToDTO(offerItem.GetAssignees()),
		Questionnaire: func() *Questionnaire {
			if offerItem.GetOptionalQuestionnaire() == nil {
				return nil
			}
			return &Questionnaire{
				Description: offerItem.GetQuestionnaire().GetDescription(),
				Questions: func() []Question {
					res := make([]Question, 0, len(offerItem.GetQuestionnaire().GetQuestions()))
					for _, q := range offerItem.GetQuestionnaire().GetQuestions() {
						res = append(res, Question{
							ID: func() *string {
								if q.GetId() == "" {
									return nil
								}
								id := q.GetId()
								return &id
							}(),
							QuestionType: model.QuestionType(q.GetQuestionType()),
							Title:        q.GetTitle(),
							ImageURL:     q.GetImageUrl(),
							Options:      q.GetOptions(),
						})
					}
					return res
				}(),
			}
		}(),
		DraftedItemInfo: func() *ItemInfo {
			if offerItem.GetDraftedItemInfo() == nil {
				return nil
			}
			return &ItemInfo{
				Name:        offerItem.GetDraftedItemInfo().GetName(),
				ContentName: offerItem.GetDraftedItemInfo().GetContentName(),
				ImageURL:    offerItem.GetDraftedItemInfo().GetImageUrl(),
				URL:         offerItem.GetDraftedItemInfo().GetUrl(),
				MinCommission: func() *Commission {
					if offerItem.GetDraftedItemInfo().GetMinCommission() == nil {
						return nil
					}
					return &Commission{
						CommissionType: convertCommissionType(offerItem.GetDraftedItemInfo().GetMinCommission().GetCommissionType()),
						CalculatedRate: float64(offerItem.GetDraftedItemInfo().GetMinCommission().GetCalculatedRate()),
					}
				}(),
				MaxCommission: func() *Commission {
					if offerItem.GetDraftedItemInfo().GetMaxCommission() == nil {
						return nil
					}
					return &Commission{
						CommissionType: convertCommissionType(offerItem.GetDraftedItemInfo().GetMaxCommission().GetCommissionType()),
						CalculatedRate: float64(offerItem.GetDraftedItemInfo().GetMaxCommission().GetCalculatedRate()),
					}
				}(),
			}
		}(),
	}
}

func SaveScheduleListPBToDTO(scheduleListPB []*offer_item.SaveSchedule) ScheduleList {
	scheduleList := make(ScheduleList, 0, len(scheduleListPB))

	for _, schedulePB := range scheduleListPB {
		var startDate *time.Time
		if schedulePB.GetStartDate() != nil {
			t := schedulePB.GetStartDate().AsTime()
			startDate = &t
		}

		var endDate *time.Time
		if schedulePB.GetEndDate() != nil {
			t := schedulePB.GetEndDate().AsTime()
			endDate = &t
		}

		scheduleList = append(scheduleList, Schedule{
			ID:           &schedulePB.Id,
			ScheduleType: ScheduleType(schedulePB.GetScheduleType()),
			StartDate:    startDate,
			EndDate:      endDate,
		})
	}
	return scheduleList
}

func SaveAssigneeListPBToDTO(assigneeList []*offer_item.SaveAssignee) AssigneeList {
	assignees := make(AssigneeList, 0, len(assigneeList))

	for _, assignee := range assigneeList {
		var declineReason *string
		if assignee.GetDeclineReason() != "" {
			s := assignee.GetDeclineReason()
			declineReason = &s
		}

		assignees = append(assignees, Assignee{
			AmebaID:       assignee.AmebaId,
			Stage:         StagePBToDTO(assignee.GetStage()),
			WritingFee:    int(assignee.GetWritingFee()),
			DeclineReason: declineReason,
			IsDeleted:     assignee.GetIsDeleted(),
		})
	}
	return assignees
}

func PostTargetPBToDTO(pb offer_item.PostTarget) PostTarget {
	switch pb {
	case offer_item.PostTarget_AMEBA:
		return PostTarget_AMEBA
	case offer_item.PostTarget_X:
		return PostTarget_X
	case offer_item.PostTarget_INSTAGRAM:
		return PostTarget_INSTAGRAM
	default:
		return PostTarget_UNDEFINED
	}
}

func StagePBToDTO(pb offer_item.Stage) Stage {
	switch pb {
	case offer_item.Stage_STAGE_BEFORE_INVITATION:
		return Stage_STAGE_BEFORE_INVITATION
	case offer_item.Stage_STAGE_INVITATION:
		return Stage_STAGE_INVITATION
	case offer_item.Stage_STAGE_LOTTERY:
		return Stage_STAGE_LOTTERY
	case offer_item.Stage_STAGE_LOTTERY_LOST:
		return Stage_STAGE_LOTTERY_LOST
	case offer_item.Stage_STAGE_SHIPMENT:
		return Stage_STAGE_SHIPMENT
	case offer_item.Stage_STAGE_DRAFT_SUBMISSION:
		return Stage_STAGE_DRAFT_SUBMISSION
	case offer_item.Stage_STAGE_PRE_EXAMINATION:
		return Stage_STAGE_PRE_EXAMINATION
	case offer_item.Stage_STAGE_PRE_REEXAMINATION:
		return Stage_STAGE_PRE_REEXAMINATION
	case offer_item.Stage_STAGE_ARTICLE_POSTING:
		return Stage_STAGE_ARTICLE_POSTING
	case offer_item.Stage_STAGE_EXAMINATION:
		return Stage_STAGE_EXAMINATION
	case offer_item.Stage_STAGE_REEXAMINATION:
		return Stage_STAGE_REEXAMINATION
	case offer_item.Stage_STAGE_PAYING:
		return Stage_STAGE_PAYING
	case offer_item.Stage_STAGE_PAYMENT_COMPLETED:
		return Stage_STAGE_PAYMENT_COMPLETED
	case offer_item.Stage_STAGE_DONE:
		return Stage_STAGE_DONE
	default:
		return Stage_STAGE_UNKNOWN
	}
}
