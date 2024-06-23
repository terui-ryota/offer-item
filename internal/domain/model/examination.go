package model

import (
	"github.com/friendsofgo/errors"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/offer-item/pkg/id"
)

//go:generate go run github.com/terui-ryota/gen-getter -type=Examination
type Examination struct {
	// ID
	id ExaminationID
	// オファー案件ID
	offerItemID OfferItemID
	// アメーバID
	amebaID AmebaID
	// 記事ID
	entryID *EntryID
	//// SNS
	//sns *SNS
	// 審査者名
	examinerName *string
	// 再審査理由
	reason *string
	// アサイニーID
	assigneeID AssigneeID
	// 記事タイプ
	entryType EntryType
	// 記事提出数
	entrySubmissionCount uint
}

type ExaminationList []*Examination

type ExaminationID string

func (e ExaminationID) String() string {
	return string(e)
}

// SNS 下書き記事投稿の際はuserIDはnilになる。本投稿の際はuserIDとsnsScreenshotURLが必須。
//
//go:generate go run github.com/terui-ryota/gen-getter -type=SNS
type SNS struct {
	// SNSユーザーID。
	userID *string
	// SNSに投稿するポストのスクリーンショットURL
	snsScreenshotURL string
}

type EntryID string

func (e EntryID) String() string {
	return string(e)
}

// TODO: SNS関連の処理はSecondリリースで実装する
func NewExamination(
	offerItemID OfferItemID,
	amebaID AmebaID,
	entryID *EntryID,
	assigneeID AssigneeID,
	entryType EntryType,
) (*Examination, error) {
	// 1stリリースではSNSの実装は行わない。その為entryIDは必ず値がある想定なので、nilの場合はエラーを返す。セカンドリリースではこのバリデーションは不要
	if entryID == nil {
		return nil, apperr.OfferItemValidationError.Wrap(errors.New("entryID is required"))
	}
	return &Examination{
		id:          ExaminationID(id.New()),
		offerItemID: offerItemID,
		amebaID:     amebaID,
		entryID:     entryID,
		assigneeID:  assigneeID,
		entryType:   entryType,
	}, nil
}

// TODO: SNS関連の処理はSecondリリースで実装する
func NewExaminationFromRepository(
	id ExaminationID,
	offerItemID OfferItemID,
	amebaID AmebaID,
	entryID *EntryID,
	examinerName,
	reason *string,
	assigneeID AssigneeID,
	entryType EntryType,
	entrySubmissionCount uint,
) *Examination {
	return &Examination{
		id:                   id,
		offerItemID:          offerItemID,
		amebaID:              amebaID,
		entryID:              entryID,
		examinerName:         examinerName,
		reason:               reason,
		assigneeID:           assigneeID,
		entryType:            entryType,
		entrySubmissionCount: entrySubmissionCount,
	}
}

func NewSNSFromRepository(snsUserID, snsScreenshotURL *string) *SNS {
	// snsの情報を送信する際、snsScreenshotURLがnilの場合はない為、snsScreenshotURLがnilの場合はsnsはnilと判断する
	if snsScreenshotURL == nil {
		return nil
	}

	return &SNS{
		userID:           snsUserID,
		snsScreenshotURL: *snsScreenshotURL,
	}
}

// 記事タイプ
type EntryType int

func (s EntryType) Int() int {
	return int(s)
}

const (
	EntryTypeUnknown EntryType = iota // 不明
	EntryTypeDraft                    // 下書き
	EntryTypeEntry                    // 本投稿
)

func (e *Examination) SetExaminationResult(isPassed bool, examinerName string, reason *string) error {
	// 後方互換のため呼び出し側で設定完了した後に有効にしてください
	// if examinerName == "" {
	// 	return apperr.OfferItemValidationError.Wrap(errors.New("examinerName is required"))
	// }
	if !isPassed {
		// 審査否認される場合は、理由は必須
		if reason == nil || *reason == "" {
			return apperr.OfferItemValidationError.Wrap(errors.New("reason is required"))
		}
	}
	e.examinerName = &examinerName
	e.reason = reason
	return nil
}
