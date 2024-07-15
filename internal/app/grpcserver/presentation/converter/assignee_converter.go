package converter

import (
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AssigneeModelToPB(m *model.Assignee) *offer_item.Assignee {
	return &offer_item.Assignee{
		Id:          m.ID().String(),
		OfferItemId: m.OfferItemID().String(),
		AmebaId:     m.AmebaID().String(),
		WritingFee:  int64(m.WritingFee()),
		Stage:       StageModelToPB(m.Stage()),
		CreatedAt:   timestamppb.New(m.CreatedAt()),
	}
}

func AssigneeCountModelToPB(m model.AssigneeCount) *offer_item.AssigneeCount {
	return &offer_item.AssigneeCount{
		Stage: StageModelToPB(m.Stage()),
		Count: int64(m.Count()),
	}
}
