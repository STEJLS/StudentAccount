package dbutils

import (
	"fmt"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
)

// GetAllFields - возвращает все направления подготовки
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

// GetAllFieldsInMap - возвращает карту где ключ код направления подготовки, а значение само направление подготовки
func GetAllFieldsInMap() map[string]*t.FieldOfStudy {
	fields := GetAllFields()

	result := make(map[string]*t.FieldOfStudy)

	for _, f := range fields {
		result[f.Code] = f
	}

	return result
}

// GetFieldsofstudyAliases - возвращает все алиасы (которкие названия) направлений подготовки
func GetFieldsofstudyAliases() map[string]int {
	rows, err := g.DB.Query("SELECT id, alias FROM fieldsofstudy")
	if err != nil {
		panic(fmt.Errorf("Ошибка. При получении алиасов направлений подготовки: %v", err.Error()))
	}

	defer rows.Close()

	shortNames := make(map[string]int)
	var name string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При получении алиасов направлений подготовки: %v", err.Error()))
		}
		shortNames[name] = id
	}

	return shortNames
}
