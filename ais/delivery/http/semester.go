package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WeCodingNow/AIS_SUG_backend/ais"
	"github.com/WeCodingNow/AIS_SUG_backend/models"
	"github.com/labstack/echo"
)

type Semester struct {
	ID        int        `json:"id"`
	Number    int        `json:"number"`
	Beginning time.Time  `json:"beginning"`
	End       *time.Time `json:"end"`

	// ControlEvents []*SemesterShortControlEvent `json:"control_events"`
}

// type SemesterShortControlEvent struct {
// 	ID               int
// 	ControlEventType *ControlEventType        `json:"type"`
// 	Discipline       *SemesterShortDiscipline `json:"discipline"`
// 	Marks            []*SemesterShortMark     `json:"marks"`
// 	Date             time.Time                `json:"date"`
// }

type SemesterShortDiscipline struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Hours int    `json:"hours"`
}

type SemesterShortMark struct {
	ID      int                   `json:"id"`
	Student *SemesterShortStudent `json:"student"`
	Date    time.Time             `json:"date"`
	Value   int                   `json:"value"`
}

type SemesterShortStudent struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	SecondName string  `json:"second_name"`
	ThirdName  *string `json:"third_name"`

	Group *SemesterShortGroup `json:"group"`
}

type SemesterShortGroup struct {
	ID     int `json:"id"`
	Number int `json:"number"`
}

func toJsonSemester(s *models.Semester) *Semester {
	// jsonControlEvents := make([]*SemesterShortControlEvent, 0)
	// for _, controlEvent := range s.ControlEvents {
	// 	jsonControlEvents = append(jsonControlEvents, toJsonSemesterShortControlEvent(controlEvent))
	// }

	return &Semester{
		ID:        s.ID,
		Number:    s.Number,
		Beginning: s.Beginning,
		End:       s.End,

		// ControlEvents: jsonControlEvents,
	}
}

// func toJsonSemesterShortControlEvent(controlEvent *models.ControlEvent) *SemesterShortControlEvent {
// 	jsonShortMarks := make([]*SemesterShortMark, 0)
// 	for _, mark := range controlEvent.Marks {
// 		jsonShortMarks = append(jsonShortMarks, toJsonSemesterShortMark(mark))
// 	}

// 	return &SemesterShortControlEvent{
// 		ID:               controlEvent.ID,
// 		Marks:            jsonShortMarks,
// 		Date:             controlEvent.Date,
// 		ControlEventType: toJsonControlEventType(controlEvent.ControlEventType),
// 		Discipline:       toJsonSemesterShortDiscipline(controlEvent.Discipline),
// 	}
// }

func toJsonSemesterShortDiscipline(discipline *models.Discipline) *SemesterShortDiscipline {
	return &SemesterShortDiscipline{
		ID:    discipline.ID,
		Name:  discipline.Name,
		Hours: discipline.Hours,
	}
}

func toJsonSemesterShortMark(mark *models.Mark) *SemesterShortMark {
	return &SemesterShortMark{
		ID:      mark.ID,
		Date:    mark.Date,
		Value:   mark.Value,
		Student: toSemesterShortStudent(mark.Student),
	}
}

func toSemesterShortStudent(student *models.Student) *SemesterShortStudent {
	return &SemesterShortStudent{
		ID:         student.ID,
		Name:       student.Name,
		SecondName: student.SecondName,
		ThirdName:  student.ThirdName,

		Group: toJsonSemesterShortGroup(student.Group),
	}
}

func toJsonSemesterShortGroup(group *models.Group) *SemesterShortGroup {
	return &SemesterShortGroup{
		ID:     group.ID,
		Number: group.Number,
	}
}

func (h *Handler) GetSemester(c echo.Context) error {
	semesterIDParam := c.Param("id")
	semesterID, err := strconv.Atoi(semesterIDParam)

	if err != nil {
		return err
	}

	semesterModel, err := h.useCase.GetSemester(c.Request().Context(), semesterID)

	if err != nil {
		if err == ais.ErrSemesterNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return c.JSON(http.StatusOK, toJsonSemester(semesterModel))
}

type manySemestersOutput struct {
	Semesters []*Semester `json:"semesters"`
}

func (h *Handler) GetAllSemesters(c echo.Context) error {
	semesters, err := h.useCase.GetAllSemesters(c.Request().Context())

	if err != nil {
		return err
	}

	semesterJSONs := make([]*Semester, 0, len(semesters))
	for _, semesterModel := range semesters {
		semesterJSONs = append(semesterJSONs, toJsonSemester(semesterModel))
	}

	return c.JSON(http.StatusOK, manySemestersOutput{Semesters: semesterJSONs})
}
