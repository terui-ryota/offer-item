package repository_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/internal/infrastructure/converter"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opencensus.io/trace"
)

func NewExaminationRepositoryImpl() repository.ExaminationRepository {
	return &ExaminationRepositoryImpl{}
}

type ExaminationRepositoryImpl struct{}

// BulkGetByOfferItemID 各アサイ二ーの最新のexaminationを取得する
func (e *ExaminationRepositoryImpl) BulkGetByOfferItemID(ctx context.Context, exec *sql.DB, offerItemID model.OfferItemID, entryType model.EntryType) (map[model.AmebaID]*model.Examination, error) {
	ctx, span := trace.StartSpan(ctx, "ExaminationRepositoryImpl.BulkGetByOfferItemID")
	defer span.End()

	// SQLではサブクエリを使った場合複雑になる為、一度条件に合う全てのExaminationを取得する。
	examinationsEntities, err := entity.Examinations(
		entity.ExaminationWhere.OfferItemID.EQ(offerItemID.String()),
		entity.ExaminationWhere.EntryType.EQ(uint(entryType)),
		qm.Load(entity.ExaminationRels.Assignee),
	).All(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[model.AmebaID]*model.Examination{}, nil
		}
		return nil, fmt.Errorf("entity.Examinations.All: %w", err)
	}

	//「offerItem」と「同じassigneeIDのexamination」が1対Nの関係になっている。「offerItem」と「createdAtが最新のassigneeIDのexamination」は1対1の関係にするために、最新のExaminationのみを取得する
	latestExaminations := make(map[string]*entity.Examination)
	assigneeCountMap := make(map[string]int)
	for _, e := range examinationsEntities {
		count := assigneeCountMap[e.AssigneeID]
		count++
		assigneeCountMap[e.AssigneeID] = count
		if latest, exists := latestExaminations[e.AssigneeID]; exists {
			if e.CreatedAt.After(latest.CreatedAt) {
				latestExaminations[e.AssigneeID] = e
			}
		} else {
			latestExaminations[e.AssigneeID] = e
		}
	}

	examinations := make(model.ExaminationList, 0, len(latestExaminations))

	// modelに変換
	for _, examinationEntity := range latestExaminations {
		examination := converter.ExaminationEntityToModel(examinationEntity, assigneeCountMap[examinationEntity.AssigneeID])
		examinations = append(examinations, examination)
	}

	// AmebaIDをキーにしたmapに変換
	examinationMap := make(map[model.AmebaID]*model.Examination)
	for _, examination := range examinations {
		examinationMap[examination.AmebaID()] = examination
	}

	return examinationMap, nil
}

func (e *ExaminationRepositoryImpl) Get(ctx context.Context, exec *sql.DB, offerItemID model.OfferItemID, assigneeID model.AssigneeID, entryType model.EntryType) (*model.Examination, error) {
	ctx, span := trace.StartSpan(ctx, "ExaminationRepositoryImpl.Get")
	defer span.End()

	entities, err := entity.Examinations(
		entity.ExaminationWhere.OfferItemID.EQ(offerItemID.String()),
		entity.ExaminationWhere.AssigneeID.EQ(assigneeID.String()),
		entity.ExaminationWhere.EntryType.EQ(uint(entryType)),
		qm.Load(entity.ExaminationRels.Assignee),
		qm.OrderBy(entity.ExaminationColumns.CreatedAt+" DESC"),
	).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.Examinations.All: %w", err)
	}

	// 該当するExaminationが存在しない場合はエラーを返す
	if len(entities) == 0 {
		return nil, apperr.OfferItemNotFoundError.Wrap(errors.New("examination not found"))
	}

	examination := converter.ExaminationEntityToModel(entities[0], len(entities))

	return examination, nil
}

func (e *ExaminationRepositoryImpl) Update(ctx context.Context, exec *sql.DB, examination *model.Examination) error {
	ctx, span := trace.StartSpan(ctx, "ExaminationRepositoryImpl.Update")
	defer span.End()

	blackList := boil.Blacklist(
		entity.ExaminationColumns.CreatedAt,
	)
	examinationEntity := converter.ExaminationModelToEntity(examination)
	if _, err := examinationEntity.Update(ctx, exec, blackList); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}
	return nil
}

func (e *ExaminationRepositoryImpl) Create(ctx context.Context, exec *sql.DB, examination *model.Examination) error {
	ctx, span := trace.StartSpan(ctx, "ExaminationRepositoryImpl.Create")
	defer span.End()

	// TODO: SecondリリースでSNSの実装を行う
	examinationEntity := converter.ExaminationModelToEntity(examination)

	if err := examinationEntity.Insert(ctx, exec, boil.Infer()); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}

	return nil
}
