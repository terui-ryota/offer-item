// RakutenITemとRakutenCommissionをつくる
package model

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

const (
	SpecialRateTrue = 1
)

// AffiliateItemDFItem はDF案件モデルです。
// Amazon/Rakuten/CAWiseの差異を吸収します。今回は楽天だけです。

//go:generate go run github.com/terui-ryota/gen-getter -type=AffiliateItemDFItem
type AffiliateItemDFItem struct {
	// 案件ID。案件を一意に識別するIDでPartnerASP毎に元となる値が異なる。
	// Amazon/ANSI, Rakuten/ShopCode+ItemId, CAWise/ProductId
	dfItemId AffiliateItemDFItemID
	// DFItemIDと同様の値が入っている。（元の値）
	// DFItemIDはPartnerASPの識別等ができるように一段クッションを入れて設定
	dfItemCode string
	// 案件名
	dfItemName string
	// 案件URL
	rawDFItemUrl string
	// トラッキングURL。PartnerASP毎に実績コンテキスト用のプレースホルダーが埋め込まれている。
	rawTrackingUrl string
	// 案件URL
	dfItemUrls URLMap
	// トラッキングURL。PartnerASP毎に実績コンテキスト用のプレースホルダーが埋め込まれている。
	trackingUrls URLMap
	// ブログ面に貼ることができる商品素材。
	material MaterialList
	// ブランド情報
	brand string
	// ショップ名
	shopName string
	// ショップコード、楽天専用項目
	shopCode string
	// 案件ID。DF案件の元となる案件情報。CAWiseの場合はItemIDを元に料率を取得する。
	itemId ItemID
	//// 卸売価格
	//RetailPrice Price
	//// 販売価格
	//SalePrice Price
	// 価格
	dfItemPriceList AffiliateItemDFItemPriceList
	// 報酬料率
	commission AffiliateItemCommissionList
	// 販売終了日
	endDate time.Time
	// 利用可能フラグ。CAWise以外は原則Active
	active bool
	// 更新日情報。gRPCサーバーでは利用せず、全文検索用のデータで利用する。
	updatedDate time.Time
	// 更新時に利用する
	hash string
	// LABOから新たにデータを更新されたものを取得するときに利用する。
	// 更新自体はどの案件も毎日更新を行うため、データ更新が行われた日付を保持する。
	modifiedDate *time.Time
	// アダルト商品かどうかのフラグ。tama側の広告表示制御で使用する。
	// レビュー情報（楽天のみ）
	review *Review
	// お気に入りフラグ
	isFavorited bool
	// 通常の料率に関する情報（2022年10月時点では楽天のみ使用。特別料率が付与されている場合のみ）
	generalRateInfo *GeneralRateInfo
	// 楽天のみ. 表示用のラベル
	salePriceLabel string
}

//go:generate go run github.com/terui-ryota/gen-getter -type=AffiliateItemPair
type AffiliateItemPair struct {
	affiliateItemItem   *AffiliateItemItem
	affiliateItemDFItem *AffiliateItemDFItem
}

func NewAffiliateItemPair(
	affiliateItemItem *AffiliateItemItem,
	affiliateItemDFItem *AffiliateItemDFItem,
) *AffiliateItemPair {
	return &AffiliateItemPair{
		affiliateItemItem:   affiliateItemItem,
		affiliateItemDFItem: affiliateItemDFItem,
	}
}

type AffiliateItemPairList []*AffiliateItemPair

//item, err := model.NewItem(
//model.ItemID(pb.GetId()),
//pb.GetImg(),
//pb.GetName(),
//minCommissionRate,
//maxCommissionRate,
//urls,
//pb.GetTieup(),
//pb.GetContentName(),
//pb.GetEnabledSelfBack(),
//pb.GetIsDf(),
//)
//
//return &Item{
//id:                id,
//img:               img,
//name:              name,
//minCommissionRate: minCommissionRate,
//maxCommissionRate: maxCommissionRate,
//urls:              urls,
//hasTieup:          hasTieup,
//contentName:       contentName,
//enabledSelfBack:   enabledSelfBack,
//isDF:              isDF,
//}, nil
//

func AffiliateItemItemToItem(affiliateiItemItem *AffiliateItemItem) (*Item, error) {
	itemCommission := affiliateiItemItem.Commission()
	rate := itemCommission.Rate()
	commission, err := NewCommission(
		CommissionTypeFixedRate,
		*rate,
	)

	itemURLS := affiliateiItemItem.ItemUrls()
	urlList := make([]*PlatformURL, 0, len(itemURLS))
	for _, u := range itemURLS {
		url, err := NewPlatformURL(
			PlatformTypeAll,
			u,
		)
		if err != nil {
			return nil, fmt.Errorf("platformURLPBToModel: %w", err)
		}
		urlList = append(urlList, url)
	}

	item, err := NewItem(
		affiliateiItemItem.ItemId(),
		affiliateiItemItem.ItemImage(),
		affiliateiItemItem.ItemName(),
		commission,
		commission,
		urlList,
		false,
		affiliateiItemItem.ContentName(),
		false,
		affiliateiItemItem.IsDF(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}
	return item, nil
}

func AffiliateItemDFItemToDFItem(affiliateItemDFItem *AffiliateItemDFItem) (*DFItem, error) {
	dfItemCommission := affiliateItemDFItem.Commission()

	var commission *Commission
	var err error
	for _, v := range dfItemCommission {
		rate := v.Rate()
		commission, err = NewCommission(
			CommissionTypeFixedRate,
			*rate,
		)
		if err != nil {
			return nil, fmt.Errorf("platformURLPBToModel: %w", err)
		}
	}

	urls := affiliateItemDFItem.DfItemUrls()
	urlList := make([]*PlatformURL, 0, len(urls))
	for _, u := range urls {
		url, err := NewPlatformURL(
			PlatformTypeAll,
			u,
		)
		if err != nil {
			return nil, fmt.Errorf("platformURLPBToModel: %w", err)
		}
		urlList = append(urlList, url)
	}

	image := ""
	// DFItemモデルは画像を1枚のみ保持します
	for _, material := range affiliateItemDFItem.Material() {
		if material.MaterialType == MaterialType_MATERIAL_TYPE_IMAGE {
			image = material.GetValue()
			break
		}
	}

	dfItem := NewDFItem(
		DFItemID(affiliateItemDFItem.DfItemId()),
		image,
		affiliateItemDFItem.DfItemName(),
		commission,
		commission,
		urlList,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}
	return dfItem, nil

}

//go:generate go run github.com/terui-ryota/gen-getter -type=AffiliateItemItem
type AffiliateItemItem struct {
	// 案件ID
	itemId ItemID
	// 案件画像
	itemImage string
	// 案件名
	itemName string
	// 案件オーナーID
	itemOwnerId string
	// 素材URL
	largeMaterialImage   string
	mediumMaterialImage1 string
	smallMaterialImage1  string
	mediumMaterialImage2 string
	mediumMaterialImage3 string
	smallMaterialImage2  string
	// 価格
	price []Price
	// 報酬料率
	commission AffiliateItemCommission
	// 案件URL（LPIDがPlaceholderの状態）
	rawItemUrl string
	// トラッキングURL（LPIDがPlaceholderの状態）
	rawTrackingUrl string
	// 案件URLs プラットフォーム毎のURL
	itemUrls URLMap
	// トラッキングURLs プラットフォーム毎のURL
	trackingUrls URLMap
	// 案件概要
	itemDescription string
	// 否認条件
	denialCondition string
	// リスティングNGキーワード
	listingNgKeyWords string
	// 提携完了までの待機日数
	approveWatiTime string
	// 広告名
	contentName string
	// DF案件フラグ
	isDF bool
}

func NewAffiliateItemItem(
	itemId ItemID,
	itemImage,
	itemName string,
	commission AffiliateItemCommission,
	contentName string,
	itemUrls URLMap,
	isDF bool,
) *AffiliateItemItem {
	return &AffiliateItemItem{
		itemId:      itemId,
		itemImage:   itemImage,
		itemName:    itemName,
		commission:  commission,
		itemUrls:    itemUrls,
		contentName: contentName,
		isDF:        isDF,
	}
}

type AffiliateItemDFItemID string

type AffiliateItemDFItemList []*AffiliateItemDFItem

// CommissionList は報酬料率リストです。
type CommissionList []Commission

//// AddStaticRate は固定料率をもった報酬料率を追加します。
//func (rrl *CommissionList) AddStaticRate(staticRate *float32) {
//	rrl.Append(Commission{
//		StaticRate: staticRate,
//	})
//}
//
//func (rrl *CommissionList) Append(rate Commission) {
//	*rrl = append(*rrl, rate)
//}

type URLMap map[PlatformType]string

func (u URLMap) Has(target string) bool {
	for _, url := range u {
		if url == target {
			return true
		}
	}
	return false
}

func (u URLMap) Get(platformType PlatformType) (string, bool) {
	v, ok := u[platformType]
	return v, ok
}

// 通常の料率に関する情報
//
//go:generate go run github.com/terui-ryota/gen-getter -type=GeneralRateInfo
type GeneralRateInfo struct {
	// 0が通常の料率、1が特別料率
	rateType int
	// 通常の料率
	rate float32
}

// NewDFItem はDF案件モデルを生成します。
func NewAffiliateItemDFItem(
	dfItemId,
	dfItemName,
	rawDfItemUrl string,
	material MaterialList,
	shopName,
	shopCode string,
	itemId ItemID,
	dfItemPriceList AffiliateItemDFItemPriceList,
	rate AffiliateItemCommissionList,
) *AffiliateItemDFItem {
	itemURL := rawDfItemUrl
	i := &AffiliateItemDFItem{
		dfItemId:        AffiliateItemDFItemID(dfItemId),
		dfItemName:      dfItemName,
		rawDFItemUrl:    itemURL,
		material:        material,
		shopName:        shopName,
		shopCode:        shopCode,
		itemId:          itemId,
		dfItemPriceList: dfItemPriceList,
		commission:      rate,
	}
	return i
}

type AffiliateItemDFItemPrice struct {
	// 卸売価格
	RetailPrice Price
	// 販売価格
	SalePrice Price
}

// 価格リスト
type AffiliateItemDFItemPriceList []AffiliateItemDFItemPrice

// 最低価格を取得
func (dpl AffiliateItemDFItemPriceList) MinPrice() AffiliateItemDFItemPrice {
	if len(dpl) == 0 {
		return AffiliateItemDFItemPrice{}
	}

	result := AffiliateItemDFItemPrice{
		SalePrice:   dpl[0].SalePrice,
		RetailPrice: dpl[0].RetailPrice,
	}
	for _, v := range dpl {
		// 販売価格で判断する
		if v.SalePrice < result.SalePrice {
			result = v
		}
	}
	return result
}

// 最高価格を取得
func (dpl AffiliateItemDFItemPriceList) MaxPrice() AffiliateItemDFItemPrice {
	if len(dpl) == 0 {
		return AffiliateItemDFItemPrice{}
	}

	result := AffiliateItemDFItemPrice{
		SalePrice:   dpl[0].SalePrice,
		RetailPrice: dpl[0].RetailPrice,
	}
	for _, v := range dpl {
		// 販売価格で判断する
		if v.SalePrice > result.SalePrice {
			result = v
		}
	}
	return result
}

// Add はDFItemListに案件を追加します。
func (dl *AffiliateItemDFItemList) Append(item *AffiliateItemDFItem) {
	list := append(*dl, item)
	*dl = list
}

//// Slice はDFItemListを指定サイズで分割します。
//func (dl AffiliateItemDFItemList) Slice(chunkSize int) []AffiliateItemDFItemList {
//	result := make([]AffiliateItemDFItemList, 0)
//	for i, end := range ifhelper.GenIndexSet(len(dl), chunkSize) {
//		result = append(result, dl[i:end])
//	}
//	return result
//}
//
//// Exist は指定IDの案件が存在するか確認します。
//func (dl AffiliateItemDFItemList) Exist(dfItemId string) bool {
//	for _, v := range dl {
//		if v.ItemId.Value() == dfItemId {
//			return true
//		}
//	}
//	return false
//}

//// GetHashes はDFItemListのID一覧を取得します。
//func (dl AffiliateItemDFItemList) GetIds() []string {
//	result := make([]string, 0)
//	for _, v := range dl {
//		result = append(result, v.DFItemId.Value())
//	}
//	return result
//}

// DFItemId はDF案件ID型です。
type AffiliateItemDFItemId string

// Value は文字列を取得します。
func (i AffiliateItemDFItemId) Value() string {
	return string(i)
}

const TAX_RATE = 0.1

// Price は価格です。
type Price int64

// WithoutTax は税抜き価格を表示します。
func (p Price) WithoutTax() int64 {
	x := big.NewFloat(float64(p.Value()))
	y := big.NewFloat(float64(1 + TAX_RATE))

	fp, _ := new(big.Float).SetPrec(1024).Quo(x, y).Float64()

	return int64(math.Ceil(fp))
}

// Value はIntを取得します。
func (p Price) Value() int64 {
	return int64(p)
}

// StockStatus は在庫情報です。
// 固定値はproto.enumsに定義します。
type StockStatus int

// Value はIntを取得します。
func (s StockStatus) Value() int {
	return int(s)
}

// ItemMaterials は案件毎の素材情報を保持します。
type ItemMaterials map[AffiliateItemDFItemId]MaterialList

// Add は素材を追加します。
func (im ItemMaterials) Add(id string, materialType int, value string) {
	itemId := AffiliateItemDFItemId(id)
	if im[itemId] == nil {
		im[itemId] = MaterialList{}
	}
	im[itemId] = append(im[itemId], Material{
		MaterialType: MaterialType(materialType),
		Value:        value,
	})
}

// Get は素材を取得します。
func (imat ItemMaterials) Get(id string) MaterialList {
	return imat[AffiliateItemDFItemId(id)]
}

// Material は素材です。
//
//go:generate go run github.com/terui-ryota/gen-getter -type=AffiliateItemMaterial
type AffiliateItemMaterial struct {
	materialType MaterialType
	value        string
}

// Key は素材のタイプと値を合わせたキーを取得します。
func (m AffiliateItemMaterial) Key() string {
	return fmt.Sprintf("%d#%s", m.materialType, m.value)
}

// AffiliateItemMaterialType はMaterialType型です。
type AffiliateItemMaterialType int

// Value はintを返します。
func (m AffiliateItemMaterialType) Value() int {
	return int(m)
}

type AffiliateItemMaterialList []AffiliateItemMaterial

// GetValues はMaterialListの値一覧を取得します。
func (ml AffiliateItemMaterialList) GetValues() []string {
	values := make([]string, 0)
	for _, v := range ml {
		values = append(values, v.value)
	}
	return values
}

//// GetValuesForHash はハッシュ生成用の値を取得します。
//func (ml AffiliateItemMaterialList) GetValuesForHash() []string {
//	values := make([]string, 0)
//	for _, v := range ml {
//		values = append(values, "material:"+strconv.Itoa(v.type.Value())+"-"+v.Value)
//	}
//	return values
//}

// Distinct は素材の重複を排除します。
//
//	func (ml *AffiliateItemMaterialList) Distinct() {
//		maps := make(map[string]AffiliateItemMaterial, len(*ml))
//		for _, v := range *ml {
//			maps[v.Key()] = v
//		}
//
//		list := MaterialList{}
//		for _, v := range maps {
//			list = append(list, v)
//		}
//		*ml = list
//	}
type Review struct {
	// レート
	Average float64
	// レビュー数
	Count int
}
