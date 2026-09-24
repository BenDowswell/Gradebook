package gradebook

import (
	"cmp"
	"errors"
	"math"
	"slices"
)

var (
	ErrEmptyField         = errors.New("empty field passed in")
	ErrStudentID          = errors.New("student ID doesn't exist")
	ErrSubjectID          = errors.New("subject ID doesn't exist")
	ErrGradeValue         = errors.New("grade must be a number between 0 and 100")
	ErrNoGradesForStudent = errors.New("no grades exist for this student yet")
	ErrNoGradeForSubject  = errors.New("grade for this subject doesn't exist for this student")
)

type Student struct {
	StudentID int
	Name      string
	School    string
}

type Subject struct {
	SubjectID int
	Name      string
}

type GradeBook struct {
	students       map[int]Student
	subjects       map[int]Subject
	grades         map[int]map[int]float64
	studentCounter int
	subjectCounter int
}

type StudentGrade struct {
	SubjectID   int
	SubjectName string
	Grade       float64
}

func NewGradeBook() *GradeBook {
	g := GradeBook{
		students:       make(map[int]Student),
		subjects:       make(map[int]Subject),
		grades:         make(map[int]map[int]float64),
		studentCounter: 1,
		subjectCounter: 1,
	}

	return &g
}

func (g *GradeBook) validateIDs(studentID, subjectID int) error {
	if _, ok := g.students[studentID]; !ok {
		return ErrStudentID
	}
	if _, ok := g.subjects[subjectID]; !ok {
		return ErrSubjectID
	}
	return nil
}

func (g *GradeBook) AddStudent(name, school string) (int, error) {
	if name == "" || school == "" {
		return 0, ErrEmptyField
	}
	id := g.studentCounter

	newStudent := Student{
		StudentID: id,
		Name:      name,
		School:    school,
	}

	g.students[id] = newStudent
	g.studentCounter++

	return id, nil
}

func (g *GradeBook) AddSubject(name string) (int, error) {
	if name == "" {
		return 0, ErrEmptyField
	}
	id := g.subjectCounter

	newSubject := Subject{
		SubjectID: id,
		Name:      name,
	}

	g.subjects[id] = newSubject
	g.subjectCounter++

	return id, nil
}

func (g *GradeBook) AddGrade(studentID, subjectID int, grade float64) error {
	if err := g.validateIDs(studentID, subjectID); err != nil {
		return err
	}

	if math.IsNaN(grade) || (grade > 100 || grade < 0) {
		return ErrGradeValue
	}
	// look up and see if grades map exists yet
	studentGrades, ok := g.grades[studentID]
	if !ok {
		studentGrades = make(map[int]float64)
		g.grades[studentID] = studentGrades
	}

	studentGrades[subjectID] = grade
	return nil
}

func (g *GradeBook) GetGrade(studentID, subjectID int) (float64, error) {
	if err := g.validateIDs(studentID, subjectID); err != nil {
		return 0, err
	}

	studentGrades, ok := g.grades[studentID]
	if !ok {
		return 0, ErrNoGradesForStudent
	}

	grade, ok := studentGrades[subjectID]
	if !ok {
		return 0, ErrNoGradeForSubject
	}

	return grade, nil
}

func (g *GradeBook) ListStudents() []Student {
	students := make([]Student, 0, len(g.students))

	for _, student := range g.students {
		students = append(students, student)
	}

	slices.SortFunc(students, func(a, b Student) int {
		return cmp.Compare(a.StudentID, b.StudentID)
	})

	return students
}

func (g *GradeBook) ListSubjects() []Subject {
	subjects := make([]Subject, 0, len(g.subjects))

	for _, subject := range g.subjects {
		subjects = append(subjects, subject)
	}

	slices.SortFunc(subjects, func(a, b Subject) int {
		return cmp.Compare(a.SubjectID, b.SubjectID)
	})

	return subjects
}

func (g *GradeBook) ListGradesForStudent(studentID int) ([]StudentGrade, error) {
	if _, ok := g.students[studentID]; !ok {
		return nil, ErrStudentID
	}

	studentGrades := g.grades[studentID]

	grades := make([]StudentGrade, 0, len(studentGrades))

	for subjectID, grade := range studentGrades {
		subject, ok := g.subjects[subjectID]
		if !ok {
			continue
		}

		grades = append(grades, StudentGrade{
			SubjectID:   subjectID,
			SubjectName: subject.Name,
			Grade:       grade,
		})
	}

	slices.SortFunc(grades, func(a, b StudentGrade) int {
		return cmp.Compare(a.SubjectID, b.SubjectID)
	})

	return grades, nil
}

func (g *GradeBook) DeleteStudent(studentID int) error {
	if _, ok := g.students[studentID]; !ok {
		return ErrStudentID
	}

	delete(g.students, studentID)
	delete(g.grades, studentID)

	return nil
}

func (g *GradeBook) DeleteSubject(subjectID int) error {
	if _, ok := g.subjects[subjectID]; !ok {
		return ErrSubjectID
	}

	delete(g.subjects, subjectID)

	for studentID, studentGrades := range g.grades {
		delete(studentGrades, subjectID)

		if len(studentGrades) == 0 {
			delete(g.grades, studentID)
		}
	}

	return nil
}
