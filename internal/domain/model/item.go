package model

import (
	"errors"
	"fmt"
	"strconv"
)

// 案件
// データの管理元は affiliate-item コンテキスト
//
//go:generate go run github.com/terui-ryota/gen-getter -type=Item
type Item struct {
	// 案件ID
	id ItemID
	// 案件画像
	img string
	// 案件名
	name string
	// 最低報酬料率
	minCommissionRate *Commission
	// 最高報酬料率
	maxCommissionRate *Commission
	// 商品URLリスト
	urls []*PlatformURL
	// 提携有無
	hasTieup bool
	// 広告主
	contentName string
	// セルフバック可能か
	enabledSelfBack bool
	// DF を持っているか
	isDF bool
}

func NewItem(
	id ItemID,
	img,
	name string,
	minCommissionRate,
	maxCommissionRate *Commission,
	urls []*PlatformURL,
	hasTieup bool,
	contentName string,
	enabledSelfBack bool,
	isDF bool,
) (*Item, error) {
	if len(id) == 0 {
		return nil, errors.New("ID should not be empty.")
	}
	if len(name) == 0 {
		return nil, errors.New("Name should not be empty.")
	}

	return &Item{
		id:                id,
		img:               img,
		name:              name,
		minCommissionRate: minCommissionRate,
		maxCommissionRate: maxCommissionRate,
		urls:              urls,
		hasTieup:          hasTieup,
		contentName:       contentName,
		enabledSelfBack:   enabledSelfBack,
		isDF:              isDF,
	}, nil
}

// ItemInfo ItemInfoは楽天などの商品情報が消されても管理面への影響を与えないために、バックエンドのDBにキャッシュするために使用します
//go:generate go run github.com/terui-ryota/gen-getter -type=ItemInfo

type ItemInfo struct {
	// オファーアイテム
	offerItemID OfferItemID
	// 商品名
	name string
	// 会社名
	contentName string
	// 商品画像URL
	imageURL string
	// 詳細を見るの遷移URL
	url string
	// 最低報酬料率
	minCommission *Commission
	// 最高報酬単価
	maxCommission *Commission
}

func NewItemInfo(
	offerItemID OfferItemID,
	name,
	contentName,
	imageURL,
	url string,
	minCommission,
	maxCommission *Commission,
) (*ItemInfo, error) {
	if err := validateItemInfo(name, contentName, imageURL, url, minCommission, maxCommission); err != nil {
		return nil, fmt.Errorf("validateItemInfo: %w", err)
	}
	return &ItemInfo{
		offerItemID:   offerItemID,
		name:          name,
		contentName:   contentName,
		imageURL:      imageURL,
		url:           url,
		minCommission: minCommission,
		maxCommission: maxCommission,
	}, nil
}

func validateItemInfo(
	name,
	contentName,
	imageURL,
	url string,
	minCommission,
	maxCommission *Commission,
) error {
	if name == "" {
		return errors.New("name is required")
	}
	if contentName == "" {
		return errors.New("contentName is required")
	}
	if imageURL == "" {
		return errors.New("imageURL is required")
	}
	if url == "" {
		return errors.New("url is required")
	}
	if minCommission != nil && !minCommission.IsValid() {
		return errors.New("minCommission is invalid")
	}
	// minCommission.CalculatedRateがmaxCommission.CalculatedRateより大きい場合はエラー
	if minCommission.CalculatedRate() > maxCommission.CalculatedRate() {
		return errors.New("minCommission is greater than maxCommission")
	}
	if maxCommission != nil && !maxCommission.IsValid() {
		return errors.New("maxCommission is invalid")
	}
	// maxCommission.CalculatedRateがminCommission.CalculatedRateより小さい場合はエラー
	if maxCommission.CalculatedRate() < minCommission.CalculatedRate() {
		return errors.New("maxCommission is less than minCommission")
	}

	return nil
}

func NewDraftedItemInfoFromRepository(
	offerItemID OfferItemID,
	name,
	contentName,
	imageURL,
	url string,
	minCommission,
	maxCommission *Commission,
) *ItemInfo {
	return &ItemInfo{
		offerItemID:   offerItemID,
		name:          name,
		contentName:   contentName,
		imageURL:      imageURL,
		url:           url,
		minCommission: minCommission,
		maxCommission: maxCommission,
	}
}

func NewItemByItemID(id ItemID) (*Item, error) {
	if len(id) == 0 {
		return nil, errors.New("ID should not be empty.")
	}

	return &Item{id: id}, nil
}

func (i *Item) IsDf() bool {
	return i.isDF
}

// 案件が存在するかをチェックする
func (i *Item) Exists() bool {
	if i != nil {
		// IDと案件名でチェック（IDと案件名は必ず設定されるため）
		if len(i.id) > 0 && len(i.name) > 0 {
			return true
		}
	}
	return false
}

// 案件マップ
type ItemMap map[ItemID]*Item

// 案件ID
type ItemID string

func (ii ItemID) String() string {
	return string(ii)
}

// 案件IDリスト
type ItemIDList []ItemID

func NewItemIDList(strList []string) ItemIDList {
	result := make(ItemIDList, 0, len(strList))
	for _, str := range strList {
		result = append(result, ItemID(str))
	}
	return result
}

func (il ItemIDList) String() []string {
	result := make([]string, 0, len(il))
	for _, i := range il {
		result = append(result, i.String())
	}
	return result
}

// DF案件ID
type DFItemID string

func (id DFItemID) String() string {
	return string(id)
}

// DF案件
// データの管理元は affiliate-item コンテキスト
//
//go:generate go run github.com/terui-ryota/gen-getter -type=DFItem
type DFItem struct {
	// DF案件ID
	id DFItemID
	// 案件画像
	img string
	// 案件名
	name string
	// 最低報酬料率
	minCommissionRate *Commission
	// 最高報酬料率
	maxCommissionRate *Commission
	// 商品URLリスト
	urls []*PlatformURL
}

func (d *DFItem) SetDFItemIDEmpty() {
	d.id = ""
}

func NewDFItemByID(id DFItemID) *DFItem {
	if id == "" {
		return nil
	}

	return &DFItem{id: id}
}

func NewDFItem(
	id DFItemID,
	img string,
	name string,
	minCommissionRate *Commission,
	maxCommissionRate *Commission,
	urls []*PlatformURL,
) *DFItem {
	return &DFItem{
		id:                id,
		img:               img,
		name:              name,
		minCommissionRate: minCommissionRate,
		maxCommissionRate: maxCommissionRate,
		urls:              urls,
	}
}

func (i *DFItem) Exists() bool {
	if i != nil {
		if i.id != "" && i.name != "" {
			return true
		}
	}

	return false
}

// BannerID は バナー ID を表現します
type BannerID string

func NewBannerID(v string) (*BannerID, error) {
	if v == "" {
		return nil, errors.New("bannerID must not be empty")
	}
	i := BannerID(v)
	return &i, nil
}

func (id BannerID) String() string {
	return string(id)
}

type Items struct {
	Item   Item
	DFItem DFItem
}

type ItemIdentifier struct {
	itemID   ItemID
	dfItemID DFItemID
}

type ItemIdentifiers []ItemIdentifier

func NewItemIdentifier(
	itemID ItemID,
	dfItemID DFItemID,
) *ItemIdentifier {
	return &ItemIdentifier{
		itemID:   itemID,
		dfItemID: dfItemID,
	}
}

func (i *ItemIdentifier) ItemID() ItemID {
	return i.itemID
}

func (i *ItemIdentifier) DFItemID() DFItemID {
	return i.dfItemID
}

// Material は素材です。
type Material struct {
	MaterialType MaterialType
	Value        string
}

func (x *Material) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// Key は素材のタイプと値を合わせたキーを取得します。
func (m Material) Key() string {
	return fmt.Sprintf("%d#%s", m.MaterialType, m.Value)
}

// MaterialType はMaterialType型です。
type MaterialType int

const (
	MaterialType_MATERIAL_TYPE_UNKNOWN     MaterialType = 0
	MaterialType_MATERIAL_TYPE_DESCRIPTION MaterialType = 1
	MaterialType_MATERIAL_TYPE_IMAGE       MaterialType = 2
)

const (
	Description = MaterialType(MaterialType_MATERIAL_TYPE_DESCRIPTION)
	Image       = MaterialType(MaterialType_MATERIAL_TYPE_IMAGE)
)

// Value はintを返します。
func (m MaterialType) Value() int {
	return int(m)
}

type MaterialList []Material

// GetValues はMaterialListの値一覧を取得します。
func (ml MaterialList) GetValues() []string {
	values := make([]string, 0)
	for _, v := range ml {
		values = append(values, v.Value)
	}
	return values
}

// GetValuesForHash はハッシュ生成用の値を取得します。
func (ml MaterialList) GetValuesForHash() []string {
	values := make([]string, 0)
	for _, v := range ml {
		values = append(values, "material:"+strconv.Itoa(v.MaterialType.Value())+"-"+v.Value)
	}
	return values
}

// Distinct は素材の重複を排除します。
func (ml *MaterialList) Distinct() {
	maps := make(map[string]Material, len(*ml))
	for _, v := range *ml {
		maps[v.Key()] = v
	}

	list := MaterialList{}
	for _, v := range maps {
		list = append(list, v)
	}
	*ml = list
}
