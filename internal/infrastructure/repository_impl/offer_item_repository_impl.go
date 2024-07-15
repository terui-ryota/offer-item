package repository_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/terui-ryota/offer-item/internal/domain/dto"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.opencensus.io/trace"

	"github.com/terui-ryota/offer-item/internal/domain/model"
	"github.com/terui-ryota/offer-item/internal/domain/repository"
	"github.com/terui-ryota/offer-item/internal/infrastructure/converter"
	"github.com/terui-ryota/offer-item/internal/infrastructure/db/entity"
	"github.com/terui-ryota/offer-item/internal/infrastructure/util/dbhelper"
)

func NewOfferItemRepositoryImpl() repository.OfferItemRepository {
	//entity.AddOfferItemHook(boil.BeforeInsertHook, func(ctx context.Context, exec boil.ContextExecutor, e *entity.OfferItem) error {
	//	requestedBy, err := common_metadata.GetRequestedByFromContext(ctx)
	//	if err != nil {
	//		return fmt.Errorf("common_metadata.GetRequestedByFromContext: %w", err)
	//	}
	//	e.CreatedBy = requestedBy.String()
	//	e.UpdatedBy = requestedBy.String()
	//	return nil
	//})
	//entity.AddOfferItemHook(boil.BeforeDeleteHook, func(ctx context.Context, exec boil.ContextExecutor, e *entity.OfferItem) error {
	//	requestedBy, err := common_metadata.GetRequestedByFromContext(ctx)
	//	if err != nil {
	//		return fmt.Errorf("common_metadata.GetRequestedByFromContext: %w", err)
	//	}
	//	e.UpdatedBy = requestedBy.String()
	//	if e.DeletedAt.Valid {
	//		e.DeletedBy = converter.ToNullString(requestedBy.String())
	//	}
	//	return nil
	//})
	return &OfferItemRepositoryImpl{}
}

type OfferItemRepositoryImpl struct{}

// // オファー案件の一覧を取得する
func (o *OfferItemRepositoryImpl) List(ctx context.Context, exec boil.ContextExecutor, condition *model.ListCondition, isClosed bool) (*model.ListOfferItemResult, error) {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.List")
	defer span.End()

	queries := make([]qm.QueryMod, 0)

	if !isClosed {
		queries = append(queries, entity.OfferItemWhere.IsClosed.EQ(false))
	}
	// データ取得前に検索結果の総数を取得する
	totalCount, err := entity.OfferItems(queries...).Count(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.OfferItems.Count: %w", err)
	}

	// リスト条件をクエリに追加する
	queries = append(queries, qm.Limit(condition.Limit()), qm.Offset(condition.Offset()))
	if len(condition.Sorts()) > 0 {
		queries = append(queries, qm.OrderBy(dbhelper.CreateOrderByClause(entity.OfferItemColumns, condition.Sorts())))
	}

	// DBからデータを取得し、モデルに変換する
	offerItems := make(model.OfferItemList, 0, condition.Limit())
	offerItemEntities, err := entity.OfferItems(queries...).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.OfferItems.All: %w", err)
	}

	for _, offerItemEntity := range offerItemEntities {
		offerItemID := model.OfferItemID(offerItemEntity.ID)
		// offerItemIDをkeyにスケジュールを取得
		scheduleEntities, err := entity.Schedules(entity.ScheduleWhere.OfferItemID.EQ(offerItemID.String())).All(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("entity.Schedules.All: %w", err)
		}
		schedules := make(model.ScheduleList, 0, len(scheduleEntities))
		// スケジュールをモデルに変換
		for _, scheduleEntity := range scheduleEntities {
			schedule := converter.ScheduleEntityToModel(scheduleEntity)
			schedules = append(schedules, schedule)
		}

		draftedItemInfoEntity, err := entity.DraftedItemInfos(entity.DraftedItemInfoWhere.OfferItemID.EQ(offerItemID.String())).One(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("entity.DraftedItemInfos.One: %w", err)
		}
		draftedItemInfo, err := converter.ConvertDraftedItemToModel(draftedItemInfoEntity)
		if err != nil {
			return nil, fmt.Errorf("converter.ConvertDraftedItemToModel: %w", err)
		}

		offerItem, err := converter.ConvertOfferItemToModel(offerItemEntity, schedules, draftedItemInfo)
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemEntityToModel: %w", err)
		}
		offerItems = append(offerItems, offerItem)
	}

	listOfferItem, err := model.NewListOfferItemResult(offerItems, int(totalCount))
	if err != nil {
		return nil, fmt.Errorf("model.NewListOfferItemResult: %w", err)
	}
	return listOfferItem, nil
}

// オファー案件を論理削除する
func (o *OfferItemRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id model.OfferItemID) error {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.Delete")
	defer span.End()

	// SQLBoiler v4.0系から指定しなくてもwhereの条件にdeleted_at is nullが追加されている
	offerItemEntity, err := entity.OfferItems(
		entity.OfferItemWhere.ID.EQ(id.String()),
	).One(ctx, tx)
	if err != nil {
		// すでに存在しない場合、何もしない
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("entity.OfferItemWhere.ID.EQ.One: %w", err)
	}

	// 論理削除
	if _, err := offerItemEntity.Delete(ctx, tx, false); err != nil {
		return fmt.Errorf("offerItemEntity.Delete: %w", err)
	}
	if _, err := entity.Schedules(entity.ScheduleWhere.OfferItemID.EQ(id.String())).DeleteAll(ctx, tx, false); err != nil {
		return fmt.Errorf("entity.Schedules.DeleteAll: %w", err)
	}
	if _, err := entity.Assignees(entity.AssigneeWhere.OfferItemID.EQ(id.String())).DeleteAll(ctx, tx, false); err != nil {
		return fmt.Errorf("entity.Assignees.DeleteAll: %w", err)
	}
	if _, err := entity.Examinations(entity.ExaminationWhere.OfferItemID.EQ(id.String())).DeleteAll(ctx, tx, false); err != nil {
		return fmt.Errorf("entity.Examinations.DeleteAll: %w", err)
	}

	return nil
}

// 検索条件を指定して、オファー案件を検索する
func (o *OfferItemRepositoryImpl) Search(ctx context.Context, exec boil.ContextExecutor, criteria *dto.SearchOfferItemCriteria, condition *model.ListCondition, isClosed bool) (*model.ListOfferItemResult, error) {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.Search")
	defer span.End()

	queries := make([]qm.QueryMod, 0)

	if criteria.NameContains != nil && *criteria.NameContains != "" {
		// 検索クエリが大文字小文字を無視するように設定
		queries = append(queries, qm.Where("LOWER("+entity.OfferItemColumns.Name+") like ?", fmt.Sprintf("%%%s%%", strings.ToLower(*criteria.NameContains))))
	}

	if criteria.ItemIDEqual != nil && criteria.ItemIDEqual.String() != "" {
		// 文字列比較が大文字小文字を無視するように設定
		queries = append(queries, qm.Where("LOWER("+entity.OfferItemColumns.ItemID+") = ?", strings.ToLower(criteria.ItemIDEqual.String())))
	}

	if criteria.DfItemIDEqual != nil && criteria.DfItemIDEqual.String() != "" {
		// 文字列比較が大文字小文字を無視するように設定
		queries = append(queries, qm.Where("LOWER("+entity.OfferItemColumns.DFItemID+") = ?", strings.ToLower(criteria.DfItemIDEqual.String())))
	}
	if !isClosed {
		queries = append(queries, entity.OfferItemWhere.IsClosed.EQ(false))
	}
	// データ取得前に検索結果の総数を取得する
	totalCount, err := entity.OfferItems(queries...).Count(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.OfferItems.Count: %w", err)
	}
	if totalCount == 0 {
		result, err := model.NewListOfferItemResult(model.OfferItemList{}, int(totalCount))
		if err != nil {
			return nil, fmt.Errorf("model.NewListOfferItemResult: %w", err)
		}
		return result, nil
	}

	// リスト条件をクエリに追加する
	queries = append(queries, qm.Limit(condition.Limit()), qm.Offset(condition.Offset()))
	if len(condition.Sorts()) > 0 {
		queries = append(queries, qm.OrderBy(dbhelper.CreateOrderByClause(entity.OfferItemColumns, condition.Sorts())))
	}

	// DBからデータを取得し、モデルに変換する
	offerItems := make(model.OfferItemList, 0, condition.Limit())
	offerItemEntities, err := entity.OfferItems(queries...).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.OfferItems.All: %w", err)
	}
	for _, offerItemEntity := range offerItemEntities {
		offerItemID := model.OfferItemID(offerItemEntity.ID)
		scheduleEntities, err := entity.Schedules(entity.ScheduleWhere.OfferItemID.EQ(offerItemID.String())).All(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("entity.Schedules.All: %w", err)
		}
		schedules := make(model.ScheduleList, 0, len(scheduleEntities))
		for _, scheduleEntity := range scheduleEntities {
			schedule := converter.ScheduleEntityToModel(scheduleEntity)
			schedules = append(schedules, schedule)
		}

		draftedItemInfoEntity, err := entity.DraftedItemInfos(entity.DraftedItemInfoWhere.OfferItemID.EQ(offerItemID.String())).One(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("entity.DraftedItemInfos.One: %w", err)
		}
		draftedItemInfo, err := converter.ConvertDraftedItemToModel(draftedItemInfoEntity)
		if err != nil {
			return nil, fmt.Errorf("converter.ConvertDraftedItemToModel: %w", err)
		}

		offerItem, err := converter.ConvertOfferItemToModel(offerItemEntity, schedules, draftedItemInfo)
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemEntityToModel: %w", err)
		}
		offerItems = append(offerItems, offerItem)
	}

	result, err := model.NewListOfferItemResult(offerItems, int(totalCount))
	if err != nil {
		return nil, fmt.Errorf("model.NewListOfferItemResult: %w", err)
	}

	return result, nil
}

// オファー案件を取得する
func (o *OfferItemRepositoryImpl) Get(ctx context.Context, exec boil.ContextExecutor, offerItemID model.OfferItemID, withLock bool) (*model.OfferItem, error) {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.Get")
	defer span.End()

	scheduleEntities, err := entity.Schedules(entity.ScheduleWhere.OfferItemID.EQ(offerItemID.String())).All(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.Schedules.All: %w", err)
	}

	schedules := make(model.ScheduleList, 0, len(scheduleEntities))
	for _, scheduleEntity := range scheduleEntities {
		schedule := converter.ScheduleEntityToModel(scheduleEntity)
		schedules = append(schedules, schedule)
	}

	draftedItemInfoEntity, err := entity.DraftedItemInfos(entity.DraftedItemInfoWhere.OfferItemID.EQ(offerItemID.String())).One(ctx, exec)
	if err != nil {
		return nil, fmt.Errorf("entity.DraftedItemInfos.One: %w", err)
	}
	draftedItemInfo, err := converter.ConvertDraftedItemToModel(draftedItemInfoEntity)
	if err != nil {
		return nil, fmt.Errorf("converter.ConvertDraftedItemToModel: %w", err)
	}

	// オファー案件取得
	offerItemQueries := make([]qm.QueryMod, 0, 2)
	offerItemQueries = append(offerItemQueries, entity.OfferItemWhere.ID.EQ(offerItemID.String()))
	if withLock {
		offerItemQueries = append(offerItemQueries, qm.For("UPDATE"))
	}
	offerItemEntity, err := entity.OfferItems(offerItemQueries...).One(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.OfferItemNotFoundError.Wrap(err)
		}
		return nil, fmt.Errorf("entity.OfferItems.One: %w", err)
	}

	offerItem, err := converter.ConvertOfferItemToModel(offerItemEntity, schedules, draftedItemInfo)
	if err != nil {
		return nil, fmt.Errorf("converter.OfferItemEntityToModel: %w", err)
	}

	return offerItem, nil
}

// オファー案件作成
func (o *OfferItemRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, offerItem *model.OfferItem) error {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.Create")
	defer span.End()

	offerItemEntity := converter.ConvertOfferItemModelToEntity(offerItem)
	if err := offerItemEntity.Insert(ctx, tx, boil.Greylist(entity.OfferItemColumns.NeedsPRMark)); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}

	draftedItemInfoEntity := converter.ConvertDraftedItemModelToEntity(offerItem.DraftedItemInfo())
	if err := draftedItemInfoEntity.Insert(ctx, tx, boil.Infer()); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}

	for _, schedule := range offerItem.Schedules() {
		scheduleEntity := converter.ScheduleModelToEntity(schedule)
		if err := scheduleEntity.Insert(ctx, tx, boil.Infer()); err != nil {
			return apperr.OfferItemInternalError.Wrap(err)
		}
	}

	return nil
}

// オファー案件を更新する
func (o *OfferItemRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, offerItem *model.OfferItem) error {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.Update")
	defer span.End()

	blackList := boil.Blacklist(
		entity.OfferItemColumns.CreatedAt,
		entity.OfferItemColumns.CreatedBy,
		entity.OfferItemColumns.DeletedAt,
		entity.OfferItemColumns.DeletedBy,
	)
	offerItemEntity := converter.ConvertOfferItemModelToEntity(offerItem)
	if _, err := offerItemEntity.Update(ctx, tx, blackList); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}

	for _, schedule := range offerItem.Schedules() {
		scheduleEntity := converter.ScheduleModelToEntity(schedule)
		if _, err := scheduleEntity.Update(ctx, tx, blackList); err != nil {
			return apperr.OfferItemInternalError.Wrap(err)
		}
	}

	draftedItemInfoEntity := converter.ConvertDraftedItemModelToEntity(offerItem.DraftedItemInfo())
	if _, err := draftedItemInfoEntity.Update(ctx, tx, blackList); err != nil {
		return apperr.OfferItemInternalError.Wrap(err)
	}
	return nil
}

func (o *OfferItemRepositoryImpl) BulkGet(ctx context.Context, exec boil.ContextExecutor, ids []model.OfferItemID, isClosed bool) (map[model.OfferItemID]*model.OfferItem, error) {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.BulkGet")
	defer span.End()

	offerItems := make(map[model.OfferItemID]*model.OfferItem, len(ids))
	amebaIDsStr := make([]string, 0, len(ids))
	for _, id := range ids {
		amebaIDsStr = append(amebaIDsStr, id.String())
	}

	queries := make([]qm.QueryMod, 0, 2)
	queries = append(queries, entity.OfferItemWhere.ID.IN(amebaIDsStr))
	if !isClosed {
		queries = append(queries, entity.OfferItemWhere.IsClosed.EQ(false))
	}
	offerItemEntities, err := entity.OfferItems(queries...).All(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("entity.OfferItems.All: %w", err)
	}
	for _, offerItemEntity := range offerItemEntities {
		offerItemID := model.OfferItemID(offerItemEntity.ID)
		scheduleEntities, err := entity.Schedules(entity.ScheduleWhere.OfferItemID.EQ(offerItemID.String())).All(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("entity.Schedules.All: %w", err)
		}
		schedules := make(model.ScheduleList, 0, len(scheduleEntities))
		for _, scheduleEntity := range scheduleEntities {
			schedule := converter.ScheduleEntityToModel(scheduleEntity)
			schedules = append(schedules, schedule)
		}

		draftedItemInfoEntity, err := entity.DraftedItemInfos(entity.DraftedItemInfoWhere.OfferItemID.EQ(offerItemID.String())).One(ctx, exec)
		if err != nil {
			return nil, fmt.Errorf("entity.DraftedItemInfos.One: %w", err)
		}
		draftedItemInfo, err := converter.ConvertDraftedItemToModel(draftedItemInfoEntity)
		if err != nil {
			return nil, fmt.Errorf("converter.ConvertDraftedItemToModel: %w", err)
		}

		offerItem, err := converter.ConvertOfferItemToModel(offerItemEntity, schedules, draftedItemInfo)
		if err != nil {
			return nil, fmt.Errorf("converter.OfferItemEntityToModel: %w", err)
		}
		offerItems[offerItemID] = offerItem
	}
	return offerItems, nil
}

// 例：sinceEndDateが2024年1月1日、untilEndDateが2024年1月7日の場合、endDateが2024年1月1日から2024年1月7日の期間内に設定してあるOfferItemのIDを取得する
func (o *OfferItemRepositoryImpl) ListIDsByEndDate(ctx context.Context, exec boil.ContextExecutor, sinceEndDate, untilEndDate time.Time) (model.OfferItemIDList, error) {
	ctx, span := trace.StartSpan(ctx, "OfferItemRepository.ListIDsByEndDate")
	defer span.End()

	schedules, err := entity.Schedules(
		qm.Where(entity.ScheduleColumns.ScheduleType+" = ?", model.ScheduleTypeInvitation.Int()),
		qm.Where(entity.ScheduleColumns.EndDate+" >= ?", sinceEndDate),
		qm.Where(entity.ScheduleColumns.EndDate+" <= ?", untilEndDate),
	).All(ctx, exec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OfferItemIDList{}, nil
		}
		return nil, apperr.OfferItemInternalError.Wrap(err)
	}

	offerItemIDs := make(model.OfferItemIDList, 0, len(schedules))
	for _, schedule := range schedules {
		offerItemIDs = append(offerItemIDs, model.OfferItemID(schedule.OfferItemID))
	}
	return offerItemIDs, nil
}
