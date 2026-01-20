package lit

import (
	"strconv"
	"strings"
)

type PgInsertUpdateQueryGenerator struct{}

func (PgInsertUpdateQueryGenerator) GenerateInsertQuery(tableName string, columnKeys []string, hasIntId bool) (string, []string) {
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

	counter := 1
	insertColumns := []string{}
	for i, k := range columnKeys {
		if hasIntId && k == "id" {
			insertQuery.WriteString("DEFAULT")
		} else {
			insertColumns = append(insertColumns, k)
			insertQuery.WriteString("$" + strconv.Itoa(counter))
			counter++
		}
		if i != totalKeys-1 {
			insertQuery.WriteString(",")
		}
	}
	insertQuery.WriteString(") RETURNING id")

	return insertQuery.String(), insertColumns
}

func (PgInsertUpdateQueryGenerator) GenerateUpdateQuery(tableName string, columnKeys []string) string {
	var updateQuery strings.Builder
	updateQuery.WriteString("UPDATE ")
	updateQuery.WriteString(tableName)
	updateQuery.WriteString(" SET ")

	totalKeys := len(columnKeys)
	for i, k := range columnKeys {
		updateQuery.WriteString(k)
		updateQuery.WriteString(" = $" + strconv.Itoa(i+1))
		if i != totalKeys-1 {
			updateQuery.WriteString(",")
		}
	}

	updateQuery.WriteString(" WHERE ")

	return updateQuery.String()
}

func pgRenumberPlaceholders(where string, offset int) string {
	if !strings.Contains(where, "$") {
		return where
	}

	var newWhere strings.Builder
	parsingIdentifier := false

	for _, c := range where {
		if c == '$' {
			parsingIdentifier = true
			newWhere.WriteRune(c)
		} else if parsingIdentifier {
			if c >= '0' && c <= '9' {
				continue
			} else {
				parsingIdentifier = false
				offset++
				newWhere.WriteString(strconv.Itoa(offset))
				newWhere.WriteRune(c)
			}
		} else {
			newWhere.WriteRune(c)
		}
	}
	if parsingIdentifier {
		offset++
		newWhere.WriteString(strconv.Itoa(offset))
	}

	return newWhere.String()
}

func pgJoinStringForIn(offset int, count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString("$" + strconv.Itoa(i+1+offset))
		if i < count-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
