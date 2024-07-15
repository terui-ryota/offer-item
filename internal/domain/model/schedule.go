//go:generate go run github.com/terui-ryota/gen-getter -type=Schedule

package model

import (
	"time"

	"github.com/friendsofgo/errors"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/offer-item/pkg/id"
)

// スケジュール
type Schedule struct {
	// スケジュールID
	id ScheduleID
	// オファー案件ID
	offerItemID OfferItemID
	// スケジュールタイプ
	scheduleType ScheduleType
	// 開始日
	startDate *time.Time
	// 終了日
	endDate *time.Time
}

// スケジュールリスト
type ScheduleList []*Schedule

// スケジュールタイプを指定し、該当するスケジュールを1つ取得する
func (sl ScheduleList) GetByScheduleType(scheduleType ScheduleType) (*Schedule, bool) {
	for _, schedule := range sl {
		if schedule.scheduleType == scheduleType {
			return schedule, true
		}
	}
	return nil, false
}

// スケジュールID
type ScheduleID string

func (si ScheduleID) String() string {
	return string(si)
}

func NewScheduleFromRepository(
	id ScheduleID,
	offerItemID OfferItemID,
	scheduleType ScheduleType,
	startDate *time.Time,
	endDate *time.Time,
) *Schedule {
	return &Schedule{
		id:           id,
		offerItemID:  offerItemID,
		scheduleType: scheduleType,
		startDate:    startDate,
		endDate:      endDate,
	}
}

func NewSchedule(
	scheduleType ScheduleType,
	startDate *time.Time,
	endDate *time.Time,
) (*Schedule, error) {
	if startDate != nil && endDate != nil {
		// スタートがエンドを超えていた場合はエラー
		if endDate.Before(*startDate) {
			return nil, apperr.OfferItemValidationError.Wrap(errors.New("endDate must be after startDate"))
		}
	}

	switch scheduleType {
	// 支払いの場合、startDateはnil、endDateの指定は必須
	case ScheduleTypePayment:
		if startDate != nil || endDate == nil {
			return nil, apperr.OfferItemValidationError.Wrap(errors.New("startDate must be nil and endDate must not be nil"))
		}
	// 参加募集、記事投稿の場合、startDateとendDateの指定は必須
	case ScheduleTypeInvitation, ScheduleTypeArticlePosting:
		if startDate == nil || endDate == nil {
			return nil, apperr.OfferItemValidationError.Wrap(errors.New("startDate and endDate must be set"))
		}
	// その他のスケジュールはstartDateとendDateがどちらもnil or どちらも指定されている必要がある
	default:
		if (startDate == nil && endDate != nil) || (startDate != nil && endDate == nil) {
			return nil, apperr.OfferItemValidationError.Wrap(errors.New("startDate and endDate must be both nil or both set"))
		}
	}

	return &Schedule{
		id:           ScheduleID(id.New()),
		scheduleType: scheduleType,
		startDate:    startDate,
		endDate:      endDate,
	}, nil
}

func (s *Schedule) SetDate(startDate, endDate *time.Time) error {
	if startDate != nil && endDate != nil {
		// スタートがエンドを超えていた場合はエラー
		if endDate.Before(*startDate) {
			return apperr.OfferItemValidationError.Wrap(errors.New("endDate must be after startDate"))
		}
	}

	switch s.scheduleType {
	// 支払いの場合、startDateはnil、endDateの指定は必須
	case ScheduleTypePayment:
		if startDate != nil || endDate == nil {
			return apperr.OfferItemValidationError.Wrap(errors.New("startDate must be nil and endDate must not be nil"))
		}
	// 参加募集、記事投稿の場合、startDateとendDateの指定は必須
	case ScheduleTypeInvitation, ScheduleTypeArticlePosting:
		if startDate == nil || endDate == nil {
			return apperr.OfferItemValidationError.Wrap(errors.New("startDate and endDate must be set"))
		}
	// その他のスケジュールはstartDateとendDateがどちらもnil or どちらも指定されている必要がある
	default:
		if (startDate == nil && endDate != nil) || (startDate != nil && endDate == nil) {
			return apperr.OfferItemValidationError.Wrap(errors.New("startDate and endDate must be both nil or both set"))
		}
	}

	s.startDate = startDate
	s.endDate = endDate
	return nil
}

// スケジュールタイプ
type ScheduleType int

func ScheduleTypeValues() []ScheduleType {
	return []ScheduleType{
		ScheduleTypeInvitation,
		ScheduleTypeLottery,
		ScheduleTypeShipment,
		ScheduleTypeDraftSubmission,
		ScheduleTypePreExamination,
		ScheduleTypeArticlePosting,
		ScheduleTypeExamination,
		ScheduleTypePayment,
	}
}

func MustScheduleTypeValues() []ScheduleType {
	return []ScheduleType{
		ScheduleTypeInvitation,
		ScheduleTypeArticlePosting,
		ScheduleTypePayment,
	}
}

const (
	ScheduleTypeUnknown         ScheduleType = iota // 不明
	ScheduleTypeInvitation                          // 参加募集
	ScheduleTypeLottery                             // 抽選
	ScheduleTypeShipment                            // 発送
	ScheduleTypeDraftSubmission                     // 下書き提出
	ScheduleTypePreExamination                      // 下書き審査
	ScheduleTypeArticlePosting                      // 記事投稿
	ScheduleTypeExamination                         // 審査
	ScheduleTypePayment                             // 支払い
)

func ScheduleEntityToModel(v int) ScheduleType {
	switch v {
	case int(ScheduleTypeInvitation):
		return ScheduleTypeInvitation
	case int(ScheduleTypeLottery):
		return ScheduleTypeLottery
	case int(ScheduleTypeShipment):
		return ScheduleTypeShipment
	case int(ScheduleTypeDraftSubmission):
		return ScheduleTypeDraftSubmission
	case int(ScheduleTypePreExamination):
		return ScheduleTypePreExamination
	case int(ScheduleTypeArticlePosting):
		return ScheduleTypeArticlePosting
	case int(ScheduleTypeExamination):
		return ScheduleTypeExamination
	case int(ScheduleTypePayment):
		return ScheduleTypePayment
	default:
		return ScheduleTypeUnknown
	}
}

func (o *OfferItem) SetSchedules(schedules []*Schedule) error {
	scheduleMap := make(map[ScheduleType]*Schedule)
	for _, s := range schedules {
		if _, ok := scheduleMap[s.scheduleType]; ok {
			// スケジュールタイプでユニークになる
			return apperr.OfferItemValidationError.Wrap(errors.New("scheduleType is duplicated"))
		}
		scheduleMap[s.scheduleType] = s
	}
	// 「参加募集」「記事投稿」「支払い」は必須の為設定されているか確認する
	for _, mustScheduleType := range MustScheduleTypeValues() {
		if scheduleMap[mustScheduleType] == nil {
			return apperr.OfferItemValidationError.Wrap(errors.New("must mustScheduleType is required"))
		}
	}

	// スケジュールタイプが不正な値の場合はエラー
	if len(scheduleMap) > len(ScheduleTypeValues()) || len(scheduleMap) < len(MustScheduleTypeValues()) {
		return apperr.OfferItemValidationError.Wrap(errors.New("scheduleType is invalid"))
	}
	o.schedules = schedules
	return nil
}

func (st ScheduleType) Int() int {
	return int(st)
}
