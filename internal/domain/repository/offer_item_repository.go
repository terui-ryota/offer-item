//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type OfferItemRepository interface {
	List(ctx context.Context, exec boil.ContextExecutor, condition *model.ListCondition, isClosed bool) (*model.ListOfferItemResult, error)
	Delete(ctx context.Context, tx *sql.Tx, id model.OfferItemID) error
	Search(ctx context.Context, exec boil.ContextExecutor, searchCriteria *dto.SearchOfferItemCriteria, condition *model.ListCondition) (*model.ListOfferItemResult, error)
	Get(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID, withLock bool) (*model.OfferItem, error)
	Create(ctx context.Context, tx *sql.Tx, offerItem *model.OfferItem) error
	Update(ctx context.Context, tx *sql.Tx, offerItem *model.OfferItem) error
	BulkGet(ctx context.Context, exec boil.ContextExecutor, ids []model.OfferItemID, isClosed bool) (map[model.OfferItemID]*model.OfferItem, error)
	ListIDsByEndDate(ctx context.Context, exec boil.ContextExecutor, sinceEndDate, untilEndDate time.Time) (model.OfferItemIDList, error)
}
