package converter

import (
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
)

func CommissionPBToModel(pb *offer_item.Commission) (*model.Commission, error) {
	commission, err := model.NewCommission(
		CommissionTypePBToModel(pb.GetCommissionType()),
		pb.GetCalculatedRate(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create commission: %w", err)
	}
	return commission, nil
}

func CommissionTypePBToModel(pb offer_item.ItemCommissionType) model.CommissionType {
	switch pb {
	case offer_item.ItemCommissionType_COMMISSION_TYPE_FIXED_RATE:
		return model.CommissionTypeFixedRate
	case offer_item.ItemCommissionType_COMMISSION_TYPE_FIXED_AMOUNT:
		return model.CommissionTypeFixedAmount
	case offer_item.ItemCommissionType_COMMISSION_TYPE_MULTI_FIXED_AMOUNT:
		return model.CommissionTypeMultiFixedAmounts
	default:
		return model.CommissionTypeUnknown
	}
}
