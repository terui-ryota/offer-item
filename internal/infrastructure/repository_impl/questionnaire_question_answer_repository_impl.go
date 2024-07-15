package repository_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.opencensus.io/trace"
)

func NewQuestionnaireQuestionAnswerRepositoryImpl(db *sql.DB) repository.QuestionnaireQuestionAnswerRepository {
	return &questionnaireQuestionAnswerRepository{
		db: db,
	}
}

type questionnaireQuestionAnswerRepository struct {
	db *sql.DB
}

// DeleteByOfferItemID implements repository.QuestionnaireQuestionAnswerRepository.
func (r *questionnaireQuestionAnswerRepository) DeleteByOfferItemID(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID) error {
	ctx, span := trace.StartSpan(ctx, "questionnaireQuestionAnswerRepository.DeleteByOfferItemID")
	defer span.End()

	if _, err := entity.QuestionnaireQuestionAnswers(
		entity.QuestionnaireQuestionAnswerWhere.OfferItemID.EQ(offerItemID.String()),
	).DeleteAll(ctx, tx); err != nil {
		return fmt.Errorf("entity.QuestionnaireQuestionAnswers.DeleteAll: %w", err)
	}
	return nil
}

// BulkGetByOfferItemIDAndAssigneeIDs implements repository.QuestionnaireQuestionAnswerRepository.
func (r *questionnaireQuestionAnswerRepository) BulkGetByOfferItemIDAndAssigneeIDs(ctx context.Context, offerItemID model.OfferItemID, assigneeIDs []model.AssigneeID) (map[model.AssigneeID]map[model.QuestionID]model.QuestionAnswer, error) {
	ctx, span := trace.StartSpan(ctx, "questionnaireQuestionAnswerRepository.BulkgetByOfferItemIDAndAssigneeIDs")
	defer span.End()

	if len(assigneeIDs) == 0 {
		return make(map[model.AssigneeID]map[model.QuestionID]model.QuestionAnswer), nil
	}
	answers, err := entity.QuestionnaireQuestionAnswers(
		entity.QuestionnaireQuestionAnswerWhere.OfferItemID.EQ(offerItemID.String()),
		entity.QuestionnaireQuestionAnswerWhere.AssigneeID.IN(func() []string {
			ids := make([]string, 0, len(assigneeIDs))
			for _, i := range assigneeIDs {
				ids = append(ids, i.String())
			}
			return ids
		}()),
	).All(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make(map[model.AssigneeID]map[model.QuestionID]model.QuestionAnswer), nil
		}
		return nil, fmt.Errorf("entity.QuestionnaireQuestionAnswers.All: %w", err)
	}
	res := make(map[model.AssigneeID]map[model.QuestionID]model.QuestionAnswer)
	for _, a := range answers {
		as, ok := res[model.AssigneeID(a.AssigneeID)]
		if !ok {
			as = make(map[model.QuestionID]model.QuestionAnswer)
		}
		as[model.QuestionID(a.QuestionnaireQuestionID)] = *model.NewQuestionAnswerFromRepository(model.AssigneeID(a.AssigneeID), model.OfferItemID(a.OfferItemID), a.QuestionnaireQuestionID, a.Answer)
		res[model.AssigneeID(a.AssigneeID)] = as
	}
	return res, nil
}

func (r *questionnaireQuestionAnswerRepository) Save(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID, assigneeID model.AssigneeID, answers []model.QuestionAnswer) error {
	ctx, span := trace.StartSpan(ctx, "questionnaireQuestionAnswerRepository.Save")
	defer span.End()

	if _, err := entity.QuestionnaireQuestionAnswers(
		entity.QuestionnaireQuestionAnswerWhere.OfferItemID.EQ(offerItemID.String()),
		entity.QuestionnaireQuestionAnswerWhere.AssigneeID.EQ(assigneeID.String()),
	).DeleteAll(ctx, tx); err != nil {
		return fmt.Errorf("entity.QuestionnaireQuestionAnswers.DeleteAll: %w", err)
	}
	for _, a := range answers {
		if err := convertQuestionnaireQuestionAnswerToEntity(&a).Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("convertQuestionnaireQuestionAnswerToEntity.Insert: %w", err)
		}
	}
	return nil
}

func convertQuestionnaireQuestionAnswerToEntity(m *model.QuestionAnswer) *entity.QuestionnaireQuestionAnswer {
	return &entity.QuestionnaireQuestionAnswer{
		AssigneeID:              m.AssigneeID().String(),
		QuestionnaireQuestionID: m.QuestionID().String(),
		OfferItemID:             m.OfferItemID().String(),
		Answer:                  m.Content(),
	}
}
