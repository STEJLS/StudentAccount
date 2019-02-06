package dbutils

import (
	"fmt"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
	u "github.com/STEJLS/StudentAccount/utils"
)

func GetAllFields() []*t.FieldOfStudy {
	rows, err := g.DB.Query(`SELECT id, id_department, name, alias, code FROM fieldsOfStudy`)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке всех направлений подготовки: %v", err.Error()))
	}

	defer rows.Close()

	result := make([]*t.FieldOfStudy, 0)

	for rows.Next() {
		field := &t.FieldOfStudy{}
		err = rows.Scan(&field.ID, &field.IDDepartment, &field.Name, &field.Alias, &field.Code)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке всех направлений подготовки: %v", err.Error()))
		}
		result = append(result, field)
	}

	return result
}

// getAllFieldsInMap - возвращает карту где ключ код направления подготовки, а значение само направление подготовки
func GetAllFieldsInMap() map[string]*t.FieldOfStudy {
	fields := GetAllFields()

	result := make(map[string]*t.FieldOfStudy)

	for _, f := range fields {
		result[f.Code] = f
	}

	return result
}

func GetAllFaculties() []*t.Faculty {
	rows, err := g.DB.Query(`SELECT id, name, shortName FROM faculties`)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке всех факултетов: %v", err.Error()))
	}

	defer rows.Close()

	result := make([]*t.Faculty, 0)

	for rows.Next() {
		f := &t.Faculty{}
		err = rows.Scan(&f.ID, &f.Name, &f.ShortName)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке всех факултетов: %v", err.Error()))
		}
		result = append(result, f)
	}

	return result
}

func GetAllDepartments() []*t.Department {
	rows, err := g.DB.Query(`SELECT id, id_faculty, name, shortName FROM departments`)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке всех кафедр: %v", err.Error()))
	}

	defer rows.Close()

	result := make([]*t.Department, 0)

	for rows.Next() {
		dep := &t.Department{}
		err = rows.Scan(&dep.ID, &dep.IDFaculty, &dep.Name, &dep.ShortName)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке всех кафедр: %v", err.Error()))
		}
		result = append(result, dep)
	}

	return result
}

// getAllDepartmentsInMap - возвращает карту где ключ имя кафедры, а значение сама кафедра
func GetAllDepartmentsInMap() map[string]*t.Department {
	deps := GetAllDepartments()

	result := make(map[string]*t.Department)

	for _, d := range deps {
		result[d.Name] = d
	}

	return result
}

// createUserStudent - уникальность number не проверяется
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

func getFacultyShortNames() map[string]bool {
	rows, err := g.DB.Query("SELECT shortname FROM faculties")
	if err != nil {
		panic(fmt.Errorf("Ошибка. При получении коротких названий факультетов: %v", err.Error()))
	}

	defer rows.Close()

	shortNames := make(map[string]bool)
	var name string

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При получении коротких названий факультетов: %v", err.Error()))
		}
		shortNames[name] = true
	}

	return shortNames
}

func getDepartmentShortNames() map[string]bool {
	rows, err := g.DB.Query("SELECT shortname FROM departments")
	if err != nil {
		panic(fmt.Errorf("Ошибка. При получении коротких названий кафедр: %v", err.Error()))
	}

	defer rows.Close()
	shortNames := make(map[string]bool)
	var name string

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При получении коротких названий кафедр: %v", err.Error()))
		}
		shortNames[name] = true
	}

	return shortNames
}

func getFieldsofstudyAliases() map[string]bool {
	rows, err := g.DB.Query("SELECT alias FROM fieldsofstudy")
	if err != nil {
		panic(fmt.Errorf("Ошибка. При получении алиасов направлений подготовки: %v", err.Error()))
	}

	defer rows.Close()

	shortNames := make(map[string]bool)
	var name string

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При получении алиасов направлений подготовки: %v", err.Error()))
		}
		shortNames[name] = true
	}

	return shortNames
}
