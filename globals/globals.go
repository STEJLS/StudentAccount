package globals

import (
	"database/sql"
	"path"
	"sync"
)

// logFileName - имя файла для логов, задается через флаг командной строки.
var LogSource string

// ConfigSource - имя файла для конфига, задается через флаг командной строки.
var ConfigSource string

// salt - соль для пароля.
var Salt = [12]byte{152, 123, 2, 1, 6, 84, 216, 35, 140, 158, 69, 128}

//DB - глобальная переменная подключения к Бд
var DB *sql.DB

// Sessions - карта для авторизации пользователей. Ключ токен, а значение - логин.
var Sessions = make(map[string]string)

// lock - Мьютекс для корректной параллельной работы с картой sessions.
var Lock = new(sync.RWMutex)

var MinLoginLength = 5
var MinPasswordLength = 6
var LoginValueName = "login"
var PasswordValueName = "password"
var NewPasswordValueName = "newPassword"

var DataDirectoryName = path.Join(".", "data")
var DocumentsDirectoryName = path.Join(DataDirectoryName, "documents")
var FOSDirectoryName = path.Join(DocumentsDirectoryName, "ФОС")
var RPDDirectoryName = path.Join(DocumentsDirectoryName, "РПД")
var PasswordFileStorageName = path.Join(DataDirectoryName, "passwords.csv")
var ArticlesDirectory = path.Join(DataDirectoryName, "articles")
