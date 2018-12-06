package contexts

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	g "github.com/STEJLS/StudentAccount/globals"
	u "github.com/STEJLS/StudentAccount/utils"
	"github.com/gocraft/web"
)

type VerifContext struct {
	*Context
}

func (c *Context) getArticle(rw web.ResponseWriter, req *web.Request) {
	stringID := req.PathParams["article_id"]
	if stringID == "" {
		rw.WriteHeader(http.StatusBadRequest)
		c.response.Message = "Укажите id статьи"
		return
	}
	id, err := strconv.Atoi(stringID)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		c.response.Message = "Укажите валидный id статьи"
		return
	}

	var localFileName string
	var realFileName string

	err = g.DB.QueryRow(`SELECT fileName, realFileName FROM articles WHERE id =$1`, id).Scan(&localFileName, &realFileName)
	if err == sql.ErrNoRows {
		rw.WriteHeader(http.StatusBadRequest)
		c.response.Message = "Статьи с указанным id не существует"
		return
	} else if err != nil {
		panic(fmt.Errorf("Ошибка. Во время выборки статьи: %v", err.Error()))
	}

	data, err := ioutil.ReadFile(path.Join(g.ArticlesDirectory, localFileName))
	if err != nil {
		panic(fmt.Errorf("Ошибка. При чтении статьи: %v", err.Error()))
	}

	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\""+"%s"+"\"", realFileName))
	rw.Header().Add("Content-type", "application/pdf")
	rw.Header().Add("Content-Length", fmt.Sprintf("%v", len(data)))

	rw.Write(data)
	c.notJSON = true
}

func (c *VerifContext) cancelArticle(rw web.ResponseWriter, req *web.Request) {
	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		c.response.Message = "id должно быть числом"
		return
	}

	var fileName string
	err = g.DB.QueryRow(`SELECT filename FROM articles WHERE id = $1`, id).Scan(&fileName)
	if err == sql.ErrNoRows {
		c.response.Message = "C указанным id статьи не существует"
		return
	} else if err != nil {
		panic(fmt.Errorf("Ошибка. Во время отклонения статьи: %v", err.Error()))
	}

	_, err = g.DB.Exec(`DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		panic(fmt.Errorf("Ошибка. Во время отклонения статьи: %v", err.Error()))
	}

	os.Remove(path.Join(g.ArticlesDirectory, fileName))

	c.response.Сompleted = true
	c.response.Body = "Статья успешно отклонена"
}

func (c *VerifContext) cancelCourse(rw web.ResponseWriter, req *web.Request) {
	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		c.response.Message = "id должно быть числом"
		return
	}

	r, err := g.DB.Exec(`UPDATE courseworks 
					   SET theme = null, confirmed = false  
					   WHERE id = $1`, id)
	if err != nil {
		panic(fmt.Errorf("Ошибка. Во время отклонения курсовой работы: %v", err.Error()))
	}

	if n, _ := r.RowsAffected(); n == 0 {
		c.response.Message = "C указанным id курсовой работы не существует"
		return
	}

	c.response.Сompleted = true
	c.response.Body = "Статья успешно отклонена"
}

func (c *VerifContext) confirmCourse(rw web.ResponseWriter, req *web.Request) {
	theme := req.FormValue("theme")
	if theme == "" {
		c.response.Message = "Необходимо указать тему"
		return
	}

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		c.response.Message = "id должно быть числом"
		return
	}

	r, err := g.DB.Exec(`UPDATE courseworks 
					   SET theme = $1, confirmed = true  
					   WHERE id = $2`, theme, id)
	if err != nil {
		panic(fmt.Errorf("Ошибка. Во время подтверждения курсовой работы: %v", err.Error()))
	}

	if n, _ := r.RowsAffected(); n == 0 {
		c.response.Message = "C указанным id курсовой работы не существует"
		return
	}

	c.response.Сompleted = true
	c.response.Body = "Курсовая работа успешно подтверждена"
}

func (c *VerifContext) confirmArticle(rw web.ResponseWriter, req *web.Request) {
	article, errStr := u.ValidateArticle(req.FormValue("name"), req.FormValue("journal"), req.FormValue("biblioRecord"), req.FormValue("type"))
	if article == nil {
		c.response.Message = errStr
		return
	}

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		c.response.Message = "id должно быть числом"
		return
	}

	r, err := g.DB.Exec(`UPDATE articles 
					   SET name = $1, journal = $2, biblioRecord = $3, type = $4, confirmed = true  
					   WHERE id = $5`,
		article.Name, article.Journal, article.BiblioRecord, article.ArticleType, id)
	if err != nil {
		panic(fmt.Errorf("Ошибка. Во время подтверждения статьи: %v", err.Error()))
	}

	if n, _ := r.RowsAffected(); n == 0 {
		c.response.Message = "C указанным id статьи не существует"
		return
	}

	c.response.Сompleted = true
	c.response.Body = "Статья успешно подтверждена"
}

func (c *VerifContext) articlesForVerif(rw web.ResponseWriter, req *web.Request) {
	rows, err := g.DB.Query(`SELECT a.id, u.fullname, (s.team||'-'|| s.teamnumber) as team, 
									a.name, a.journal, a.bibliorecord, a.type FROM articles as a
							JOIN students s ON s.id = a.id_student
							JOIN users u ON u.id_student = s.id
							WHERE confirmed = false 
							AND a.id_student IN (
								SELECT id FROM students where id_field in(
								SELECT id FROM fieldsofstudy WHERE id_department = $1))`,
		c.user.IDDepartment)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке статей: %v", err.Error()))
	}
	defer rows.Close()

	articlesInfo := make([]*struct {
		ID           int
		FIO          string
		Team         string
		Name         string
		Journal      string
		BiblioRecord string
		ArticlType   string
	}, 0)

	for rows.Next() {
		articleInfo := new(struct {
			ID           int
			FIO          string
			Team         string
			Name         string
			Journal      string
			BiblioRecord string
			ArticlType   string
		})

		err = rows.Scan(&articleInfo.ID, &articleInfo.FIO, &articleInfo.Team, &articleInfo.Name,
			&articleInfo.Journal, &articleInfo.BiblioRecord, &articleInfo.ArticlType)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке статей: %v", err.Error()))
		}

		articlesInfo = append(articlesInfo, articleInfo)
	}

	c.response.Body = articlesInfo
	c.response.Сompleted = true
}

func (c *VerifContext) coursesForVerif(rw web.ResponseWriter, req *web.Request) {
	rows, err := g.DB.Query(`SELECT subj.name, cw.id, u.fullname, (s.team||'-'|| s.teamnumber) as team, 
								cw.semester, cw.theme, cw.head, cw.rating FROM courseworks as cw
							JOIN students s ON s.id = cw.id_student
							JOIN users u ON u.id_student = s.id
							JOIN subjects subj ON subj.id = cw.id_subject
							WHERE confirmed = false AND cw.theme IS NOT NULL
							AND cw.id_student IN (
								SELECT id FROM students where id_field in(
								SELECT id FROM fieldsofstudy WHERE id_department = $1))`,
		c.user.IDDepartment)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке практик: %v", err.Error()))
	}
	defer rows.Close()

	practicesInfo := make([]*struct {
		ID       int
		Subject  string
		FIO      string
		Team     string
		Semester int
		Theme    string
		Head     string
		Rating   int
	}, 0)

	for rows.Next() {
		practiceInfo := new(struct {
			ID       int
			Subject  string
			FIO      string
			Team     string
			Semester int
			Theme    string
			Head     string
			Rating   int
		})

		err = rows.Scan(&practiceInfo.Subject, &practiceInfo.ID, &practiceInfo.FIO, &practiceInfo.Team, &practiceInfo.Semester,
			&practiceInfo.Theme, &practiceInfo.Head, &practiceInfo.Rating)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке практик: %v", err.Error()))
		}

		practicesInfo = append(practicesInfo, practiceInfo)
	}

	c.response.Body = practicesInfo
	c.response.Сompleted = true
}
