package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/presentation/converter"
	"github.com/terui-ryota/offer-item/internal/application/usecase"
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	offer_item "github.com/terui-ryota/protofiles/go/offer_item"
)

func NewOfferItemHandler(offerItemUsecase usecase.OfferItemUsecase, assigneeUsecase usecase.AssigneeUsecase) offer_item.OfferItemHandlerServer {
	return &offerItemHandler{
		offerItemUsecase: offerItemUsecase,
		assigneeUsecase:  assigneeUsecase,
	}
}

type offerItemHandler struct {
	offerItemUsecase usecase.OfferItemUsecase
	assigneeUsecase  usecase.AssigneeUsecase
	offer_item.UnimplementedOfferItemHandlerServer
}

func (h *offerItemHandler) Decline(ctx context.Context, req *offer_item.DeclineRequest) (*offer_item.DeclineResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())
	amebaID := model.AmebaID(req.GetAmebaId())
	reason := req.GetDeclineReason()

	if err := h.assigneeUsecase.Decline(ctx, offerItemID, amebaID, reason); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.Decline: %w", err)
	}

	return &offer_item.DeclineResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) Invitation(ctx context.Context, req *offer_item.InvitationRequest) (*offer_item.InvitationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())
	amebaID := model.AmebaID(req.GetAmebaId())

	questionAnswers := make(map[model.QuestionID]string)
	for _, a := range req.GetQuestionAnswers() {
		questionAnswers[model.QuestionID(a.GetQuestionId())] = a.GetContent()
	}

	if err := h.assigneeUsecase.Invitation(ctx, offerItemID, amebaID, req.InvitationFlag, questionAnswers); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.Invitation: %w", err)
	}

	return &offer_item.InvitationResponse{
		Request: req,
	}, nil
}

// BulkGetQuestionnaireQuestionAnswers implements offer_item_v2.OfferItemHandlerServer.
func (h *offerItemHandler) BulkGetQuestionnaireQuestionAnswers(ctx context.Context, req *offer_item.BulkGetQuestionnaireQuestionAnswersRequest) (*offer_item.BulkGetQuestionnaireQuestionAnswersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}
	answers, err := h.assigneeUsecase.BulkGetQuestionnaireQuestionAnswers(ctx, model.OfferItemID(req.OfferItemId), func() []model.AmebaID {
		ids := make([]model.AmebaID, 0, len(req.GetAmebaIds()))
		for _, a := range req.GetAmebaIds() {
			ids = append(ids, model.AmebaID(a))
		}
		return ids
	}())
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.BulkGetQuestionnaireQuestionAnswers: %w", err)
	}
	return &offer_item.BulkGetQuestionnaireQuestionAnswersResponse{
		Request: req,
		AmebaIdAnswersMap: func() map[string]*offer_item.BulkGetQuestionnaireQuestionAnswersResponse_QuestionAnswerMap {
			m := make(map[string]*offer_item.BulkGetQuestionnaireQuestionAnswersResponse_QuestionAnswerMap)
			for amebaID, answers := range answers {
				as := &offer_item.BulkGetQuestionnaireQuestionAnswersResponse_QuestionAnswerMap{
					QuestionIdAnswerMap: func() map[string]*offer_item.QuestionAnswer {
						res := make(map[string]*offer_item.QuestionAnswer)
						for _, answer := range answers {
							res[answer.QuestionID().String()] = converter.QuestionAnswerModelToPB(&answer)
						}
						return res
					}(),
				}
				m[amebaID.String()] = as
			}
			return m
		}(),
	}, nil
}

func (h *offerItemHandler) GetAssigneeByAmebaIDOfferItemID(ctx context.Context, req *offer_item.GetAssigneeByAmebaIDOfferItemIDRequest) (*offer_item.GetAssigneeByAmebaIDOfferItemIDResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	amebaID := model.AmebaID(req.GetAmebaId())
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	assignee, err := h.assigneeUsecase.GetAssigneeByAmebaIDOfferItemID(ctx, amebaID, offerItemID)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.GetAssignee: %w", err)
	}

	return &offer_item.GetAssigneeByAmebaIDOfferItemIDResponse{
		Request:  req,
		Assignee: converter.AssigneeModelToPB(assignee),
	}, nil
}

func (h *offerItemHandler) FinishedShipment(ctx context.Context, req *offer_item.FinishedShipmentRequest) (*offer_item.FinishedShipmentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	if err := h.assigneeUsecase.FinishedShipment(ctx, offerItemID); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.FinishedShipment: %w", err)
	}

	return &offer_item.FinishedShipmentResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) CompletedOfferItem(ctx context.Context, req *offer_item.CompletedOfferItemRequest) (*offer_item.CompletedOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	if err := h.assigneeUsecase.CompletedOfferItem(ctx, offerItemID); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.CompletedOfferItem: %w", err)
	}

	return &offer_item.CompletedOfferItemResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) PaymentCompleted(ctx context.Context, req *offer_item.PaymentCompletedRequest) (*offer_item.PaymentCompletedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())
	amebaIDs := make([]model.AmebaID, 0, len(req.GetAmebaIds()))
	for _, amebaID := range req.GetAmebaIds() {
		amebaIDs = append(amebaIDs, model.AmebaID(amebaID))
	}

	if err := h.assigneeUsecase.PaymentCompleted(ctx, offerItemID, amebaIDs); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.PaymentCompleted: %w", err)
	}

	return &offer_item.PaymentCompletedResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) HealthCheck(ctx context.Context, req *offer_item.HealthCheckReq) (*offer_item.HealthCheckRes, error) {
	fmt.Print("============HealthCheck===============")
	// TODO: HealthCheckの実装を追加
	return &offer_item.HealthCheckRes{
		Request: req,
		Num:     1,
	}, nil
}

// オファー案件を作成する
func (h *offerItemHandler) SaveOfferItem(ctx context.Context, req *offer_item.SaveOfferItemRequest) (*offer_item.SaveOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// DTOに変換する
	offerItemDTO := dto.SaveOfferItemPBToDTO(req.GetOfferItem())

	// 作成する
	if err := h.offerItemUsecase.SaveOfferItem(ctx, offerItemDTO); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.CreateOfferItem: %w", err)
	}

	return &offer_item.SaveOfferItemResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) GetOfferItem(ctx context.Context, req *offer_item.GetOfferItemRequest) (*offer_item.GetOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	offerItem, err := h.offerItemUsecase.GetOfferItem(ctx, offerItemID)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.GetOfferItem: %w", err)
	}
	offerItemPB, err := converter.OfferItemModelToPB(offerItem)
	if err != nil {
		return nil, fmt.Errorf("converter.OfferItemModelToPB: %w", err)
	}

	return &offer_item.GetOfferItemResponse{
		Request:   req,
		OfferItem: offerItemPB,
	}, nil
}

// オファー案件一覧取得
func (h *offerItemHandler) ListOfferItem(ctx context.Context, req *offer_item.ListOfferItemRequest) (*offer_item.ListOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	condition, err := converter.ListConditionPBToModel(req.GetCondition())
	if err != nil {
		return nil, fmt.Errorf("converter.ListConditionPBToModel: %w", err)
	}

	// オファー案件一覧を取得
	result, err := h.offerItemUsecase.ListOfferItem(ctx, condition)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.ListOfferItem: %w", err)
	}

	// protoに変換する
	offerItemPBs := make([]*offer_item.OfferItem, 0, len(result.OfferItems()))
	for _, offerItem := range result.OfferItems() {
		offerItemPB, err := converter.OfferItemModelToPB(offerItem)
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemModelToPB: %w", err)
		}
		offerItemPBs = append(offerItemPBs, offerItemPB)
	}

	return &offer_item.ListOfferItemResponse{
		Request:    req,
		OfferItems: offerItemPBs,
		Result:     converter.ListResultModelToPB(result.ListResult()),
	}, nil
}

// オファー案件を削除する
func (h *offerItemHandler) DeleteOfferItem(ctx context.Context, req *offer_item.DeleteOfferItemRequest) (*offer_item.DeleteOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	// 削除する
	if err := h.offerItemUsecase.DeleteOfferItem(ctx, offerItemID); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.DeleteOfferItem: %w", err)
	}

	return &offer_item.DeleteOfferItemResponse{Request: req}, nil
}

// オファー案件を検索する
func (h *offerItemHandler) SearchOfferItem(ctx context.Context, req *offer_item.SearchOfferItemRequest) (*offer_item.SearchOfferItemResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// 検索条件を設定
	searchCriteria := &dto.SearchOfferItemCriteria{}
	if req.OptionalItemId != nil {
		itemID := model.ItemID(req.GetItemId())
		searchCriteria.ItemIDEqual = &itemID
	}

	if req.OptionalDfItemId != nil {
		dfItemID := model.DFItemID(req.GetDfItemId())
		searchCriteria.DfItemIDEqual = &dfItemID
	}

	if req.OptionalOfferItemName != nil {
		offerItemName := req.GetOfferItemName()
		searchCriteria.NameContains = &offerItemName
	}

	condition, err := converter.ListConditionPBToModel(req.GetCondition())
	if err != nil {
		return nil, fmt.Errorf("converter.ListConditionPBToModel: %w", err)
	}

	searchResult, err := h.offerItemUsecase.SearchOfferItem(ctx, searchCriteria, condition)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.SearchOfferItem: %w", err)
	}

	// protoに変換する
	offerItemPBs := make([]*offer_item.OfferItem, 0, len(searchResult.OfferItems()))
	for _, offerItem := range searchResult.OfferItems() {
		offerItemPB, err := converter.OfferItemModelToPB(offerItem)
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemModelToPB: %w", err)
		}
		offerItemPBs = append(offerItemPBs, offerItemPB)
	}

	return &offer_item.SearchOfferItemResponse{
		Request:    req,
		OfferItems: offerItemPBs,
		Result:     converter.ListResultModelToPB(searchResult.ListResult()),
	}, nil
}

func (h *offerItemHandler) ListAssigneeOfferItemPair(ctx context.Context, req *offer_item.ListAssigneeOfferItemPairRequest) (*offer_item.ListAssigneeOfferItemPairResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	amebaID := model.AmebaID(req.GetAmebaId())

	listAssigneeOfferItemPairs, err := h.offerItemUsecase.ListAssigneeOfferItemPair(ctx, amebaID)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.BulkGetAssigneesByAmebaID: %w", err)
	}
	// protoに変換する
	assigneeOfferItemPairPBs := make([]*offer_item.AssigneeOfferItemPair, 0, len(listAssigneeOfferItemPairs))
	for _, assigneeOfferItemPair := range listAssigneeOfferItemPairs {
		offerItemPB, err := converter.OfferItemModelToPB(assigneeOfferItemPair.OfferItem())
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemModelToPB: %w", err)
		}
		assigneeOfferItemPairPB := &offer_item.AssigneeOfferItemPair{
			Assignee:  converter.AssigneeModelToPB(assigneeOfferItemPair.Assignee()),
			OfferItem: offerItemPB,
		}
		assigneeOfferItemPairPBs = append(assigneeOfferItemPairPBs, assigneeOfferItemPairPB)
	}

	return &offer_item.ListAssigneeOfferItemPairResponse{
		Request:               req,
		AssigneeOfferItemPair: assigneeOfferItemPairPBs,
	}, nil
}

// GetQuestionnaire implements offer_item_v2.OfferItemHandlerServer.
func (h *offerItemHandler) GetQuestionnaire(ctx context.Context, req *offer_item.GetQuestionnaireRequest) (*offer_item.GetQuestionnaireResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}
	q, err := h.offerItemUsecase.GetQuestionnaire(ctx, model.OfferItemID(req.GetOfferItemId()))
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.GetQuestionnaire: %w", err)
	}
	return &offer_item.GetQuestionnaireResponse{
		Request:       req,
		Questionnaire: converter.QuestionnaireModelToPB(q),
	}, nil
}

func (h *offerItemHandler) ListAssignee(ctx context.Context, req *offer_item.ListAssigneeRequest) (*offer_item.ListAssigneeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())
	stage := model.Stage(req.GetStage())

	assignees, err := h.assigneeUsecase.ListAssignee(ctx, offerItemID, stage)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.ListAssignee: %w", err)
	}

	// protoに変換する
	assigneePBs := make([]*offer_item.Assignee, 0, len(assignees))
	for _, assignee := range assignees {
		assigneePBs = append(assigneePBs, converter.AssigneeModelToPB(assignee))
	}

	return &offer_item.ListAssigneeResponse{
		Request:   req,
		Assignees: assigneePBs,
	}, nil
}

func (h *offerItemHandler) ListAssigneeUnderExamination(ctx context.Context, req *offer_item.ListAssigneeUnderExaminationRequest) (*offer_item.ListAssigneeUnderExaminationResponse, error) {
	assignees, err := h.assigneeUsecase.ListAssigneeUnderExamination(ctx)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.ListAssigneeUnderExamination: %w", err)
	}

	// protoに変換する
	assigneePBs := make([]*offer_item.Assignee, 0, len(assignees))
	for _, assignee := range assignees {
		assigneePBs = append(assigneePBs, converter.AssigneeModelToPB(assignee))
	}

	return &offer_item.ListAssigneeUnderExaminationResponse{
		Request:   req,
		Assignees: assigneePBs,
	}, nil
}

func (h *offerItemHandler) ListStageAssigneeCount(ctx context.Context, req *offer_item.ListStageAssigneeCountRequest) (*offer_item.ListStageAssigneeCountResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	assigneeCounts, err := h.assigneeUsecase.ListAssigneeCount(ctx, offerItemID)
	if err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.ListAssigneeCount: %w", err)
	}

	// protoに変換する
	assigneeCountsPB := make([]*offer_item.AssigneeCount, 0, len(assigneeCounts))
	for _, assigneeCount := range assigneeCounts {
		assigneeCountPB := converter.AssigneeCountModelToPB(assigneeCount)
		assigneeCountsPB = append(assigneeCountsPB, assigneeCountPB)
	}

	return &offer_item.ListStageAssigneeCountResponse{
		Request:        req,
		AssigneeCounts: assigneeCountsPB,
	}, nil
}

func (h *offerItemHandler) InviteOffer(ctx context.Context, req *offer_item.InviteOfferRequest) (*offer_item.InviteOfferResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())

	if err := h.assigneeUsecase.InviteOffer(ctx, offerItemID); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.InviteOffer: %w", err)
	}

	return &offer_item.InviteOfferResponse{
		Request: req,
	}, nil
}

func (h *offerItemHandler) UploadLotteryResults(ctx context.Context, req *offer_item.UploadLotteryResultsRequest) (*offer_item.UploadLotteryResultsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, apperr.OfferItemValidationError.Wrap(err)
	}

	// 古いフィールドを参照したら怒るようにする(後続処理を行わずにエラー)
	// adminを上げた後、protofilesを削除し、この行を削除する
	//nolint:staticcheck
	if len(req.GetMapLotteryResult()) > 0 {
		return nil, apperr.OfferItemValidationError.Wrap(errors.New("please don't specify this field coz this is deprecated"))
	}

	// モデルに変換する
	offerItemID := model.OfferItemID(req.GetOfferItemId())
	mapLotteryResult := make(map[model.AmebaID]model.LotteryResult, len(req.GetMapLotteryResultWithShippingData()))
	for amebaID, lotteryResult := range req.GetMapLotteryResultWithShippingData() {
		var janCode *string
		if lotteryResult.GetOptionalJanCode() != nil {
			janCode = func() *string {
				s := lotteryResult.GetJanCode()
				return &s
			}()
		}
		lr := model.NewLotteryResult(lotteryResult.GetIsPassedLottery(), lotteryResult.GetShippingData(), janCode)
		mapLotteryResult[model.AmebaID(amebaID)] = *lr
	}

	if err := h.assigneeUsecase.UploadLotteryResults(ctx, offerItemID, mapLotteryResult); err != nil {
		return nil, fmt.Errorf("h.offerItemUsecase.UploadLotteryResult: %w", err)
	}

	return &offer_item.UploadLotteryResultsResponse{
		Request: req,
	}, nil
}
