package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/rivo/uniseg"

	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/offer-item/pkg/id"
)

// アサイニー
//
//go:generate go run github.com/terui-ryota/gen-getter -type=Assignee
type Assignee struct {
	// アサイニーID
	id AssigneeID
	// 依頼案件ID
	offerItemID OfferItemID
	// アメーバID
	amebaID AmebaID
	// 執筆報酬
	writingFee int
	// ステージ
	stage Stage
	// 辞退理由
	declineReason *string
	// 作成日時
	createdAt time.Time
	// 発送情報
	shippingData []string
	// 発送した商品のJANコード
	janCode *string
}

func NewAssignee(
	offerItemID OfferItemID,
	amebaID AmebaID,
	writingFee int,
	stage Stage,
) (*Assignee, error) {
	if writingFee < 0 {
		return nil, apperr.OfferItemValidationError.Wrap(errors.New("writingFee must be greater than 0"))
	}
	if amebaID == "" {
		return nil, fmt.Errorf("amebaID must be set")
	}
	return &Assignee{
		id:          AssigneeID(id.New()),
		offerItemID: offerItemID,
		amebaID:     amebaID,
		writingFee:  writingFee,
		stage:       stage,
	}, nil
}

func NewAssigneeFromRepository(
	id AssigneeID,
	offerItemID OfferItemID,
	amebaID AmebaID,
	writingFee int,
	stage Stage,
	declineReason *string,
	createdAt time.Time,
) *Assignee {
	return &Assignee{
		id:            id,
		offerItemID:   offerItemID,
		amebaID:       amebaID,
		writingFee:    writingFee,
		stage:         stage,
		declineReason: declineReason,
		createdAt:     createdAt,
	}
}

type LotteryResult struct {
	isPassedLottery bool         // 選考を通過したかどうか
	shippingData    ShippingData // 選考結果インポート時に `発送した商品`というフィールド名のもの(順序保証)。選考以降でサンプルありの場合、nil以外が入る
	janCode         *string
}

type ShippingData []string

func (lr *LotteryResult) IsPassedLottery() bool {
	return lr.isPassedLottery
}

func (lr *LotteryResult) ShippingData() []string {
	return lr.shippingData
}

func (lr *LotteryResult) JanCode() *string {
	return lr.janCode
}

func NewLotteryResult(isPassedLottery bool, shippingData []string, janCode *string) *LotteryResult {
	if shippingData == nil { // リクエストによるので、ここでnilなら空配列を入れておく
		shippingData = []string{}
	}
	return &LotteryResult{
		isPassedLottery: isPassedLottery,
		shippingData:    shippingData,
		janCode:         janCode,
	}
}

func (a *Assignee) SetWritingFee(v int) error {
	if v < 0 {
		return apperr.OfferItemValidationError.Wrap(errors.New("writingFee must be greater than 0"))
	}
	a.writingFee = v
	return nil
}

// ステージを「抽選」から「発送」に変更する
func (a *Assignee) SetStageShipment() error {
	if a.Stage() != StageLottery {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageLottery"))
	}
	a.stage = StageShipment
	return nil
}

// ステージを「抽選」からからオファーアイテムの設定項目を確認し、適切なステージに変更する
func (a *Assignee) ChangeStageByLotteryResult(offerItem *OfferItem, shippingData []string, janCode *string) error {
	if a.Stage() != StageLottery {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageLottery"))
	}
	switch {
	case offerItem.HasSample(): // サンプルがある場合。発送情報を入れる
		a.stage = StageShipment
		a.shippingData = shippingData
		a.janCode = janCode
	case offerItem.NeedsPreliminaryReview(): // 事前審査が必要
		a.stage = StageDraftSubmission
	default:
		a.stage = StageArticlePosting // それ以外は記事提出
	}
	return nil
}

// ステージを「終了」に変更する。また、案件の辞退理由を設定する
func (a *Assignee) SetStageDoneByDecline(declineReason string) error {
	if declineReason == "" {
		return apperr.OfferItemValidationError.Wrap(errors.New("declineReason must be set"))
	}
	if uniseg.GraphemeClusterCount(declineReason) > 128 {
		return apperr.OfferItemValidationError.Wrap(errors.New("declineReason must be less than 128 characters"))
	}
	a.declineReason = &declineReason
	a.SetStageDone()
	return nil
}

// ステージを「抽選」から「抽選落ち」に変更する
func (a *Assignee) SetStageLotteryLost() error {
	if a.Stage() != StageLottery {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageLottery"))
	}
	a.stage = StageLotteryLost
	return nil
}

// ステージを「参加募集前」から「参加募集」に変更する
func (a *Assignee) SetStageInvitation() error {
	if a.Stage() != StageBeforeInvitation {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageBeforeInvitation"))
	}
	a.stage = StageInvitation
	return nil
}

// 審査を通過している場合はステージを「記事提出」に、通過していない場合「下書き再審査」に変更する
func (a *Assignee) PreExamination(isPass bool) error {
	if a.Stage() != StagePreExamination {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StagePreExamination"))
	}
	if isPass {
		a.stage = StageArticlePosting
	} else {
		a.stage = StagePreReexamination
	}
	return nil
}

// 審査を通過している場合はステージを「支払い中」に、通過していない場合「記事再審査」に変更する
func (a *Assignee) Examination(isPass bool) error {
	if a.Stage() != StageExamination {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageExamination"))
	}
	if isPass {
		a.stage = StagePaying
	} else {
		a.stage = StageReexamination
	}
	return nil
}

// ステージを「支払い中」から「支払い完了」に変更する
func (a *Assignee) SetStagePaymentCompleted() error {
	if a.Stage() != StagePaying {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StagePaying"))
	}
	a.stage = StagePaymentCompleted
	return nil
}

// ステージを「発送」から「下書き提出」もしくは「記事提出」に変更する
func (a *Assignee) FinishedShipment(needsPreliminaryReview bool) error {
	if a.Stage() != StageShipment {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageShipment"))
	}

	// 事前審査が必須の場合はステージを「下書き審査」、必須ではない場合を「記事提出」に変更する
	if needsPreliminaryReview {
		a.stage = StageDraftSubmission
	} else {
		a.stage = StageArticlePosting
	}
	return nil
}

// ステージを「終了」に変更する。以前のステージの制限はしない。
func (a *Assignee) SetStageDone() {
	a.stage = StageDone
}

// ChangeStageByDraftSubmission はステージを「下書き提出」or「下書き再審査」から「下書き審査」に変更する
func (a *Assignee) ChangeStageByDraftSubmission() error {
	if a.Stage() != StageDraftSubmission && a.Stage() != StagePreReexamination {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageDraftSubmission or StagePreReexamination"))
	}
	a.stage = StagePreExamination
	return nil
}

// ChangeStageByEntrySubmission はステージが記事提出or記事再審査から事後審査がある場合はステージを「記事審査」、ない場合を「支払い中」に変更する
func (a *Assignee) ChangeStageByEntrySubmission(needsAfterReview bool) error {
	if a.Stage() != StageArticlePosting && a.Stage() != StageReexamination {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageArticlePosting or StageReexamination"))
	}

	// 事後審査がありの場合はステージを「記事審査」、ない場合を「支払い中」に変更する
	if needsAfterReview {
		a.stage = StageExamination
	} else {
		a.stage = StagePaying
	}
	return nil
}

// ステージを「参加募集」から「終了」に変更する
func (a *Assignee) SetStageDoneFromInvitation() error {
	if a.Stage() != StageInvitation {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageInvitation"))
	}
	a.stage = StageDone
	return nil
}

// ステージを「参加募集」からオファーアイテムの設定項目を確認し、適切なステージに変更する
func (a *Assignee) Invitation(offerItem *OfferItem) error {
	if a.Stage() != StageInvitation {
		return apperr.OfferItemValidationError.Wrap(errors.New("stage must be StageInvitation"))
	}
	switch {
	case offerItem.HasLottery():
		a.stage = StageLottery // 抽選がある場合は抽選ステージ
	case offerItem.HasSample(): // サンプルがある
		a.stage = StageShipment
	case offerItem.NeedsPreliminaryReview(): // 事前審査が必要
		a.stage = StageDraftSubmission
	default:
		a.stage = StageArticlePosting // それ以外は記事提出
	}
	return nil
}

// 指定されたステージに強制的にステージを変更する
func (a *Assignee) SetStageAll(s Stage) {
	a.stage = s
}

// AssigneeList アサイニーリスト
type AssigneeList []*Assignee

// アサイニーID
type AssigneeID string

func (ai AssigneeID) String() string {
	return string(ai)
}

// ステージ
type Stage int

func (s Stage) Int() int {
	return int(s)
}

//go:generate go run github.com/terui-ryota/gen-getter -type=AssigneeCount
type AssigneeCount struct {
	// ステージ
	stage Stage
	// カウント
	count int
}

func NewAssigneeCountFromRepository(stage Stage, count int) AssigneeCount {
	return AssigneeCount{
		stage: stage,
		count: count,
	}
}

const (
	StageUnknown          Stage = iota // 不明
	StageBeforeInvitation              // 参加募集前
	StageInvitation                    // 参加募集
	StageLottery                       // 抽選
	StageLotteryLost                   // 抽選落ち
	StageShipment                      // 発送
	StageDraftSubmission               // 下書き提出
	StagePreExamination                // 下書き審査
	StagePreReexamination              // 下書き再審査
	StageArticlePosting                // 記事提出
	StageExamination                   // 記事審査
	StageReexamination                 // 記事再審査
	StagePaying                        // 支払い中
	StagePaymentCompleted              // 支払い完了
	StageDone                          // 終了(辞退の場合も含む)
)
