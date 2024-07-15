//go:generate go run github.com/terui-ryota/gen-getter -type=Commission

package model

import (
	"fmt"

	"github.com/terui-ryota/offer-item/pkg/apperr"
)

func NewAffiliateItemCommission(rate, amount, staticRate *float32, description string) AffiliateItemCommission {
	c := AffiliateItemCommission{
		rate:        rate,
		amount:      amount,
		staticRate:  staticRate,
		description: description,
	}

	// CommissionType は報酬料率のタイプ（定率・定額）を設定します。
	c.setRateType()

	return c
}

// setRate は報酬料率タイプを設定します
func (c *AffiliateItemCommission) setRateType() {
	if c.isMultipleVC() {
		// Priceのみが設定されている（複数件VC）
		c.commissionType = CommissionTypeMultiFixedAmounts
	} else if c.hasAmount() {
		// POINT（定額）に値が入っている場合
		c.commissionType = CommissionTypeFixedAmount
	} else if c.hasRate() {
		// POINT(料率)に値が入っている場合
		c.commissionType = CommissionTypeFixedRate
	} else {
		// いずれも設定されてない場合は不明タイプで返却
		c.commissionType = CommissionTypeUnknown
	}
}

func (c AffiliateItemCommission) hasAmount() bool {
	return c.amount != nil && *c.amount != 0
}

func (c AffiliateItemCommission) hasRate() bool {
	return c.rate != nil && *c.rate != 0
}

func (c AffiliateItemCommission) isMultipleVC() bool {
	return c.description != "" && (c.rate == nil && c.amount == nil && c.amountPoint == nil && c.ratePoint == nil)
}

// CommissionList は報酬料率リストです。
type AffiliateItemCommissionList []AffiliateItemCommission

//go:generate go run github.com/terui-ryota/gen-getter -type=AffiliateItemCommission
type AffiliateItemCommission struct {
	// 料率
	rate *float32
	// 定額単価
	amount *float32
	// 料率（ポイント）
	ratePoint *float32
	// 定額単価（ポイント）
	amountPoint *float32
	// 価格
	price *float32
	// Amazon/楽天にて使用する。
	staticRate *float32
	// 概要
	description string
	// 報酬料率タイプ
	commissionType CommissionType
}

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
