package gradebook

import (
	"errors"
	"math"
	"testing"
)

func TestNewGradeBook(t *testing.T) {
	g := NewGradeBook()

	if g.Students == nil {
		t.Error("expected Students map to be initialized, got nil")
	}
	if g.Subjects == nil {
		t.Error("expected Subjects map to be initialized, got nil")
	}
	if g.Grades == nil {
		t.Error("expected Grades map to be initialized, got nil")
	}
	if g.studentCounter != 1 {
		t.Errorf("expected studentCounter to start at 1, got %d", g.studentCounter)
	}
	if g.subjectCounter != 1 {
		t.Errorf("expected subjectCounter to start at 1, got %d", g.subjectCounter)
	}
}

func TestAddStudent(t *testing.T) {
	tests := []struct {
		name        string
		studentName string
		schoolName  string
		expectedErr error
	}{
		{"Valid student", "Alice", "Springfield High", nil},
		{"Empty student name", "", "Springfield High", ErrEmptyField},
		{"Empty school name", "Bob", "", ErrEmptyField},
		{"Both empty", "", "", ErrEmptyField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb := NewGradeBook()
			id, err := gb.AddStudent(tt.studentName, tt.schoolName)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil {
				if id != 1 {
					t.Errorf("expected first student ID to be 1, got %d", id)
				}
				if gb.Students[id].Name != tt.studentName {
					t.Errorf("expected name %s, got %s", tt.studentName, gb.Students[id].Name)
				}
			}
		})
	}
}

func TestAddSubject(t *testing.T) {
	tests := []struct {
		name        string
		subjectName string
		expectedErr error
	}{
		{"Valid subject", "Math", nil},
		{"Empty subject name", "", ErrEmptyField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb := NewGradeBook()
			id, err := gb.AddSubject(tt.subjectName)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil {
				if id != 1 {
					t.Errorf("expected first subject ID to be 1, got %d", id)
				}
				if gb.Subjects[id].Name != tt.subjectName {
					t.Errorf("expected subject name %s, got %s", tt.subjectName, gb.Subjects[id].Name)
				}
			}
		})
	}
}

func TestAddGrade(t *testing.T) {
	tests := []struct {
		name        string
		studentID   int
		subjectID   int
		grade       float64
		expectedErr error
	}{
		{"Valid grade", 1, 1, 95.5, nil},
		{"Invalid student ID", 99, 1, 95.5, ErrStudentID},
		{"Invalid subject ID", 1, 99, 95.5, ErrSubjectID},
		{"Grade too low", 1, 1, -1.0, ErrGradeValue},
		{"Grade too high", 1, 1, 101.0, ErrGradeValue},
		{"Grade is NaN", 1, 1, math.NaN(), ErrGradeValue},
		{"Grade exactly 0", 1, 1, 0.0, nil},
		{"Grade exactly 100", 1, 1, 100.0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb := NewGradeBook()
			// Setup required state
			_, _ = gb.AddStudent("Alice", "High School")
			_, _ = gb.AddSubject("Math")

			err := gb.AddGrade(tt.studentID, tt.subjectID, tt.grade)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil {
				if gb.Grades[tt.studentID][tt.subjectID] != tt.grade {
					t.Errorf("expected stored grade to be %f, got %f", tt.grade, gb.Grades[tt.studentID][tt.subjectID])
				}
			}
		})
	}
}

func TestGetGrade(t *testing.T) {
	tests := []struct {
		name          string
		reqStudentID  int
		reqSubjectID  int
		expectedGrade float64
		expectedErr   error
	}{
		{"Valid retrieval", 1, 1, 88.0, nil},
		{"Invalid student ID", 99, 1, 0, ErrStudentID},
		{"Invalid subject ID", 1, 99, 0, ErrSubjectID},
		{"Student has no grades at all", 2, 1, 0, ErrNoGradesForStudent},
		{"Student has no grade for subject", 1, 2, 0, ErrNoGradeForSubject},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gb := NewGradeBook()

			// Setup state: Student 1 has Math (Subject 1) grade, but no Science (Subject 2) grade.
			// Student 2 exists but has no grades.
			_, _ = gb.AddStudent("Alice", "High School") // ID 1
			_, _ = gb.AddStudent("Bob", "High School")   // ID 2
			_, _ = gb.AddSubject("Math")                 // ID 1
			_, _ = gb.AddSubject("Science")              // ID 2

			_ = gb.AddGrade(1, 1, 88.0)

			grade, err := gb.GetGrade(tt.reqStudentID, tt.reqSubjectID)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if grade != tt.expectedGrade {
				t.Errorf("expected grade %f, got %f", tt.expectedGrade, grade)
			}
		})
	}
}

func TestAddGrade_OverwriteBehavior(t *testing.T) {
	g := NewGradeBook()

	studentID, _ := g.AddStudent("Alice", "High School")
	subjectID, _ := g.AddSubject("Math")

	// 1. Add the initial grade
	err := g.AddGrade(studentID, subjectID, 85.0)
	if err != nil {
		t.Fatalf("unexpected error on first AddGrade: %v", err)
	}

	// 2. Add a new grade for the exact same student and subject
	err = g.AddGrade(studentID, subjectID, 95.5)
	if err != nil {
		t.Fatalf("unexpected error on second AddGrade (overwrite): %v", err)
	}

	// 3. Verify that the grade was successfully updated to the new value
	finalGrade, err := g.GetGrade(studentID, subjectID)
	if err != nil {
		t.Fatalf("unexpected error on GetGrade: %v", err)
	}

	if finalGrade != 95.5 {
		t.Errorf("expected grade to be overwritten to 95.5, but got %f", finalGrade)
	}
}

func TestAddStudent_IDSequencing(t *testing.T) {
	g := NewGradeBook()

	// Define the sequence of students we want to add
	students := []struct {
		name       string
		school     string
		expectedID int
	}{
		{"Alice", "Springfield High", 1},
		{"Bob", "Springfield High", 2},
		{"Charlie", "Shelbyville High", 3},
	}

	for _, s := range students {
		id, err := g.AddStudent(s.name, s.school)
		if err != nil {
			t.Fatalf("unexpected error adding student %s: %v", s.name, err)
		}

		// 1. Verify the returned ID matches the expected sequence
		if id != s.expectedID {
			t.Errorf("for student %s: expected ID %d, got %d", s.name, s.expectedID, id)
		}

		// 2. Verify the student was actually saved in the map under this new ID
		storedStudent, exists := g.Students[id]
		if !exists {
			t.Fatalf("expected student to be saved at ID %d, but nothing was found", id)
		}

		// 3. Verify the saved data matches
		if storedStudent.Name != s.name {
			t.Errorf("expected saved student name to be %s, got %s", s.name, storedStudent.Name)
		}
	}

	// 4. Finally, verify the internal counter state is ready for the 4th student
	if g.studentCounter != 4 {
		t.Errorf("expected internal studentCounter to be 4, got %d", g.studentCounter)
	}
}

func TestListStudents(t *testing.T) {
	g := NewGradeBook()

	_, _ = g.AddStudent("Charlie", "Shelbyville High") // ID 1
	_, _ = g.AddStudent("Alice", "Springfield High")   // ID 2
	_, _ = g.AddStudent("Bob", "Springfield High")     // ID 3

	students := g.ListStudents()

	if len(students) != 3 {
		t.Fatalf("expected 3 students, got %d", len(students))
	}

	if students[0].StudentID != 1 {
		t.Errorf("expected first student ID to be 1, got %d", students[0].StudentID)
	}

	if students[1].StudentID != 2 {
		t.Errorf("expected second student ID to be 2, got %d", students[1].StudentID)
	}

	if students[2].StudentID != 3 {
		t.Errorf("expected third student ID to be 3, got %d", students[2].StudentID)
	}
}

func TestListStudentsEmptyGradeBook(t *testing.T) {
	g := NewGradeBook()

	students := g.ListStudents()

	if len(students) != 0 {
		t.Errorf("expected 0 students, got %d", len(students))
	}
}

func TestListSubjects(t *testing.T) {
	g := NewGradeBook()

	_, _ = g.AddSubject("Science") // ID 1
	_, _ = g.AddSubject("Math")    // ID 2
	_, _ = g.AddSubject("English") // ID 3

	subjects := g.ListSubjects()

	if len(subjects) != 3 {
		t.Fatalf("expected 3 subjects, got %d", len(subjects))
	}

	if subjects[0].SubjectID != 1 {
		t.Errorf("expected first subject ID to be 1, got %d", subjects[0].SubjectID)
	}

	if subjects[1].SubjectID != 2 {
		t.Errorf("expected second subject ID to be 2, got %d", subjects[1].SubjectID)
	}

	if subjects[2].SubjectID != 3 {
		t.Errorf("expected third subject ID to be 3, got %d", subjects[2].SubjectID)
	}
}

func TestListSubjectsEmptyGradeBook(t *testing.T) {
	g := NewGradeBook()

	subjects := g.ListSubjects()

	if len(subjects) != 0 {
		t.Errorf("expected 0 subjects, got %d", len(subjects))
	}
}

func TestListGradesForStudent(t *testing.T) {
	tests := []struct {
		name          string
		studentID     int
		setupGrades   bool
		expectedCount int
		expectedErr   error
	}{
		{"Invalid student ID", 99, false, 0, ErrStudentID},
		{"Student has no grades", 1, false, 0, nil},
		{"Student has grades", 1, true, 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGradeBook()

			studentID, _ := g.AddStudent("Alice", "High School")
			mathID, _ := g.AddSubject("Math")
			scienceID, _ := g.AddSubject("Science")

			if tt.setupGrades {
				_ = g.AddGrade(studentID, scienceID, 92)
				_ = g.AddGrade(studentID, mathID, 88)
			}

			grades, err := g.ListGradesForStudent(tt.studentID)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if len(grades) != tt.expectedCount {
				t.Fatalf("expected %d grades, got %d", tt.expectedCount, len(grades))
			}

			if tt.setupGrades {
				if grades[0].SubjectID != mathID {
					t.Errorf("expected first subject ID to be %d, got %d", mathID, grades[0].SubjectID)
				}

				if grades[0].SubjectName != "Math" {
					t.Errorf("expected first subject name to be Math, got %s", grades[0].SubjectName)
				}

				if grades[0].Grade != 88 {
					t.Errorf("expected first grade to be 88, got %f", grades[0].Grade)
				}

				if grades[1].SubjectID != scienceID {
					t.Errorf("expected second subject ID to be %d, got %d", scienceID, grades[1].SubjectID)
				}

				if grades[1].SubjectName != "Science" {
					t.Errorf("expected second subject name to be Science, got %s", grades[1].SubjectName)
				}

				if grades[1].Grade != 92 {
					t.Errorf("expected second grade to be 92, got %f", grades[1].Grade)
				}
			}
		})
	}
}

func TestDeleteStudent(t *testing.T) {
	g := NewGradeBook()

	studentID, _ := g.AddStudent("Alice", "High School")
	subjectID, _ := g.AddSubject("Math")

	_ = g.AddGrade(studentID, subjectID, 88)

	err := g.DeleteStudent(studentID)
	if err != nil {
		t.Fatalf("unexpected error deleting student: %v", err)
	}

	if _, ok := g.Students[studentID]; ok {
		t.Errorf("expected student ID %d to be deleted", studentID)
	}

	if _, ok := g.Grades[studentID]; ok {
		t.Errorf("expected grades for student ID %d to be deleted", studentID)
	}
}

func TestDeleteStudentUnknownID(t *testing.T) {
	g := NewGradeBook()

	err := g.DeleteStudent(99)

	if !errors.Is(err, ErrStudentID) {
		t.Errorf("expected error %v, got %v", ErrStudentID, err)
	}
}

func TestDeleteSubjectRemovesSubjectAndGradeButKeepsOtherGrades(t *testing.T) {
	g := NewGradeBook()

	studentID, _ := g.AddStudent("Alice", "High School")
	mathID, _ := g.AddSubject("Math")
	scienceID, _ := g.AddSubject("Science")

	_ = g.AddGrade(studentID, mathID, 88)
	_ = g.AddGrade(studentID, scienceID, 92)

	err := g.DeleteSubject(mathID)
	if err != nil {
		t.Fatalf("unexpected error deleting subject: %v", err)
	}

	if _, ok := g.Subjects[mathID]; ok {
		t.Errorf("expected subject ID %d to be deleted", mathID)
	}

	studentGrades, ok := g.Grades[studentID]
	if !ok {
		t.Fatalf("expected student ID %d to still have grades", studentID)
	}

	if _, ok := studentGrades[mathID]; ok {
		t.Errorf("expected grade for deleted subject ID %d to be removed", mathID)
	}

	grade, ok := studentGrades[scienceID]
	if !ok {
		t.Fatalf("expected grade for subject ID %d to remain", scienceID)
	}

	if grade != 92 {
		t.Errorf("expected Science grade to remain 92, got %f", grade)
	}
}

func TestDeleteSubjectCleansUpEmptyStudentGrades(t *testing.T) {
	g := NewGradeBook()

	studentID, _ := g.AddStudent("Alice", "High School")
	subjectID, _ := g.AddSubject("Math")

	_ = g.AddGrade(studentID, subjectID, 88)

	err := g.DeleteSubject(subjectID)
	if err != nil {
		t.Fatalf("unexpected error deleting subject: %v", err)
	}

	if _, ok := g.Grades[studentID]; ok {
		t.Errorf("expected grades map for student ID %d to be removed", studentID)
	}
}

func TestDeleteSubjectUnknownID(t *testing.T) {
	g := NewGradeBook()

	err := g.DeleteSubject(99)

	if !errors.Is(err, ErrSubjectID) {
		t.Errorf("expected error %v, got %v", ErrSubjectID, err)
	}
}

func TestDeleteSubjectRemovesGradeFromAllStudents(t *testing.T) {
	g := NewGradeBook()

	aliceID, _ := g.AddStudent("Alice", "High School")
	bobID, _ := g.AddStudent("Bob", "High School")

	mathID, _ := g.AddSubject("Math")
	scienceID, _ := g.AddSubject("Science")

	_ = g.AddGrade(aliceID, mathID, 88)
	_ = g.AddGrade(aliceID, scienceID, 91)

	_ = g.AddGrade(bobID, mathID, 76)
	_ = g.AddGrade(bobID, scienceID, 84)

	err := g.DeleteSubject(mathID)
	if err != nil {
		t.Fatalf("unexpected error deleting subject: %v", err)
	}

	if _, ok := g.Subjects[mathID]; ok {
		t.Errorf("expected subject ID %d to be deleted", mathID)
	}

	aliceGrades := g.Grades[aliceID]

	if _, ok := aliceGrades[mathID]; ok {
		t.Errorf("expected Math grade to be removed for Alice")
	}

	if grade, ok := aliceGrades[scienceID]; !ok {
		t.Errorf("expected Alice's Science grade to remain")
	} else if grade != 91 {
		t.Errorf("expected Alice's Science grade to be 91, got %f", grade)
	}

	bobGrades := g.Grades[bobID]

	if _, ok := bobGrades[mathID]; ok {
		t.Errorf("expected Math grade to be removed for Bob")
	}

	if grade, ok := bobGrades[scienceID]; !ok {
		t.Errorf("expected Bob's Science grade to remain")
	} else if grade != 84 {
		t.Errorf("expected Bob's Science grade to be 84, got %f", grade)
	}
}

func TestDeleteStudentOnlyRemovesTargetStudent(t *testing.T) {
	g := NewGradeBook()

	aliceID, _ := g.AddStudent("Alice", "High School")
	bobID, _ := g.AddStudent("Bob", "High School")

	mathID, _ := g.AddSubject("Math")

	_ = g.AddGrade(aliceID, mathID, 88)
	_ = g.AddGrade(bobID, mathID, 92)

	err := g.DeleteStudent(aliceID)
	if err != nil {
		t.Fatalf("unexpected error deleting student: %v", err)
	}

	if _, ok := g.Students[aliceID]; ok {
		t.Errorf("expected Alice to be deleted")
	}

	if _, ok := g.Grades[aliceID]; ok {
		t.Errorf("expected Alice's grades to be deleted")
	}

	if _, ok := g.Students[bobID]; !ok {
		t.Errorf("expected Bob to remain")
	}

	bobGrades, ok := g.Grades[bobID]
	if !ok {
		t.Fatalf("expected Bob's grades to remain")
	}

	if grade, ok := bobGrades[mathID]; !ok {
		t.Errorf("expected Bob's Math grade to remain")
	} else if grade != 92 {
		t.Errorf("expected Bob's Math grade to be 92, got %f", grade)
	}
}
