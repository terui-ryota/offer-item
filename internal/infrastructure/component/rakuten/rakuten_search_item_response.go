package rakuten

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/terui-ryota/offer-item/internal/domain/model"
)

type AvailableFlag int

const (
	ENABLE AvailableFlag = iota
	AVAILABLE
)

type RakutenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// IsInValidKeyword はkeyword間違いか判定
func (re RakutenError) IsInValidItemCode() bool {
	return re.ErrorDescription == "itemCode is not valid"
}

// IsInValidKeyword はkeyword間違いか判定
func (re RakutenError) IsInValidKeyword() bool {
	return re.ErrorDescription == "keyword is not valid"
}

// IsOverLimitKeywords は検索件数が最大長を超えているか判定
func (re RakutenError) IsOverLimitKeywords() bool {
	return re.ErrorDescription == "keyword must be under 128 length"
}

// IsTooManyRequest はTooManyRequestで返ってきているか確認します。
func (re RakutenError) IsTooManyRequest() bool {
	return re.Error == "too_many_requests"
}

type ItemResult struct {
	Count     int            `json:"count"`
	Page      int            `json:"page"`
	First     int            `json:"first"`
	Last      int            `json:"last"`
	Hits      int            `json:"hits"`
	Carrier   int            `json:"carrier"`
	PageCount int            `json:"pageCount"`
	Items     []ItemResponse `json:"Items"`
	// GenreInformationとTagInformationは返却がないため空データとする
	// GenreInformation []interface{} `json:"GenreInformation"`
	// TagInformation   []interface{} `json:"TagInformation"`
}

type ItemResponse struct {
	Item `json:"Item"`
}

type Item struct {
	ItemName           string             `json:"itemName"`
	Catchcopy          string             `json:"catchcopy"`
	ItemCode           string             `json:"itemCode"`
	ItemPrice          int                `json:"itemPrice"`
	ItemCaption        string             `json:"itemCaption"`
	ItemURL            string             `json:"itemUrl"`
	AffiliateURL       string             `json:"affiliateUrl"`
	ImageFlag          int                `json:"imageFlag"`
	SmallImageUrls     []RakutenItemImage `json:"smallImageUrls"`
	MediumImageUrls    []RakutenItemImage `json:"mediumImageUrls"`
	Availability       int                `json:"availability"`
	TaxFlag            int                `json:"taxFlag"`
	PostageFlag        int                `json:"postageFlag"`
	CreditCardFlag     int                `json:"creditCardFlag"`
	ShopOfTheYearFlag  int                `json:"shopOfTheYearFlag"`
	ShipOverseasFlag   int                `json:"shipOverseasFlag"`
	ShipOverseasArea   string             `json:"shipOverseasArea"`
	AsurakuFlag        int                `json:"asurakuFlag"`
	AsurakuClosingTime string             `json:"asurakuClosingTime"`
	AsurakuArea        string             `json:"asurakuArea"`
	AffiliateRate      float64            `json:"affiliateRate"`
	StartTime          string             `json:"startTime"`
	EndTime            string             `json:"endTime"`
	ReviewCount        int                `json:"reviewCount"`
	ReviewAverage      float64            `json:"reviewAverage"`
	PointRate          int                `json:"pointRate"`
	PointRateStartTime string             `json:"pointRateStartTime"`
	PointRateEndTime   string             `json:"pointRateEndTime"`
	GiftFlag           int                `json:"giftFlag"`
	ShopName           string             `json:"shopName"`
	ShopCode           string             `json:"shopCode"`
	ShopURL            string             `json:"shopUrl"`
	ShopAffiliateURL   string             `json:"shopAffiliateUrl"`
	GenreID            string             `json:"genreId"`
	TagIds             []int              `json:"tagIds"`
}

type RakutenItemImage struct {
	ImageURL string `json:"imageUrl"`
}

// https://github.com/ca-media-nantes/pick/backend/affiliate-item/issues/207 対応
func (rii RakutenItemImage) UrlWithoutPath() string {
	// AffiliateURLをパース
	parsedUrl, err := url.Parse(rii.ImageURL)
	if err != nil {
		return rii.ImageURL
	}

	// Pathの先頭に `/` がなければ足す
	if !strings.HasPrefix(parsedUrl.Path, "/") {
		return fmt.Sprintf("%s://%s/%s", parsedUrl.Scheme, parsedUrl.Host, parsedUrl.Path)
	}
	return fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, parsedUrl.Path)
}

// ItemUrl は商品URLを返却します。
// 楽天の場合、ItemURLからはAffiliateURLが返却される仕様のため、商品URLを返却します。
// https://rakuten-webservice.tumblr.com/post/121801720357/%E6%A5%BD%E5%A4%A9%E5%B8%82%E5%A0%B4%E5%95%86%E5%93%81%E6%A4%9C%E7%B4%A2api%E6%A5%BD%E5%A4%A9%E5%B8%82%E5%A0%B4%E3%83%A9%E3%83%B3%E3%82%AD%E3%83%B3%E3%82%B0api%E3%81%AE%E4%BB%95%E6%A7%98%E5%A4%89%E6%9B%B4%E3%81%AB%E3%81%A4%E3%81%84%E3%81%A6
// Mobile/PC等の複数URLが返却されますが、MobileURLはPCから参照できないようなので、PCURLを選択します。
func (i Item) ParsedItemUrl() (string, error) {
	// AffiliateURLをパース
	parsedUrl, err := url.Parse(i.AffiliateURL)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to parse url. url:%s", i.AffiliateURL))
	}

	// PC URLを取得
	query := parsedUrl.Query()
	pcUrl := query.Get("pc")

	// AffiliateUrlのpcパラメタにない場合、ItemURLはitemUrlフィールドにあるためそちらを返す
	if pcUrl == "" {
		return i.ItemURL, nil
	}

	// QueryパラメータとしてEscapeされているため、Unescapeする。
	unescapeUrl, err := url.PathUnescape(pcUrl)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to unescape url. pcurl:%s", pcUrl))
	}

	return unescapeUrl, nil
}

func (i Item) UrlWithPlaceHolder() (string, error) {
	// AffiliateURLをパース
	parsedUrl, err := url.Parse(i.AffiliateURL)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Failed to parse url. url:%s", i.AffiliateURL))
	}

	// 任意パラメーターを置き換え
	urlPath := strings.Split(parsedUrl.Path, "/")
	if len(urlPath) < 2 {
		// ShopCodeが取得できない場合はエラーとして処理する
		return "", errors.Wrap(err, fmt.Sprintf("Failed to parase url path. url:%s", parsedUrl.Path))
	}

	// hgc を ichiba に変換する（ディープリンク対応）
	firstPath := urlPath[1] // 取得したPathは /hgc/{楽天アフィリエイトID}/ となっているため[1]を指定
	if firstPath == "hgc" {
		firstPath = "ichiba"
	}

	// プレースホルダーを指定したURLPathを設定する
	urlPaths := []string{firstPath, "${rakuten_afid}", "${partner_id}${click_id_prefix}_${click_id}"}
	pathString := strings.Join(urlPaths, "/")
	parsedUrl.RawPath = pathString

	// 楽天側でPick経由の流入を確認できるように指定されたパラメータをつける
	{
		q := parsedUrl.Query()
		q.Set("rafmid", "0001")

		parsedUrl.RawQuery = q.Encode()
	}

	// URLパスをエスケープしないためにurl.String()を利用しない
	return fmt.Sprintf("%s://%s/%s?%s", parsedUrl.Scheme, parsedUrl.Host, parsedUrl.RawPath, parsedUrl.RawQuery), nil
}

// 在庫情報加工
//func (i *ItemResponse) MakeStockInfo() model.StockStatus {
//	switch AvailableFlag(i.Item.Availability) {
//	case AVAILABLE:
//		return model.InStock
//	case ENABLE:
//		return model.OutOfStock
//	default:
//		return model.Unknown
//	}
//}

// MakeMaterial はMaterialを生成します。
func (i *ItemResponse) MakeMaterial() model.MaterialList {
	materials := model.MaterialList{}

	// description
	if i.Item.ItemCaption != "" {
		materials = append(materials, model.Material{MaterialType: model.Description, Value: i.Item.ItemCaption})
	}
	// TODO CatchCopy どうするのかフロントの方に合わせよう。
	if i.Item.Catchcopy != "" {
		materials = append(materials, model.Material{MaterialType: model.Description, Value: i.Item.Catchcopy})
	}

	// 商品画像がある場合のみ
	if i.Item.ImageFlag == 1 {
		// 大きい画像を優先的に利用する #154
		if len(i.Item.MediumImageUrls) > 0 {
			// MediumSize
			for _, url := range i.Item.MediumImageUrls {
				if url.ImageURL != "" {
					materials = append(materials, model.Material{MaterialType: model.Description, Value: url.UrlWithoutPath()})
				}
			}
		} else if len(i.Item.SmallImageUrls) > 0 {
			// SmallSize
			for _, url := range i.Item.SmallImageUrls {
				if url.ImageURL != "" {
					materials = append(materials, model.Material{MaterialType: model.Image, Value: url.UrlWithoutPath()})
				}
			}
		}
	}

	// 重複排除処理
	// materials.Distinct()

	return materials
}
