package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	grpcCong "github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/offer-item/internal/application/service"
	"github.com/terui-ryota/offer-item/internal/common/txhelper"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"go.opencensus.io/trace"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
)

type AssigneeUsecase interface {
	ListAssignee(ctx context.Context, offerItemID model.OfferItemID, stage model.Stage) (model.AssigneeList, error)
	ListAssigneeUnderExamination(ctx context.Context) (model.AssigneeList, error)
	ListAssigneeCount(ctx context.Context, offerItemID model.OfferItemID) ([]model.AssigneeCount, error)
	InviteOffer(ctx context.Context, offerItemID model.OfferItemID) error
	UploadLotteryResults(ctx context.Context, offerItemID model.OfferItemID, mapLotteryResult map[model.AmebaID]model.LotteryResult) error
	PaymentCompleted(ctx context.Context, offerItemID model.OfferItemID, amebaIDs []model.AmebaID) error
	CompletedOfferItem(ctx context.Context, offerItemID model.OfferItemID) error
	FinishedShipment(ctx context.Context, offerItemID model.OfferItemID) error
	GetAssigneeByAmebaIDOfferItemID(ctx context.Context, amebaID model.AmebaID, offerItemID model.OfferItemID) (*model.Assignee, error)
	BulkGetQuestionnaireQuestionAnswers(ctx context.Context, offerItemID model.OfferItemID, amebaIDs []model.AmebaID) (map[model.AmebaID]map[model.QuestionID]model.QuestionAnswer, error)
	Invitation(ctx context.Context, offerItemID model.OfferItemID, amebaID model.AmebaID, accepted bool, questionAnswers map[model.QuestionID]string) error
	Decline(ctx context.Context, offerItemID model.OfferItemID, amebaID model.AmebaID, declineReason string) error
}

func NewAssigneeUsecase(
	db *sql.DB,
	config *grpcCong.GRPCConfig,
	assigneeRepository repository.AssigneeRepository,
	offerItemRepository repository.OfferItemRepository,
	questionnaireRepository repository.QuestionnaireRepository,
	questionnaireQuestionAnswerRepository repository.QuestionnaireQuestionAnswerRepository,
) AssigneeUsecase {
	return &assigneeUsecaseImpl{
		db:                                    db,
		config:                                config,
		assigneeRepository:                    assigneeRepository,
		offerItemRepository:                   offerItemRepository,
		questionnaireRepository:               questionnaireRepository,
		questionnaireQuestionAnswerRepository: questionnaireQuestionAnswerRepository,
	}
}

type assigneeUsecaseImpl struct {
	db                                    *sql.DB
	config                                *grpcCong.GRPCConfig
	assigneeRepository                    repository.AssigneeRepository
	offerItemRepository                   repository.OfferItemRepository
	questionnaireRepository               repository.QuestionnaireRepository
	questionnaireQuestionAnswerRepository repository.QuestionnaireQuestionAnswerRepository
	offerItemService                      service.OfferItemService
}

// BulkGetQuestionnaireQuestionAnswers implements AssigneeUsecase.
func (a *assigneeUsecaseImpl) BulkGetQuestionnaireQuestionAnswers(ctx context.Context, offerItemID model.OfferItemID, amebaIDs []model.AmebaID) (map[model.AmebaID]map[model.QuestionID]model.QuestionAnswer, error) {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.BulkGetQuestionnaireQuestionAnswers")
	defer span.End()

	assigneeIDs := make([]model.AssigneeID, 0, len(amebaIDs))
	assigneeIDAmebaIDMap := make(map[model.AssigneeID]model.AmebaID)
	for _, aid := range amebaIDs {
		a, err := a.assigneeRepository.GetByAmebaIDOfferItemID(ctx, a.db, aid, offerItemID)
		if err != nil {
			if !errors.Is(err, apperr.OfferItemNotFoundError) {
				return nil, fmt.Errorf("a.assigneeRepository.GetByAmebaIDOfferItemID: %w", err)
			}
			continue
		}
		assigneeIDs = append(assigneeIDs, a.ID())
		assigneeIDAmebaIDMap[a.ID()] = aid
	}
	answers, err := a.questionnaireQuestionAnswerRepository.BulkGetByOfferItemIDAndAssigneeIDs(ctx, offerItemID, assigneeIDs)
	if err != nil {
		return nil, fmt.Errorf("a.questionnaireQuestionAnswerRepository.BulkGetByOfferItemIDAndAssigneeIDs: %w", err)
	}
	res := make(map[model.AmebaID]map[model.QuestionID]model.QuestionAnswer)
	for k, v := range answers {
		if aid, ok := assigneeIDAmebaIDMap[k]; ok {
			res[aid] = v
		}
	}
	return res, nil
}

// 支払い完了ステージに変更する
func (a *assigneeUsecaseImpl) PaymentCompleted(ctx context.Context, offerItemID model.OfferItemID, amebaIDs []model.AmebaID) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.PaymentCompleted")
	defer span.End()

	assigneeList, err := a.assigneeRepository.ListUnderPaying(ctx, a.db, offerItemID, amebaIDs)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
	}

	if err := txhelper.WithTransaction(ctx, a.db, func(tx *sql.Tx) error {
		for _, assignee := range assigneeList {
			if err := assignee.SetStagePaymentCompleted(); err != nil {
				return fmt.Errorf("assignee.SetStageLottery: %w", err)
			}

			if err := a.assigneeRepository.Update(ctx, tx, assignee); err != nil {
				return fmt.Errorf("o.assigneeRepository.Update: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}
	return nil
}

// オファー案件の全アサイニーのステージを完了に変更する
func (a *assigneeUsecaseImpl) CompletedOfferItem(ctx context.Context, offerItemID model.OfferItemID) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.CompletedOfferItem")
	defer span.End()

	assigneeList, err := a.assigneeRepository.ListByOfferItemID(ctx, a.db, offerItemID)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
	}
	if err := txhelper.WithTransaction(ctx, a.db, func(tx *sql.Tx) error {
		for _, assignee := range assigneeList {
			assignee.SetStageDone()

			if err := a.assigneeRepository.Update(ctx, tx, assignee); err != nil {
				return fmt.Errorf("o.assigneeRepository.Update: %w", err)
			}
		}

		offerItem, err := a.offerItemRepository.Get(ctx, tx, offerItemID, true)
		if err != nil {
			return fmt.Errorf("o.offerItemRepository.Get: %w", err)
		}
		offerItem.SetIsClosed(true)

		if err := a.offerItemRepository.Update(ctx, tx, offerItem); err != nil {
			return fmt.Errorf("o.offerItemRepository.Update: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	return nil
}

// stageを「発送中」から「下書き提出」もしくは「記事提出」に変更する
func (a *assigneeUsecaseImpl) FinishedShipment(ctx context.Context, offerItemID model.OfferItemID) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.FinishedShipment")
	defer span.End()

	assigneeList, err := a.assigneeRepository.ListByOfferItemIDStage(ctx, a.db, offerItemID, model.StageShipment)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
	}

	offerItem, err := a.offerItemRepository.Get(ctx, a.db, offerItemID, false)
	if err != nil {
		return fmt.Errorf("o.offerItemRepository.Get: %w", err)
	}

	if err := txhelper.WithTransaction(ctx, a.db, func(tx *sql.Tx) error {
		for _, assignee := range assigneeList {
			if err := assignee.FinishedShipment(offerItem.NeedsPreliminaryReview()); err != nil {
				return fmt.Errorf("assignee.SetStageShipmentFinished: %w", err)
			}

			if err := a.assigneeRepository.Update(ctx, tx, assignee); err != nil {
				return fmt.Errorf("o.assigneeRepository.Update: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	return nil
}

// amebaIDに紐づくアサイニーを取得する
func (a *assigneeUsecaseImpl) GetAssigneeByAmebaIDOfferItemID(ctx context.Context, amebaID model.AmebaID, offerItemID model.OfferItemID) (*model.Assignee, error) {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.GetAssigneeByAmebaIDOfferItemID")
	defer span.End()

	assignee, err := a.assigneeRepository.GetByAmebaIDOfferItemID(ctx, a.db, amebaID, offerItemID)
	if err != nil {
		return nil, fmt.Errorf("o.assigneeRepository.GetByAmebaID: %w", err)
	}
	return assignee, nil
}

// ブロガーが案件に参加するかどうかを決定する
func (a *assigneeUsecaseImpl) Invitation(ctx context.Context, offerItemID model.OfferItemID, amebaID model.AmebaID, accepted bool, questionAnswers map[model.QuestionID]string) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.Invitaion")
	defer span.End()

	// オファーのアイテムの設定項目により、ステージの遷移先が変わる為取得する
	offerItem, err := a.offerItemRepository.Get(ctx, a.db, offerItemID, false)
	if err != nil {
		return fmt.Errorf("o.offerItemRepository.Get: %w", err)
	}
	questionnaire, err := a.questionnaireRepository.Get(ctx, a.db, offerItemID, false)
	if err != nil {
		if !errors.Is(err, apperr.OfferItemNotFoundError) {
			return fmt.Errorf("a.questionnaireRepository.Get: %w", err)
		}
	}
	if err := txhelper.WithTransaction(ctx, a.db, func(tx *sql.Tx) error {
		assignee, err := a.assigneeRepository.GetByAmebaIDOfferItemID(ctx, tx, amebaID, offerItemID)
		if err != nil {
			return fmt.Errorf("o.assigneeRepository.GetByAmebaIDOfferItemID: %w", err)
		}
		if accepted {
			if err := assignee.Invitation(offerItem); err != nil {
				return fmt.Errorf("assignee.Invitation: %w", err)
			}
			if questionnaire != nil {
				answers, err := model.NewQuestionAnswers(assignee.ID(), *questionnaire, questionAnswers)
				if err != nil {
					return apperr.OfferItemValidationError.Wrap(err)
				}
				if err := a.questionnaireQuestionAnswerRepository.Save(ctx, tx, offerItemID, assignee.ID(), answers); err != nil {
					return fmt.Errorf("a.questionnaireQuestionAnswerRepository.Save: %w", err)
				}
			}

			//// ステージが抽選ではない場合案件に参加が決定する為メールを飛ばす
			//if offerItem.IsOfferDetailMailSent() && assignee.Stage() != model.StageLottery {
			//	// メールを送信する
			//	if err := a.queueAdapter.BulkSendMailQueue(ctx, model.AssigneeList{assignee}, offerItemID, "amebapick_offer_item_v2_lottery_is_passed"); err != nil {
			//		return fmt.Errorf("a.queueAdapter.BulkSendMailQueue: %w", err)
			//	}
			//}
		} else {
			if err := assignee.SetStageDoneFromInvitation(); err != nil {
				return fmt.Errorf("assignee.SetStageDoneFromInvitation: %w", err)
			}
		}
		if err := a.assigneeRepository.Update(ctx, tx, assignee); err != nil {
			return fmt.Errorf("o.assigneeRepository.Update: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	return nil
}

func (a *assigneeUsecaseImpl) Decline(ctx context.Context, offerItemID model.OfferItemID, amebaID model.AmebaID, declineReason string) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.Decline")
	defer span.End()

	// アサイニーを取得
	assignee, err := a.assigneeRepository.GetByAmebaIDOfferItemID(ctx, a.db, amebaID, offerItemID)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.GetByAmebaIDOfferItemID: %w", err)
	}

	// アサイニーのステージを更新し辞退理由を設定
	if err := assignee.SetStageDoneByDecline(declineReason); err != nil {
		return fmt.Errorf("assignee.SetStageDoneByDecline: %w", err)
	}

	if err := a.assigneeRepository.Update(ctx, a.db, assignee); err != nil {
		return fmt.Errorf("o.assigneeRepository.Update: %w", err)
	}
	return nil
}

// 抽選結果を元にステージを更新する
func (a *assigneeUsecaseImpl) UploadLotteryResults(ctx context.Context, offerItemID model.OfferItemID, mapLotteryResult map[model.AmebaID]model.LotteryResult) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.UploadLotteryResults")
	defer span.End()

	assigneeList, err := a.assigneeRepository.ListByOfferItemIDStage(ctx, a.db, offerItemID, model.StageLottery)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
	}

	offerItem, err := a.offerItemRepository.Get(ctx, a.db, offerItemID, false)
	if err != nil {
		return fmt.Errorf("o.offerItemRepository.Get: %w", err)
	}

	// アイテム情報を付与する
	if err = a.offerItemService.AddItemInfo(ctx, model.OfferItemList{offerItem}); err != nil {
		return fmt.Errorf("o.offerItemService.AddItemInfo: %w", err)
	}

	var isPassedAssignees model.AssigneeList
	if err := txhelper.WithTransaction(ctx, a.db, func(tx *sql.Tx) error {
		for _, assignee := range assigneeList {
			amebaID := assignee.AmebaID()

			// 抽選結果に含まれていないアサイニーはスキップする
			if _, ok := mapLotteryResult[amebaID]; !ok {
				continue
			}

			// 抽選を通過した場合はオファーアイテムの設定を見て適切なステージに、落選した場合は抽選落ちステージに変更する
			if lotteryResult := mapLotteryResult[amebaID]; lotteryResult.IsPassedLottery() {
				if err := assignee.ChangeStageByLotteryResult(offerItem, lotteryResult.ShippingData(), lotteryResult.JanCode()); err != nil {
					return fmt.Errorf("assignee.SetStageShipment: %w", err)
				}

				// 抽選を通過したアサイニーを配列に追加
				isPassedAssignees = append(isPassedAssignees, assignee)
			} else {
				if err := assignee.SetStageLotteryLost(); err != nil {
					return fmt.Errorf("assignee.SetStageLotteryLost: %w", err)
				}
			}
			if err := a.assigneeRepository.Update(ctx, a.db, assignee); err != nil {
				return fmt.Errorf("o.assigneeRepository.Update: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	return nil
}

// ステージを参加募集前から参加中に変更する
func (a *assigneeUsecaseImpl) InviteOffer(ctx context.Context, offerItemID model.OfferItemID) error {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.InviteOffer")
	defer span.End()

	assigneeList, err := a.assigneeRepository.ListByOfferItemIDStage(ctx, a.db, offerItemID, model.StageBeforeInvitation)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
	}

	// TODO: アサイニーが増えた場合Timeoutする可能性がある為、BulkUpdateの実装を検討する
	if err := txhelper.WithTransaction(ctx, a.db, func(tx *sql.Tx) error {
		for _, assignee := range assigneeList {
			if err := assignee.SetStageInvitation(); err != nil {
				return fmt.Errorf("assignee.SetStageInvitation: %w", err)
			}
			if err := a.assigneeRepository.Update(ctx, tx, assignee); err != nil {
				return fmt.Errorf("o.assigneeRepository.Update: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	//var adsTemplateCode string
	//offerItem, err := a.offerItemRepository.Get(ctx, a.db, offerItemID, false)
	//if err != nil {
	//	return fmt.Errorf("o.offerItemRepository.Get: %w", err)
	//}
	//
	//// メールを送信しない場合はreturnする
	//if !offerItem.IsInvitationMailSent() {
	//	return nil
	//}
	//
	//if offerItem.HasLottery() {
	//	adsTemplateCode = "amebapick_offer_item_v2_invitation_lottery"
	//} else {
	//	adsTemplateCode = "amebapick_offer_item_v2_invitation_without_lottery"
	//}
	//
	//// メールを送信する
	//if err := a.queueAdapter.BulkSendMailQueue(ctx, assigneeList, offerItemID, adsTemplateCode); err != nil {
	//	return fmt.Errorf("a.queueAdapter.BulkSendMailQueue: %w", err)
	//}

	return nil
}

// ステージに紐づくアサイニーの数を取得する
func (a *assigneeUsecaseImpl) ListAssigneeCount(ctx context.Context, offerItemID model.OfferItemID) ([]model.AssigneeCount, error) {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.ListAssigneeCount")
	defer span.End()

	result, err := a.assigneeRepository.ListCount(ctx, a.db, offerItemID)
	if err != nil {
		return nil, fmt.Errorf("o.assigneeRepository.ListCountByOfferItemID: %w", err)
	}
	return result, nil
}

// 下書き審査、記事審査中のアサイニー一覧を取得する
func (a *assigneeUsecaseImpl) ListAssigneeUnderExamination(ctx context.Context) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.ListAssigneeUnderExamination")
	defer span.End()

	result, err := a.assigneeRepository.ListUnderExamination(ctx, a.db)
	if err != nil {
		return nil, fmt.Errorf("o.assigneeRepository.List: %w", err)
	}
	return result, nil
}

// オファー案件IDとステージに紐づくアサイニー一覧を取得する
func (a *assigneeUsecaseImpl) ListAssignee(ctx context.Context, offerItemID model.OfferItemID, stage model.Stage) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "assigneeUsecaseImpl.ListAssignee")
	defer span.End()

	result, err := a.assigneeRepository.ListByOfferItemIDStage(ctx, a.db, offerItemID, stage)
	if err != nil {
		return nil, fmt.Errorf("o.assigneeRepository.List: %w", err)
	}

	return result, nil
}
