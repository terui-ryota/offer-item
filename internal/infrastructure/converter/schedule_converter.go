package converter

import (
	"time"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	null "github.com/volatiletech/null/v8"
)

func ScheduleEntityToModel(e *entity.Schedule) *model.Schedule {
	var startDate *time.Time
	var endDate *time.Time

	if e.StartDate.Valid {
		startDate = &e.StartDate.Time
	}
	if e.EndDate.Valid {
		endDate = &e.EndDate.Time
	}
	return model.NewScheduleFromRepository(
		model.ScheduleID(e.ID),
		model.OfferItemID(e.OfferItemID),
		model.ScheduleEntityToModel(int(e.ScheduleType)),
		startDate,
		endDate,
	)
}

func ScheduleModelToEntity(m *model.Schedule) *entity.Schedule {
	var startDate *time.Time
	var endDate *time.Time
	if m.StartDate() != nil {
		startDate = m.StartDate()
	}
	if m.EndDate() != nil {
		endDate = m.EndDate()
	}
	return &entity.Schedule{
		ID:           m.ID().String(),
		OfferItemID:  m.OfferItemID().String(),
		ScheduleType: uint(m.ScheduleType().Int()),
		StartDate:    null.TimeFromPtr(startDate),
		EndDate:      null.TimeFromPtr(endDate),
	}
}
