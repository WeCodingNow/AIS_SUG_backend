package models

// нужно превращать *Group в int(group_id) у хендлера чтобы избежать рекурсии (студент->группа->этот же студент->...)
type Student struct {
	ID         int
	Name       string
	SecondName string
	ThirdName  *string

	*Group
	*Residence
	Contacts []*Contact
	Marks    []*Mark
}
