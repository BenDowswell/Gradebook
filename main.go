package student

import "errors"

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
	Students map[int]Student
	Subjects map[int]Subject
	Grades   map[int]map[int]float64
}

func NewGradeBook() *GradeBook {
	students := make(map[int]Student)
	subjects := make(map[int]Subject)
	grades := make(map[int]map[int]float64)

	g := GradeBook{
		Students: students,
		Subjects: subjects,
		Grades:   grades,
	}

	return &g
}

func (g *GradeBook) AddStudent(student, school string) {
	// work out the next available id  come back to this later when we want to delete students potentially
	id := len(g.Students) + 1

	newStudent := Student{
		StudentID: id,
		Name:      student,
		School:    school,
	}

	g.Students[id] = newStudent
}

func (g *GradeBook) AddSubject(name string) {
	// work out the next availble id  might need to come back to this later if we decide to delete subjects unlikely tho
	id := len(g.Subjects) + 1

	newSubject := Subject{
		SubjectID: id,
		Name:      name,
	}

	g.Subjects[id] = newSubject
}

func (g *GradeBook) AddGrade(studentID, subjectID int, grade float64) error {
	_, ok := g.Students[studentID]
	if !ok {
		return errors.New("Student ID doesnt exist please create student")
	}

	_, ok = g.Subjects[subjectID]
	if !ok {
		return errors.New("Subject  ID doesnt exist please create subject")
	}

	if grade > 100 || grade < 0 {
		return errors.New("Grade is incorrect value must be between 0 and 100 ")
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
