// Code generated by gen-getter. DO NOT EDIT.
package model

func (e *Examination) ID() ExaminationID {
	return e.id
}
func (e *Examination) OfferItemID() OfferItemID {
	return e.offerItemID
}
func (e *Examination) AmebaID() AmebaID {
	return e.amebaID
}
func (e *Examination) EntryID() *EntryID {
	return e.entryID
}
func (e *Examination) ExaminerName() *string {
	return e.examinerName
}
func (e *Examination) Reason() *string {
	return e.reason
}
func (e *Examination) AssigneeID() AssigneeID {
	return e.assigneeID
}
func (e *Examination) EntryType() EntryType {
	return e.entryType
}
func (e *Examination) EntrySubmissionCount() uint {
	return e.entrySubmissionCount
}
