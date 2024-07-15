package repository_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/internal/infrastructure/converter"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opencensus.io/trace"

	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
)

func NewAssigneeRepositoryImpl() repository.AssigneeRepository {
	//	entity.AddAssigneeHook(boil.BeforeInsertHook, func(ctx context.Context, exec boil.ContextExecutor, e *entity.Assignee) error {
	//		requestedBy, err := common_metadata.GetRequestedByFromContext(ctx)
	//		if err != nil {
	//			return fmt.Errorf("common_metadata.GetRequestedByFromContext: %w", err)
	//		}
	//		e.CreatedBy = requestedBy.String()
	//		e.UpdatedBy = requestedBy.String()
	//		return nil
	//	})
	//	entity.AddAssigneeHook(boil.BeforeUpdateHook, func(ctx context.Context, exec boil.ContextExecutor, e *entity.Assignee) error {
	//		requestedBy, err := common_metadata.GetRequestedByFromContext(ctx)
	//		if err != nil {
	//			return fmt.Errorf("common_metadata.GetRequestedByFromContext: %w", err)
	//		}
	//		e.UpdatedBy = requestedBy.String()
	//		return nil
	//	})
	return &AssigneeRepositoryImpl{}
}

type AssigneeRepositoryImpl struct{}

func (a *AssigneeRepositoryImpl) ListByOfferItemIDAmebaIDs(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID, amebaIDs []model.AmebaID, withLock bool) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListByOfferItemIDAmebaIDs")
	defer span.End()

	amebaIDStrings := make([]string, len(amebaIDs))
	for i, id := range amebaIDs {
		amebaIDStrings[i] = id.String()
	}

	queries := make([]qm.QueryMod, 0)
	queries = append(queries, entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()))
	queries = append(queries, entity.AssigneeWhere.AmebaID.IN(amebaIDStrings))
	if withLock {
		queries = append(queries, qm.For("UPDATE"))
	}

	assigneeEntities, err := entity.Assignees(queries...).All(ctx, tx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.OfferItemNotFoundError.Wrap(err)
		}
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}

	assignees := make([]*model.Assignee, len(assigneeEntities))
	for i, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		assignees[i] = assignee
	}

	return assignees, nil
}

func (a *AssigneeRepositoryImpl) BulkGetByOfferItemIDAmebaIDs(ctx context.Context, db *sql.DB, offerItemID model.OfferItemID, amebaIDs []model.AmebaID, withLock bool) (map[model.AmebaID]*model.Assignee, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.BulkGetByOfferItemIDAmebaIDs")
	defer span.End()

	amebaIDStrings := make([]string, len(amebaIDs))
	for i, id := range amebaIDs {
		amebaIDStrings[i] = id.String()
	}

	queries := make([]qm.QueryMod, 0)
	queries = append(queries, entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()))
	queries = append(queries, entity.AssigneeWhere.AmebaID.IN(amebaIDStrings))
	if withLock {
		queries = append(queries, qm.For("UPDATE"))
	}

	assigneeEntities, err := entity.Assignees(queries...).All(ctx, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.OfferItemNotFoundError.Wrap(err)
		}
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}
	assigneeMap := make(map[model.AmebaID]*model.Assignee)
	for _, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		// AmebaIDとOfferItemIDの組み合わせでユニークなキーを作成
		assigneeMap[assignee.AmebaID()] = assignee
	}
	return assigneeMap, nil
}

func (a *AssigneeRepositoryImpl) Update(ctx context.Context, exec boil.ContextExecutor, assignee *model.Assignee) error {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.Update")
	defer span.End()

	blackList := boil.Blacklist(
		entity.AssigneeColumns.CreatedAt,
		entity.AssigneeColumns.CreatedBy,
	)
	assigneeEntity := converter.AssigneeModelToEntity(assignee)

	if _, err := assigneeEntity.Update(ctx, exec, blackList); err != nil {
		return fmt.Errorf("entity.Assignees.Update: %w", err)
	}

	return nil
}

// アサイニーを作成する
func (a *AssigneeRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, assignee *model.Assignee) error {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.Create")
	defer span.End()

	assigneeEntity := converter.AssigneeModelToEntity(assignee)
	if err := assigneeEntity.Insert(ctx, tx, boil.Infer()); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}
	return nil
}

// 指定されたOfferItemIDとStageに紐づくAssigneeを取得する
func (a *AssigneeRepositoryImpl) ListByOfferItemIDStage(ctx context.Context, db boil.ContextExecutor, offerItemID model.OfferItemID, stage model.Stage) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListByOfferItemIDStage")
	defer span.End()

	assigneeEntities, err := entity.Assignees(entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()), entity.AssigneeWhere.Stage.EQ(uint(stage))).All(ctx, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AssigneeList{}, nil
		}
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}

	assignees := make([]*model.Assignee, len(assigneeEntities))
	for i, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		assignees[i] = assignee
	}

	return assignees, nil
}

// 下書き審査、記事審査ステージのアサイニーを取得する
func (a *AssigneeRepositoryImpl) ListUnderExamination(ctx context.Context, db boil.ContextExecutor) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListUnderExamination")
	defer span.End()

	assigneeEntities, err := entity.Assignees(
		qm.Expr(entity.AssigneeWhere.Stage.EQ(uint(model.StagePreExamination)), qm.Or2(entity.AssigneeWhere.Stage.EQ(uint(model.StageExamination)))),
	).All(ctx, db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AssigneeList{}, nil
		}
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}

	assignees := make([]*model.Assignee, len(assigneeEntities))
	for i, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		assignees[i] = assignee
	}

	return assignees, nil
}

// ステージごとのアサイニー数を取得する
func (a *AssigneeRepositoryImpl) ListCount(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID) ([]model.AssigneeCount, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListCount")
	defer span.End()

	// クエリ作成
	mods := []qm.QueryMod{
		qm.Select(
			entity.AssigneeColumns.Stage,
			"COUNT(*) AS count",
		),
		qm.From(entity.TableNames.Assignee),
		entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()),
		qm.GroupBy(
			entity.AssigneeColumns.Stage,
		),
	}

	type StageAssigneeCount struct {
		Stage uint `boil:"stage"`
		Count int  `boil:"count"`
	}

	var records []StageAssigneeCount

	if err := entity.NewQuery(mods...).Bind(ctx, exec, &records); err != nil {
		return nil, fmt.Errorf("entity.NewQuery.Bind: %w", err)
	}

	assigneeCounts := make([]model.AssigneeCount, 0, len(records))
	for _, record := range records {
		assigneeCount := model.NewAssigneeCountFromRepository(model.Stage(record.Stage), record.Count)
		assigneeCounts = append(assigneeCounts, assigneeCount)
	}
	return assigneeCounts, nil
}

// ステージが「支払い中のアサイニーを取得する
func (a *AssigneeRepositoryImpl) ListUnderPaying(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID, amebaIDs []model.AmebaID) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListUnderPaying")
	defer span.End()

	queries := make([]qm.QueryMod, 0, 3)
	queries = append(
		queries,
		entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()),
		entity.AssigneeWhere.AmebaID.IN(amebaIDsToSliceString(amebaIDs)),
		entity.AssigneeWhere.Stage.EQ(uint(model.StagePaying)),
	)

	assigneeEntities, err := entity.Assignees(queries...).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}

	assignees := make(model.AssigneeList, 0, len(assigneeEntities))
	for _, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		assignees = append(assignees, assignee)
	}
	return assignees, nil
}

func amebaIDsToSliceString(amebaIDs []model.AmebaID) []string {
	amebaIDStrings := make([]string, len(amebaIDs))
	for i, id := range amebaIDs {
		amebaIDStrings[i] = id.String()
	}
	return amebaIDStrings
}

// アサイニーを取得する
func (a *AssigneeRepositoryImpl) ListByOfferItemID(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListByOfferItemID")
	defer span.End()

	assigneeEntities, err := entity.Assignees(entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String())).All(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AssigneeList{}, nil
		}
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}

	assignees := make([]*model.Assignee, len(assigneeEntities))
	for i, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		assignees[i] = assignee
	}
	return assignees, nil
}

// アメーバIDに紐づくアサイニー、オファーアイテムを取得する
func (a *AssigneeRepositoryImpl) GetByAmebaIDOfferItemID(ctx context.Context, exec boil.ContextExecutor, amebaID model.AmebaID, offerItemID model.OfferItemID) (*model.Assignee, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.GetByAmebaIDOfferItemID")
	defer span.End()

	assigneeEntity, err := entity.Assignees(
		entity.AssigneeWhere.AmebaID.EQ(amebaID.String()),
		entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()),
	).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.OfferItemNotFoundError.Wrap(err)
		}
		return nil, fmt.Errorf("entity.Assignees.One: %w", err)
	}

	assignee := converter.AssigneeEntityToModel(assigneeEntity)
	return assignee, nil
}

// アメーバIDに紐づくアサイニー、オファーアイテムを取得する
func (a *AssigneeRepositoryImpl) ListByAmebaID(ctx context.Context, exec boil.ContextExecutor, amebaID model.AmebaID) (model.AssigneeList, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.ListByAmebaID")
	defer span.End()

	assigneeEntities, err := entity.Assignees(entity.AssigneeWhere.AmebaID.EQ(amebaID.String())).All(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AssigneeList{}, nil
		}
		return nil, fmt.Errorf("entity.Assignees.All: %w", err)
	}

	assignees := make([]*model.Assignee, len(assigneeEntities))
	for i, assigneeEntity := range assigneeEntities {
		assignee := converter.AssigneeEntityToModel(assigneeEntity)
		assignees[i] = assignee
	}
	return assignees, nil
}

// アサイニーを取得する
func (a *AssigneeRepositoryImpl) Get(ctx context.Context, exec boil.ContextExecutor, assigneeID model.AssigneeID) (*model.Assignee, error) {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.Get")
	defer span.End()

	assigneeEntity, err := entity.Assignees(entity.AssigneeWhere.ID.EQ(assigneeID.String())).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.OfferItemNotFoundError.Wrap(err)
		}
		return nil, fmt.Errorf("entity.FindAssignee: %w", err)
	}

	assignee := converter.AssigneeEntityToModel(assigneeEntity)
	return assignee, nil
}

func (a *AssigneeRepositoryImpl) DeleteByOfferItemIDAndAmebaID(ctx context.Context, tx *sql.Tx, offerItemID model.OfferItemID, amebaID model.AmebaID) error {
	ctx, span := trace.StartSpan(ctx, "AssigneeRepositoryImpl.DeleteByOfferItemIDAndAmebaID")
	defer span.End()

	_, err := entity.Assignees(
		entity.AssigneeWhere.OfferItemID.EQ(offerItemID.String()),
		entity.AssigneeWhere.AmebaID.EQ(amebaID.String()),
	).DeleteAll(ctx, tx, false)
	if err != nil {
		return fmt.Errorf("entity.Assignees.DeleteAll: %w", err)
	}
	return nil
}
