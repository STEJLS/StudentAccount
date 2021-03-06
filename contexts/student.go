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

	rows, err := g.DB.Query(`SELECT subjects.name, marks.rating, subjects.passtype, marks.semester, marks.repass
			  			   FROM students
						   JOIN marks ON students.id = marks.id_student
						   JOIN subjects ON marks.id_subject = subjects.id
						   WHERE students.id = $1`, c.user.IDStudent)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке оценок пользователя(login - %v ): %v", c.user.Login, err.Error()))
	}
	defer rows.Close()

	marks := make([]*t.ResponseMarks, 0)

	for rows.Next() {
		mark := t.ResponseMarks{}
		err = rows.Scan(&mark.Subject, &mark.Rating, &mark.PassType, &mark.Semester, &mark.Repass)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке оценок пользователя(login - %v ): %v", c.user.Login, err.Error()))
		}

		marks = append(marks, &mark)
	}

	c.response.Body = marks
	c.response.Completed = true
}

func getStudentInfo(c *studentContext, rw web.ResponseWriter, req *web.Request) {
	row := g.DB.QueryRow(`SELECT s.number, u.fullname, s.team, s.teamnumber,
	fi.name, fi.code, fi.profile, fi.level,
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
		FieldProfile        string
		level               int
		DepartmentName      string
		DepartmentShortName string
		FacultyName         string
		FacultyShortName    string
	}

	err := row.Scan(&userInfo.Number, &userInfo.FullName, &userInfo.Team, &userInfo.TeamNumber, &userInfo.FieldName,
		&userInfo.FieldCode, &userInfo.FieldProfile, &userInfo.level, &userInfo.DepartmentName, &userInfo.DepartmentShortName, &userInfo.FacultyName,
		&userInfo.FacultyShortName)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке информации о пользователе(login - %v ): %v", c.user.Login, err.Error()))
	}

	c.response.Body = userInfo
	c.response.Completed = true
}

func addArticle(c *studentContext, rw web.ResponseWriter, req *web.Request) {
	name := req.FormValue("name")
	journal := req.FormValue("journal")
	biblioRecord := req.FormValue("biblioRecord")
	articleType := req.FormValue("type")

	article, errStr := u.ValidateArticle(name, journal, biblioRecord, articleType)
	if article == nil {
		c.response.Message = errStr
		return
	}

	file, header, err := req.FormFile("article")
	isRecieved := true
	realName := ""
	if err != nil {
		isRecieved = false
	}

	fileName := ""
	if isRecieved {
		realName = header.Filename
		if filepath.Ext(header.Filename) != ".pdf" {
			c.response.Message = "Файл должен иметь формат pdf"
			return
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			c.response.Message = "Файл не получен"
			return
		}

		fileName = u.GenerateToken()

		err = ioutil.WriteFile(path.Join(g.ArticlesDirectory, fileName), data, 0666)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При записи статьи на диск: %v", err.Error()))
		}
	}

	_, err = g.DB.Exec(`INSERT INTO articles(name, journal, biblioRecord, type, fileName, realFileName,id_student) 
			 VALUES($1, $2, $3, $4, $5, $6, $7)`, article.Name, article.Journal, article.BiblioRecord, article.ArticleType, fileName, realName, c.user.IDStudent)
	if err != nil {
		if isRecieved {
			os.Remove(path.Join(g.ArticlesDirectory, fileName))
		}

		if pe, ok := err.(*pq.Error); ok { // Нарушение уникальности
			if pe.Code == "23505" {
				c.response.Message = "Статья с таким названием у вас уже есть"
				return
			}
		} else {
			panic(fmt.Errorf("Ошибка. При записи статьи в БД: %v", err.Error()))
		}

	}

	c.response.Completed = true
	c.response.Message = "Статья добавлена"
}

func addCourseWorkName(c *studentContext, rw web.ResponseWriter, req *web.Request) {
	theme := req.FormValue("theme")
	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		c.response.Message = "Укажите верный id курсовой работы"
		return
	}

	if theme == "" {
		c.response.Message = "Укажите верную тему курсовой работы"
		return
	}

	var confirmed bool
	err = g.DB.QueryRow(`SELECT confirmed FROM courseworks WHERE id = $1`, id).Scan(&confirmed)
	if err == sql.ErrNoRows {
		c.response.Message = "Курсовой работы с указанным id не существует"
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
		return
	}

	c.response.Completed = true
	c.response.Message = "Тема курсовой работы успешно добавлена"
}

func (c *studentContext) getPractices(rw web.ResponseWriter, req *web.Request) {
	rows, err := g.DB.Query(`SELECT semester, name, head, company,
							 'c '||to_char(begin_date, 'DD-MM-YYYY') || ' по ' || to_char(end_date, 'DD-MM-YYYY'), rating
							FROM practicis WHERE id_student = $1 ORDER BY semester`, c.user.IDStudent)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке практик: %v", err.Error()))
	}
	defer rows.Close()

	practicesInfo := make([]*struct {
		Semester int
		Name     string
		Head     string
		Company  string
		Date     string
		Rating   int
	}, 0)

	for rows.Next() {
		practiceInfo := new(struct {
			Semester int
			Name     string
			Head     string
			Company  string
			Date     string
			Rating   int
		})

		err = rows.Scan(&practiceInfo.Semester, &practiceInfo.Name, &practiceInfo.Head, &practiceInfo.Company,
			&practiceInfo.Date, &practiceInfo.Rating)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке практик: %v", err.Error()))
		}

		practicesInfo = append(practicesInfo, practiceInfo)
	}

	c.response.Body = practicesInfo
	c.response.Completed = true
}

func (c *studentContext) getArticles(rw web.ResponseWriter, req *web.Request) {
	rows, err := g.DB.Query(`SELECT id, name, journal, bibliorecord, type, filename, confirmed 
							 FROM articles
							 WHERE id_student = $1`, c.user.IDStudent)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке статей: %v", err.Error()))
	}
	defer rows.Close()

	articlesInfo := make([]*struct {
		ID           int
		Name         string
		Journal      string
		BiblioRecord string
		ArticlType   string
		FileName     string
		Confirmed    bool
	}, 0)

	for rows.Next() {
		articleInfo := new(struct {
			ID           int
			Name         string
			Journal      string
			BiblioRecord string
			ArticlType   string
			FileName     string
			Confirmed    bool
		})

		err = rows.Scan(&articleInfo.ID, &articleInfo.Name, &articleInfo.Journal, &articleInfo.BiblioRecord,
			&articleInfo.ArticlType, &articleInfo.FileName, &articleInfo.Confirmed)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке статей: %v", err.Error()))
		}

		articlesInfo = append(articlesInfo, articleInfo)
	}

	c.response.Body = articlesInfo
	c.response.Completed = true
}

func (c *studentContext) getCourses(rw web.ResponseWriter, req *web.Request) {
	rows, err := g.DB.Query(`SELECT subj.name, cw.id, u.fullname, (s.team||'-'|| s.teamnumber) as team, 
								cw.semester, cw.theme, cw.head, cw.rating, cw.confirmed FROM courseworks as cw
							JOIN students s ON s.id = cw.id_student
							JOIN users u ON u.id_student = s.id
							JOIN subjects subj ON subj.id = cw.id_subject							
							WHERE cw.id_student = $1`, c.user.IDStudent)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке курсовых работ: %v", err.Error()))
	}
	defer rows.Close()

	coursesInfo := make([]*struct {
		ID        int
		Subject   string
		FIO       string
		Team      string
		Semester  int
		Theme     sql.NullString
		Head      string
		Rating    int
		Confirmed bool
	}, 0)

	for rows.Next() {
		courseInfo := new(struct {
			ID        int
			Subject   string
			FIO       string
			Team      string
			Semester  int
			Theme     sql.NullString
			Head      string
			Rating    int
			Confirmed bool
		})

		err = rows.Scan(&courseInfo.Subject, &courseInfo.ID, &courseInfo.FIO, &courseInfo.Team, &courseInfo.Semester,
			&courseInfo.Theme, &courseInfo.Head, &courseInfo.Rating, &courseInfo.Confirmed)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке курсовых работ: %v", err.Error()))
		}

		coursesInfo = append(coursesInfo, courseInfo)
	}

	c.response.Body = coursesInfo
	c.response.Completed = true
}

func (c *studentContext) getFOSandRPDList(rw web.ResponseWriter, req *web.Request) {
	rows, err := g.DB.Query(`SELECT t1.name, t1.id, t2.id  from (SELECT name, id  from documents where     id_faculty = $1 and
		id_department = $2 and
		id_field  = (SELECT id_field FROM students where id = $3) and type = 0) as t1
	join (SELECT name, id  from documents where     id_faculty = $4 and
		id_department = $5 and
		id_field  = (SELECT id_field FROM students where id = $6) and type = 1) as t2 on t2.name = t1.name`,
		c.user.IDFaculty, c.user.IDDepartment, c.user.IDStudent, c.user.IDFaculty, c.user.IDDepartment, c.user.IDStudent)

	if err != nil {
		if err == sql.ErrNoRows {
			c.response.Completed = true
			c.response.Message = "Для вашего направления подготовки еще нет документов"
			return
		}
		panic(fmt.Errorf("Ошибка. При выборке списка ФОС и РПД документов: %v", err.Error()))
	}

	result := make([]*struct {
		Name  string
		FosID int
		RpdID int
	}, 0)

	for rows.Next() {
		docInfo := new(struct {
			Name  string
			FosID int
			RpdID int
		})

		err = rows.Scan(&docInfo.Name, &docInfo.FosID, &docInfo.RpdID)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке ФОС или РПД произошла ошибка: %v", err.Error()))
		}

		result = append(result, docInfo)
	}

	if len(result) == 0 {
		c.response.Message = "Для вашего направления подготовки еще нет документов"
	}

	c.response.Body = result
	c.response.Completed = true
}

func (c *studentContext) getDocument(rw web.ResponseWriter, req *web.Request) {
	id, err := strconv.Atoi(req.PathParams["document_id"])
	if err != nil {
		c.response.Message = "Укажите верный id документа"
		return
	}

	var path string
	err = g.DB.QueryRow(`SELECT path FROM documents WHERE id = $1`, id).Scan(&path)
	if err != nil {
		if err == sql.ErrNoRows {
			c.response.Message = "С указанным id файла не существует"
			return
		}
		panic(fmt.Errorf("Ошибка. При выборке пути документа с id = %v: %v", id, err.Error()))
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При чтении статьи: %v", err.Error()))
	}

	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\""+"%s"+"\"", filepath.Base(path)))
	rw.Header().Add("Content-type", "application/pdf")
	rw.Header().Add("Content-Length", fmt.Sprintf("%v", len(data)))

	rw.Write(data)
}
