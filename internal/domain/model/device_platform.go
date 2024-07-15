package model

import (
	"errors"
)

// プラットフォームタイプ
type PlatformType int

const (
	PlatformTypeUnknown PlatformType = iota // 不明
	PlatformTypeAll                         // 全て
	PlatformTypeAndroid                     // Android
	PlatformTypeIOS                         // iOS
)

func (pt PlatformType) Int() int {
	return int(pt)
}

func (pt PlatformType) EQ(other PlatformType) bool {
	return pt == other
}

func (pt PlatformType) IsValid() bool {
	return pt != PlatformTypeUnknown
}

// プラットフォーム別URL
type PlatformURL struct {
	// プラットフォームタイプ
	platformType PlatformType
	// URL
	url string
}

func NewPlatformURL(platformType PlatformType, url string) (*PlatformURL, error) {
	if !platformType.IsValid() {
		return nil, errors.New("Platform Type is not valid.")
	}
	if len(url) == 0 {
		return nil, errors.New("URL should not be empty.")
	}
	return &PlatformURL{
		platformType: platformType,
		url:          url,
	}, nil
}

func (p *PlatformURL) PlatformType() PlatformType {
	return p.platformType
}

func (p *PlatformURL) URL() string {
	return p.url
}
