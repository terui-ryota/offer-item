package converter

import (
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
)

func ItemPBToModel(pb *offer_item.Item) (*model.Item, error) {
	minCommissionRate, err := CommissionPBToModel(pb.GetMinCommissionRate())
	if err != nil {
		return nil, fmt.Errorf("failed to create min commission rate. : %w", err)
	}
	maxCommissionRate, err := CommissionPBToModel(pb.GetMaxCommissionRate())
	if err != nil {
		return nil, fmt.Errorf("failed to create max commission rate: %w", err)
	}

	urls, err := func(us []*offer_item.URL) ([]*model.PlatformURL, error) {
		list := make([]*model.PlatformURL, 0, len(us))
		for _, u := range us {
			url, err := platformURLPBToModel(u)
			if err != nil {
				return nil, fmt.Errorf("platformURLPBToModel: %w", err)
			}
			list = append(list, url)
		}
		return list, nil
	}(pb.GetUrls())
	if err != nil {
		return nil, fmt.Errorf("failed to create urls: %w", err)
	}

	item, err := model.NewItem(
		model.ItemID(pb.GetId()),
		pb.GetImg(),
		pb.GetName(),
		minCommissionRate,
		maxCommissionRate,
		urls,
		pb.GetTieup(),
		pb.GetContentName(),
		pb.GetEnabledSelfBack(),
		pb.GetIsDf(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return item, nil
}

//func DFItemPBToModel(pb *offer_item.DFItem) (*model.DFItem, error) {
//	image := ""
//
//	// DFItemモデルは画像を1枚のみ保持します
//	for _, material := range pb.GetMaterials() {
//		if material.MaterialType == offer_item.MaterialType_MATERIAL_TYPE_IMAGE {
//			image = material.GetValue()
//			break
//		}
//	}
//
//	minCommissionRate, err := CommissionPBToModel(pb.GetMinCommissionRate())
//	if err != nil {
//		return nil, fmt.Errorf("CommissionPBToModel: %w", err)
//	}
//
//	maxCommissionRate, err := CommissionPBToModel(pb.GetMaxCommissionRate())
//	if err != nil {
//		return nil, fmt.Errorf("CommissionPBToModel: %w", err)
//	}
//
//	urls := make([]*model.PlatformURL, 0, len(pb.GetUrls()))
//	for _, u := range pb.GetUrls() {
//		url, err := platformURLPBToModel(u)
//		if err != nil {
//			return nil, fmt.Errorf("platformURLPBToModel: %w", err)
//		}
//
//		urls = append(urls, url)
//	}
//
//	dfItem, err := model.NewDFItem(
//		model.DFItemID(pb.GetId()),
//		image,
//		pb.GetName(),
//		minCommissionRate,
//		maxCommissionRate,
//		urls,
//	), nil
//	if err != nil {
//		return nil, fmt.Errorf("failed to create df item: %w", err)
//	}
//
//	return dfItem, nil
//}

func platformURLPBToModel(pb *offer_item.URL) (*model.PlatformURL, error) {
	platFormURL, err := model.NewPlatformURL(
		platformTypePBToModel(pb.GetPlatformType()),
		pb.GetUrl(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform url: %w", err)
	}
	return platFormURL, nil
}

func platformTypePBToModel(pb offer_item.PlatformType) model.PlatformType {
	switch pb {
	case offer_item.PlatformType_ALL:
		return model.PlatformTypeAll
	case offer_item.PlatformType_ANDROID:
		return model.PlatformTypeAndroid
	case offer_item.PlatformType_IOS:
		return model.PlatformTypeIOS
	default:
		return model.PlatformTypeUnknown
	}
}
