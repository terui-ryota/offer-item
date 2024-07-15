package model

import (
	"testing"
	"time"
)

func TestAssignee_SetStageDoneByDecline(t *testing.T) {
	type fields struct {
		id            AssigneeID
		offerItemID   OfferItemID
		amebaID       AmebaID
		writingFee    int
		stage         Stage
		declineReason *string
		createdAt     time.Time
	}
	type args struct {
		declineReason string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "異常系。 0 文字",
			fields: fields{},
			args: args{
				declineReason: "",
			},
			wantErr: true,
		},
		{
			name:   "正常系。半角 128 文字",
			fields: fields{},
			args: args{
				declineReason: func() string {
					var r string
					for range 128 {
						r += "a"
					}
					return r
				}(),
			},
			wantErr: false,
		},
		{
			name:   "異常系。半角 129 文字",
			fields: fields{},
			args: args{
				declineReason: func() string {
					var r string
					for range 129 {
						r += "a"
					}
					return r
				}(),
			},
			wantErr: true,
		},
		{
			name:   "正常系。マルチバイト 128 文字",
			fields: fields{},
			args: args{
				declineReason: func() string {
					var r string
					for range 128 {
						r += "👨‍👩‍👧‍👦"
					}
					return r
				}(),
			},
			wantErr: false,
		},
		{
			name:   "異常系。マルチバイト 129 文字",
			fields: fields{},
			args: args{
				declineReason: func() string {
					var r string
					for range 129 {
						r += "👨‍👩‍👧‍👦"
					}
					return r
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assignee{
				id:            tt.fields.id,
				offerItemID:   tt.fields.offerItemID,
				amebaID:       tt.fields.amebaID,
				writingFee:    tt.fields.writingFee,
				stage:         tt.fields.stage,
				declineReason: tt.fields.declineReason,
				createdAt:     tt.fields.createdAt,
			}
			if err := a.SetStageDoneByDecline(tt.args.declineReason); (err != nil) != tt.wantErr {
				t.Errorf("Assignee.SetStageDoneByDecline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssignee_ChangeStageByDraftSubmission(t *testing.T) {
	type fields struct {
		id            AssigneeID
		offerItemID   OfferItemID
		amebaID       AmebaID
		writingFee    int
		stage         Stage
		declineReason *string
		createdAt     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    Stage
	}{
		{
			name: "正常系。下書き提出から、下書き審査に変更される",
			fields: fields{
				stage: StageDraftSubmission,
			},
			wantErr: false,
			want:    StagePreExamination,
		},
		{
			name: "正常系。下書き再審査から、下書き審査に変更される",
			fields: fields{
				stage: StagePreReexamination,
			},
			wantErr: false,
			want:    StagePreExamination,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assignee{
				id:            tt.fields.id,
				offerItemID:   tt.fields.offerItemID,
				amebaID:       tt.fields.amebaID,
				writingFee:    tt.fields.writingFee,
				stage:         tt.fields.stage,
				declineReason: tt.fields.declineReason,
				createdAt:     tt.fields.createdAt,
			}
			if err := a.ChangeStageByDraftSubmission(); (err != nil) != tt.wantErr {
				t.Errorf("Assignee.ChangeStageByDraftSubmission() error = %v, wantErr %v", err, tt.wantErr)
			}
			if a.stage != tt.want {
				t.Errorf("Assignee.stage error = %v, want %v", a.stage, tt.want)
			}
		})
	}
}

func TestAssignee_ChangeStageByEntrySubmission(t *testing.T) {
	type fields struct {
		id            AssigneeID
		offerItemID   OfferItemID
		amebaID       AmebaID
		writingFee    int
		stage         Stage
		declineReason *string
		createdAt     time.Time
	}
	type args struct {
		needsAfterReview bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Stage
	}{
		{
			name: "正常系。記事提出 x 事後審査あり、から記事審査に変更される",
			fields: fields{
				stage: StageArticlePosting,
			},
			args: args{
				needsAfterReview: true,
			},
			wantErr: false,
			want:    StageExamination,
		},
		{
			name: "正常系。記事提出 x 事後審査なし、から支払い中に変更される",
			fields: fields{
				stage: StageArticlePosting,
			},
			args: args{
				needsAfterReview: false,
			},
			wantErr: false,
			want:    StagePaying,
		},
		{
			name: "正常系。記事再審査 x 事後審査あり、から記事審査に変更される",
			fields: fields{
				stage: StageReexamination,
			},
			args: args{
				needsAfterReview: true,
			},
			wantErr: false,
			want:    StageExamination,
		},
		{
			name: "正常系。記事再審査 x 事後審査なし、から支払い中に変更される",
			fields: fields{
				stage: StageReexamination,
			},
			args: args{
				needsAfterReview: false,
			},
			wantErr: false,
			want:    StagePaying,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Assignee{
				id:            tt.fields.id,
				offerItemID:   tt.fields.offerItemID,
				amebaID:       tt.fields.amebaID,
				writingFee:    tt.fields.writingFee,
				stage:         tt.fields.stage,
				declineReason: tt.fields.declineReason,
				createdAt:     tt.fields.createdAt,
			}
			if err := a.ChangeStageByEntrySubmission(tt.args.needsAfterReview); (err != nil) != tt.wantErr {
				t.Errorf("Assignee.ChangeStageByEntrySubmission() error = %v, wantErr %v", err, tt.wantErr)
			}
			if a.stage != tt.want {
				t.Errorf("Assignee.stage error = %v, want %v", a.stage, tt.want)
			}
		})
	}
}
