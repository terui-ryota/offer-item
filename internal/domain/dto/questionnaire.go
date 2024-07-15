package dto

import (
	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/protofiles/go/offer_item"
)

type Questionnaire struct {
	Description string
	Questions   []Question
}

type Question struct {
	ID           *string
	QuestionType model.QuestionType
	Title        string
	ImageURL     string
	Options      []string
}

func convertQuestionType(pb offer_item.Questionnaire_QuestionType) model.QuestionType {
	switch pb {
	case offer_item.Questionnaire_QUESTION_TYPE_RADIO:
		return model.QuestionTypeRadio
	case offer_item.Questionnaire_QUESTION_TYPE_TEXT:
		return model.QuestionTypeText
	default:
		return model.QuestionTypeUnknown
	}
}

func ConvertQuestionnaire(pb *offer_item.Questionnaire) *Questionnaire {
	return &Questionnaire{
		Description: pb.GetDescription(),
		Questions: func() []Question {
			qs := make([]Question, 0, len(pb.GetQuestions()))
			for _, q := range pb.GetQuestions() {
				qs = append(qs, Question{
					ID: func() *string {
						if q.GetId() == "" {
							return nil
						}
						id := q.GetId()
						return &id
					}(),
					QuestionType: convertQuestionType(q.GetQuestionType()),
					Title:        q.GetTitle(),
					ImageURL:     q.GetImageUrl(),
					Options:      q.GetOptions(),
				})
			}
			return qs
		}(),
	}
}
