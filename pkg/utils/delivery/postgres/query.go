package postgres

import (
	"fmt"
	"strings"
)

func prefixFields(table, fieldsString string) string {
	fields := strings.Split(fieldsString, ",")

	for i := range fields {
		fields[i] = fmt.Sprintf("%s.%s", table, fields[i])
	}

	return strings.Join(fields, ",")
}

func MakeJoinQuery(leftTable, leftFields, leftKey, rightTable, rightKey, rightID string) string {
	leftFields = prefixFields(leftTable, leftFields)
	// leftFields := strings.Split(leftFieldsString, ",")

	// for i := range leftFields {
	// 	leftFields[i] = fmt.Sprintf("%s.%s", leftTable, leftFields[i])
	// }
	// leftFieldsString = strings.Join(leftFields, ",")

	return fmt.Sprintf(
		"SELECT %s FROM %s JOIN %s ON %s.%s = %s.%s WHERE %s.%s = $1",
		leftFields, leftTable, rightTable, leftTable, leftKey, rightTable, rightKey, rightTable, rightID,
	)
}

// SELECT Группа.id,Группа.номер FROM Группа__Семестр
// 	JOIN Группа  ON Группа.id  = Группа__Семестр.id_группы
// 	JOIN Семестр ON Семестр.id = Группа__Семестр.id_семестра
// 	WHERE Группа.id = 1;

// returns left table rows
func MakeManyToManyJoinQuery(
	leftTable, leftFields, leftKey, leftMiddleKey,
	rightTable, rightKey, rightMiddleKey,
	middleTable string,
) string {
	leftFields = prefixFields(leftTable, leftFields)

	return fmt.Sprintf(
		`SELECT %s FROM %s
			JOIN %s ON %s.%s = %s.%s
			JOIN %s ON %s.%s = %s.%s
			WHERE %s.%s = $1`,
		leftFields, middleTable,
		leftTable, leftTable, leftKey, middleTable, leftMiddleKey,
		rightTable, rightTable, rightKey, middleTable, rightMiddleKey,
		leftTable, leftKey,
	)
}
