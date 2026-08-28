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
