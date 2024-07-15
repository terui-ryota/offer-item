package converter

import (
	"github.com/volatiletech/null/v8"
)

func ToNullString(s string) null.String {
	if len(s) == 0 {
		return null.StringFromPtr(nil)
	}
	return null.StringFrom(s)
}
