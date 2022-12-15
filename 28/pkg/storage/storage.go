package storage

import "study/pkg/student"

// Storage интерфейс хранилища с методами put, getAll
type Storage interface {
	Put()
	GetAll()
}

// StudentsStorage структура хранилища и его методы
type StudentsStorage struct {
	students map[string]*student.Student
}

func NewStudentsStorage() *StudentsStorage {
	return &StudentsStorage{
		students: make(map[string]*student.Student),
	}
}

func (ss *StudentsStorage) Put(name string, student *student.Student) {
	ss.students[name] = student
}

func (ss *StudentsStorage) GetAll() map[string]*student.Student {
	var result = make(map[string]*student.Student)
	for k, v := range ss.students {
		result[k] = v
	}
	return result
}
