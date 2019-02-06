package contexts

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	dbu "github.com/STEJLS/StudentAccount/dbutils"
	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
	u "github.com/STEJLS/StudentAccount/utils"
	"github.com/gocraft/web"
)

// Context - базовый контекст
type Context struct {
	response t.ResponseMessage
	user     t.User
}

// UserContext - контекст пользователя
type UserContext struct {
	*Context
}

// getAndValidateAuthData - достает из запроса логин, пароль и проверяет их на правила
func (c *AuthContext) getAndValidateAuthData(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.login = req.FormValue(g.LoginValueName)
	password := req.FormValue(g.PasswordValueName)

	if len(password) < g.MinPasswordLength {
		log.Printf("Инфо. Попытка авторизация с неверным паролем(password - `%v`)\n", c.login)
		c.response.Message = fmt.Sprintf("Длина пароля не может быть менее %v символов.", g.MinLoginLength)
		return
	}

	if len(c.login) < g.MinLoginLength {
		log.Printf("Инфо. Попытка авторизация с некорректным логином(login - `%v`)\n", c.login)
		c.response.Message = fmt.Sprintf("Длина логина не может быть менее %v символов.", g.MinLoginLength)
		return
	}

	c.password = u.GenerateMD5hash(password)

	next(rw, req)
}

// userRequire - middleware, которое требует пользователя
func (c *Context) userRequire(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	cookie, err := req.Cookie("token")

	if err == http.ErrNoCookie {
		c.response.Message = "Необходимо авторизоваться"
		return
	}

	if err != nil {
		panic(fmt.Errorf("Ошибка. При чтении cookie: " + err.Error()))
	}

	if cookie.Value == "" {
		http.SetCookie(rw, &http.Cookie{Name: "token", Expires: time.Now().UTC()})
		c.response.Message = "Необходимо авторизоваться"
		return
	}

	g.Lock.RLock()
	login, ok := g.Sessions[cookie.Value]
	g.Lock.RUnlock()

	if !ok || login == "" {
		http.SetCookie(rw, &http.Cookie{Name: "token", Expires: time.Now().UTC()})
		c.response.Message = "Необходимо авторизоваться"
		return
	}

	row := g.DB.QueryRow(`SELECT id, role, login, fullName, isActivated, id_faculty, id_department, id_student
						FROM users WHERE login = $1`, login)
	var user t.User
	err = row.Scan(&user.ID, &user.Role, &user.Login, &user.FullName, &user.IsActivated,
		&user.IDFaculty, &user.IDDepartment, &user.IDStudent)

	if err == sql.ErrNoRows { // Пользователя с таким логином нет.
		log.Printf("Ошибка. По логину из токена нет пользователя (login - %v )\n", login)
		c.response.Message = "Пользователя с таким логином не существует."
		return
	}

	if err != nil { // Ошибка при работе с БД
		panic(fmt.Errorf("Ошибка. При поиске в БД пользователя(login - %v ): %v", login, err.Error()))
	}

	c.user = user

	next(rw, req)
}

// adminRequire - middleware, которое требует прав админа
func (c *Context) adminRequire(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if c.user.Role != t.Admin {
		c.response.Message = "Недостаточно прав."
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	next(rw, req)
}

// verifRequire - middleware, которое требует прав верификатора
func (c *Context) verifRequire(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	if c.user.Role != t.Verificator {
		c.response.Message = "Недостаточно прав."
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	next(rw, req)
}

// panicHandler - middleware, которое обрабатывает панику
func (c *Context) panicHandler(rw web.ResponseWriter, req *web.Request, err interface{}) {
	c.response.Сompleted = false
	c.response.Body = nil
	c.response.Message = "Неполадки на сервере, повторите попытку позже."

	if pe, ok := err.(error); ok {
		log.Println(pe.Error())
	}

	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(u.ConvertToJSON(c.response))
}

// toJSON - middleware, которое отвечает за отдачу ответа в формате json
func (c *Context) toJSON(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	next(rw, req)

	if rw.Size() != 0 || rw.StatusCode() == http.StatusNotModified { // отдача статики
		return
	}

	rw.Header().Add("Content-type", "application/json;")
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	_, err := rw.Write(u.ConvertToJSON(c.response))
	if err != nil {
		panic(fmt.Errorf("Ошибка. При отдачи: %v", err.Error()))
	}
}

func (c *Context) getDepartments(rw web.ResponseWriter, req *web.Request) {
	c.response.Body = dbu.GetAllDepartments()
	c.response.Сompleted = true
}

func (c *Context) getFaculties(rw web.ResponseWriter, req *web.Request) {
	c.response.Body = dbu.GetAllFaculties()
	c.response.Сompleted = true
}
func (c *Context) notFound(rw web.ResponseWriter, req *web.Request) {
	c.response.Message = "Указанной страницы не существует"
	c.response.Сompleted = false
}

func GetRoots() *web.Router {
	rootPath, _ := os.Getwd()

	rootRouter := web.New(Context{}).
		Middleware((*Context).toJSON).
		Error((*Context).panicHandler).
		NotFound((*Context).notFound).
		Middleware(web.LoggerMiddleware).
		Middleware(web.StaticMiddleware(path.Join(rootPath, "public"), web.StaticOption{IndexFile: "index.html", Prefix: "/"}))

	rootRouter.Subrouter(AuthContext{}, "/account").
		Middleware((*AuthContext).getAndValidateAuthData).
		Post("/login", (*AuthContext).Login)

	rootRouter.Subrouter(AuthContext{}, "/account").
		Middleware((*AuthContext).userRequire).
		Post("/changePassword", ChangePassword).
		Post("/logout", (*AuthContext).Logout)

	adminRouter := rootRouter.Subrouter(UserContext{}, "/admin")
	adminRouter.Middleware((*UserContext).userRequire).
		Middleware((*UserContext).adminRequire).
		Post("verif", createVerif).
		Post("/faculties", addFacultiesFromCSV).
		Post("/departments", addDepartmentsFromCSV).
		Post("/fieldsOfStudy", addFieldsOfStudyFromCSV).
		Post("/students", addStudentsFromCSV).
		Post("/practicis", addPracticisFromCSV).
		Post("/сourseWorks", addCourseWorksFromCSV).
		Get("/tempPasswords", getTempPasswords).
		Get("/departments", (*UserContext).getDepartments).
		Get("/faculties", (*UserContext).getFaculties)

	verificatorRouter := rootRouter.Subrouter(VerifContext{}, "/verif")
	verificatorRouter.Middleware((*VerifContext).userRequire).
		Middleware((*VerifContext).verifRequire).
		Get("/departments", (*VerifContext).getDepartments).
		Get("/faculties", (*VerifContext).getFaculties).
		Get("/article/:article_id", (*VerifContext).getArticle).
		Get("/articlesForVerif", (*VerifContext).articlesForVerif).
		Get("/coursesForVerif", (*VerifContext).coursesForVerif).
		Post("/article", (*VerifContext).confirmArticle).
		Post("/cancelArticle", (*VerifContext).cancelArticle).
		Post("/course", (*VerifContext).confirmCourse).
		Post("/cancelCourse", (*VerifContext).cancelCourse)

	userRouter := rootRouter.Subrouter(studentContext{}, "/student")
	userRouter.Middleware((*studentContext).userRequire).
		Middleware((*studentContext).StudentRequire).
		Get("/marks", getStudentMarks).
		Get("/info", getStudentInfo).
		Get("/article/:article_id", (*studentContext).getArticle).
		Get("/practices", (*studentContext).getPractices).
		Get("/articles", (*studentContext).getArticles).
		Get("/courses", (*studentContext).getCourses).
		Post("/article", addArticle).
		Post("/courseWork", addCourseWorkName)
	return rootRouter
}
