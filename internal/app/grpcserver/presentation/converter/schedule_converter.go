package converter

import (
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ScheduleModelToPB(m *model.Schedule) *offer_item.Schedule {
	var optionalStartDate *offer_item.Schedule_StartDate
	var optionalEndDate *offer_item.Schedule_EndDate

	if m.StartDate() != nil {
		optionalStartDate = &offer_item.Schedule_StartDate{
			StartDate: timestamppb.New(*m.StartDate()),
		}
	}

	if m.EndDate() != nil {
		optionalEndDate = &offer_item.Schedule_EndDate{
			EndDate: timestamppb.New(*m.EndDate()),
		}
	}

	return &offer_item.Schedule{
		Id:                string(m.ID()),
		OfferItemId:       string(m.OfferItemID()),
		ScheduleType:      ScheduleTypeModelToPB(m.ScheduleType()),
		OptionalStartDate: optionalStartDate,
		OptionalEndDate:   optionalEndDate,
	}
}

func ScheduleTypeModelToPB(m model.ScheduleType) offer_item.ScheduleType {
	switch m {
	case model.ScheduleTypeInvitation:
		return offer_item.ScheduleType_SCHEDULE_TYPE_INVITATION
	case model.ScheduleTypeLottery:
		return offer_item.ScheduleType_SCHEDULE_TYPE_LOTTERY
	case model.ScheduleTypeShipment:
		return offer_item.ScheduleType_SCHEDULE_TYPE_SHIPMENT
	case model.ScheduleTypeDraftSubmission:
		return offer_item.ScheduleType_SCHEDULE_TYPE_DRAFT_SUBMISSION
	case model.ScheduleTypePreExamination:
		return offer_item.ScheduleType_SCHEDULE_TYPE_PRE_EXAMINATION
	case model.ScheduleTypeArticlePosting:
		return offer_item.ScheduleType_SCHEDULE_TYPE_ARTICLE_POSTING
	case model.ScheduleTypeExamination:
		return offer_item.ScheduleType_SCHEDULE_TYPE_EXAMINATION
	case model.ScheduleTypePayment:
		return offer_item.ScheduleType_SCHEDULE_TYPE_PAYMENT
	default:
		return offer_item.ScheduleType_SCHEDULE_TYPE_UNKNOWN
	}
}

func ScheduleTypeDTOToModel(v dto.ScheduleType) model.ScheduleType {
	switch v {
	case dto.ScheduleType_SCHEDULE_TYPE_INVITATION:
		return model.ScheduleTypeInvitation
	case dto.ScheduleType_SCHEDULE_TYPE_LOTTERY:
		return model.ScheduleTypeLottery
	case dto.ScheduleType_SCHEDULE_TYPE_SHIPMENT:
		return model.ScheduleTypeShipment
	case dto.ScheduleType_SCHEDULE_TYPE_DRAFT_SUBMISSION:
		return model.ScheduleTypeDraftSubmission
	case dto.ScheduleType_SCHEDULE_TYPE_PRE_EXAMINATION:
		return model.ScheduleTypePreExamination
	case dto.ScheduleType_SCHEDULE_TYPE_ARTICLE_POSTING:
		return model.ScheduleTypeArticlePosting
	case dto.ScheduleType_SCHEDULE_TYPE_EXAMINATION:
		return model.ScheduleTypeExamination
	case dto.ScheduleType_SCHEDULE_TYPE_PAYMENT:
		return model.ScheduleTypePayment
	default:
		return model.ScheduleTypeUnknown
	}
}
