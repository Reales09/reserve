package mapper

import (
	"central_reserve/internal/domain/entities"
)

func TableToTable(table entities.Table) entities.Table {
	return table
}

func TableSliceToTableSlice(tables []entities.Table) []entities.Table {
	return tables
}
