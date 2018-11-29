package contexts

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
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
	stringID := req.FormValue("id")
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

func (c *VerifContext) confirmArticle(rw web.ResponseWriter, req *web.Request) {

	article, errStr := u.ValidateArticle(req.FormValue("name"), req.FormValue("journal"), req.FormValue("biblioRecord"), req.FormValue("type"))
	if article == nil {
		rw.WriteHeader(http.StatusBadRequest)
		c.response.Message = errStr
		return
	}

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
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
		rw.WriteHeader(http.StatusBadRequest)
		c.response.Message = "C указанным id статьи не существует"
		return
	}

	c.response.Сompleted = true
	c.response.Body = "Статья успешно подтверждена"

}
