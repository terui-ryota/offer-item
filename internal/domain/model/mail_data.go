package model

import "errors"

// メール送信データ
//
//go:generate go run github.com/terui-ryota/gen-getter -type=MailData
type MailData struct {
	// メールテンプレートコード
	adsTemplateCode string
	// 送信データ
	params *MailParams
}

func NewMailData(
	adsTemplateCode string,
	mailParams *MailParams,
) (*MailData, error) {
	if adsTemplateCode == "" {
		return nil, errors.New("template Code should not be empty")
	}

	return &MailData{
		adsTemplateCode: adsTemplateCode,
		params:          mailParams,
	}, nil
}

type MailParams struct {
	offerItem   *OfferItem
	assignee    *Assignee
	examination *Examination
}

// メール送信データ
//
//go:generate go run github.com/terui-ryota/gen-getter -type=MailParams
func NewMailParams(
	offerItem *OfferItem,
	assignee *Assignee,
	examination *Examination,
) *MailParams {
	return &MailParams{
		offerItem:   offerItem,
		assignee:    assignee,
		examination: examination,
	}
}
