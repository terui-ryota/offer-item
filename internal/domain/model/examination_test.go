package model

import (
	"testing"

	"github.com/volatiletech/null/v8"
)

func TestExamination_SetExaminationResult(t *testing.T) {
	type fields struct {
		id                   ExaminationID
		offerItemID          OfferItemID
		amebaID              AmebaID
		entryID              *EntryID
		examinerName         *string
		reason               *string
		assigneeID           AssigneeID
		entryType            EntryType
		entrySubmissionCount uint
	}
	type args struct {
		isPassed     bool
		examinerName string
		reason       *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "正常系。 isPassed=false",
			fields: fields{},
			args: args{
				isPassed:     false,
				examinerName: "サイバー太郎",
				reason:       null.StringFrom("xxxな理由でNG").Ptr(),
			},
			wantErr: false,
		},
		{
			name:   "異常系。 isPassed=false x NG 理由が空文字",
			fields: fields{},
			args: args{
				isPassed:     false,
				examinerName: "サイバー太郎",
				reason:       null.StringFrom("").Ptr(),
			},
			wantErr: true,
		},
		{
			name:   "異常系。 isPassed=false x NG理由が nil",
			fields: fields{},
			args: args{
				isPassed:     false,
				examinerName: "サイバー太郎",
				reason:       nil,
			},
			wantErr: true,
		},
		// 後方互換用のコメントアウトを解除後に有効にしてください
		// {
		// 	name:   "異常系。 isPassed=false x 審査者が空文字",
		// 	fields: fields{},
		// 	args: args{
		// 		isPassed:     false,
		// 		examinerName: "",
		// 		reason:       null.StringFrom("xxxな理由でNG").Ptr(),
		// 	},
		// 	wantErr: true,
		// },
		{
			name:   "正常系。 isPassed=true",
			fields: fields{},
			args: args{
				isPassed:     true,
				examinerName: "サイバー太郎",
				reason:       nil,
			},
			wantErr: false,
		},
		// 後方互換用のコメントアウトを解除後に有効にしてください
		// {
		// 	name:   "異常系。 isPassed=true x 審査者が空文字",
		// 	fields: fields{},
		// 	args: args{
		// 		isPassed:     true,
		// 		examinerName: "",
		// 		reason:       null.StringFrom("OK").Ptr(),
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Examination{
				id:                   tt.fields.id,
				offerItemID:          tt.fields.offerItemID,
				amebaID:              tt.fields.amebaID,
				entryID:              tt.fields.entryID,
				examinerName:         tt.fields.examinerName,
				reason:               tt.fields.reason,
				assigneeID:           tt.fields.assigneeID,
				entryType:            tt.fields.entryType,
				entrySubmissionCount: tt.fields.entrySubmissionCount,
			}
			if err := e.SetExaminationResult(tt.args.isPassed, tt.args.examinerName, tt.args.reason); (err != nil) != tt.wantErr {
				t.Errorf("Examination.SetExaminationResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
