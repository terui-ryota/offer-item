package dto

//func MapExaminationResultPBToDTO(mapExaminationResultsPB map[string]*offer_item.ExaminationResult) map[string]*ExaminationResultDTO {
//	resultsDTOMap := make(map[string]*ExaminationResultDTO)
//
//	for amebaIDStr, v := range mapExaminationResultsPB {
//		var reason *string
//		if v.GetReason() != "" {
//			tmpReason := v.GetReason()
//			reason = &tmpReason
//		}
//
//		resultsDTOMap[amebaIDStr] = &ExaminationResultDTO{
//			IsPassed:     v.GetIsPassed(),
//			ExaminerName: v.GetExaminerName(),
//			Reason:       reason,
//		}
//	}
//
//	return resultsDTOMap
//}

type ExaminationResultDTO struct {
	IsPassed     bool
	ExaminerName string
	Reason       *string
	EntryID      *string
	SNS          *SNS
}

type SNS struct {
	UserID        *string
	ScreenshotURL string
}
