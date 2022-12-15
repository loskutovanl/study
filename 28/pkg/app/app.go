package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"study/pkg/storage"
	"study/pkg/student"
)

// App структура приложения и его методы
type App struct {
	repository storage.Storage
}

func NewApp(repository storage.Storage) *App {
	return &App{
		repository: repository,
	}
}

func (a *App) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	space := " "
	hyphen := 45

	newStorage := storage.NewStudentsStorage()

	// инструкция для пользователя
	fmt.Println("Инструкция:")
	fmt.Println("Вводите строки в формате: имя возраст оценка")
	fmt.Println("Для завершения работы подайте сигнал Ctrl + D")
	fmt.Println(strings.Repeat("-", hyphen))
	fmt.Println("Строки:")

	// ввод пользователя в бесконечном цикле до введения Ctrl+d
	for scanner.Scan() {
		studentInfo := strings.Split(scanner.Text(), space)

		if student.IsStudentValid(studentInfo) {
			newStudent := student.NewStudent(studentInfo)
			newStorage.Put(studentInfo[0], newStudent)
		}
	}

	fmt.Println(strings.Repeat("-", hyphen))
	fmt.Println("Студенты из хранилища:")
	allStudents := newStorage.GetAll()

	for _, v := range allStudents {
		fmt.Println(v.PrintStudent())
	}
}
