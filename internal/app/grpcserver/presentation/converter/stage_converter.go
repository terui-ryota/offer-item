package converter

import (
	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
)

func StageDTOToModel(dtoStage dto.Stage) model.Stage {
	switch dtoStage {
	case dto.Stage_STAGE_UNKNOWN:
		return model.StageUnknown
	case dto.Stage_STAGE_BEFORE_INVITATION:
		return model.StageBeforeInvitation
	case dto.Stage_STAGE_INVITATION:
		return model.StageInvitation
	case dto.Stage_STAGE_LOTTERY:
		return model.StageLottery
	case dto.Stage_STAGE_LOTTERY_LOST:
		return model.StageLotteryLost
	case dto.Stage_STAGE_SHIPMENT:
		return model.StageShipment
	case dto.Stage_STAGE_DRAFT_SUBMISSION:
		return model.StageDraftSubmission
	case dto.Stage_STAGE_PRE_EXAMINATION:
		return model.StagePreExamination
	case dto.Stage_STAGE_PRE_REEXAMINATION:
		return model.StagePreReexamination
	case dto.Stage_STAGE_ARTICLE_POSTING:
		return model.StageArticlePosting
	case dto.Stage_STAGE_EXAMINATION:
		return model.StageExamination
	case dto.Stage_STAGE_REEXAMINATION:
		return model.StageReexamination
	case dto.Stage_STAGE_PAYING:
		return model.StagePaying
	case dto.Stage_STAGE_PAYMENT_COMPLETED:
		return model.StagePaymentCompleted
	case dto.Stage_STAGE_DONE:
		return model.StageDone
	default:
		return model.StageUnknown // 未知の値に対処するためにデフォルトを設定します。
	}
}

func StageModelToPB(stage model.Stage) offer_item.Stage {
	switch stage {
	case model.StageBeforeInvitation:
		return offer_item.Stage_STAGE_BEFORE_INVITATION
	case model.StageInvitation:
		return offer_item.Stage_STAGE_INVITATION
	case model.StageLottery:
		return offer_item.Stage_STAGE_LOTTERY
	case model.StageLotteryLost:
		return offer_item.Stage_STAGE_LOTTERY_LOST
	case model.StageShipment:
		return offer_item.Stage_STAGE_SHIPMENT
	case model.StageDraftSubmission:
		return offer_item.Stage_STAGE_DRAFT_SUBMISSION
	case model.StagePreExamination:
		return offer_item.Stage_STAGE_PRE_EXAMINATION
	case model.StagePreReexamination:
		return offer_item.Stage_STAGE_PRE_REEXAMINATION
	case model.StageArticlePosting:
		return offer_item.Stage_STAGE_ARTICLE_POSTING
	case model.StageExamination:
		return offer_item.Stage_STAGE_EXAMINATION
	case model.StageReexamination:
		return offer_item.Stage_STAGE_REEXAMINATION
	case model.StagePaying:
		return offer_item.Stage_STAGE_PAYING
	case model.StagePaymentCompleted:
		return offer_item.Stage_STAGE_PAYMENT_COMPLETED
	case model.StageDone:
		return offer_item.Stage_STAGE_DONE
	default:
		return offer_item.Stage_STAGE_UNKNOWN
	}
}
