package utils

import (
	"fmt"
	"strings"
)

func MakeInsertQuery(tableName string, fields ...string) string {
	var leftPart strings.Builder
	var rightPart strings.Builder

	leftPart.WriteString(fmt.Sprintf("INSERT INTO %s(", tableName))
	rightPart.WriteString("VALUES(")

	for i, field := range fields {
		leftPart.WriteString(field)
		rightPart.WriteString(fmt.Sprintf("$%d", i+1))

		if i != len(fields)-1 {
			leftPart.WriteString(",")
			rightPart.WriteString(",")
		}
	}

	leftPart.WriteString(") ")
	rightPart.WriteString(")")

	leftPart.WriteString(rightPart.String())

	return leftPart.String()
}

func MakeInsertQueryReturningModel(tableName string, fields ...string) string {
	var queryBuilder strings.Builder

	queryBuilder.WriteString(MakeInsertQuery(tableName, fields...))
	queryBuilder.WriteString(" RETURNING id,")

	for i, field := range fields {
		queryBuilder.WriteString(field)

		if i != len(fields)-1 {
			queryBuilder.WriteString(",")
		}
	}

	return queryBuilder.String()
}

func MakeInsertQueryFromString(tableName string, fieldsString string) string {
	fields := strings.Split(fieldsString, ",")

	return MakeInsertQuery(tableName, fields...)
}
