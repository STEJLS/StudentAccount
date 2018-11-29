package csv

import (
	"encoding/csv"
	"io"
	"log"
)

// ReadFaculties - создает слайс факультетов из csv файла
func ReadFaculties(r io.Reader) []*FacultyCSV {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Ошибка. При парсинге csv файла: %v", err.Error())
		return nil
	}

	result := make([]*FacultyCSV, 0, len(records))

	for _, v := range records {
		if f := facultyFromCSVLine(v); f != nil {
			result = append(result, f)
		} else {
			return nil
		}
	}

	return result
}

// ReadDepartments - создает слайс кафедр из csv файла
func ReadDepartments(r io.Reader) []*DepartmentCSV {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Ошибка. При парсинге csv файла: %v", err.Error())
		return nil
	}

	result := make([]*DepartmentCSV, 0, len(records))

	for _, v := range records {
		if d := departmentFromCSVLine(v); d != nil {
			result = append(result, d)
		} else {
			return nil
		}
	}

	return result
}

// departmentName,name,code,alias,level
// ReadFieldsOfStudy - создает слайс направлений подготовки из csv файла
func ReadFieldsOfStudy(r io.Reader) []*FieldOfStudyCSV {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Ошибка. При парсинге csv файла: %v", err.Error())
		return nil
	}

	result := make([]*FieldOfStudyCSV, 0, len(records))

	for _, v := range records {
		if field := fieldOfStudyFromCSVLine(v); field != nil {
			result = append(result, field)
		} else {
			return nil
		}
	}

	return result
}

// ReadStudents - создает слайс студентов из csv файла
func ReadStudents(r io.Reader) []*StudentLineCSV {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Ошибка. При парсинге csv файла: %v", err.Error())
		return nil
	}

	result := make([]*StudentLineCSV, 0, len(records))

	for _, v := range records {
		if s := studentFromCSVLine(v); s != nil {
			result = append(result, s)
		} else {
			return nil
		}
	}

	return result
}

// ReadPracticis - создает слайс практик из csv файла
func ReadPracticis(r io.Reader) []*PracticeCSV {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Ошибка. При парсинге csv файла: %v", err.Error())
		return nil
	}

	result := make([]*PracticeCSV, 0, len(records))

	for _, v := range records {
		if p := practiceFromCSVLine(v); p != nil {
			result = append(result, p)
		} else {
			return nil
		}
	}

	return result
}

// ReadPracticis - создает слайс практик из csv файла
func ReadCourseWorks(r io.Reader) []*CourseWorksCSV {
	csvReader := csv.NewReader(r)

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("Ошибка. При парсинге csv файла: %v", err.Error())
		return nil
	}

	result := make([]*CourseWorksCSV, 0, len(records))

	for _, v := range records {
		if cw := courseWorkFromCSVLine(v); cw != nil {
			result = append(result, cw)
		} else {
			return nil
		}
	}

	return result
}
