package rakuten

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlWithPlaceHolder(t *testing.T) {
	t.Run("'hgc' should be replaced with 'ichiba'", func(t *testing.T) {
		item := Item{
			AffiliateURL: "https://hb.afl.rakuten.co.jp/hgc/dummy_id/?pc=https%3A%2F%2Fitem.rakuten.co.jp%2Fropepicnic%2Fr41443%2F&m=http%3A%2F%2Fm.rakuten.co.jp%2Fropepicnic%2Fi%2F10012667%2F",
		}

		expectedUrl := "https://hb.afl.rakuten.co.jp/ichiba/${rakuten_afid}/${partner_id}${click_id_prefix}_${click_id}?pc=https%3A%2F%2Fitem.rakuten.co.jp%2Fropepicnic%2Fr41443%2F&m=http%3A%2F%2Fm.rakuten.co.jp%2Fropepicnic%2Fi%2F10012667%2F&rafmid=0001"
		actualUrl, err := item.UrlWithPlaceHolder()
		result, parseErr := sameURLStrings(expectedUrl, actualUrl)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		assert.Equal(t, true, result)
		assert.Equal(t, nil, err)
	})

	t.Run("first segment of path should not be modified unless it's 'hgc'", func(t *testing.T) {
		item := Item{
			AffiliateURL: "https://hb.afl.rakuten.co.jp/ghc/dummy_id/?pc=https%3A%2F%2Fitem.rakuten.co.jp%2Fropepicnic%2Fr41443%2F&m=http%3A%2F%2Fm.rakuten.co.jp%2Fropepicnic%2Fi%2F10012667%2F",
		}

		expectedUrl := "https://hb.afl.rakuten.co.jp/ghc/${rakuten_afid}/${partner_id}${click_id_prefix}_${click_id}?pc=https%3A%2F%2Fitem.rakuten.co.jp%2Fropepicnic%2Fr41443%2F&m=http%3A%2F%2Fm.rakuten.co.jp%2Fropepicnic%2Fi%2F10012667%2F&rafmid=0001"
		actualUrl, err := item.UrlWithPlaceHolder()
		result, parseErr := sameURLStrings(expectedUrl, actualUrl)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		assert.Equal(t, true, result)
		assert.Equal(t, nil, err)
	})
}

func sameURLStrings(a, b string) (bool, error) {
	urlA, err := url.Parse(a)
	if err != nil {
		return false, err
	}
	urlB, err := url.Parse(b)
	if err != nil {
		return false, err
	}

	if urlA.Path != urlB.Path {
		return false, nil
	}

	for key, values := range urlA.Query() {
		valuesB, ok := urlB.Query()[key]
		if !ok {
			return false, nil
		}
		if len(values) != len(valuesB) {
			return false, nil
		}
		for i, v := range values {
			if v != valuesB[i] {
				return false, nil
			}
		}
	}

	return true, nil
}
