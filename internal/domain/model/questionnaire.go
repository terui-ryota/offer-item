package model

import (
	"fmt"

	"github.com/terui-ryota/offer-item/pkg/id"
)

//go:generate go run github.com/terui-ryota/gen-getter -type=Questionnaire
type Questionnaire struct {
	offerItemID OfferItemID
	description string
	questions   []Question
}

func (q *Questionnaire) SetDescription(v string) error {
	q.description = v
	return nil
}

func (q *Questionnaire) SetQuestions(v []Question) error {
	q.questions = v
	return nil
}

func NewQuestionnaire(
	offerItemID OfferItemID,
	description string,
	questions []Question,
) (*Questionnaire, error) {
	if description == "" {
		return nil, fmt.Errorf("description must not be empty")
	}
	if err := validateQuestions(questions); err != nil {
		return nil, fmt.Errorf("validateQuestions: %w", err)
	}
	return &Questionnaire{
		offerItemID: offerItemID,
		description: description,
		questions:   questions,
	}, nil
}

func NewQuestionnaireFromRepository(
	offerItemID OfferItemID,
	description string,
	questions []Question,
) *Questionnaire {
	return &Questionnaire{
		offerItemID: offerItemID,
		description: description,
		questions:   questions,
	}
}

func validateQuestions(questions []Question) error {
	if len(questions) == 0 {
		return fmt.Errorf("len of questions must not be 0")
	}
	m := func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, q := range questions {
			m[q.title] = struct{}{}
		}
		return m
	}()
	if len(m) != len(questions) {
		return fmt.Errorf("question title must be unique in questions: %+v", questions)
	}
	return nil
}

type QuestionType int

const (
	QuestionTypeUnknown QuestionType = iota
	QuestionTypeRadio
	QuestionTypeText
)

type QuestionID string

func NewQuestionType(i int) QuestionType {
	switch i {
	case int(QuestionTypeRadio):
		return QuestionTypeRadio
	case int(QuestionTypeText):
		return QuestionTypeText
	default:
		return QuestionTypeUnknown
	}
}

func (i QuestionID) String() string {
	return string(i)
}

//go:generate go run github.com/terui-ryota/gen-getter -type=Question
type Question struct {
	id           QuestionID
	offerItemID  OfferItemID
	questionType QuestionType
	title        string
	imageURL     string
	options      []string
}

func (q *Question) SetOptions(questionType QuestionType, options ...string) error {
	if err := validateOptions(questionType, options); err != nil {
		return fmt.Errorf("validateOptions: %w, %+v, %+v", err, *q, options)
	}
	q.questionType = questionType
	q.options = options
	return nil
}

func (q *Question) SetTitle(v string) error {
	q.title = v
	return nil
}

func (q *Question) SetImageURL(v string) error {
	q.imageURL = v
	return nil
}

func NewQuestion(
	offerItemID OfferItemID,
	questionType QuestionType,
	title,
	imageURL string,
	options []string,
) (*Question, error) {
	if questionType == QuestionTypeUnknown {
		return nil, fmt.Errorf("questionType must not be unknown")
	}
	if title == "" {
		return nil, fmt.Errorf("title must not be empty")
	}
	// 1st リリースで画像の連携は行わないためバリデーションを無効化
	// if imageURL == "" {
	// 	return nil, fmt.Errorf("imageURL must not be empty")
	// }
	if err := validateOptions(questionType, options); err != nil {
		return nil, fmt.Errorf("validateOptions: %w", err)
	}
	return &Question{
		id:           QuestionID(id.New()),
		offerItemID:  offerItemID,
		questionType: questionType,
		title:        title,
		imageURL:     imageURL,
		options:      options,
	}, nil
}

func NewQuestionFromRepository(
	id QuestionID,
	offerItemID OfferItemID,
	questionType QuestionType,
	title,
	imageURL string,
	options []string,
) *Question {
	return &Question{
		id:           id,
		offerItemID:  offerItemID,
		questionType: questionType,
		title:        title,
		imageURL:     imageURL,
		options:      options,
	}
}

func validateOptions(t QuestionType, options []string) error {
	switch t {
	case QuestionTypeRadio:
		if len(options) == 0 {
			return fmt.Errorf("len of options must not be 0, type: %d", t)
		}

		optionsMap := make(map[string]struct{})
		for _, option := range options {
			if _, ok := optionsMap[option]; ok {
				// 重複がある場合はエラーを返却する
				return fmt.Errorf("option name must be unique in options: %+v", options)
			} else {
				optionsMap[option] = struct{}{}
			}
		}

		return nil
	case QuestionTypeText:
		if len(options) != 0 {
			return fmt.Errorf("len of options must be 0, type: %d", t)
		}
		return nil
	default:
		return fmt.Errorf("unknown type: %d", t)
	}
}

//go:generate go run github.com/terui-ryota/gen-getter -type=QuestionAnswer
type QuestionAnswer struct {
	assigneeID  AssigneeID
	questionID  QuestionID
	offerItemID OfferItemID
	content     string
}

func NewQuestionAnswers(
	assigneeID AssigneeID,
	questionnaire Questionnaire,
	answers map[QuestionID]string,
) ([]QuestionAnswer, error) {
	if len(questionnaire.questions) != len(answers) {
		return nil, fmt.Errorf("len of questions and len of answers are not equal: %d %d", len(questionnaire.questions), len(answers))
	}
	res := make([]QuestionAnswer, 0, len(answers))
	for _, q := range questionnaire.questions {
		if answer, ok := answers[q.id]; !ok {
			return nil, fmt.Errorf("questionID does not contain: %s", q.id)
		} else {
			q, err := newQuestionAnswer(
				assigneeID,
				q,
				answer,
			)
			if err != nil {
				return nil, fmt.Errorf("newQuestionAnswer: %w", err)
			}
			res = append(res, *q)
		}
	}
	return res, nil
}

func newQuestionAnswer(
	assigneeID AssigneeID,
	question Question,
	content string,
) (*QuestionAnswer, error) {
	m := func() map[string]struct{} {
		m := make(map[string]struct{})
		for _, o := range question.options {
			m[o] = struct{}{}
		}
		return m
	}()
	switch question.questionType {
	case QuestionTypeRadio:
		if _, ok := m[content]; !ok {
			return nil, fmt.Errorf("unknown option: %+v, %v", question, content)
		}
	case QuestionTypeText:
		// noop
	default:
		return nil, fmt.Errorf("unknown type: %+v", question)
	}
	return &QuestionAnswer{
		assigneeID:  assigneeID,
		questionID:  question.id,
		offerItemID: question.offerItemID,
		content:     content,
	}, nil
}

func NewQuestionAnswerFromRepository(
	assigneeID AssigneeID,
	offerItemID OfferItemID,
	questionID string,
	content string,
) *QuestionAnswer {
	return &QuestionAnswer{
		assigneeID:  assigneeID,
		questionID:  QuestionID(questionID),
		offerItemID: offerItemID,
		content:     content,
	}
}
