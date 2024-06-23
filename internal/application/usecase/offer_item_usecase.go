package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/dto"
)

type OfferItemUsecase interface {
	//SaveOfferItem(ctx context.Context, offerItemDTO *dto.OfferItemDTO) error
	SaveOfferItem(ctx context.Context, offerItemDTO *dto.OfferItemDTO) error
}

func NewOfferItemUsecase(
	db *sql.DB,
	// offerItemRepository repository.OfferItemRepository,
	// assigneeRepository repository.AssigneeRepository,
	// questionnaireRepository repository.QuestionnaireRepository,
	// questionnaireQuestionAnswerRepository repository.QuestionnaireQuestionAnswerRepository,
	// affiliateItemAdapter adapter.AffiliateItemAdapter,
	// affiliatorAdapter adapter.AffiliatorAdapter,
	// examinationRepository repository.ExaminationRepository,
	// validationConfig *config.ValidationConfig,
) OfferItemUsecase {
	return &offerItemUsecaseImpl{
		db: db,
		//offerItemRepository:                   offerItemRepository,
		//assigneeRepository:                    assigneeRepository,
		//questionnaireRepository:               questionnaireRepository,
		//questionnaireQuestionAnswerRepository: questionnaireQuestionAnswerRepository,
		//affiliateItemAdapter:                  affiliateItemAdapter,
		//affiliatorAdapter:                     affiliatorAdapter,
		//examinationRepository:                 examinationRepository,
		//validationConfig:                      validationConfig,
	}
}

type offerItemUsecaseImpl struct {
	db *sql.DB
	//offerItemRepository                   repository.OfferItemRepository
	//assigneeRepository                    repository.AssigneeRepository
	//questionnaireRepository               repository.QuestionnaireRepository
	//questionnaireQuestionAnswerRepository repository.QuestionnaireQuestionAnswerRepository
	//affiliateItemAdapter                  adapter.AffiliateItemAdapter
	//affiliatorAdapter                     adapter.AffiliatorAdapter
	//examinationRepository                 repository.ExaminationRepository
	//validationConfig                      *config.ValidationConfig
}

func (o *offerItemUsecaseImpl) SaveOfferItem(ctx context.Context, offerItemDTO *dto.OfferItemDTO) error {

	fmt.Println(offerItemDTO)
	//ctx, span := trace.StartSpan(ctx, "offerItemUsecaseImpl.SaveOfferItem")
	//defer span.End()

	return nil
}
