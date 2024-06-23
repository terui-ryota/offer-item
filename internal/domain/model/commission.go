//go:generate go run github.com/terui-ryota/gen-getter -type=Commission

package model

import (
	"fmt"

	"github.com/terui-ryota/offer-item/pkg/apperr"
)

// 報酬
type Commission struct {
	// 報酬タイプ
	commissionType CommissionType
	// 計算済みの報酬料率
	calculatedRate float32
}

func NewCommission(
	commissionType CommissionType,
	calculatedRate float32,
) (*Commission, error) {
	newCommission := &Commission{
		commissionType: commissionType,
		calculatedRate: calculatedRate,
	}
	if !newCommission.IsValid() {
		return nil, apperr.OfferItemValidationError.Wrap(fmt.Errorf(fmt.Sprintf("invalid commission: %+v", newCommission)))
	}

	return newCommission, nil
}

func (c *Commission) IsValid() bool {
	// 報酬タイプが不明でない場合かつ計算済みの報酬料率が0より大きい場合
	if c.commissionType != CommissionTypeUnknown && c.calculatedRate > 0 {
		return true
		// 報酬タイプが不明でない場合かつ計算済みの報酬料率が0の場合
	} else if c.commissionType == CommissionTypeUnknown && c.calculatedRate == 0 {
		return true
	}
	return false
}

// 報酬タイプ
type CommissionType int

const (
	CommissionTypeUnknown           CommissionType = iota // 不明
	CommissionTypeFixedRate                               // 定率
	CommissionTypeFixedAmount                             // 定額
	CommissionTypeMultiFixedAmounts                       // 複数件VC
)

func (ct CommissionType) Int() int {
	return int(ct)
}

func ConvertCommissionType(commissionTypeInt int) CommissionType {
	switch commissionTypeInt {
	case int(CommissionTypeFixedRate):
		return CommissionTypeFixedRate
	case int(CommissionTypeFixedAmount):
		return CommissionTypeFixedAmount
	case int(CommissionTypeMultiFixedAmounts):
		return CommissionTypeMultiFixedAmounts
	default:
		return CommissionTypeUnknown
	}
}
