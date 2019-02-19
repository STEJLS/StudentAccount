package FOSandRPDparser

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/STEJLS/StudentAccount/dbutils"
	g "github.com/STEJLS/StudentAccount/globals"
)

type addDocToDBFunc func(string, string)

type parser struct {
	facultyShortNames     map[string]int
	departmentShortNames  map[string]int
	fieldsOfStudyAliases  map[string]int
	currentFacultyName    string
	currentDepartmentName string
	currentFieldName      string
	currentFacultyID      int
	currentDepartmentID   int
	currentFieldID        int
	action                addDocToDBFunc
	stmt                  *sql.Stmt
}

func NewParser() *parser {
	return &parser{
		facultyShortNames:    dbutils.GetFacultyShortNames(),
		departmentShortNames: dbutils.GetDepartmentShortNames(),
		fieldsOfStudyAliases: dbutils.GetFieldsofstudyAliases(),
	}
}

// Parse - парсит файлы ФОС и РПД, которые лежат в папках,
// определенных в глобальных переменных FOSDirectoryName и RPDDirectoryName.
// В случае успеха возвращает true, иначе false.
func (p *parser) Parse() bool {

	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для парсинга ФОС и РПД: %v", err.Error()))
	}

	stmt, err := transaction.Prepare(`INSERT INTO documents (id_faculty, id_department, id_field, type, name, path)
	VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для парсинга ФОС и РПД: %v", err.Error()))
		}
		panic(fmt.Errorf("Ошибка. При создании подготовленного выражения для парсинга ФОС и РПД: %v", err.Error()))
	}

	p.stmt = stmt

	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
			panic(r)
		}
	}()

	p.action = p.addFOS
	isParsedFOS := p.parse(g.FOSDirectoryName)
	p.action = p.addRPD
	isParsedRPD := p.parse(g.RPDDirectoryName)

	if err = stmt.Close(); err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для парсинга ФОС и РПД: %v", err.Error()))
		}
		log.Printf("Ошибка. При закрытии подготовленного выражения для парсинга ФОС и РПД: %v", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		panic(fmt.Errorf("Ошибка. При подтверждении транзакции для парсинга ФОС и РПД: %v", err.Error()))
	}

	if !isParsedFOS && !isParsedRPD {
		return false
	}

	return true
}

func (p *parser) parse(directoryName string) bool {
	docDir, err := os.Open(directoryName) //ФОС или РПД
	if err != nil {
		if os.IsNotExist(err) { // такой папки нет
			log.Printf("Ошибка. При парсинге ФОС или РПД: не могу открыть папку %v. %v", directoryName, err.Error())
			return false
		}

		panic(fmt.Errorf("Ошибка. При открытии папки(%v) для парсинга ФОС или РПД: %v", directoryName, err.Error()))
	}

	facultyDirNames, err := docDir.Readdirnames(0)
	docDir.Close()
	if err != nil {
		log.Printf("Ошибка. При парсинге ФОС или РПД: не могу прочитать список файлов в папке %v.", directoryName)
		return false
	}

	p.processFaculties(directoryName, facultyDirNames)

	return true
}

func (p *parser) processFaculties(parentDir string, facultiesDirNames []string) {
	for _, facultyDirName := range facultiesDirNames { // Факультеты
		id, ok := p.facultyShortNames[facultyDirName]
		if !ok { // Проверияем есть ли такой факультет
			log.Printf("Ошибка. При парсинге ФОС или РПД: факультета с таким названием(%v) не существует.", facultyDirName)
			continue
		}

		p.currentFacultyName = facultyDirName
		p.currentFacultyID = id

		facultyDirPath := path.Join(parentDir, facultyDirName)
		facultyDir, err := os.Open(facultyDirPath)
		if err != nil {
			log.Printf("Ошибка. При открытии папки факультета(%v) для парсинга ФОС или РПД: %v", facultyDirPath, err.Error())
			continue
		}

		departmentDirNames, err := facultyDir.Readdirnames(0)
		facultyDir.Close()
		if err != nil {
			log.Printf("Ошибка. При парсинге ФОС или РПД: не могу прочитать список файлов в папке %v.", facultyDirPath)
			continue
		}
		p.processDepartments(facultyDirPath, departmentDirNames)
	}
}

func (p *parser) processDepartments(parentDir string, departmentsDirNames []string) {
	for _, departmentDirName := range departmentsDirNames { // Кафедры
		id, ok := p.departmentShortNames[departmentDirName]
		if !ok { // Проверияем есть ли такая кафедра
			log.Printf("Ошибка. При парсинге ФОС или РПД: кафедры с таким названием(%v) не существует.", departmentDirName)
			continue
		}

		p.currentDepartmentName = departmentDirName
		p.currentDepartmentID = id

		departmentDirPath := path.Join(parentDir, departmentDirName)
		departmentDir, err := os.Open(departmentDirPath)
		if err != nil {
			log.Printf("Ошибка. При открытии папки кафедры(%v) для парсинга ФОС или РПД: %v", departmentDirPath, err.Error())
			continue
		}

		fieldDirNames, err := departmentDir.Readdirnames(0)
		departmentDir.Close()
		if err != nil {
			log.Printf("Ошибка. При получении списка папок в папке %v для парсинга ФОС или РПД: %v", departmentDirPath, err.Error())
			continue
		}

		p.processField(departmentDirPath, fieldDirNames)
	}
}

func (p *parser) processField(parentDir string, fieldsDirNames []string) {
	for _, fieldDirName := range fieldsDirNames { //Направление подготовки
		id, ok := p.fieldsOfStudyAliases[fieldDirName]
		if !ok { // Проверияем есть ли такое направление подготовки
			log.Printf("Ошибка. При парсинге ФОС или РПД: направления подготовки с таким названием(%v) не существует.", fieldDirName)
			continue
		}

		p.currentFieldName = fieldDirName
		p.currentFieldID = id

		fieldDirPath := path.Join(parentDir, fieldDirName)
		files, err := ioutil.ReadDir(fieldDirPath)
		if err != nil {
			log.Printf("Ошибка. При парсинге ФОС или РПД: не могу получить список файлов в папке(%v): %v", fieldDirName, err.Error())
			continue
		}

		for _, file := range files {
			if !file.IsDir() {
				p.action(fieldDirPath, file.Name())
			}
		}
	}
}

func (p *parser) addFOS(parentDir string, fileName string) {
	_, err := p.stmt.Exec(p.currentFacultyID, p.currentDepartmentID, p.currentFieldID, 0,
		strings.TrimRight(fileName, filepath.Ext(fileName)),
		path.Join(parentDir, fileName))

	if err != nil {
		log.Printf("Ошибка. При вставке в БД нового ФОС документа(%v): %v", path.Join(parentDir, fileName), err.Error())
	}
}

func (p *parser) addRPD(parentDir string, fileName string) {
	_, err := p.stmt.Exec(p.currentFacultyID, p.currentDepartmentID,
		p.currentFieldID, 1, strings.TrimRight(fileName, filepath.Ext(fileName)),
		path.Join(parentDir, fileName))

	if err != nil {
		log.Printf("Ошибка. При вставке в БД нового РПД документа(%v): %v", path.Join(parentDir, fileName), err.Error())
	}
}
