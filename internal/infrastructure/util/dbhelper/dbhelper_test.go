package dbhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testEntityColumns = struct {
	ID        string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func TestCreateOrderByClause(t *testing.T) {
	sorts := func() []*model.Sort {
		list := make([]*model.Sort, 0, 3)
		sort, _ := model.NewSort("id", false)
		list = append(list, sort)
		sort, _ = model.NewSort("created_at", true)
		list = append(list, sort)
		sort, _ = model.NewSort("updated_at", false)
		list = append(list, sort)
		return list
	}()
	expected := "id, created_at desc, updated_at"
	result := CreateOrderByClause(testEntityColumns, sorts)
	assert.Equal(t, expected, result)

	// 存在しないカラムを指定されても結果が変わらないことを確認
	sorts = append(sorts,
		func() *model.Sort {
			sort, _ := model.NewSort("test", true)
			return sort
		}(),
	)
	result = CreateOrderByClause(testEntityColumns, sorts)
	assert.Equal(t, expected, result)
}

func TestGetColumnsFromEntityColumns(t *testing.T) {
	expected := []string{"id", "created_at", "updated_at"}
	assert.EqualValues(t, expected, GetColumnsFromEntityColumns(testEntityColumns))
}
