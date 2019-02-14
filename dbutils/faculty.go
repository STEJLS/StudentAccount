package dbutils

import (
	"fmt"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
)

// GetAllFaculties - возврщает все факультеты
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

// GetFacultyShortNames - возвращает все короткие названия факультетов
func GetFacultyShortNames() map[string]int {
	rows, err := g.DB.Query("SELECT id, shortname FROM faculties")
	if err != nil {
		panic(fmt.Errorf("Ошибка. При получении коротких названий факультетов: %v", err.Error()))
	}

	defer rows.Close()

	shortNames := make(map[string]int)
	var name string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При получении коротких названий факультетов: %v", err.Error()))
		}
		shortNames[name] = id
	}

	return shortNames
}
