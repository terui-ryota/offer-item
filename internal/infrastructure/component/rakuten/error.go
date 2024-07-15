package rakuten

import (
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
)

func NewASPTooManyRequestError(itemId model.ItemID) ASPTooManyRequestError {
	return ASPTooManyRequestError{
		itemId: itemId,
	}
}

type ASPTooManyRequestError struct {
	itemId model.ItemID
}

func (e ASPTooManyRequestError) Error() string {
	return fmt.Sprintf("ASP API returned error. itemId:%s", e.itemId.String())
}
