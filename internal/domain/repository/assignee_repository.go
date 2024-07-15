//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package repository

import (
	"context"
	"database/sql"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type AssigneeRepository interface {
	ListByOfferItemIDAmebaIDs(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID, amebaIDs []model.AmebaID, withLock bool) (model.AssigneeList, error)
	Update(ctx context.Context, exec boil.ContextExecutor, assignee *model.Assignee) error
	Create(ctx context.Context, tx *sql.Tx, assignee *model.Assignee) error
	BulkGetByOfferItemIDAmebaIDs(ctx context.Context, db *sql.DB, offerItemID model.OfferItemID, amebaIDs []model.AmebaID, withLock bool) (map[model.AmebaID]*model.Assignee, error)
	ListByOfferItemIDStage(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID, stage model.Stage) (model.AssigneeList, error)
	ListUnderExamination(ctx context.Context, exec boil.ContextExecutor) (model.AssigneeList, error)
	ListCount(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID) ([]model.AssigneeCount, error)
	ListUnderPaying(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID, amebaIDs []model.AmebaID) (model.AssigneeList, error)
	ListByOfferItemID(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID) (model.AssigneeList, error)
	GetByAmebaIDOfferItemID(ctx context.Context, exec boil.ContextExecutor, amebaID model.AmebaID, offerItemID model.OfferItemID) (*model.Assignee, error)
	ListByAmebaID(ctx context.Context, exec boil.ContextExecutor, amebaID model.AmebaID) (model.AssigneeList, error)
	Get(ctx context.Context, exec boil.ContextExecutor, assigneeID model.AssigneeID) (*model.Assignee, error)
	DeleteByOfferItemIDAndAmebaID(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID, amebaID model.AmebaID) error
}
