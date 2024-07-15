//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package repository

import (
	"context"
	"database/sql"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type QuestionnaireRepository interface {
	Get(ctx context.Context, exec boil.ContextExecutor, id model.OfferItemID, withLock bool) (*model.Questionnaire, error)
	BulkGet(ctx context.Context, exec boil.ContextExecutor, ids []model.OfferItemID) (map[model.OfferItemID]model.Questionnaire, error)
	Save(ctx context.Context, tx *sql.Tx, questionnaire model.Questionnaire) error
	Delete(ctx context.Context, tx *sql.Tx, id model.OfferItemID) error
}
