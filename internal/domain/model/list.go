package model

import "errors"

// リスト取得結果
type ListResult struct {
	// 取得データ数
	count int
	// データ総数
	totalCount int
}

func NewListResult(count, totalCount int) (*ListResult, error) {
	if count < 0 {
		return nil, errors.New("Count should not be less than 0.")
	}
	if totalCount < 0 {
		return nil, errors.New("Total Count should not be less than 0.")
	}

	return &ListResult{
		count:      count,
		totalCount: totalCount,
	}, nil
}

func (l *ListResult) Count() int {
	return l.count
}

func (l *ListResult) TotalCount() int {
	return l.totalCount
}

// リスト取得条件
type ListCondition struct {
	// 読み飛ばしデータ数
	offset int
	// データ取得上限数
	limit int
	// ソート設定リスト
	sorts []*Sort
}

func NewListCondition(offset, limit int, sorts []*Sort) (*ListCondition, error) {
	if offset < 0 {
		return nil, errors.New("Offset should not be less than 0.")
	}
	if limit < 0 {
		return nil, errors.New("Limit should not be less than 0.")
	}

	return &ListCondition{
		offset: offset,
		limit:  limit,
		sorts:  sorts,
	}, nil
}

func (l *ListCondition) Offset() int {
	return l.offset
}

func (l *ListCondition) SetOffset(offset int) {
	if offset >= 0 {
		l.offset = offset
	}
}

func (l *ListCondition) Limit() int {
	return l.limit
}

func (l *ListCondition) Sorts() []*Sort {
	return l.sorts
}

// ソート設定
type Sort struct {
	// カラム名
	orderBy string
	// 降順フラグ
	desc bool
}

func NewSort(orderBy string, desc bool) (*Sort, error) {
	if len(orderBy) == 0 {
		return nil, errors.New("Order By should not be empty.")
	}

	return &Sort{
		orderBy: orderBy,
		desc:    desc,
	}, nil
}

func (s *Sort) OrderBy() string {
	return s.orderBy
}

func (s *Sort) Desc() bool {
	return s.desc
}
