package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation/converter"
	"github.com/terui-ryota/offer-item/internal/common/txhelper"
	"github.com/terui-ryota/offer-item/internal/domain/adapter"
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/offer-item/pkg/id"
	"go.opencensus.io/trace"
)

type OfferItemUsecase interface {
	SaveOfferItem(ctx context.Context, offerItemDTO *dto.OfferItemDTO) error
	GetOfferItem(ctx context.Context, offerItemID model.OfferItemID) (*model.OfferItem, error)
	ListOfferItem(ctx context.Context, condition *model.ListCondition) (*model.ListOfferItemResult, error)
}

func NewOfferItemUsecase(
	db *sql.DB,
	offerItemRepository repository.OfferItemRepository,
	assigneeRepository repository.AssigneeRepository,
	questionnaireRepository repository.QuestionnaireRepository,
	questionnaireQuestionAnswerRepository repository.QuestionnaireQuestionAnswerRepository,
	affiliateItemAdapter adapter.AffiliateItemAdapter,
	// affiliatorAdapter adapter.AffiliatorAdapter,
	examinationRepository repository.ExaminationRepository,
	validationConfig *config.ValidationConfig,
) OfferItemUsecase {
	return &offerItemUsecaseImpl{
		db:                                    db,
		offerItemRepository:                   offerItemRepository,
		assigneeRepository:                    assigneeRepository,
		questionnaireRepository:               questionnaireRepository,
		questionnaireQuestionAnswerRepository: questionnaireQuestionAnswerRepository,
		affiliateItemAdapter:                  affiliateItemAdapter,
		//affiliatorAdapter:                     affiliatorAdapter,
		examinationRepository: examinationRepository,
		validationConfig:      validationConfig,
	}
}

type offerItemUsecaseImpl struct {
	db                                    *sql.DB
	offerItemRepository                   repository.OfferItemRepository
	assigneeRepository                    repository.AssigneeRepository
	questionnaireRepository               repository.QuestionnaireRepository
	questionnaireQuestionAnswerRepository repository.QuestionnaireQuestionAnswerRepository
	affiliateItemAdapter                  adapter.AffiliateItemAdapter
	//affiliatorAdapter                     adapter.AffiliatorAdapter
	examinationRepository repository.ExaminationRepository
	validationConfig      *config.ValidationConfig
}

// offer-item、schedule、assignee、questionnaireの作成・更新を行う
func (o *offerItemUsecaseImpl) SaveOfferItem(ctx context.Context, offerItemDTO *dto.OfferItemDTO) error {
	ctx, span := trace.StartSpan(ctx, "offerItemUsecaseImpl.SaveOfferItem")
	defer span.End()

	// 案件の取得を行う
	var (
		itemID   model.ItemID
		dfItemID model.DFItemID
	)

	// ItemIDが存在しない場合はエラーを返す
	if offerItemDTO.ItemID == "" {
		return apperr.OfferItemValidationError.Wrap(errors.New("ItemID is required"))
	}
	itemID = model.ItemID(offerItemDTO.ItemID)

	if offerItemDTO.DfItemID != nil {
		dfItemID = model.DFItemID(*offerItemDTO.DfItemID)
	}

	if len(offerItemDTO.Assignees) > o.validationConfig.MaxInputAssigneeListNum {
		return apperr.OfferItemValidationError.Wrap(errors.New("assignees input over MAX_ASSIGNEE_INPUT_NUM"))
	}

	items, err := o.affiliateItemAdapter.GetItems(ctx, *model.NewItemIdentifier(itemID, dfItemID))
	if err != nil {
		return fmt.Errorf("o.affiliateItemAdapter.GetItems: %w. Item ID: %s, DF Item ID: %s", err, itemID.String(), dfItemID.String())
	}

	//amebaIDs := offerItemDTO.Assignees.GetAmebaIDs()
	//affiliatorIDMap, err := o.affiliatorAdapter.BulkGetAffiliatorIDsByAmebaIDs(ctx, amebaIDs)
	//if err != nil {
	//	return fmt.Errorf("o.affiliatorAdapter.BulkGetAffiliatorIDsByAmebaIDs: %w", err)
	//}

	if err := txhelper.WithTransaction(ctx, o.db, func(tx *sql.Tx) error {
		var offerItemID model.OfferItemID
		AssigneesDTOs := offerItemDTO.Assignees
		// オファーアイテムIDの存在が存在する場合更新処理を行う
		if offerItemDTO.ID != nil {
			offerItemID = model.OfferItemID(*offerItemDTO.ID)
			// スケジュールIDの存在を確認する
			for _, schedule := range offerItemDTO.Schedules {
				if schedule.ID == nil || *schedule.ID == "" {
					return apperr.OfferItemValidationError.Wrap(errors.New("ScheduleID is required"))
				}
			}

			// 行ロックする際、idでレコードを取得する
			offerItem, err := o.offerItemRepository.Get(ctx, tx, offerItemID, true)
			if err != nil {
				return fmt.Errorf("o.offerItemRepository.Get: %w", err)
			}

			if err := offerItem.SetItem(&items.Item); err != nil {
				return fmt.Errorf("offerItem.SetItem: %w", err)
			}

			// MEMO: DFITemIDが設定されていたオファー案件を編集しDFItemを削除した場合、フロントからDFItemIDはnilでリクエストされる。
			// その場合、DFItemIDを明示的に空にする。（空にしない場合、offerItem.DfItem().ID()が更新されずに値が入った状態になってしまう）
			if offerItemDTO.DfItemID == nil {
				if offerItem.DfItem().Exists() {
					offerItem.DfItem().SetDFItemIDEmpty()
				}
			} else {
				offerItem.SetDFItem(&items.DFItem)
			}

			if err := setOfferItemFields(offerItem, offerItemDTO); err != nil {
				return fmt.Errorf("setOfferItemFields: %w", err)
			}

			// OfferItem、Scheduleを更新
			if err := o.offerItemRepository.Update(ctx, tx, offerItem); err != nil {
				return fmt.Errorf("o.offerItemRepository.Update: %w", err)
			}

			if offerItemDTO.Questionnaire != nil {
				// アンケートが設定されていれば新規作成または更新
				q, err := o.questionnaireRepository.Get(ctx, tx, offerItem.ID(), true)
				if err != nil {
					if !errors.Is(err, apperr.OfferItemNotFoundError) {
						return fmt.Errorf("o.questionnaireRepository.Get: %w", err)
					}
					q, err = createQuestionnaire(offerItemID, *offerItemDTO.Questionnaire)
					if err != nil {
						return fmt.Errorf("createQuestionnaire: %w", err)
					}
				} else {
					q, err = updateQuestionnaire(*q, *offerItemDTO.Questionnaire)
					if err != nil {
						return fmt.Errorf("o.updateQuestionnaire: %w", err)
					}
				}
				if err := o.questionnaireRepository.Save(ctx, tx, *q); err != nil {
					return fmt.Errorf("o.questionnaireRepository.Save: %w", err)
				}
			} else {
				if err := o.questionnaireQuestionAnswerRepository.DeleteByOfferItemID(ctx, tx, offerItem.ID()); err != nil {
					return fmt.Errorf("o.questionnaireQuestionAnswerRepository.DeleteByOfferItemID: %w", err)
				}
				if err := o.questionnaireRepository.Delete(ctx, tx, offerItem.ID()); err != nil {
					return fmt.Errorf("o.questionnaireRepository.Delete: %w", err)
				}
			}

			// AmebaIDの存在を確認する
			var amebaIDs []model.AmebaID
			for _, AssigneeDTO := range AssigneesDTOs {
				if AssigneeDTO.AmebaID == "" {
					return apperr.OfferItemValidationError.Wrap(errors.New("AmebaID is required"))
				}
				amebaIDs = append(amebaIDs, model.AmebaID(AssigneeDTO.AmebaID))
			}
			assigneeMap, err := o.assigneeRepository.BulkGetByOfferItemIDAmebaIDs(ctx, o.db, offerItemID, amebaIDs, true)
			if err != nil {
				return fmt.Errorf("o.assigneeRepository.BulkGetByOfferItemIdAmebaIDs: %w", err)
			}
			for i := range AssigneesDTOs {
				assigneeDTO := AssigneesDTOs[i]
				if assignee, ok := assigneeMap[model.AmebaID(assigneeDTO.AmebaID)]; ok {
					if assigneeDTO.IsDeleted {
						if assignee.Stage() != model.StageBeforeInvitation {
							return apperr.OfferItemValidationError.Wrap(errors.New("Assignee stage must be StageBeforeInvitation"))
						}
						if err := o.assigneeRepository.DeleteByOfferItemIDAndAmebaID(ctx, tx, offerItemID, model.AmebaID(assigneeDTO.AmebaID)); err != nil {
							return fmt.Errorf("o.assigneeRepository.DeleteByOfferItemIDAmebaID: %w", err)
						}
						continue
					}

					// リクエストされたステージが「下書き審査」または「審査」の場合、examinationにデータが存在することを確認する
					var entryType model.EntryType
					if assigneeDTO.Stage == dto.Stage_STAGE_PRE_EXAMINATION || assigneeDTO.Stage == dto.Stage_STAGE_EXAMINATION {
						if assigneeDTO.Stage == dto.Stage_STAGE_PRE_EXAMINATION {
							entryType = model.EntryTypeDraft
						}
						if assigneeDTO.Stage == dto.Stage_STAGE_EXAMINATION {
							entryType = model.EntryTypeEntry
						}

						// examinationが存在しない場合はエラーを返す
						_, err := o.examinationRepository.Get(ctx, o.db, offerItemID, assignee.ID(), entryType)
						if errors.Is(err, apperr.OfferItemNotFoundError) {
							return apperr.OfferItemNotFoundError.Wrap(errors.New("if stage is pre-examination or examination, examination must exist"))
						}
						if err != nil {
							return fmt.Errorf("o.examinationRepository.Get: %w", err)
						}
					}

					if err := setAssigneeFields(assignee, &assigneeDTO); err != nil {
						return fmt.Errorf("setAssigneeFields: %w", err)
					}

					// Assigneeを更新
					if err := o.assigneeRepository.Update(ctx, tx, assignee); err != nil {
						return fmt.Errorf("o.assigneeRepository.Update: %w", err)
					}
				} else {
					// Assigneeが追加された場合
					//itemID := offerItemDTO.ItemID
					//affiliatorID, ok := affiliatorIDMap[model.AmebaID(assigneeDTO.AmebaID)]
					//
					//// affiliatorIDが存在しない場合はエラーを返す
					//if !ok {
					//	logger.FromContext(ctx).Warn("missing affiliatorID", zap.String("ameba_id", assigneeDTO.AmebaID))
					//	return apperr.OfferItemValidationError.Wrap(errors.New("missing affiliatorID"))
					//}
					assignee, err := model.NewAssignee(
						offerItem.ID(),
						model.AmebaID(assigneeDTO.AmebaID),
						assigneeDTO.WritingFee,
						converter.StageDTOToModel(assigneeDTO.Stage),
					)
					if err != nil {
						return fmt.Errorf("model.NewAssignee: %w", err)
					}
					if err := o.assigneeRepository.Create(ctx, tx, assignee); err != nil {
						return fmt.Errorf("o.assigneeRepository.Create: %w", err)
					}
					//// 提携処理
					//if !items.Item.HasTieup() {
					//	continue
					//}
					//if err = o.affiliateItemAdapter.ApplyTieupItem(ctx, model.ItemID(itemID), affiliatorID); err != nil {
					//	return fmt.Errorf("o.affiliateItemAdapter.ApplyTieupItem: %w", err)
					//}
					//if err := o.affiliateItemAdapter.ApproveTieupItem(ctx, model.ItemID(itemID), affiliatorID); err != nil {
					//	return fmt.Errorf("o.affilaiteItemAdapter.ApproveTieupItem: %w", err)
					//}
				}
			}
		} else {
			var schedules model.ScheduleList
			for _, scheduleDTO := range offerItemDTO.Schedules {
				schedule, err := model.NewSchedule(
					converter.ScheduleTypeDTOToModel(scheduleDTO.ScheduleType),
					scheduleDTO.StartDate,
					scheduleDTO.EndDate,
				)
				if err != nil {
					return fmt.Errorf("model.NewSchedule: %w", err)
				}
				schedules = append(schedules, schedule)
			}

			offerItemID = model.OfferItemID(id.New())

			draftedItemInfoMinCommission, err := model.NewCommission(
				model.CommissionType(offerItemDTO.DraftedItemInfo.MinCommission.CommissionType),
				float32(offerItemDTO.DraftedItemInfo.MinCommission.CalculatedRate),
			)
			if err != nil {
				return fmt.Errorf("failed to create minCommission: %w", err)
			}

			draftedItemInfoMaxCommission, err := model.NewCommission(
				model.CommissionType(offerItemDTO.DraftedItemInfo.MaxCommission.CommissionType),
				float32(offerItemDTO.DraftedItemInfo.MaxCommission.CalculatedRate),
			)
			if err != nil {
				return fmt.Errorf("failed to create maxCommission: %w", err)
			}

			draftedItemInfo, err := model.NewItemInfo(
				offerItemID,
				offerItemDTO.DraftedItemInfo.Name,
				offerItemDTO.DraftedItemInfo.ContentName,
				offerItemDTO.DraftedItemInfo.ImageURL,
				offerItemDTO.DraftedItemInfo.URL,
				draftedItemInfoMinCommission,
				draftedItemInfoMaxCommission,
			)
			if err != nil {
				return fmt.Errorf("convertOfferItemDTOToItemInfo: %w", err)
			}

			offerItem, err := model.NewOfferItem(
				offerItemID,
				offerItemDTO.Name,
				&items.Item,
				&items.DFItem,
				offerItemDTO.CouponBannerID,
				offerItemDTO.SpecialRate,
				offerItemDTO.SpecialAmount,
				offerItemDTO.HasSample,
				offerItemDTO.NeedsPreliminaryReview,
				offerItemDTO.NeedsAfterReview,
				offerItemDTO.NeedsPRMark,
				offerItemDTO.PostRequired,
				converter.PostTargetDTOToModel(offerItemDTO.PostTarget),
				offerItemDTO.HasCoupon,
				offerItemDTO.HasSpecialCommission,
				offerItemDTO.HasLottery,
				offerItemDTO.ProductFeatures,
				offerItemDTO.CautionaryPoints,
				offerItemDTO.ReferenceInfo,
				offerItemDTO.OtherInfo,
				offerItemDTO.IsInvitationMailSent,
				offerItemDTO.IsOfferDetailMailSent,
				offerItemDTO.IsPassedPreliminaryReviewMailSent,
				offerItemDTO.IsFailedPreliminaryReviewMailSent,
				offerItemDTO.IsArticlePostMailSent,
				offerItemDTO.IsPassedAfterReviewMailSent,
				offerItemDTO.IsFailedAfterReviewMailSent,
				offerItemDTO.IsClosed,
				schedules,
				draftedItemInfo,
			)
			if err != nil {
				return fmt.Errorf("model.NewOfferItem: %w", err)
			}
			if err := o.offerItemRepository.Create(ctx, tx, offerItem); err != nil {
				return fmt.Errorf("o.offerItemRepository.Create: %w", err)
			}
			if offerItemDTO.Questionnaire != nil {
				// アンケートが設定されていれば新規作成
				q, err := createQuestionnaire(offerItem.ID(), *offerItemDTO.Questionnaire)
				if err != nil {
					return fmt.Errorf("createQuestionnaire: %w", err)
				}
				if err := o.questionnaireRepository.Save(ctx, tx, *q); err != nil {
					return fmt.Errorf("o.questionnaireRepository.Save: %w", err)
				}
			}
			// アサイニーインサート
			for _, assigneeDTO := range AssigneesDTOs {
				//itemID := offerItemDTO.ItemID
				//affiliatorID, ok := affiliatorIDMap[model.AmebaID(assigneeDTO.AmebaID)]
				//// affiliatorIDが存在しない場合はエラーを返す
				//if !ok {
				//	logger.FromContext(ctx).Warn("missing affiliatorID", zap.String("ameba_id", assigneeDTO.AmebaID))
				//	return apperr.OfferItemValidationError.Wrap(errors.New("missing affiliatorID"))
				//}

				// ステージが「事前審査」または「審査」の場合、審査結果をアップロードした場合nil pointerでエラーになる可能性がある為バリデーションを行う
				if assigneeDTO.Stage == dto.Stage_STAGE_PRE_EXAMINATION || assigneeDTO.Stage == dto.Stage_STAGE_EXAMINATION {
					return apperr.OfferItemValidationError.Wrap(errors.New("stage must not be pre-examination or examination"))
				}
				assignee, err := model.NewAssignee(
					offerItem.ID(),
					model.AmebaID(assigneeDTO.AmebaID),
					assigneeDTO.WritingFee,
					converter.StageDTOToModel(assigneeDTO.Stage),
				)
				if err != nil {
					return fmt.Errorf("model.NewAssignee: %w", err)
				}
				if err := o.assigneeRepository.Create(ctx, tx, assignee); err != nil {
					return fmt.Errorf("o.assigneeRepository.Create: %w", err)
				}
				// 提携処理
				if !items.Item.HasTieup() {
					continue
				}
				//if err = o.affiliateItemAdapter.ApplyTieupItem(ctx, model.ItemID(itemID), affiliatorID); err != nil {
				//	return fmt.Errorf("o.affiliateItemAdapter.ApplyTieupItem: %w", err)
				//}
				//if err := o.affiliateItemAdapter.ApproveTieupItem(ctx, model.ItemID(itemID), affiliatorID); err != nil {
				//	return fmt.Errorf("o.affilaiteItemAdapter.ApproveTieupItem: %w", err)
				//}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("txhelper.WithTransaction: %w", err)
	}

	return nil
}

func createQuestionnaire(offerItemID model.OfferItemID, input dto.Questionnaire) (*model.Questionnaire, error) {
	qs := make([]model.Question, 0, len(input.Questions))
	for _, q := range input.Questions {
		q, err := model.NewQuestion(
			offerItemID,
			q.QuestionType,
			q.Title,
			q.ImageURL,
			q.Options,
		)
		if err != nil {
			return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("model.NewQuestion: %w", err))
		}
		qs = append(qs, *q)
	}
	questionnaire, err := model.NewQuestionnaire(
		offerItemID,
		input.Description,
		qs,
	)
	if err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("model.NewQuestionnaire: %w", err))
	}
	return questionnaire, nil
}

func updateQuestionnaire(q model.Questionnaire, input dto.Questionnaire) (*model.Questionnaire, error) {
	if err := q.SetDescription(input.Description); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("%w", err))
	}
	questionMap := func() map[string]model.Question {
		m := make(map[string]model.Question)
		for _, q := range q.Questions() {
			m[q.ID().String()] = q
		}
		return m
	}()
	qs := make([]model.Question, 0, len(input.Questions))
	idmap := map[string]struct{}{}
	for _, qi := range input.Questions {
		// NOTE: フロントがQuestionをコピーしたときにidを空白にするのが難しい
		//       そのためQuestionのIDを先勝ちにし、あとからきたものをid=''にして新規作成する
		//       そのためフロント(backend/admin/templates/offer_item/edit.html)でのJSONEditorの順序変更を不可能(disable_array_reorder: true)にしている
		if qi.ID != nil {
			_, ok := idmap[*qi.ID]
			if ok {
				qi.ID = nil
			} else {
				idmap[*qi.ID] = struct{}{}
			}
		}

		if qi.ID == nil {
			q, err := model.NewQuestion(
				q.OfferItemID(),
				qi.QuestionType,
				qi.Title,
				qi.ImageURL,
				qi.Options,
			)
			if err != nil {
				return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("model.NewQuestion: %w", err))
			}
			qs = append(qs, *q)
		} else {
			q, ok := questionMap[*qi.ID]
			if !ok {
				return nil, fmt.Errorf("not found: %s", *qi.ID)
			}
			if err := q.SetOptions(qi.QuestionType, qi.Options...); err != nil {
				return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("q.SetOptions: %w", err))
			}
			if err := q.SetTitle(qi.Title); err != nil {
				return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("q.SetTitle: %w", err))
			}
			if err := q.SetImageURL(qi.ImageURL); err != nil {
				return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("q.SetImageURL: %w", err))
			}
			qs = append(qs, q)
		}
	}
	if err := q.SetQuestions(qs); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf("q.SetQuestions: %w", err))
	}
	return &q, nil
}

func setOfferItemFields(offerItem *model.OfferItem, d *dto.OfferItemDTO) error {
	if err := offerItem.SetName(d.Name); err != nil {
		return fmt.Errorf("offerItem.SetName: %w", err)
	}
	if err := offerItem.SetSpecialCommission(d.HasSpecialCommission, d.SpecialRate, d.SpecialAmount, offerItem.Item().MinCommissionRate().CommissionType()); err != nil {
		return fmt.Errorf("offerItem.SetSpecialCommission: %w", err)
	}
	offerItem.SetHasSample(d.HasSample)
	offerItem.SetNeedsPreliminaryReview(d.NeedsPreliminaryReview)
	offerItem.SetNeedsAfterReview(d.NeedsAfterReview)
	offerItem.SetNeedsPRMark(d.NeedsPRMark)
	offerItem.SetPostRequired(d.PostRequired)
	offerItem.SetPostTarget(converter.PostTargetDTOToModel(d.PostTarget))
	if err := offerItem.SetCoupon(d.HasCoupon, d.CouponBannerID); err != nil {
		return fmt.Errorf("offerItem.SetCoupon: %w", err)
	}
	offerItem.SetHasLottery(d.HasLottery)
	if err := offerItem.SetProductFeatures(d.ProductFeatures); err != nil {
		return fmt.Errorf("offerItem.SetProductFeatures: %w", err)
	}
	if err := offerItem.SetCautionaryPoints(d.CautionaryPoints); err != nil {
		return fmt.Errorf("offerItem.SetCautionaryPoints: %w", err)
	}
	if err := offerItem.SetReferenceInfo(d.ReferenceInfo); err != nil {
		return fmt.Errorf("offerItem.SetReferenceInfo: %w", err)
	}
	if err := offerItem.SetOtherInfo(d.OtherInfo); err != nil {
		return fmt.Errorf("offerItem.SetOtherInfo: %w", err)
	}
	offerItem.SetIsInvitationMailSent(d.IsInvitationMailSent)
	offerItem.SetIsOfferDetailMailSent(d.IsOfferDetailMailSent)
	offerItem.SetIsArticlePostMailSent(d.IsArticlePostMailSent)
	offerItem.SetIsPassedPreliminaryReviewMailSent(d.IsPassedPreliminaryReviewMailSent)
	offerItem.SetIsFailedPreliminaryReviewMailSent(d.IsFailedPreliminaryReviewMailSent)
	offerItem.SetIsPassedAfterReviewMailSent(d.IsPassedAfterReviewMailSent)
	offerItem.SetIsFailedAfterReviewMailSent(d.IsFailedAfterReviewMailSent)
	schedulesMap := func() map[model.ScheduleID]model.Schedule {
		m := make(map[model.ScheduleID]model.Schedule)
		for _, s := range offerItem.Schedules() {
			m[s.ID()] = *s
		}
		return m
	}()
	schedules := make([]*model.Schedule, 0, len(d.Schedules))
	for _, scheduleDTO := range d.Schedules {
		s, ok := schedulesMap[model.ScheduleID(*scheduleDTO.ID)]
		if !ok {
			return fmt.Errorf("notfound: %s", *scheduleDTO.ID)
		}

		if err := s.SetDate(scheduleDTO.StartDate, scheduleDTO.EndDate); err != nil {
			return fmt.Errorf("s.SetDate: %w", err)
		}
		schedules = append(schedules, &s)
	}

	if err := offerItem.SetSchedules(schedules); err != nil {
		return fmt.Errorf("offerItem.SetSchedules: %w", err)
	}

	draftedItemInfoMinCommission, err := model.NewCommission(
		model.CommissionType(d.DraftedItemInfo.MinCommission.CommissionType),
		float32(d.DraftedItemInfo.MinCommission.CalculatedRate),
	)
	if err != nil {
		return fmt.Errorf("failed to create minCommission: %w", err)
	}

	draftedItemInfoMaxCommission, err := model.NewCommission(
		model.CommissionType(d.DraftedItemInfo.MaxCommission.CommissionType),
		float32(d.DraftedItemInfo.MaxCommission.CalculatedRate),
	)
	if err != nil {
		return fmt.Errorf("failed to create maxCommission: %w", err)
	}

	if err := offerItem.SetDraftedItemInfo(d.DraftedItemInfo.Name, d.DraftedItemInfo.ContentName, d.DraftedItemInfo.ImageURL, d.DraftedItemInfo.URL, draftedItemInfoMinCommission, draftedItemInfoMaxCommission, offerItem.ID()); err != nil {
		return fmt.Errorf("offerItem.SetdraftedItemInfo: %w", err)
	}

	var dfItemID *model.DFItemID
	if offerItem.DfItem().Exists() {
		optDFITemID := offerItem.DfItem().ID()
		dfItemID = &optDFITemID
	}

	// TODO: bannerIDは今後複数対応を行う
	var bannerIDs []model.BannerID
	if d.CouponBannerID != nil && *d.CouponBannerID != "" {
		bannerIDs = append(bannerIDs, model.BannerID(*d.CouponBannerID))
	}

	offerItem.SetPickInfo(model.ItemID(d.ItemID), dfItemID, bannerIDs)

	return nil
}

func setAssigneeFields(assignee *model.Assignee, d *dto.Assignee) error {
	if err := assignee.SetWritingFee(d.WritingFee); err != nil {
		return fmt.Errorf("assignee.SetWritingFee: %w", err)
	}

	assignee.SetStageAll(converter.StageDTOToModel(d.Stage))
	return nil
}

func (o *offerItemUsecaseImpl) addItemInfo(ctx context.Context, offerItems model.OfferItemList) error {
	itemIdentifiers := offerItems.ItemIdentifiers()

	// 案件ID、DF案件IDが設定されているオファー案件がない場合、何もしない
	if len(itemIdentifiers) == 0 {
		return nil
	}

	// 案件情報を取得
	itemMap, err := o.affiliateItemAdapter.BulkGetItems(ctx, itemIdentifiers)
	if err != nil {
		return fmt.Errorf("o.affiliateItemAdapter.BulkGetItems: %w", err)
	}

	// 取得した情報を適用する
	for _, offerItem := range offerItems {
		var dfItemID model.DFItemID
		if offerItem.DfItem() != nil {
			dfItemID = offerItem.DfItem().ID()
		}

		// 案件情報を適用
		if items, ok := itemMap[*model.NewItemIdentifier(offerItem.Item().ID(), dfItemID)]; ok {
			if err := offerItem.SetItem(&items.Item); err != nil {
				return fmt.Errorf("offerItem.SetItem: %w", err)
			}
			if items.DFItem.Exists() {
				offerItem.SetDFItem(&items.DFItem)
			}
		}
	}

	return nil
}

func (o *offerItemUsecaseImpl) GetOfferItem(ctx context.Context, offerItemID model.OfferItemID) (*model.OfferItem, error) {
	ctx, span := trace.StartSpan(ctx, "offerItemUsecaseImpl.GetOfferItem")
	defer span.End()

	offerItem, err := o.offerItemRepository.Get(ctx, o.db, offerItemID, false)
	if err != nil {
		return nil, fmt.Errorf("o.offerItemRepository.Get: %w", err)
	}

	// アイテム情報を付与する
	if err = o.addItemInfo(ctx, model.OfferItemList{offerItem}); err != nil {
		return nil, fmt.Errorf("o.addItemInfo: %w", err)
	}

	return offerItem, nil
}

// オファー案件一覧を取得する
func (o *offerItemUsecaseImpl) ListOfferItem(ctx context.Context, condition *model.ListCondition) (*model.ListOfferItemResult, error) {
	ctx, span := trace.StartSpan(ctx, "offerItemUsecaseImpl.ListOfferItem")
	defer span.End()

	result, err := o.offerItemRepository.List(ctx, o.db, condition, true)
	if err != nil {
		return nil, fmt.Errorf("o.offerItemRepository.List: %w", err)
	}

	// アイテム情報を付与する
	if err = o.addItemInfo(ctx, result.OfferItems()); err != nil {
		return nil, fmt.Errorf("o.addItemInfo: %w", err)
	}

	return result, nil
}
