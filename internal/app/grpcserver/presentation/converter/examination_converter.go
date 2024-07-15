package converter

import (
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
)

func ExaminationModelToPB(m *model.Examination) *offer_item.Examination {
	var optionalEntryID *offer_item.Examination_EntryId
	if m.EntryID() != nil {
		optionalEntryID = &offer_item.Examination_EntryId{
			EntryId: m.EntryID().String(),
		}
	}

	/*
		// TODO: SNS実装時に処理を追加する
		var optionalSNS *offer_item.Examination_Sns
		if m.Sns() != nil {
			optionalSNS = &offer_item.Examination_Sns{Sns: SnsModelToPB(m.Sns())}
		}
	*/

	var optionalReason *offer_item.Examination_Reason
	if m.Reason() != nil {
		optionalReason = &offer_item.Examination_Reason{
			Reason: *m.Reason(),
		}
	}

	return &offer_item.Examination{
		OfferItemId:     m.OfferItemID().String(),
		AmebaId:         m.AmebaID().String(),
		OptionalEntryId: optionalEntryID,
		// OptionalSns:     optionalSNS,
		OptionalReason: optionalReason,
		OptionalExaminerName: func() *offer_item.Examination_ExaminerName {
			if m.ExaminerName() == nil {
				return nil
			}
			return &offer_item.Examination_ExaminerName{
				ExaminerName: *m.ExaminerName(),
			}
		}(),
		EntrySubmissionCount: uint32(m.EntrySubmissionCount()),
	}
}

func SnsModelToPB(m *model.SNS) *offer_item.SNS {
	var userId *string
	if m.UserID() != nil {
		userId = m.UserID()
	}
	return &offer_item.SNS{
		OptionalUserId: &offer_item.SNS_UserId{
			UserId: *userId,
		},
		ScreenshotUrl: m.SnsScreenshotURL(),
	}
}

type EntryType int32
