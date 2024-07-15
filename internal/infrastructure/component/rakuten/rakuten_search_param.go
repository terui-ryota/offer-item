//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=mock/mock_$GOFILE -package=mock_$GOPACKAGE
package rakuten

import "net/url"

type RakutenSearchParam interface {
	UrlValues() url.Values
}
