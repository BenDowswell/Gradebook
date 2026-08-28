package gradebook

import (
	"errors"
	"math"
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
	Students       map[int]Student
	Subjects       map[int]Subject
	Grades         map[int]map[int]float64
	studentCounter int
	subjectCounter int
}

func NewGradeBook() *GradeBook {
	g := GradeBook{
		Students:       make(map[int]Student),
		Subjects:       make(map[int]Subject),
		Grades:         make(map[int]map[int]float64),
		studentCounter: 1,
		subjectCounter: 1,
	}

	return &g
}

func (g *GradeBook) validateIDs(studentID, subjectID int) error {
	if _, ok := g.Students[studentID]; !ok {
		return ErrStudentID
	}
	if _, ok := g.Subjects[subjectID]; !ok {
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

	g.Students[id] = newStudent
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

	g.Subjects[id] = newSubject
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
	studentGrades, ok := g.Grades[studentID]
	if !ok {
		studentGrades = make(map[int]float64)
		g.Grades[studentID] = studentGrades
	}

	studentGrades[subjectID] = grade
	return nil
}

func (g *GradeBook) GetGrade(studentID, subjectID int) (float64, error) {
	if err := g.validateIDs(studentID, subjectID); err != nil {
		return 0, err
	}

	studentGrades, ok := g.Grades[studentID]
	if !ok {
		return 0, ErrNoGradesForStudent
	}

	grade, ok := studentGrades[subjectID]
	if !ok {
		return 0, ErrNoGradeForSubject
	}

	return grade, nil
}
