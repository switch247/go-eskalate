package main

import (
	"fmt"
)

type subject struct {
	name  string
	marks float64
	grade string
}

var subjects = []subject{}
var student string
var numberOfSubjects int
var DB = map[string][]subject{}
var choice string

func main() {
	for {
		fmt.Println("Welcome to Grade Calculator!")
		validateInput("Enter your name:", &student, "")

		fmt.Println("Hello, " + student + "!")

		validateInput("Enter the number of subjects you take:", &numberOfSubjects, "")

		getSubjects()
		fmt.Printf("Subjects: %v\n", subjects)
		var averageGrade = calculateAverage(subjects)
		fmt.Printf("Average Grade: %v\n", calculateGrade(averageGrade))
		DB[student] = append(DB[student], subjects...)
		subjects = []subject{}
		fmt.Println("=====================================")

		validateInput("Enter 'exit' to quit or any other key to continue:", &choice, "")
		if choice == "exit" {
			break
		}
	}
}

func getSubjects() {
	for i := 0; i < numberOfSubjects; i++ {
		var subjectName string
		var subjectMarks float64
		validateInput("Enter the name of the subject:", &subjectName, "")
		subjectMarks = validateMarks("Enter the marks you scored in " + subjectName + ":")

		subjects = append(subjects, subject{name: subjectName, marks: subjectMarks, grade: calculateGrade(subjectMarks)})
	}
}

func validateInput(prompt string, input interface{}, errMsg string) {
	for {
		fmt.Println(prompt)
		_, err := fmt.Scanln(input)

		if err == nil {
			break
		}
		fmt.Printf("Invalid input. Please enter a valid %T.\n", input)
		if errMsg != "" {
			fmt.Println(errMsg)
		}
	}
}

func validateMarks(prompt string) float64 {
	for {
		var input float64
		fmt.Print(prompt)
		_, err := fmt.Scanln(&input)
		if err == nil && input >= 0 && input <= 100 {
			return input
		}
		fmt.Println("Invalid input. Please enter a valid float64 between 0-100.")
	}
}

func calculateGrade(marks float64) string {
	switch {
	case marks >= 90:
		return "A+"
	case marks >= 85:
		return "A"
	case marks >= 80:
		return "A-"
	case marks >= 75:
		return "B+"
	case marks >= 70:
		return "B"
	case marks >= 65:
		return "B-"
	case marks >= 60:
		return "C+"
	case marks >= 50:
		return "C"
	case marks >= 45:
		return "D"
	default:
		return "F"
	}
}

func calculateAverage(x []subject) float64 {
	total := 0.0
	for _, value := range x {
		total += value.marks
	}
	return total / float64(len(x))
}
