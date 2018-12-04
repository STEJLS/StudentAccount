package contexts

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
	u "github.com/STEJLS/StudentAccount/utils"

	"github.com/gocraft/web"
)

type AuthContext struct {
	*Context
	login    string
	password string
}

// Login - авторизация пользователя
func (c *AuthContext) Login(rw web.ResponseWriter, req *web.Request) {
	row := g.DB.QueryRow(`SELECT id, role, login, password, fullName, isActivated, id_faculty, id_department, id_student 
						FROM users WHERE login = $1`, c.login)
	var user t.User
	err := row.Scan(&user.ID, &user.Role, &user.Login, &user.Password, &user.FullName,
		&user.IsActivated, &user.IDFaculty, &user.IDDepartment, &user.IDStudent)

	if err == sql.ErrNoRows { // Пользователя с таким логином нет.
		log.Printf("Инфо. Попытка авторизации по несуществующему логину (login - %v )\n", c.login)
		c.response.Message = "Пользователя с таким логином не существует."
		return
	}

	if err != nil { // Ошибка при работе с БД
		panic(fmt.Errorf("Ошибка. При поиске в БД пользователя(login - %v ): %v", c.login, err.Error()))
	}

	if c.password != user.Password { // Пользователь ввел неверный пароль
		log.Printf("Инфо. Попытка авторизация с неверным паролем(login - %v )\n", c.login)
		c.response.Message = "Неверный пароль."
		return
	}

	token := u.GenerateToken()
	g.Sessions[token] = c.login

	http.SetCookie(rw, &http.Cookie{Name: "token", Value: token, Path: "/"})

	log.Printf("Инфо. Пользователь %v авторизовался.", c.login)
	user.Password = ""
	c.response.Сompleted = true
	c.response.Message = "Вы успешно авторизовались в системе"
	c.response.Body = user
}

// Logout - выход пользователя из системы
func (c *AuthContext) Logout(rw web.ResponseWriter, req *web.Request) {
	cookie, err := req.Cookie("token")

	if err == http.ErrNoCookie {
		c.response.Message = "Необходимо авторизоваться"
		return
	}

	if err != nil {
		panic(fmt.Errorf("Ошибка. При чтении cookie: " + err.Error()))
	}

	http.SetCookie(rw, &http.Cookie{Name: "token", Expires: time.Now().UTC()})

	if cookie.Value == "" {
		c.response.Message = "Необходимо авторизоваться"
		return
	}

	g.Lock.RLock()
	delete(g.Sessions, cookie.Value)
	g.Lock.RUnlock()
	c.response.Сompleted = true
	c.response.Message = "Вы успешно вышли из системы"
}

func ChangePassword(c *AuthContext, rw web.ResponseWriter, req *web.Request) {
	newPassword := req.FormValue(g.NewPasswordValueName)
	oldPassword := req.FormValue(g.PasswordValueName)
	if len(newPassword) < g.MinPasswordLength && len(oldPassword) < g.MinPasswordLength {
		c.response.Message = fmt.Sprintf("Длина пароля не может быть менее %v символов.", g.MinPasswordLength)
		return
	}

	r, err := g.DB.Exec(`UPDATE users SET password = $1, isActivated = true WHERE id = $2 AND password = $3`,
		u.GenerateMD5hash(newPassword), c.user.ID, u.GenerateMD5hash(oldPassword))
	if err != nil {
		panic(fmt.Errorf("Ошибка. При изменения пароля пользователю(login - %v ): %v", c.user.Login, err.Error()))
	}

	countAffected, err := r.RowsAffected()
	if err != nil {
		panic(fmt.Errorf("Ошибка. При изменения пароля пользователю(login - %v ): %v", c.user.Login, err.Error()))
	}

	if countAffected == 0 {
		c.response.Message = "Старый пароль указан неверно"
		return
	}

	c.response.Сompleted = true
	c.response.Message = "Ваш пароль успешно изменен"
}
