package dbutils

import (
	"fmt"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
)

// GetAllDepartments - возвращает все кафедры
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

// GetAllDepartmentsInMap - возвращает карту где ключ имя кафедры, а значение сама кафедра
func GetAllDepartmentsInMap() map[string]*t.Department {
	deps := GetAllDepartments()

	result := make(map[string]*t.Department)

	for _, d := range deps {
		result[d.Name] = d
	}

	return result
}

// GetDepartmentShortNames - возвращает все короткие имена кафедр
func GetDepartmentShortNames() map[string]bool {
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
