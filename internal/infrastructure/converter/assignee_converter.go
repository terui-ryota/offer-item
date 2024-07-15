package converter

import (
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	null "github.com/volatiletech/null/v8"
)

func AssigneeEntityToModel(e *entity.Assignee) *model.Assignee {
	var declineReason *string
	if e.DeclineReason.Valid {
		declineReason = &e.DeclineReason.String
	}

	return model.NewAssigneeFromRepository(
		model.AssigneeID(e.ID),
		model.OfferItemID(e.OfferItemID),
		model.AmebaID(e.AmebaID),
		e.WritingFee,
		model.Stage(e.Stage),
		declineReason,
		e.CreatedAt,
	)
}

func AssigneeModelToEntity(m *model.Assignee) *entity.Assignee {
	return &entity.Assignee{
		ID:            m.ID().String(),
		OfferItemID:   string(m.OfferItemID()),
		AmebaID:       string(m.AmebaID()),
		WritingFee:    m.WritingFee(),
		Stage:         uint(m.Stage().Int()),
		DeclineReason: null.StringFromPtr(m.DeclineReason()),
	}
}
