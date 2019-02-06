package dbutils

import (
	"fmt"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
	u "github.com/STEJLS/StudentAccount/utils"
)

// CreateUserStudent - создает пользователя студента. Уникальность number не проверяется
func CreateUserStudent(idField int, number string, team string, teamNumber int, durationOfStudy int,
	fullName string, IDFaculty int, IDDepartment int) (int, string) {
	var studentID int // Создаем студента
	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для создания студента: %v", err.Error()))
	}

	err = transaction.QueryRow(`INSERT INTO students(id_field, number, team, teamNumber, durationOfStudy)
						VALUES ($1, $2, $3, $4, $5) RETURNING id`, idField, number, team, teamNumber, durationOfStudy).
		Scan(&studentID)

	if err != nil {
		transaction.Rollback()
		panic(fmt.Errorf("Ошибка. При создании студента с логином = %v: %v", number, err.Error()))
	}

	tempPassword := u.GenerateTempPassword() // Создаем пользователя
	_, err = transaction.Exec(`INSERT INTO users(role, login, password, fullname, id_faculty, id_department, id_student) 
				 VALUES($1, $2, $3, $4, $5, $6, $7)`, t.StudentRole, number, u.GenerateMD5hash(tempPassword), fullName, IDFaculty, IDDepartment, studentID)
	if err != nil {
		transaction.Rollback()
		panic(fmt.Errorf("Ошибка. При создании пользователя с логином = %v: %v", number, err.Error()))
	}

	transaction.Commit()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При закрытии транзакции для создания студента: %v", err.Error()))
	}

	return studentID, tempPassword
}
