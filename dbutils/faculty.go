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
func GetFacultyShortNames() map[string]bool {
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
