package postgres

import (
	"fmt"
	"strings"
)

func MakeJoinQuery(leftTable, leftFieldsString, leftKey, rightTable, rightKey, rightID string) string {
	leftFields := strings.Split(leftFieldsString, ",")

	for i := range leftFields {
		leftFields[i] = fmt.Sprintf("%s.%s", leftTable, leftFields[i])
	}
	leftFieldsString = strings.Join(leftFields, ",")

	return fmt.Sprintf(
		"SELECT %s FROM %s JOIN %s ON %s.%s = %s.%s WHERE %s.%s = $1",
		leftFieldsString, leftTable, rightTable, leftTable, leftKey, rightTable, rightKey, rightTable, rightID)
}
