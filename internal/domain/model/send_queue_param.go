package model

import (
	"encoding/json"
	"fmt"
)

//go:generate go run github.com/terui-ryota/gen-getter -type=SendQueueParam
type SendQueueParam struct {
	offerItemID     string
	assigneeID      string
	adsTemplateCode string
}

func NewSendQueueParam(offerItemID, assigneeID, adsTemplateCode string) *SendQueueParam {
	return &SendQueueParam{
		offerItemID:     offerItemID,
		assigneeID:      assigneeID,
		adsTemplateCode: adsTemplateCode,
	}
}

type SendQueueParams []*SendQueueParam

func GetChunkSendQueueParams(assigneeList AssigneeList, offerItemID OfferItemID, adsTemplateCode string) []SendQueueParams {
	sendQueueParams := make(SendQueueParams, 0, len(assigneeList))
	for _, assignee := range assigneeList {
		param := NewSendQueueParam(offerItemID.String(), assignee.ID().String(), adsTemplateCode)
		sendQueueParams = append(sendQueueParams, param)
	}

	// PublishBatchで送信するために、sendQueueParamsを10個ずつに分割する
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sns#Client.PublishBatch
	chunkSize := 10
	var chunkSendQueueParams []SendQueueParams
	for i := 0; i < len(sendQueueParams); i += chunkSize {
		end := i + chunkSize
		// もし 'end'が'sendQueueParams'の長さを超えてしまった場合は、'end'を'sendQueueParams'の長さに設定する
		if end > len(sendQueueParams) {
			end = len(sendQueueParams)
		}
		chunkSendQueueParams = append(chunkSendQueueParams, sendQueueParams[i:end])
	}

	return chunkSendQueueParams
}

type Message interface {
	Messages() ([]string, error)
}

func (s SendQueueParams) Messages() ([]string, error) {
	messages := make([]string, 0, len(s))
	for _, v := range s {
		d, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal: %w", err)
		}
		messages = append(messages, string(d))
	}
	return messages, nil
}

// MarshalJSON エクスポートされていないフィールドをJSONに変換するために、MarshalJSONを実装する
func (s *SendQueueParam) MarshalJSON() ([]byte, error) {
	result, err := json.Marshal(struct {
		OfferItemID     string `json:"offer_item_id"`
		AssigneeID      string `json:"assignee_id"`
		ADSTemplateCode string `json:"ads_template_code"`
	}{
		OfferItemID:     s.offerItemID,
		AssigneeID:      s.assigneeID,
		ADSTemplateCode: s.adsTemplateCode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SendQueueParam: %w", err)
	}
	return result, nil
}
