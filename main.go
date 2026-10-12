package main

import (
	"gradebook/console"
	"gradebook/gradebook"
)

func main() {
	book := gradebook.NewGradeBook()
	app := console.NewApp(book)

	app.MainMenu()
}
