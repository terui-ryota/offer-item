//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package repository

import (
	"context"
	"database/sql"

	"github.com/terui-ryota/offer-item/internal/domain/model"
)

type QuestionnaireQuestionAnswerRepository interface {
	Save(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID, assigneeID model.AssigneeID, answers []model.QuestionAnswer) error
	BulkGetByOfferItemIDAndAssigneeIDs(ctx context.Context, offerItemID model.OfferItemID, assigneeIDs []model.AssigneeID) (map[model.AssigneeID]map[model.QuestionID]model.QuestionAnswer, error)
	DeleteByOfferItemID(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID) error
}
