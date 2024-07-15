package repository_impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/offer-item/pkg/logger"
	null "github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opencensus.io/trace"
)

func NewQuestionnaireRepositoryImpl() repository.QuestionnaireRepository {
	return &questionnaireRepositoryImpl{}
}

type questionnaireRepositoryImpl struct{}

// BulkGet implements repository.QuestionnaireRepository.
func (*questionnaireRepositoryImpl) BulkGet(ctx context.Context, exec boil.ContextExecutor, ids []model.OfferItemID) (map[model.OfferItemID]model.Questionnaire, error) {
	ctx, span := trace.StartSpan(ctx, "questionnaireRepositoryImpl.BulkGet")
	defer span.End()

	if len(ids) == 0 {
		return make(map[model.OfferItemID]model.Questionnaire), nil
	}
	idStrs := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrs = append(idStrs, id.String())
	}
	questionnaires, err := entity.Questionnaires(
		entity.QuestionnaireWhere.OfferItemID.IN(idStrs),
	).All(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make(map[model.OfferItemID]model.Questionnaire), nil
		}
		return nil, fmt.Errorf("entity.Questionnaires.All: %w", err)
	}
	questions, err := entity.QuestionnaireQuestions(
		entity.QuestionnaireQuestionWhere.OfferItemID.IN(idStrs),
	).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.QuestionnaireQuestions.All: %w", err)
	}
	questionsMap := make(map[string][]*entity.QuestionnaireQuestion)
	for _, q := range questions {
		qs, ok := questionsMap[q.OfferItemID]
		if !ok {
			qs = make([]*entity.QuestionnaireQuestion, 0)
		}
		qs = append(qs, q)
		questionsMap[q.OfferItemID] = qs
	}
	res := make(map[model.OfferItemID]model.Questionnaire)
	for _, q := range questionnaires {
		m := convertQuestionnaireToModel(ctx, q, questionsMap[q.OfferItemID])
		res[model.OfferItemID(q.OfferItemID)] = *m
	}
	return res, nil
}

// Delete implements repository.QuestionnaireRepository.
func (*questionnaireRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id model.OfferItemID) error {
	ctx, span := trace.StartSpan(ctx, "questionnaireRepositoryImpl.Delete")
	defer span.End()

	qs, err := func() ([]*entity.QuestionnaireQuestion, error) {
		qs, err := entity.QuestionnaireQuestions(
			entity.QuestionnaireQuestionWhere.OfferItemID.EQ(id.String()),
		).All(ctx, tx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return make([]*entity.QuestionnaireQuestion, 0), nil
			}
			return nil, fmt.Errorf("entity.QuestionnaireQuestions.All: %w", err)
		}
		return qs, nil
	}()
	if err != nil {
		return fmt.Errorf("listQuestionnaireQuestions: %w", err)
	}
	if len(qs) != 0 {
		if _, err := entity.QuestionnaireQuestions(
			entity.QuestionnaireQuestionWhere.OfferItemID.EQ(id.String()),
		).DeleteAll(ctx, tx); err != nil {
			return fmt.Errorf("entity.QuestionnaireQuestions.DeleteAll: %w", err)
		}
	}
	if _, err := entity.Questionnaires(
		entity.QuestionnaireWhere.OfferItemID.EQ(id.String()),
	).DeleteAll(ctx, tx); err != nil {
		return fmt.Errorf("entity.Questionnaires.DeleteAll: %w", err)
	}
	return nil
}

// Get implements repository.QuestionnaireRepository.
func (*questionnaireRepositoryImpl) Get(ctx context.Context, exec boil.ContextExecutor, id model.OfferItemID, withLock bool) (*model.Questionnaire, error) {
	ctx, span := trace.StartSpan(ctx, "questionnaireRepositoryImpl.Get")
	defer span.End()

	mods := make([]qm.QueryMod, 0)
	mods = append(mods, entity.QuestionnaireWhere.OfferItemID.EQ(id.String()))
	if withLock {
		mods = append(mods, qm.For("UPDATE"))
	}
	questionnaire, err := entity.Questionnaires(mods...).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.OfferItemNotFoundError
		}
		return nil, fmt.Errorf("entity.Questionnaires: %w", err)
	}
	questions, err := entity.QuestionnaireQuestions(
		entity.QuestionnaireQuestionWhere.OfferItemID.EQ(id.String()),
	).All(ctx, exec)
	if errors.Is(err, sql.ErrNoRows) {
		questions = make(entity.QuestionnaireQuestionSlice, 0)
	}
	return convertQuestionnaireToModel(ctx, questionnaire, questions), nil
}

func convertQuestionnaireToModel(ctx context.Context, questionnaire *entity.Questionnaire, questions []*entity.QuestionnaireQuestion) *model.Questionnaire {
	sort.Slice(questions, func(i, j int) bool {
		return questions[i].Priority < questions[j].Priority
	})
	qs := make([]model.Question, 0, len(questions))
	for _, q := range questions {
		qs = append(qs, *model.NewQuestionFromRepository(
			model.QuestionID(q.ID),
			model.OfferItemID(q.OfferItemID),
			model.NewQuestionType(q.Type),
			q.Title,
			q.Image,
			func() []string {
				if !q.AnswerOptions.Valid {
					return nil
				}
				res := make([]string, 0)
				if err := q.AnswerOptions.Unmarshal(&res); err != nil {
					logger.FromContext(ctx).Errorf("failed to marshal: %w", err)
				}
				return res
			}(),
		))
	}
	return model.NewQuestionnaireFromRepository(
		model.OfferItemID(questionnaire.OfferItemID),
		questionnaire.Description,
		qs,
	)
}

// Save implements repository.QuestionnaireRepository.
func (*questionnaireRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, m model.Questionnaire) error {
	ctx, span := trace.StartSpan(ctx, "questionnaireRepositoryImpl.Save")
	defer span.End()

	questionnaire, questions := convertQuestionnaireToEntity(m)
	if _, err := entity.Questionnaires(entity.QuestionnaireWhere.OfferItemID.EQ(m.OfferItemID().String())).DeleteAll(ctx, tx); err != nil {
		return fmt.Errorf("entity.Questionnaires.DeleteAll: %w", err)
	}
	if err := questionnaire.Insert(ctx, tx, boil.Infer()); err != nil {
		return fmt.Errorf("questionnaire.Insert: %w", err)
	}
	if _, err := entity.QuestionnaireQuestions(entity.QuestionnaireQuestionWhere.OfferItemID.EQ(m.OfferItemID().String())).DeleteAll(ctx, tx); err != nil {
		return fmt.Errorf("entity.QuestionnaireQuestions.DeleteAll: %w", err)
	}
	for _, q := range questions {
		if err := q.Insert(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("q.Insert: %w", err)
		}
	}
	return nil
}

func convertQuestionnaireToEntity(m model.Questionnaire) (entity.Questionnaire, []*entity.QuestionnaireQuestion) {
	questionnaire := entity.Questionnaire{
		OfferItemID: m.OfferItemID().String(),
		Description: m.Description(),
	}
	questions := make([]*entity.QuestionnaireQuestion, 0, len(m.Questions()))
	for idx, q := range m.Questions() {
		questions = append(questions, &entity.QuestionnaireQuestion{
			ID:          q.ID().String(),
			OfferItemID: q.OfferItemID().String(),
			Title:       q.Title(),
			Type:        int(q.QuestionType()),
			Image:       q.ImageURL(),
			Priority:    idx + 1,
			AnswerOptions: func() null.JSON {
				if len(q.Options()) == 0 {
					return null.JSONFromPtr(nil)
				}
				bs, err := json.Marshal(q.Options())
				if err != nil {
					logger.Default().Errorf("json.Marshal: %w", err)
					return null.JSONFromPtr(nil)
				} else {
					return null.JSONFromPtr(&bs)
				}
			}(),
		})
	}
	return questionnaire, questions
}
