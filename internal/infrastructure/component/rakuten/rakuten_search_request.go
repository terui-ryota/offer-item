package rakuten

import (
	"net/http"
	"net/url"
)

func NewRakutenRequest(url string, urlValue url.Values) *RakutenRequest {
	return &RakutenRequest{
		Url:      url,
		urlValue: urlValue,
	}
}

type RakutenRequest struct {
	Url      string
	urlValue url.Values
}

func (rr *RakutenRequest) GetRequest() (*http.Request, error) {
	return genGetRequest(rr.Url, rr.urlValue)
}

func genGetRequest(url string, query url.Values) (*http.Request, error) {
	urlstr := url + "?" + query.Encode()
	return http.NewRequest(http.MethodGet, urlstr, nil)
}
