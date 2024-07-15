package rakuten

import (
	"net/url"
	"reflect"
	"strings"
)

const (
	// fieldWideRange 検索対象が広い
	fieldWideRange = "0"
	// availabilityAllItem すべての商品
	availabilityAllItem = "0"
)

func NewItemSearchByKeywordParam(applicationId, format string,
	shopCode, keyword string, category *string, page string, hits string,
) RakutenSearchParam {
	p := &ItemSearchParam{
		ApplicationID: applicationId,
		Format:        format,
		ShopCode:      shopCode,
		Keyword:       keyword,
		Page:          page,
		Hits:          hits,
		Field:         fieldWideRange,
		Availability:  availabilityAllItem,
	}

	if category != nil {
		p.GenreID = *category
	}

	return p
}

func NewItemSearchByItemCodeParam(applicationId, format, itemCode string) RakutenSearchParam {
	return &ItemSearchParam{
		ApplicationID: applicationId,
		//AffiliateID:   affiliateId,
		Format:       format,
		ItemCode:     itemCode,
		Availability: availabilityAllItem,
	}
}

// アイテム検索パラメータ
type ItemSearchParam struct {
	ApplicationID string `param:"applicationId"`
	//AffiliateID             string `param:"affiliateId"`
	Format                  string `param:"format"`
	Callback                string `param:"callback"`
	Elements                string `param:"elements"`
	FormatVersion           string `param:"formatVersion"`
	Keyword                 string `param:"keyword"`
	ShopCode                string `param:"shopCode"`
	ItemCode                string `param:"itemCode"`
	GenreID                 string `param:"genreId"`
	TagID                   string `param:"tagId"`
	Hits                    string `param:"hits"`
	Page                    string `param:"page"`
	Sort                    string `param:"sort"`
	MinPrice                string `param:"minPrice"`
	MaxPrice                string `param:"maxPrice"`
	Availability            string `param:"availability"`
	Field                   string `param:"field"`
	Carrier                 string `param:"carrier"`
	ImageFlag               string `param:"imageFlag"`
	OrFlag                  string `param:"orFlag"`
	NGKeyword               string `param:"NGKeyword"`
	PurchaseType            string `param:"purchaseType"`
	ShipOverseasFlag        string `param:"shipOverseasFlag"`
	ShipOverseasArea        string `param:"shipOverseasArea"`
	AsurakuFlag             string `param:"asurakuFlag"`
	AsurakuArea             string `param:"asurakuArea"`
	PointRateFlag           string `param:"pointRateFlag"`
	PointRate               string `param:"pointRate"`
	PostageFlag             string `param:"postageFlag"`
	CreditCardFlag          string `param:"creditCardFlag"`
	GiftFlag                string `param:"giftFlag"`
	HasReviewFlag           string `param:"hasReviewFlag"`
	MaxAffiliateRate        string `param:"maxAffiliateRate"`
	MinAffiliateRate        string `param:"minAffiliateRate"`
	HasMovieFlag            string `param:"hasMovieFlag"`
	PamphletFlag            string `param:"pamphletFlag"`
	AppointDeliveryDateFlag string `param:"appointDeliveryDateFlag"`
	GenreInformationFlag    string `param:"genreInformationFlag"`
	TagInformationFlag      string `param:"tagInformationFlag"`
}

func (i *ItemSearchParam) UrlValues() url.Values {
	return convertParamToUrlValues(i)
}

// ConvertParamToUrlValues はparamタグからurl.valuesを生成します。
func convertParamToUrlValues(rwsParam interface{}) url.Values {
	query := url.Values{}
	rv := reflect.ValueOf(rwsParam).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		name := f.Name
		key := f.Tag.Get("param")
		// String or StringSliceのみをパラメーターの対象とする
		if value, ok := rv.FieldByName(name).Interface().([]string); ok {
			if len(value) > 0 {
				query.Set(key, strings.Join(value, ","))
			}
		} else if value, ok := rv.FieldByName(name).Interface().(string); ok {
			if value != "" {
				query.Set(key, value)
			}
		}
	}
	return query
}
