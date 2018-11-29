package contexts

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
	u "github.com/STEJLS/StudentAccount/utils"
	"github.com/gocraft/web"
	"github.com/lib/pq"
)

//Контекст студента
type studentContext struct {
	*Context
}

// StudentRequire - middleware, которое требует прав студента
func (c *studentContext) StudentRequire(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if c.user.Role != t.StudentRole {
		c.response.Message = "Недостаточно прав."
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	next(rw, req)
}

func getStudentMarks(c *studentContext, rw web.ResponseWriter, req *web.Request) {

	rows, err := g.DB.Query(`SELECT departments.shortname, subjects.name, marks.rating, subjects.passtype, marks.semester
			  			   FROM students
						   JOIN marks ON students.id = marks.id_student
						   JOIN subjects ON marks.id_subject = subjects.id
						   JOIN departments ON subjects.id_department = departments.id
						   WHERE students.id = $1`, c.user.IDStudent)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке оценок пользователя(login - %v ): %v", c.user.Login, err.Error()))
	}
	defer rows.Close()

	marks := make([]*t.ResponseMarks, 0)

	for rows.Next() {
		mark := t.ResponseMarks{}
		err = rows.Scan(&mark.Department, &mark.Subject, &mark.Rating, &mark.PassType, &mark.Semester)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке оценок пользователя(login - %v ): %v", c.user.Login, err.Error()))
		}

		marks = append(marks, &mark)
	}

	c.response.Body = marks
	c.response.Сompleted = true
}

func getStudentInfo(c *studentContext, rw web.ResponseWriter, req *web.Request) {
	row := g.DB.QueryRow(`SELECT s.number, u.fullname, s.team, s.teamnumber,
	fi.name, fi.code, fi.level,
							d.name, d.shortname, fa.name, fa.shortname 
							FROM students s
							JOIN users u ON u.id_student=s.id
							JOIN fieldsofstudy fi ON fi.id = s.id_field
							JOIN departments d ON d.id = u.id_department
							JOIN faculties fa ON fa.id = u.id_faculty
							WHERE s.id = $1`, c.user.IDStudent)

	var userInfo struct {
		Number              string
		FullName            string
		Team                string
		TeamNumber          int
		FieldName           string
		FieldCode           string
		level               int
		DepartmentName      string
		DepartmentShortName string
		FacultyName         string
		FacultyShortName    string
	}

	err := row.Scan(&userInfo.Number, &userInfo.FullName, &userInfo.Team, &userInfo.TeamNumber, &userInfo.FieldName,
		&userInfo.FieldCode, &userInfo.level, &userInfo.DepartmentName, &userInfo.DepartmentShortName, &userInfo.FacultyName,
		&userInfo.FacultyShortName)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке информации о пользователе(login - %v ): %v", c.user.Login, err.Error()))
	}

	c.response.Body = userInfo
	c.response.Сompleted = true
}

func addArticle(c *studentContext, rw web.ResponseWriter, req *web.Request) {
	name := req.FormValue("name")
	journal := req.FormValue("journal")
	biblioRecord := req.FormValue("biblioRecord")
	articleType := req.FormValue("type")

	file, header, err := req.FormFile("article")
	if err != nil {
		c.response.Message = "Файл не получен"
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if filepath.Ext(header.Filename) != ".pdf" {
		c.response.Message = "Файл должен иметь формат pdf"
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	article, errStr := u.ValidateArticle(name, journal, biblioRecord, articleType)
	if article == nil {
		c.response.Message = errStr
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := u.GenerateToken()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.response.Message = "Файл не получен"
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile(path.Join(g.ArticlesDirectory, fileName), data, 0666)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При записи статьи на диск: %v", err.Error()))
	}

	_, err = g.DB.Exec(`INSERT INTO articles(name, journal, biblioRecord, type, fileName, realFileName,id_student) 
			 VALUES($1, $2, $3, $4, $5, $6, $7)`, article.Name, article.Journal, article.BiblioRecord, article.ArticleType, fileName, header.Filename, c.user.IDStudent)
	if err != nil {
		os.Remove(path.Join(g.ArticlesDirectory, fileName))

		if pe, ok := err.(*pq.Error); ok { // Нарушение уникальности
			if pe.Code == "23505" {
				c.response.Message = "Статья с таким названием у вас уже есть"
				return
			}
		} else {
			panic(fmt.Errorf("Ошибка. При записи статьи в БД: %v", err.Error()))
		}

	}

	c.response.Сompleted = true
	c.response.Message = "Статья добавлена"
}

func addCourseWorkName(c *studentContext, rw web.ResponseWriter, req *web.Request) {
	theme := req.FormValue("theme")
	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		c.response.Message = "Укажите верный id курсовой работы"
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if theme == "" {
		c.response.Message = "Укажите верную тему курсовой работы"
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var confirmed bool
	err = g.DB.QueryRow(`SELECT confirmed FROM courseworks WHERE id = $1`, id).Scan(&confirmed)
	if err == sql.ErrNoRows {
		c.response.Message = "Курсовой работы с указанным id не существует"
		rw.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке курсовой работы с id = %v в БД: %v", id, err.Error()))
	}

	if !confirmed {
		_, err := g.DB.Exec(`UPDATE courseworks SET theme = $1 WHERE id = $2`, theme, id)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При обновлении темы курсовой работы в БД: %v", err.Error()))
		}
	} else {
		c.response.Message = "Нельзя изменить название подтвержденной курсовой работы"
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	c.response.Сompleted = true
	c.response.Message = "Тема курсовой работы успешно добавлена"
}
