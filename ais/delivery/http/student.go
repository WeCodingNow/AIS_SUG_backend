package http

import (
	"net/http"
	"strconv"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type Student struct {
	ID int `json:"id"`
	// *Group    `json:"group"`
	Name       string  `json:"name"`
	SecondName string  `json:"second_name"`
	ThirdName  *string `json:"third_name"`
	// Residence  *ShortResidence `json:"residence"`
	Contacts []*FilteredContact `json:"contacts"`
}

type FilteredContact struct {
	ID          int          `json:"id"`
	ContactType *ContactType `json:"contact_type"`
	Def         string       `json:"def"`
}

func filterContact(contact *models.Contact) *FilteredContact {
	return &FilteredContact{
		contact.ID,
		toJsonContactType(contact.ContactType),
		contact.Def,
	}
}

type ShortStudent struct {
	ID         int     `json:"id"`
	GroupID    int     `json:"group_id"`
	Name       string  `json:"name"`
	SecondName string  `json:"second_name"`
	ThirdName  *string `json:"third_name"`
}

func toJsonShortStudent(student *models.Student) *ShortStudent {
	return &ShortStudent{
		ID:         student.ID,
		GroupID:    student.Group.ID,
		Name:       student.Name,
		SecondName: student.SecondName,
		ThirdName:  student.ThirdName,
	}
}

func toJsonStudent(student *models.Student) *Student {
	contactJSONs := make([]*FilteredContact, 0, len(student.Contacts))

	for _, contact := range student.Contacts {
		contactJSONs = append(contactJSONs, filterContact(contact))
	}

	// residenceJSON := toJsonShortResidence(student.Residence)

	return &Student{
		ID: student.ID,
		// Group: toJsonShortGroup(student.Group)
		// GroupID:    student.Group.ID,
		Name:       student.Name,
		SecondName: student.SecondName,
		ThirdName:  student.ThirdName,
		// Residence:  residenceJSON,
		Contacts: contactJSONs,
	}
}

func (h *Handler) GetStudent(c echo.Context) error {
	studentIDParam := c.Param("id")
	studentID, err := strconv.Atoi(studentIDParam)

	if err != nil {
		return err
	}

	studentModel, err := h.useCase.GetStudent(c.Request().Context(), studentID)

	if err != nil {
		if err == ais.ErrStudentNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonStudent(studentModel))
}

type manyStudentsOutput struct {
	Students []*Student `json:"students"`
}

func (h *Handler) GetAllStudents(c echo.Context) error {
	studentModels, err := h.useCase.GetAllStudents(c.Request().Context())

	if err != nil {
		return err
	}

	students := make([]*Student, 0, len(studentModels))
	for _, studentModel := range studentModels {
		students = append(students, toJsonStudent(studentModel))
	}

	return c.JSON(http.StatusOK, manyStudentsOutput{Students: students})
}
