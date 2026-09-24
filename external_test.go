package gradebook_test

import (
	"errors"
	"gradebook"
	"testing"
)

func TestAddStudent(t *testing.T) {
	tests := []struct {
		name        string
		studentName string
		schoolName  string
		expectedErr error
	}{
		{"Valid student", "Alice", "Springfield High", nil},
		{"Empty student name", "", "Springfield High", gradebook.ErrEmptyField},
		{"Empty school name", "Bob", "", gradebook.ErrEmptyField},
		{"Both empty", "", "", gradebook.ErrEmptyField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gradebook.NewGradeBook()
			id, err := g.AddStudent(tt.studentName, tt.schoolName)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if err == nil {
				if id != 1 {
					t.Errorf("expected first student ID to be 1, got %d", id)
				}
				if g.ListStudents()[0].Name != tt.studentName {
					t.Errorf("expected name %s, got %s", tt.studentName, g.ListStudents()[0].Name)
				}
			}
		})
	}
}
