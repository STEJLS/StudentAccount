package contexts

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/STEJLS/StudentAccount/csv"
	dbu "github.com/STEJLS/StudentAccount/dbutils"
	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
	u "github.com/STEJLS/StudentAccount/utils"
	"github.com/gocraft/web"
	"github.com/lib/pq"
)

func addFacultiesFromCSV(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	file, _, err := req.FormFile("csvFile")

	if err != nil {
		c.response.Message = "Файл не получен"
		return
	}

	faculties := csv.ReadFaculties(file)

	if len(faculties) == 0 {
		c.response.Message = "Получены некорректные данные"
		return
	}

	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для создания факультета: %v", err.Error()))
	}
	stmt, err := transaction.Prepare(`INSERT INTO faculties(name, shortName) Values($1,$2)`)
	if err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления новых факультетов: %v", err.Error()))
		}
		panic(fmt.Errorf("Ошибка. При создании подготовленного выражения: %v", err.Error()))
	}

	for i, f := range faculties {
		_, err = stmt.Exec(f.Name, f.ShortName)
		if err != nil {
			if errt := transaction.Rollback(); errt != nil {
				panic(fmt.Errorf("Ошибка. При откате транзакции для добавления новых факультетов: %v", err.Error()))
			}
			if pe, ok := err.(*pq.Error); ok { // Нарушение уникальности, // Нет такой кафедры
				if pe.Code == "23505" {
					c.response.Message = fmt.Sprintf("Запись под номером %v нарушает уникальность", i+1)
				}
				c.response.Body = f
				return
			} else {
				panic(fmt.Errorf("Ошибка. При вставке новых факультетов: %v", err.Error()))
			}
		}
	}

	if err = stmt.Close(); err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления новых факультетов: %v", err.Error()))
		}
		log.Printf("Ошибка. При закрытии подготовленного выражения: %v", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		panic(fmt.Errorf("Ошибка. При подтверждении транзакции для добавления новых факультетов: %v", err.Error()))
	}

	c.response.Сompleted = true
	c.response.Message = "Факультеты успешно добавлены"
}
func addDepartmentsFromCSV(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	file, _, err := req.FormFile("csvFile")
	if err != nil {
		c.response.Message = "Файл не получен"
		return
	}

	departments := csv.ReadDepartments(file)

	if len(departments) == 0 {
		c.response.Message = "Получены некорректные данные"
		return
	}

	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для создания кафедр: %v", err.Error()))
	}
	stmt, err := transaction.Prepare(`INSERT INTO departments(name, shortName, id_faculty) Values($1,$2,
		(SELECT id FROM faculties where name = $3))`)
	if err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления новых кафедр: %v", err.Error()))
		}
		panic(fmt.Errorf("Ошибка. При создании подготовленного выражения: %v", err.Error()))
	}

	for i, d := range departments {
		_, err = stmt.Exec(d.Name, d.ShortName, d.FacultyName)
		if err != nil {
			if errt := transaction.Rollback(); errt != nil {
				panic(fmt.Errorf("Ошибка. При откате транзакции для добавления новых кафедр: %v", err.Error()))
			}
			if pe, ok := err.(*pq.Error); ok {
				if pe.Code == "23505" {
					c.response.Message = fmt.Sprintf("Запись под номером %v нарушает уникальность", i+1)
				}
				if pe.Code == "23502" {
					c.response.Message = fmt.Sprintf("В записи под номером %v указан неверный факультет", i+1)
				}
				c.response.Body = d
				return
			} else {
				panic(fmt.Errorf("Ошибка. При вставке новых кафедр: %v", err.Error()))
			}
		}
	}

	if err = stmt.Close(); err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления новых кафедр: %v", err.Error()))
		}
		log.Printf("Ошибка. При закрытии подготовленного выражения: %v", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		panic(fmt.Errorf("Ошибка. При подтверждении транзакции для добавления новых кафедр: %v", err.Error()))
	}

	c.response.Сompleted = true
	c.response.Message = "Кафедры успешно добавлены"
}

// addFieldsOfStudyFromCSV - добавляет в систему направления подготовки
func addFieldsOfStudyFromCSV(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	file, _, err := req.FormFile("csvFile")

	if err != nil {
		c.response.Message = "Файл не получен"
		return
	}

	fieldsOfStudy := csv.ReadFieldsOfStudy(file)

	if len(fieldsOfStudy) == 0 {
		c.response.Message = "Получены некорректные данные"
		return
	}

	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для создания напрвлений подготовки: %v", err.Error()))
	}

	stmt, err := transaction.Prepare(`INSERT INTO fieldsOfStudy(id_department, name, code, alias, level) VALUES( 
		(SELECT id from departments where name = $1),$2,$3,$4,$5)`)
	if err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавлении направления подготовки: %v", err.Error()))
		}
		panic(fmt.Errorf("Ошибка. При создании подготовленного выражения: %v", err.Error()))
	}

	for i, f := range fieldsOfStudy {
		_, err = stmt.Exec(f.DepartmentName, f.Name, f.Code, f.Alias, f.Level)
		if err != nil {
			if errt := transaction.Rollback(); errt != nil {
				panic(fmt.Errorf("Ошибка. При откате транзакции для добавлении направления подготовки: %v", err.Error()))
			}
			if pe, ok := err.(*pq.Error); ok { // Нарушение уникальности, // Нет такой кафедры
				if pe.Code == "23505" {
					c.response.Message = fmt.Sprintf("Запись под номером %v нарушает уникальность", i+1)
				}
				if pe.Code == "23502" {
					c.response.Message = fmt.Sprintf("В записи под номером %v указана неверная кафедра", i+1)
				}
				c.response.Body = f
				return
			} else {
				panic(fmt.Errorf("Ошибка. При вставке новых факультетов: %v", err.Error()))
			}
		}
	}

	if err = stmt.Close(); err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавлении направления подготовки: %v", err.Error()))
		}
		log.Printf("Ошибка. При закрытии подготовленного выражения: %v", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		panic(fmt.Errorf("Ошибка. При подтверждении транзакции для создания напрвлений подготовки: %v", err.Error()))
	}

	c.response.Сompleted = true
	c.response.Message = "Направления подготовки успешно добавлены"
}

func addStudentsFromCSV(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	file, _, err := req.FormFile("csvFile")

	if err != nil {
		c.response.Message = "Файл не получен"
		return
	}

	passwordFile, passwordWriter := u.InitPasswordCSVWriter()
	defer passwordFile.Close()

	studentLines := csv.ReadStudents(file)

	if len(studentLines) == 0 {
		c.response.Message = "Получены некорректные данные"
		return
	}

	fields := dbu.GetAllFieldsInMap()
	departments := dbu.GetAllDepartmentsInMap()

	for i, item := range studentLines { // проверка направления подготовки и кафедры
		_, ok := fields[item.FieldOfStudy.Code]
		if !ok {
			c.response.Message = fmt.Sprintf("Запись на строке %v имеет неверное направление подготовки", i+1)
			return
		}

		_, ok = departments[item.FieldOfStudy.DepartmentName]
		if !ok {
			c.response.Message = fmt.Sprintf("Запись на строке %v имеет неверное название кафедры", i+1)
			return
		}
	}

	//мапим по номеру студента
	studentsMap := make(map[string][]*csv.StudentLineCSV)
	for _, s := range studentLines {
		_, ok := studentsMap[s.Student.Number]
		if !ok {
			studentsMap[s.Student.Number] = make([]*csv.StudentLineCSV, 0, 15)
		}
		studentsMap[s.Student.Number] = append(studentsMap[s.Student.Number], s)
	}

	for number, s := range studentsMap {
		var currentStudentID int
		var password string
		//Ищем/Создаем студента
		err := g.DB.QueryRow(`SELECT id FROM students WHERE number = $1`, number).Scan(&currentStudentID)
		if err == sql.ErrNoRows {
			currentStudentID, password = dbu.CreateUserStudent(fields[s[0].FieldOfStudy.Code].ID, s[0].Student.Number, s[0].Student.Team, s[0].Student.TeamNumber,
				s[0].Student.DurationOfStudy, s[0].Student.FullName, departments[s[0].FieldOfStudy.DepartmentName].IDFaculty, departments[s[0].FieldOfStudy.DepartmentName].ID)
			//Записываем новый пароль в файл
			passwordWriter.Write([]string{number, s[0].Student.FullName,
				fmt.Sprintf("%v-%v", s[0].Student.Team, s[0].Student.TeamNumber), password})
			passwordWriter.Flush()
			if err := passwordWriter.Error(); err != nil {
				panic(fmt.Errorf("Ошибка. При записи временного пароля: %v", err.Error()))
			}
		} else if err != nil {
			panic(fmt.Errorf("Ошибка. При поиске студента с номером - %v : %v", number, err.Error()))
		}

		//Проходим по всем записям контретного студента
		for _, r := range s {
			transaction, err := g.DB.Begin()
			if err != nil {
				panic(fmt.Errorf("Ошибка. При создании транзакции для создания студента: %v", err.Error()))
			}

			var currentSubject t.Subject // Находим/создаем предмет
			err = transaction.QueryRow(`SELECT id,id_department, id_field, name, passType FROM subjects
								  WHERE id_department = $1 AND id_field = $2 AND name = $3 AND passType = $4`,
				fields[r.FieldOfStudy.Code].IDDepartment,
				fields[r.FieldOfStudy.Code].ID, r.Subject.Name, r.Subject.PassType).
				Scan(&currentSubject.ID, &currentSubject.IDDepartment,
					&currentSubject.IDField, &currentSubject.Name, &currentSubject.PassType)
			if err == sql.ErrNoRows {
				err = transaction.QueryRow(`INSERT INTO subjects(id_department, id_field, name, passType) 
									  VALUES ($1, $2, $3, $4) RETURNING id;`, fields[r.FieldOfStudy.Code].IDDepartment,
					fields[r.FieldOfStudy.Code].ID, r.Subject.Name, r.Subject.PassType).
					Scan(&currentSubject.ID)
				if err != nil {
					transaction.Rollback()
					panic(fmt.Errorf("Ошибка. При добавлении предмета: %v", err.Error()))
				}
			} else if err != nil {
				transaction.Rollback()
				panic(fmt.Errorf("Ошибка. При поиске предмета: %v", err.Error()))
			}

			//Добавляем успеваемость
			_, err = transaction.Exec(`INSERT INTO marks(id_student, id_subject, rating, semester, repass)
							VALUES($1, $2, $3, $4, $5)`, currentStudentID, currentSubject.ID, r.Mark.Rating, r.Mark.Semester, r.Mark.Repass)
			if err != nil {
				transaction.Rollback()
				panic(fmt.Errorf("Ошибка. При поиске предмета: %v", err.Error()))
			}

			transaction.Commit()
			if err != nil {
				panic(fmt.Errorf("Ошибка. При закрытии транзакции для создания студента: %v", err.Error()))
			}
		}
	}

	c.response.Сompleted = true
	c.response.Message = "Добавление успешно завершено."
}

// addPracticisFromCSV - добавляет в систему направления подготовки
func addPracticisFromCSV(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	file, _, err := req.FormFile("csvFile")
	if err != nil {
		c.response.Message = "Файл не получен"
		return
	}

	practicis := csv.ReadPracticis(file)

	if len(practicis) == 0 {
		c.response.Message = "Получены некорректные данные"
		return
	}

	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для создания практик: %v", err.Error()))
	}

	stmt, err := transaction.Prepare(`INSERT INTO practicis(id_student, semester, name, head, company, begin_date, end_date, rating) 
								      VALUES( (SELECT id from students where number = $1),$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления практик: %v", err.Error()))
		}
		panic(fmt.Errorf("Ошибка. При создании подготовленного выражения: %v", err.Error()))
	}

	for i, p := range practicis {
		_, err = stmt.Exec(p.StudentNumber, p.Semester, p.Name, p.Head, p.Company, p.BeginDate, p.EndDate, p.Rating)
		if err != nil {
			if errt := transaction.Rollback(); errt != nil {
				panic(fmt.Errorf("Ошибка. При откате транзакции для добавлении практик: %v", err.Error()))
			}
			if pe, ok := err.(*pq.Error); ok { // Нарушение уникальности, // Нет такой кафедры
				if pe.Code == "23505" {
					c.response.Message = fmt.Sprintf("Запись под номером %v нарушает уникальность", i+1)
				}
				if pe.Code == "23502" {
					c.response.Message = fmt.Sprintf("В записи под номером %v указан неверный номер студента", i+1)
				}
				c.response.Body = p
				return
			} else {
				panic(fmt.Errorf("Ошибка. При вставке новых практик: %v", err.Error()))
			}
		}
	}

	if err = stmt.Close(); err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления практик: %v", err.Error()))
		}
		log.Printf("Ошибка. При закрытии подготовленного выражения: %v", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		panic(fmt.Errorf("Ошибка. При подтверждении транзакции для создания практик: %v", err.Error()))
	}

	c.response.Сompleted = true
	c.response.Message = "Практики успешно добавлены"
}

// addCourseWorksFromCSV - добавляет в систему курсовые работы
func addCourseWorksFromCSV(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	file, _, err := req.FormFile("csvFile")

	if err != nil {
		c.response.Message = "Файл не получен"
		return
	}

	courseWorks := csv.ReadCourseWorks(file)

	if len(courseWorks) == 0 {
		c.response.Message = "Получены некорректные данные"
		return
	}

	transaction, err := g.DB.Begin()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При создании транзакции для создания курсовых работ: %v", err.Error()))
	}

	stmt, err := transaction.Prepare(`INSERT INTO courseworks(id_student,id_subject,semester,head,rating) 
									  VALUES( (SELECT id FROM students WHERE number = $1), 
											  (SELECT id FROM subjects WHERE id_field = 
												(SELECT id FROM fieldsOfStudy WHERE id_department = 
														(SELECT id FROM departments WHERE name = $2)
												 AND code = $3) 
											   AND name=$4),
											$5, $6, $7)`)
	if err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления курсовых работ: %v", err.Error()))
		}
		panic(fmt.Errorf("Ошибка. При создании подготовленного выражения: %v", err.Error()))
	}

	for i, cw := range courseWorks {
		_, err = stmt.Exec(cw.StudentNumber, cw.DepartmentName, cw.FieldCode, cw.SubjectName, cw.Semester, cw.Head, cw.Rating)
		if err != nil {
			if errt := transaction.Rollback(); errt != nil {
				panic(fmt.Errorf("Ошибка. При откате транзакции для добавлении курсовых работ: %v", err.Error()))
			}
			if pe, ok := err.(*pq.Error); ok { // Нарушение уникальности, // Нет такой кафедры
				if pe.Code == "23505" {
					c.response.Message = fmt.Sprintf("Запись под номером %v нарушает уникальность", i+1)
				}
				if pe.Code == "23502" {
					c.response.Message = fmt.Sprintf("В записи под номером %v указан неверный внешний ключ", i+1)
				}
				c.response.Body = cw
				return
			} else {
				panic(fmt.Errorf("Ошибка. При вставке новых курсовых работ: %v", err.Error()))
			}
		}
	}

	if err = stmt.Close(); err != nil {
		if err = transaction.Rollback(); err != nil {
			panic(fmt.Errorf("Ошибка. При откате транзакции для добавления курсовых работ: %v", err.Error()))
		}
		log.Printf("Ошибка. При закрытии подготовленного выражения: %v", err.Error())
	}

	if err = transaction.Commit(); err != nil {
		panic(fmt.Errorf("Ошибка. При подтверждении транзакции для создания курсовых работ: %v", err.Error()))
	}

	c.response.Сompleted = true
	c.response.Message = "Курсовые работы успешно добавлены"
}

func getTempPasswords(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	data, err := ioutil.ReadFile(g.PasswordFileStorageName)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При чтении файла с паролем: %v", err.Error()))
	}

	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\""+"tempPasswords-%s.csv"+"\"", time.Now().Format("2006-01-02_15:04:05")))
	rw.Header().Add("Content-type", "text/csv")
	rw.Header().Add("Content-Length", fmt.Sprintf("%v", len(data)))

	rw.Write(data)
}

func createVerif(c *UserContext, rw web.ResponseWriter, req *web.Request) {
	verif, errStr := u.ValidateVerif(req.FormValue("login"), req.FormValue("password"), req.FormValue("fullName"), req.FormValue("id_faculty"), req.FormValue("id_department"))
	if verif == nil {
		c.response.Message = errStr
		return
	}

	err := g.DB.QueryRow(`INSERT INTO users(login, password, fullname, id_faculty, id_department, role, isActivated)
	 VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`, verif.Login, verif.Password, verif.FullName, verif.IDFaculty, verif.IDDepartment, verif.Role, verif.IsActivated).
		Scan(&verif.ID)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При добавлении верификатора: %v", err.Error()))
	}

	verif.Password = ""

	c.response.Сompleted = true
	c.response.Message = "Верификатор успешно создан"
	c.response.Body = verif
}
