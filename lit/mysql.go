package lit

import (
	"strings"
)

type MySqlInsertUpdateQueryGenerator struct{}

func (MySqlInsertUpdateQueryGenerator) GenerateInsertQuery(tableName string, columnKeys []string, hasIntId bool) (string, []string) {
	var insertQuery strings.Builder

	insertQuery.WriteString("INSERT INTO ")
	insertQuery.WriteString(tableName)
	insertQuery.WriteString(" (")

	totalKeys := len(columnKeys)
	for i, k := range columnKeys {
		insertQuery.WriteString(k)
		if i != totalKeys-1 {
			insertQuery.WriteString(",")
		}
	}

	insertQuery.WriteString(") VALUES (")

	insertColumns := []string{}
	for i, k := range columnKeys {
		if hasIntId && k == "id" {
			insertQuery.WriteString("NULL")
		} else {
			insertColumns = append(insertColumns, k)
			insertQuery.WriteString("?")
		}
		if i != totalKeys-1 {
			insertQuery.WriteString(",")
		}
	}
	insertQuery.WriteString(")")

	return insertQuery.String(), insertColumns
}

func (MySqlInsertUpdateQueryGenerator) GenerateUpdateQuery(tableName string, columnKeys []string) string {
	var updateQuery strings.Builder
	updateQuery.WriteString("UPDATE ")
	updateQuery.WriteString(tableName)
	updateQuery.WriteString(" SET ")

	totalKeys := len(columnKeys)
	for i, k := range columnKeys {
		updateQuery.WriteString(k)
		updateQuery.WriteString(" = ?")
		if i != totalKeys-1 {
			updateQuery.WriteString(",")
		}
	}

	updateQuery.WriteString(" WHERE ")

	return updateQuery.String()
}

func mysqlJoinStringForIn(count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString("?")
		if i < count-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
