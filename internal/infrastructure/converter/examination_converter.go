package converter

import (
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	null "github.com/volatiletech/null/v8"
)

func ExaminationEntityToModel(e *entity.Examination, count int) *model.Examination {
	var entryID *model.EntryID
	if e.EntryID.Valid {
		tmpEntryID := model.EntryID(e.EntryID.String)
		entryID = &tmpEntryID
	}

	/*
		// TODO: SNS実装時に処理を追加する
		 var snsUserID *string
		 if e.SNSUserID.Valid {
			snsUserID = &e.SNSUserID.String
		}
	*/

	/*
		var snsScreenshotURL *string
		if e.SNSScreenshotURL.Valid {
			snsScreenshotURL = &e.SNSScreenshotURL.String
		}
	*/

	return model.NewExaminationFromRepository(
		model.ExaminationID(e.ID),
		model.OfferItemID(e.OfferItemID),
		model.AmebaID(e.R.Assignee.AmebaID),
		entryID,
		/*
			snsUserID,
			snsScreenshotURL,
		*/
		e.ExaminerName.Ptr(),
		e.Reason.Ptr(),
		model.AssigneeID(e.AssigneeID),
		model.EntryType(e.EntryType),
		uint(count),
	)
}

func ExaminationModelToEntity(examination *model.Examination) entity.Examination {
	var entryID *string
	if examination.EntryID() != nil {
		tmpEntryID := examination.EntryID().String()
		entryID = &tmpEntryID
	}

	/*
		// TODO: SNS実装時に処理を追加する
		 var snsScreenshotURL *string
		 if examination.Sns() != nil {
			tmpScreenshotURL := examination.Sns().SnsScreenshotURL()
			snsScreenshotURL = &tmpScreenshotURL
		 }
	*/

	return entity.Examination{
		ID:          examination.ID().String(),
		OfferItemID: examination.OfferItemID().String(),
		AssigneeID:  examination.AssigneeID().String(),
		EntryID:     null.StringFromPtr(entryID),
		/*
			SNSUserID:        null.StringFromPtr(examination.Sns().UserID()),
			SNSScreenshotURL: null.StringFromPtr(snsScreenshotURL),
		*/
		ExaminerName: null.StringFromPtr(examination.ExaminerName()),
		Reason:       null.StringFromPtr(examination.Reason()),
		EntryType:    uint(examination.EntryType()),
	}
}
