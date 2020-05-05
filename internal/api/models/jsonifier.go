package models

type JSONModelType int
type JSONRefTable = map[JSONModelType]bool
type JSONMap = map[string]interface{}

type RecursiveJSONifiable interface {
	GetRefType() JSONModelType
	GetJsonMap(refs JSONRefTable) JSONMap
}

const (
	CathedraT JSONModelType = iota + 1
	SemesterT
	GroupT
	StudentT
	ControlEventT
	ControlEventTypeT
	MarkT
	ContactTypeT
	ContactT
	DisciplineT
	ResidenceT
)
