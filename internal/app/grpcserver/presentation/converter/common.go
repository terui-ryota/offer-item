package converter

import (
	"fmt"

	//"github.com/ca-media-nantes/pick/protofiles/go/common"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/common"
	"github.com/terui-ryota/protofiles/go/offer_item"
)

func ListConditionPBToModel(pb *common.ListCondition) (*model.ListCondition, error) {
	sorts, err := func(ss []*common.Sort) ([]*model.Sort, error) {
		list := make([]*model.Sort, 0, len(ss))
		for _, s := range ss {
			sort, err := SortPBToModel(s)
			if err != nil {
				return nil, fmt.Errorf("SortPBToModel: %w", err)
			}
			list = append(list, sort)
		}
		return list, nil
	}(pb.GetSort())
	if err != nil {
		return nil, err
	}
	lc, err := model.NewListCondition(int(pb.GetOffset()), int(pb.GetLimit()), sorts)
	if err != nil {
		return nil, fmt.Errorf("model.NewListCondition: %w", err)
	}
	return lc, nil
}

func SortPBToModel(pb *common.Sort) (*model.Sort, error) {
	sort, err := model.NewSort(pb.GetOrderBy(), DescPBToModel(pb.GetOrdering()))
	if err != nil {
		return nil, fmt.Errorf("model.NewSort: %w", err)
	}
	return sort, nil
}

func DescPBToModel(pb common.Ordering) bool {
	return pb != common.Ordering_ASC
}

func ListResultModelToPB(m *model.ListResult) *common.ListResult {
	return &common.ListResult{
		Count:      uint32(m.Count()),
		TotalCount: uint32(m.TotalCount()),
	}
}

func CommissionModelToPB(m *model.Commission) *offer_item.Commission {
	return &offer_item.Commission{
		CommissionType: CommissionTypeModelToPB(m.CommissionType()),
		CalculatedRate: m.CalculatedRate(),
	}
}

func CommissionTypeModelToPB(m model.CommissionType) offer_item.ItemCommissionType {
	switch m {
	case model.CommissionTypeFixedRate:
		return offer_item.ItemCommissionType_COMMISSION_TYPE_FIXED_RATE
	case model.CommissionTypeFixedAmount:
		return offer_item.ItemCommissionType_COMMISSION_TYPE_FIXED_AMOUNT
	case model.CommissionTypeMultiFixedAmounts:
		return offer_item.ItemCommissionType_COMMISSION_TYPE_MULTI_FIXED_AMOUNT
	default:
		return offer_item.ItemCommissionType_COMMISSION_TYPE_UNKNOWN
	}
}

func PlatformURLsModelToPB(m []*model.PlatformURL) []*offer_item.URL {
	urlPBs := make([]*offer_item.URL, 0, len(m))
	for _, url := range m {
		urlPBs = append(urlPBs, PlatformURLModelToPB(url))
	}
	return urlPBs
}

func PlatformURLModelToPB(m *model.PlatformURL) *offer_item.URL {
	return &offer_item.URL{
		PlatformType: PlatformTypeModelToPB(m.PlatformType()),
		Url:          m.URL(),
	}
}

func PlatformTypeModelToPB(m model.PlatformType) offer_item.PlatformType {
	switch m {
	case model.PlatformTypeAll:
		return offer_item.PlatformType_ALL
	case model.PlatformTypeAndroid:
		return offer_item.PlatformType_ANDROID
	case model.PlatformTypeIOS:
		return offer_item.PlatformType_IOS
	default:
		return offer_item.PlatformType_PLATFORM_TYPE_UNKNOWN
	}
}
