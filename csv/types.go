package csv

import (
	"log"
	"strconv"
	"strings"
)

// Faculty - факультет
type FacultyCSV struct {
	Name      string
	ShortName string
}

// Department - кафедра
type DepartmentCSV struct {
	FacultyName string
	Name        string
	ShortName   string
}

// FieldOfStudy - направление подготовки
type FieldOfStudyCSV struct {
	DepartmentName string
	Name           string
	Code           string
	Profile        string
	Alias          string
	Level          int
}

// Mark - запись об оценке
type MarkCSV struct {
	Rating   int
	Semester int
	Repass   bool
}

type StudentCSV struct {
	Number          string
	FullName        string
	Team            string
	TeamNumber      int
	DurationOfStudy int
}

type SubjectCSV struct {
	Name     string
	PassType int
}

type StudentLineCSV struct {
	Student      *StudentCSV
	FieldOfStudy *FieldOfStudyCSV
	Subject      *SubjectCSV
	Mark         *MarkCSV
}

type PracticeCSV struct {
	StudentNumber string
	Semester      int
	Name          string
	Head          string
	Company       string
	BeginDate     string
	EndDate       string
	Rating        int
}

type CourseWorksCSV struct {
	StudentNumber  string
	DepartmentName string
	FieldCode      string
	SubjectName    string
	Head           string
	Rating         int
	Semester       int
}

// facultyFromCSVLine - принимает на вход масcив строк, которые
// представляют собой колонки и возвращает объект Faculty
func facultyFromCSVLine(columns []string) *FacultyCSV {
	if len(columns) != 2 {
		return nil
	}

	return &FacultyCSV{strings.TrimSpace(columns[0]), strings.TrimSpace(columns[1])}
}

// departmentFromCSVLine - принимает на вход масcив строк, которые
// представляют собой колонки и возвращает объект Department
func departmentFromCSVLine(columns []string) *DepartmentCSV {
	if len(columns) != 3 {
		return nil
	}

	return &DepartmentCSV{strings.TrimSpace(columns[0]), strings.TrimSpace(columns[1]), strings.TrimSpace(columns[2])}
}

// departmentFromCSVLine - принимает на вход масcив строк, которые
// представляют собой колонки и возвращает объект FieldsOfStudy
func fieldOfStudyFromCSVLine(columns []string) *FieldOfStudyCSV {
	if len(columns) != 6 {
		return nil
	}
	level, err := strconv.Atoi(strings.TrimSpace(columns[5]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(level - %v ): %v", columns[4], err.Error())
		return nil
	}

	return &FieldOfStudyCSV{strings.TrimSpace(columns[0]), strings.TrimSpace(columns[1]), strings.TrimSpace(columns[2]),
		strings.TrimSpace(columns[3]), strings.TrimSpace(columns[4]), level}
}

// StudentFromCSVLine - принимает на вход масcив строк, которые
// представляют собой колонки и возвращает объект Student
func studentFromCSVLine(columns []string) *StudentLineCSV {
	if len(columns) != 12 {
		return nil
	}

	teamNumber, err := strconv.Atoi(strings.TrimSpace(columns[4]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(teamNumber - %v ): %v", columns[4], err.Error())
		return nil
	}
	durationOfStudy, err := strconv.Atoi(strings.TrimSpace(columns[6]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(durationOfStudy - %v ): %v", columns[6], err.Error())
		return nil
	}

	passType, err := strconv.Atoi(strings.TrimSpace(columns[8]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(passType - %v ): %v", columns[8], err.Error())
		return nil
	}

	rating, err := strconv.Atoi(strings.TrimSpace(columns[9]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(rating - %v ): %v", columns[9], err.Error())
		return nil
	}

	semester, err := strconv.Atoi(strings.TrimSpace(columns[10]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(semester - %v ): %v", columns[10], err.Error())
		return nil
	}

	repass, err := strconv.ParseBool(strings.TrimSpace(columns[11]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(repass - %v ): %v", columns[11], err.Error())
		return nil
	}

	student := &StudentCSV{Number: columns[0],
		FullName:        columns[1],
		Team:            columns[3],
		TeamNumber:      teamNumber,
		DurationOfStudy: durationOfStudy,
	}

	field := &FieldOfStudyCSV{
		Code:           columns[2],
		Alias:          columns[3],
		DepartmentName: columns[5],
	}

	subject := &SubjectCSV{
		Name:     columns[7],
		PassType: passType,
	}

	mark := &MarkCSV{
		Rating:   rating,
		Repass:   repass,
		Semester: semester,
	}

	return &StudentLineCSV{student, field, subject, mark}
}

func practiceFromCSVLine(columns []string) *PracticeCSV {
	if len(columns) != 8 {
		return nil
	}
	semester, err := strconv.Atoi(strings.TrimSpace(columns[1]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(semester - %v ): %v", columns[1], err.Error())
		return nil
	}
	rating, err := strconv.Atoi(strings.TrimSpace(columns[7]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(rating - %v ): %v", columns[7], err.Error())
		return nil
	}
	return &PracticeCSV{strings.TrimSpace(columns[0]),
		semester,
		strings.TrimSpace(columns[2]),
		strings.TrimSpace(columns[3]),
		strings.TrimSpace(columns[4]),
		strings.TrimSpace(columns[5]),
		strings.TrimSpace(columns[6]),
		rating,
	}
}

func courseWorkFromCSVLine(columns []string) *CourseWorksCSV {
	if len(columns) != 7 {
		return nil
	}
	rating, err := strconv.Atoi(strings.TrimSpace(columns[5]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(rating - %v ): %v", columns[5], err.Error())
		return nil
	}
	semester, err := strconv.Atoi(strings.TrimSpace(columns[6]))
	if err != nil {
		log.Printf("Ошибка. Припарсинге строки в цифру(semester - %v ): %v", columns[7], err.Error())
		return nil
	}
	return &CourseWorksCSV{strings.TrimSpace(columns[0]),
		strings.TrimSpace(columns[1]),
		strings.TrimSpace(columns[2]),
		strings.TrimSpace(columns[3]),
		strings.TrimSpace(columns[4]),
		rating,
		semester,
	}
}
