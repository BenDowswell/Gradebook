package console

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"gradebook/gradebook"
)

type App struct {
	book   *gradebook.GradeBook
	reader *bufio.Scanner
	out    io.Writer
}

func NewApp(book *gradebook.GradeBook) *App {
	return &App{
		book:   book,
		reader: bufio.NewScanner(os.Stdin),
		out:    os.Stdout,
	}
}

func (a *App) readLine(prompt string) (string, bool) {
	fmt.Fprint(a.out, prompt)

	if !a.reader.Scan() {
		return "", false
	}

	return strings.TrimSpace(a.reader.Text()), true
}

func (a *App) MainMenu() {
	for {
		fmt.Fprintln(a.out, "\n=== Gradebook ===")
		fmt.Fprintln(a.out, "1) Add")
		fmt.Fprintln(a.out, "2) Delete")
		fmt.Fprintln(a.out, "3) List")
		fmt.Fprintln(a.out, "4) Lookup")
		fmt.Fprintln(a.out, "x) Exit")

		choice, ok := a.readLine("Choose an option: ")
		if !ok {
			return
		}

		switch strings.ToLower(choice) {
		case "1":
			a.AddMenu()
		case "2":
			a.DeleteMenu()
		case "3":
			a.ListMenu()
		case "4":
			a.LookupMenu()
		case "x":
			fmt.Fprintln(a.out, "Goodbye!")
			return
		default:
			fmt.Fprintln(a.out, "Not an option, try again.")
		}
	}
}

func (a *App) AddMenu() {
	for {
		fmt.Fprintln(a.out, "\n=== Add ===")
		fmt.Fprintln(a.out, "1) Add student")
		fmt.Fprintln(a.out, "2) Add subject")
		fmt.Fprintln(a.out, "3) Add grade")
		fmt.Fprintln(a.out, "x) Back")

		choice, ok := a.readLine("Choose an option: ")
		if !ok {
			return
		}

		switch strings.ToLower(choice) {
		case "1":
			a.addStudent()
		case "2":
			a.addSubject()
		case "3":
			a.addGrade()
		case "x":
			return
		default:
			fmt.Fprintln(a.out, "Not an option, try again.")
		}
	}
}

func (a *App) DeleteMenu() {
	for {
		fmt.Fprintln(a.out, "\n=== Delete ===")
		fmt.Fprintln(a.out, "1) Delete student")
		fmt.Fprintln(a.out, "2) Delete subject")
		fmt.Fprintln(a.out, "x) Back")

		choice, ok := a.readLine("Choose an option: ")
		if !ok {
			return
		}

		switch strings.ToLower(choice) {
		case "1":
			a.deleteStudent()
		case "2":
			a.deleteSubject()
		case "x":
			return
		default:
			fmt.Fprintln(a.out, "Not an option, try again.")
		}
	}
}

func (a *App) ListMenu() {
	for {
		fmt.Fprintln(a.out, "\n=== List ===")
		fmt.Fprintln(a.out, "1) List students")
		fmt.Fprintln(a.out, "2) List subjects")
		fmt.Fprintln(a.out, "3) List student grades")
		fmt.Fprintln(a.out, "x) Back")

		choice, ok := a.readLine("Choose an option: ")
		if !ok {
			return
		}

		switch strings.ToLower(choice) {
		case "1":
			a.listStudents()
		case "2":
			a.listSubjects()
		case "3":
			a.listStudentGrades()
		case "x":
			return
		default:
			fmt.Fprintln(a.out, "Not an option, try again.")
		}
	}
}

func (a *App) LookupMenu() {
	for {
		fmt.Fprintln(a.out, "\n=== Lookup ===")
		fmt.Fprintln(a.out, "1) Get student")
		fmt.Fprintln(a.out, "2) Get subject")
		fmt.Fprintln(a.out, "x) Back")

		choice, ok := a.readLine("Choose an option: ")
		if !ok {
			return
		}

		switch strings.ToLower(choice) {
		case "1":
			a.getStudent()
		case "2":
			a.getSubject()
		case "x":
			return
		default:
			fmt.Fprintln(a.out, "Not an option, try again.")
		}
	}
}

func (a *App) addStudent() {
	name, ok := a.readRequired("Student name: ")
	if !ok {
		return
	}

	school, ok := a.readRequired("School: ")
	if !ok {
		return
	}

	id, err := a.book.AddStudent(name, school)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	student, err := a.book.GetStudent(id)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintf(a.out, "\nStudent added! ID: %d\n", id)
	fmt.Fprintf(a.out, "Name: %s\nSchool: %s\n",
		student.Name, student.School)
}

func (a *App) addSubject() {
	name, ok := a.readRequired("Subject name: ")
	if !ok {
		return
	}

	id, err := a.book.AddSubject(name)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	subject, err := a.book.GetSubject(id)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintf(a.out, "\nSubject added! ID: %d\n", id)
	fmt.Fprintf(a.out, "Name: %s\n", subject.Name)
}

func (a *App) addGrade() {
	studentID, ok := a.readInt("Student ID: ")
	if !ok {
		return
	}

	subjectID, ok := a.readInt("Subject ID: ")
	if !ok {
		return
	}

	grade, ok := a.readFloat("Grade (0-100): ")
	if !ok {
		return
	}

	if err := a.book.AddGrade(studentID, subjectID, grade); err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	savedGrade, err := a.book.GetGrade(studentID, subjectID)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintf(a.out,
		"Grade saved! Student %d, Subject %d: %.2f\n",
		studentID, subjectID, savedGrade)
}

func (a *App) deleteStudent() {
	id, ok := a.readInt("Student ID to delete: ")
	if !ok {
		return
	}

	if err := a.book.DeleteStudent(id); err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintln(a.out, "Student deleted successfully.")
}

func (a *App) deleteSubject() {
	id, ok := a.readInt("Subject ID to delete: ")
	if !ok {
		return
	}

	if err := a.book.DeleteSubject(id); err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintln(a.out, "Subject deleted successfully.")
}

func (a *App) listStudents() {
	students := a.book.ListStudents()

	if len(students) == 0 {
		fmt.Fprintln(a.out, "No students found.")
		return
	}

	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tName\tSchool")
	fmt.Fprintln(w, "--\t----\t------")

	for _, student := range students {
		fmt.Fprintf(w, "%d\t%s\t%s\n",
			student.StudentID,
			student.Name,
			student.School,
		)
	}

	w.Flush()
}

func (a *App) listSubjects() {
	subjects := a.book.ListSubjects()

	if len(subjects) == 0 {
		fmt.Fprintln(a.out, "No subjects found.")
		return
	}

	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tSubject")
	fmt.Fprintln(w, "--\t-------")

	for _, subject := range subjects {
		fmt.Fprintf(w, "%d\t%s\n",
			subject.SubjectID,
			subject.Name,
		)
	}

	w.Flush()
}

func (a *App) listStudentGrades() {
	id, ok := a.readInt("Student ID: ")
	if !ok {
		return
	}

	grades, err := a.book.ListGradesForStudent(id)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	if len(grades) == 0 {
		fmt.Fprintln(a.out, "No grades found for this student.")
		return
	}

	w := tabwriter.NewWriter(a.out, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "Subject ID\tSubject\tGrade")
	fmt.Fprintln(w, "----------\t-------\t-----")

	for _, grade := range grades {
		fmt.Fprintf(w, "%d\t%s\t%.2f\n",
			grade.SubjectID,
			grade.SubjectName,
			grade.Grade,
		)
	}

	w.Flush()
}

func (a *App) getStudent() {
	id, ok := a.readInt("Student ID: ")
	if !ok {
		return
	}

	student, err := a.book.GetStudent(id)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintf(a.out, "\nStudent ID: %d\n", student.StudentID)
	fmt.Fprintf(a.out, "Name: %s\n", student.Name)
	fmt.Fprintf(a.out, "School: %s\n", student.School)
}

func (a *App) getSubject() {
	id, ok := a.readInt("Subject ID: ")
	if !ok {
		return
	}

	subject, err := a.book.GetSubject(id)
	if err != nil {
		fmt.Fprintln(a.out, "Error:", err)
		return
	}

	fmt.Fprintf(a.out, "\nSubject ID: %d\n", subject.SubjectID)
	fmt.Fprintf(a.out, "Name: %s\n", subject.Name)
}

func (a *App) readInt(prompt string) (int, bool) {
	for {
		input, ok := a.readLine(prompt)
		if !ok {
			return 0, false
		}

		value, err := strconv.Atoi(input)
		if err == nil {
			return value, true
		}

		fmt.Fprintln(a.out, "Please enter a whole number.")
	}
}

func (a *App) readFloat(prompt string) (float64, bool) {
	for {
		input, ok := a.readLine(prompt)
		if !ok {
			return 0, false
		}

		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Fprintln(a.out, "Please enter a valid number.")
			continue
		}

		if math.IsNaN(value) || math.IsInf(value, 0) ||
			value < 0 || value > 100 {
			fmt.Fprintln(a.out, "Grade must be a number between 0 and 100.")
			continue
		}

		return value, true
	}
}

func (a *App) readRequired(prompt string) (string, bool) {
	for {
		value, ok := a.readLine(prompt)
		if !ok {
			return "", false
		}

		if value != "" {
			return value, true
		}

		fmt.Fprintln(a.out, "This field cannot be empty. Try again.")
	}
}
