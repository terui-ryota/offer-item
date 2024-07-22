package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/common/txhelper"
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"go.opencensus.io/trace"
)

type ExaminationUsecase interface {
	BulkGetExaminations(ctx context.Context, offerItemID model.OfferItemID, entryType model.EntryType) (map[model.AmebaID]*model.Examination, error)
	UploadExaminationResults(ctx context.Context, offerItemID model.OfferItemID, entryType model.EntryType, examinationResultMap map[string]*dto.ExaminationResultDTO) error
	GetExaminationByAssigneeIDOfferItemID(ctx context.Context, offerItemID model.OfferItemID, assigneeID model.AssigneeID, entryType model.EntryType) (*model.Examination, error)
	Submission(ctx context.Context, offerItemID model.OfferItemID, amebaID model.AmebaID, entryType model.EntryType, entryID *model.EntryID) error
}

func NewExaminationUsecase(
	db *sql.DB,
	examinationRepository repository.ExaminationRepository,
	assigneeRepository repository.AssigneeRepository,
	offerItemRepository repository.OfferItemRepository,
) ExaminationUsecase {
	return &ExaminationUsecaseImpl{
		db:                    db,
		examinationRepository: examinationRepository,
		assigneeRepository:    assigneeRepository,
		offerItemRepository:   offerItemRepository,
	}
}

type ExaminationUsecaseImpl struct {
	db                    *sql.DB
	examinationRepository repository.ExaminationRepository
	assigneeRepository    repository.AssigneeRepository
	offerItemRepository   repository.OfferItemRepository
}

// AmebaIDをkeyにしたmapを取得する
func (e *ExaminationUsecaseImpl) BulkGetExaminations(ctx context.Context, offerItemID model.OfferItemID, entryType model.EntryType) (map[model.AmebaID]*model.Examination, error) {
	ctx, span := trace.StartSpan(ctx, "ExaminationUsecaseImpl.BulkGetExaminations")
	defer span.End()

	result, err := e.examinationRepository.BulkGetByOfferItemID(ctx, e.db, offerItemID, entryType)
	if err != nil {
		return nil, fmt.Errorf("u.examinationRepository.BulkGetByOfferItemID: %w", err)
	}
	return result, nil
}

// 下書き審査、記事審査結果を元にステージを更新する
func (e *ExaminationUsecaseImpl) UploadExaminationResults(ctx context.Context, offerItemID model.OfferItemID, entryType model.EntryType, examinationResultMap map[string]*dto.ExaminationResultDTO) error {
	ctx, span := trace.StartSpan(ctx, "ExaminationUsecaseImpl.UploadExaminationResults")
	defer span.End()

	var err error

	prePassedAssignees := make([]*model.Assignee, 0, len(examinationResultMap))
	preFailedAssignees := make([]*model.Assignee, 0, len(examinationResultMap))
	passedAssignees := make([]*model.Assignee, 0, len(examinationResultMap))
	failedAssignees := make([]*model.Assignee, 0, len(examinationResultMap))

	offerItem, err := e.offerItemRepository.Get(ctx, e.db, offerItemID, false)
	if err != nil {
		return fmt.Errorf("u.offerItemRepository.Get: %w", err)
	}

	if err = txhelper.WithTransaction(ctx, e.db, func(tx *sql.Tx) error {
		switch entryType {
		// 下書き審査の場合(下書き再審査も含む)
		case model.EntryTypeDraft:
			// ステージが「下書き審査」のアサイニーを取得
			assigneeList, err := e.assigneeRepository.ListByOfferItemIDStage(ctx, e.db, offerItemID, model.StagePreExamination)
			if err != nil {
				return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
			}
			for _, assignee := range assigneeList {
				amebaID := assignee.AmebaID()
				examinationResult, ok := examinationResultMap[amebaID.String()]
				if !ok {
					// examinationResultMapにamebaIDが存在しない場合は次のアサイニーに進む
					continue
				}

				// 下書き審査を通過した場合はステージを「記事投稿」に、通過していない場合「下書き再審査」に変更する
				if err := assignee.PreExamination(examinationResult.IsPassed); err != nil {
					return fmt.Errorf("assignee.SetStagePreReexamination: %w", err)
				}

				// 審査結果を設定する
				preExamination, err := e.examinationRepository.Get(ctx, e.db, offerItemID, assignee.ID(), model.EntryTypeDraft)
				if err != nil {
					return fmt.Errorf("u.examinationRepository.GetExamination: %w", err)
				}
				if err = preExamination.SetExaminationResult(examinationResult.IsPassed, examinationResult.ExaminerName, examinationResult.Reason); err != nil {
					return fmt.Errorf("preExamination.SetExaminationResult: %w", err)
				}
				if err = e.examinationRepository.Update(ctx, e.db, preExamination); err != nil {
					return fmt.Errorf("u.examinationRepository.Update: %w", err)
				}

				if examinationResult.IsPassed {
					prePassedAssignees = append(prePassedAssignees, assignee)
				} else {
					preFailedAssignees = append(preFailedAssignees, assignee)
				}

				// アサイニーのステージを更新
				if err = e.assigneeRepository.Update(ctx, e.db, assignee); err != nil {
					return fmt.Errorf("o.assigneeRepository.Update: %w", err)
				}
			}
		case model.EntryTypeEntry:
			// ステージが「記事投稿」のアサイニーを取得
			assigneeList, err := e.assigneeRepository.ListByOfferItemIDStage(ctx, e.db, offerItemID, model.StageExamination)
			if err != nil {
				return fmt.Errorf("o.assigneeRepository.ListByOfferItemID: %w", err)
			}
			for _, assignee := range assigneeList {
				amebaID := assignee.AmebaID()
				examinationResult, ok := examinationResultMap[amebaID.String()]
				if !ok {
					// examinationResultMapにamebaIDが存在しない場合は次のアサイニーに進む
					continue
				}

				// 記事審査を通過した場合はステージを「支払い中」に、通過していない場合「記事再投稿」に変更する
				if err = assignee.Examination(examinationResult.IsPassed); err != nil {
					return fmt.Errorf("assignee.SetStagePreReexamination: %w", err)
				}

				// 審査結果を設定する
				examination, err := e.examinationRepository.Get(ctx, e.db, offerItemID, assignee.ID(), model.EntryTypeEntry)
				if err != nil {
					return fmt.Errorf("u.examinationRepository.GetExamination: %w", err)
				}
				if err = examination.SetExaminationResult(examinationResult.IsPassed, examinationResult.ExaminerName, examinationResult.Reason); err != nil {
					return fmt.Errorf("preExamination.SetExaminationResult: %w", err)
				}
				if err = e.examinationRepository.Update(ctx, e.db, examination); err != nil {
					return fmt.Errorf("u.examinationRepository.Update: %w", err)
				}

				if examinationResult.IsPassed {
					if offerItem.IsPassedAfterReviewMailSent() {
						passedAssignees = append(passedAssignees, assignee)
					}
				} else {
					if offerItem.IsFailedAfterReviewMailSent() {
						failedAssignees = append(failedAssignees, assignee)
					}
				}
				// ステージを更新
				if err = e.assigneeRepository.Update(ctx, e.db, assignee); err != nil {
					return fmt.Errorf("o.assigneeRepository.Update: %w", err)
				}
			}
		default:
			return apperr.OfferItemValidationError.Wrap(fmt.Errorf("entryType:%d is invalid", entryType))
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	//for adsTemplateCode, assignees := range map[string][]*model.Assignee{
	//	"amebapick_offer_item_v2_pre_examination_ok": prePassedAssignees,
	//	"amebapick_offer_item_v2_pre_examination_ng": preFailedAssignees,
	//	"amebapick_offer_item_v2_examination_ok":     passedAssignees,
	//	"amebapick_offer_item_v2_examination_ng":     failedAssignees,
	//} {
	//	if len(assignees) == 0 {
	//		continue
	//	}

	//err = e.queueAdapter.BulkSendMailQueue(ctx, assignees, offerItemID, adsTemplateCode)
	//if err != nil {
	//	return fmt.Errorf("a.queueAdapter.BulkSendMailQueue: %w", err)
	//}
	//}

	return nil
}

func (e *ExaminationUsecaseImpl) GetExaminationByAssigneeIDOfferItemID(ctx context.Context, offerItemID model.OfferItemID, assigneeID model.AssigneeID, entryType model.EntryType) (*model.Examination, error) {
	ctx, span := trace.StartSpan(ctx, "ExaminationUsecaseImpl.GetExaminationByAssigneeIDOfferItemID")
	defer span.End()

	result, err := e.examinationRepository.Get(ctx, e.db, offerItemID, assigneeID, entryType)
	if err != nil {
		return nil, fmt.Errorf("u.examinationRepository.Get: %w", err)
	}
	return result, nil
}

// TODO: SecondリリースでSNSの実装を行う
// 記事投稿、下書き投稿を行う
func (e *ExaminationUsecaseImpl) Submission(ctx context.Context, offerItemID model.OfferItemID, amebaID model.AmebaID, entryType model.EntryType, entryID *model.EntryID) error {
	ctx, span := trace.StartSpan(ctx, "ExaminationUsecaseImpl.Submission")
	defer span.End()

	assignee, err := e.assigneeRepository.GetByAmebaIDOfferItemID(ctx, e.db, amebaID, offerItemID)
	if err != nil {
		return fmt.Errorf("o.assigneeRepository.GetByAmebaIDOfferItemID: %w", err)
	}

	offerItem, err := e.offerItemRepository.Get(ctx, e.db, offerItemID, false)
	if err != nil {
		return fmt.Errorf("u.offerItemRepository.Get: %w", err)
	}

	//var sendMailFlag bool
	//// ステージが記事提出かつ記事投稿メールを送る場合はsendMailFlagをtrueに変更する
	//if assignee.Stage() == model.StageArticlePosting && offerItem.IsArticlePostMailSent() {
	//	sendMailFlag = true
	//}

	examination, err := model.NewExamination(
		offerItemID,
		amebaID,
		entryID,
		assignee.ID(),
		entryType,
	)
	if err != nil {
		return fmt.Errorf("model.NewExamination: %w", err)
	}

	if err := txhelper.WithTransaction(ctx, e.db, func(tx *sql.Tx) error {
		if err := e.examinationRepository.Create(ctx, e.db, examination); err != nil {
			return fmt.Errorf("u.examinationRepository.Create: %w", err)
		}

		if entryType == model.EntryTypeDraft {
			if err := assignee.ChangeStageByDraftSubmission(); err != nil {
				return fmt.Errorf("assignee.ChangeStageByDraftSubmission: %w", err)
			}
		} else {
			if err := assignee.ChangeStageByEntrySubmission(offerItem.NeedsAfterReview()); err != nil {
				return fmt.Errorf("assignee.ChangeStageByEntrySubmission: %w", err)
			}
		}

		if err := e.assigneeRepository.Update(ctx, e.db, assignee); err != nil {
			return fmt.Errorf("o.assigneeRepository.Update: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	//if sendMailFlag {
	//	if err := e.queueAdapter.BulkSendMailQueue(ctx, model.AssigneeList{assignee}, offerItemID, "amebapick_offer_item_v2_after_article_posted"); err != nil {
	//		return fmt.Errorf("a.queueAdapter.BulkSendMailQueue: %w", err)
	//	}
	//}

	return nil
}
