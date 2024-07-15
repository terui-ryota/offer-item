package dbhelper

import (
	"reflect"

	"github.com/terui-ryota/offer-item/internal/domain/model"
)

// カラム情報とソート設定リストから、OrderBy 句を生成する
// カラム情報は、db/entity の {entity名}Columns が渡されることを想定する
func CreateOrderByClause(entityColumns interface{}, sorts []*model.Sort) (orderByClause string) {
	columns := GetColumnsFromEntityColumns(entityColumns)
	exists := func(list []string, target string) bool {
		for _, elem := range list {
			if elem == target {
				return true
			}
		}
		return false
	}

	for _, sort := range sorts {
		if exists(columns, sort.OrderBy()) {
			if orderByClause != "" {
				orderByClause = orderByClause + ", "
			}
			orderByClause = orderByClause + sort.OrderBy()
			if sort.Desc() {
				orderByClause = orderByClause + " desc"
			}
		}
	}
	return
}

// カラム情報から string スライスを生成する
func GetColumnsFromEntityColumns(entityColumns interface{}) (columns []string) {
	values := reflect.ValueOf(entityColumns)
	for i := 0; i < values.NumField(); i++ {
		value := values.Field(i)
		switch value.Kind() {
		case reflect.String:
			columns = append(columns, value.String())
		}
	}
	return
}
