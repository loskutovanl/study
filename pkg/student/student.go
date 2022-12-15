package student

import (
	"strconv"
	"strings"
)

// Student структура студента и ее методы
type Student struct {
	name  string
	age   int
	grade int
}

func NewStudent(studentInfo []string) *Student {
	name := studentInfo[0]
	age, _ := strconv.Atoi(studentInfo[1])
	grade, _ := strconv.Atoi(studentInfo[2])

	student := Student{
		name:  name,
		age:   age,
		grade: grade,
	}

	return &student
}

func IsStudentValid(studentInfo []string) bool {
	result := true

	if len(studentInfo) == 3 {
		_, err := strconv.Atoi(studentInfo[1])
		if err != nil {
			result = false
		}

		_, err = strconv.Atoi(studentInfo[2])
		if err != nil {
			result = false
		}
	} else {
		result = false
	}

	return result
}

func (s *Student) PrintStudent() string {
	return strings.Join([]string{s.name, strconv.Itoa(s.age), strconv.Itoa(s.grade)}, " ")
}
