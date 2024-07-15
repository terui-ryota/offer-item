//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package repository

import (
	"context"
	"database/sql"

	"github.com/terui-ryota/offer-item/internal/domain/model"
)

type ExaminationRepository interface {
	BulkGetByOfferItemID(ctx context.Context, db *sql.DB, offerItemID model.OfferItemID, entryType model.EntryType) (map[model.AmebaID]*model.Examination, error)
	Get(ctx context.Context, db *sql.DB, offerItemID model.OfferItemID, assigneeID model.AssigneeID, entryType model.EntryType) (*model.Examination, error)
	Update(ctx context.Context, exec *sql.DB, examination *model.Examination) error
	Create(ctx context.Context, db *sql.DB, examination *model.Examination) error
}
