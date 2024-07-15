//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package adapter

import (
	"context"

	"github.com/terui-ryota/offer-item/internal/domain/model"
)

type AffiliateItemAdapter interface {
	GetItems(ctx context.Context, itemIdentifier model.ItemIdentifier) (*model.Items, error)
	BulkGetItems(ctx context.Context, itemIdentifiers model.ItemIdentifiers) (map[model.ItemIdentifier]model.Items, error)
}
