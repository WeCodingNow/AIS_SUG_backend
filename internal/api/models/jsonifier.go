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
	BacklogT
)

func withDontWant(oldRefs JSONRefTable, types ...JSONModelType) JSONRefTable {
	if oldRefs == nil {
		oldRefs = make(JSONRefTable)
	}

	newRefs := make(JSONRefTable, len(oldRefs)+len(types))

	for t, ref := range oldRefs {
		newRefs[t] = ref
	}

	for _, t := range types {
		newRefs[t] = true
	}

	return newRefs
}
