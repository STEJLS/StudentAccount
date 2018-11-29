package utils

import (
	"crypto/md5"
	"encoding/csv"

	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"

	"github.com/STEJLS/StudentAccount/XMLconfig"
	g "github.com/STEJLS/StudentAccount/globals"
	t "github.com/STEJLS/StudentAccount/types"
)

// InitFlags - инициализирует флаги командной строки.
func InitFlags() {
	flag.StringVar(&g.LogSource, "log_source", "log.txt", "Source for log file")
	flag.StringVar(&g.ConfigSource, "config_source", "config.xml", "Source for config file")
	flag.Parse()
}

// InitLogger - инициализирует логгер.
func InitLogger() *os.File {
	logfile, err := os.OpenFile(g.LogSource, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Ошибка. Файл логов (%q) не открылся: ", g.LogSource, err)
	}

	log.SetOutput(logfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return logfile
}

// connectTog.DB - устанавливет соединение с БД и инициализирует глобальные переменные
func ConnectToDB(dbInfo XMLconfig.DataBase) {
	var err error
	connStr := fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%v sslmode=%v", dbInfo.User, dbInfo.Password, dbInfo.DBname,
		dbInfo.Host, strconv.Itoa(dbInfo.Port), dbInfo.SSLmode)

	if g.DB == nil {
		g.DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Фатал. При подключении к серверу БД(%v:%v): ", dbInfo.Host, dbInfo.Port) + err.Error())
		}
	}

	err = g.DB.Ping()
	if err != nil {
		log.Fatalln("Фатал. При пинге сервера БД: " + err.Error())
	}

	log.Printf("Инфо. Подключение к базе данных установлено.")
}

func GenerateMD5hash(password string) string {
	md5hash := md5.Sum([]byte(password))
	temp := make([]byte, 0, md5.Size+len(g.Salt))
	for _, item := range md5hash {
		temp = append(temp, item)
	}
	for _, item := range g.Salt {
		temp = append(temp, item)
	}

	return fmt.Sprintf("%x", md5.Sum(temp))
}

// initgDB - инициализирует базу данных начальными значениями
func InitDB(login string, password string) {
	_, err := g.DB.Exec(`INSERT INTO users (login, password, role, isActivated) 
					   VALUES ($1, $2, 0, true) ON CONFLICT (login) 
					   DO UPDATE SET password = $2`, login, GenerateMD5hash(password))
	if err != nil {
		log.Fatalf("Фатал.Ошибка при создании админа: %v", err.Error())
	}
}

// generateToken - генерирует уникальный токен для авторизации.
func GenerateToken() string {
	token, err := uuid.NewV4()

	if err != nil {
		panic(fmt.Errorf("Ошибка. При генерации токена авторизации: " + err.Error()))
	}

	return token.String()
}

// +generateTempPassword - создает временный пароль длиной 8
func GenerateTempPassword() string {
	return GenerateToken()[0:8]
}

// ConvertToJSON - Преобразует данные в json
func ConvertToJSON(data interface{}) []byte {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка. При маршалинге в json результата: %v\n", err.Error())
		panic(errors.New("Неполадки на сервере, повторите попытку позже"))
	}

	return jsonData
}

func InitPasswordCSVWriter() (*os.File, *csv.Writer) {
	passwordFile, err := os.OpenFile(g.PasswordFileStorageName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При открытии файла с временными паролями: %v", err.Error()))
	}

	passwordWriter := csv.NewWriter(passwordFile)
	if err := passwordWriter.Error(); err != nil {
		panic(fmt.Errorf("Ошибка. При создании csv writer для временных паролей: %v", err.Error()))
	}

	return passwordFile, passwordWriter
}

func InitFiles() {
	if _, err := os.Stat(g.DataDirectoryName); os.IsNotExist(err) {
		err = os.Mkdir(g.DataDirectoryName, 0666)
		if err != nil {
			log.Fatalf("Ошибка. При создании папки для данных с именем - %v: %v", g.DataDirectoryName, err.Error())
		}
	}

	if _, err := os.Stat(g.PasswordFileStorageName); os.IsNotExist(err) {
		file, err := os.Create(g.PasswordFileStorageName)
		if err != nil {
			log.Fatalf("Ошибка. При создании файла для паролей с именем- %v: %v", g.PasswordFileStorageName, err.Error())
		}
		file.Close()
		if err != nil {
			log.Fatalf("Ошибка. При закрытии файла для паролей с именем- %v: %v", g.PasswordFileStorageName, err.Error())
		}
	}

	if _, err := os.Stat(g.ArticlesDirectory); os.IsNotExist(err) {
		err = os.Mkdir(g.ArticlesDirectory, 0666)
		if err != nil {
			log.Fatalf("Ошибка. При создании папки для данных с именем - %v: %v", g.DataDirectoryName, err.Error())
		}
	}

}

func ValidateArticle(name string, journal string, biblioRecord string, articleType string) (*t.Article, string) {
	if name == "" {
		return nil, "Укажите название статьи"
	}

	if journal == "" {
		return nil, "Укажите название журнала"
	}

	if biblioRecord == "" {
		return nil, "Укажите библиографическую ссылку"
	}

	if articleType == "" {
		return nil, "Укажите тип журнала"
	}

	switch articleType {
	case "Конференция":
	case "Конференция РИНЦ":
	case "Статья ВАК":
	case "Статья Scopus":
	case "Статья Web of Science":
	case "Статья РИНЦ":
		break
	default:
		return nil, "Укажите верный тип журнала"
	}

	return &t.Article{Name: name,
			Journal:      journal,
			BiblioRecord: biblioRecord,
			ArticleType:  articleType},
		""
}

func ValidateVerif(login string, password string, fullName string, IDFaculty string, IDDepartment string) (*t.User, string) {
	if len(password) < g.MinPasswordLength {
		return nil, fmt.Sprintf("Длина пароля не может быть менее %v символов.", g.MinPasswordLength)
	}
	if len(login) < g.MinLoginLength {
		return nil, fmt.Sprintf("Длина логина не может быть менее %v символов.", g.MinLoginLength)
	}

	if fullName == "" {
		return nil, "Имя не может быть пустой строкой"
	}
	idF, err := strconv.Atoi(IDFaculty)
	if err != nil {
		return nil, "id факультета должно быть числом"
	}

	var exist bool

	err = g.DB.QueryRow(`SELECT EXISTS(SELECT * FROM users WHERE login = $1)`, login).Scan(&exist)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке пользователя с логином = %v в БД: %v", login, err.Error()))
	}

	if exist {
		return nil, "Указанный логин занят"
	}

	err = g.DB.QueryRow(`SELECT EXISTS(SELECT * FROM faculties WHERE id = $1)`, idF).Scan(&exist)
	if err != nil {
		panic(fmt.Errorf("Ошибка. При выборке факультета с id = %v в БД: %v", idF, err.Error()))
	}

	if !exist {
		return nil, "Указанного факультета не существует"
	}

	var IDD *int

	if IDDepartment != "" {
		idD, err := strconv.Atoi(IDDepartment)
		if err != nil {
			return nil, "id кафедры должно быть числом"
		}

		err = g.DB.QueryRow(`SELECT EXISTS(SELECT * FROM departments WHERE id = $1)`, idD).Scan(&exist)
		if err != nil {
			panic(fmt.Errorf("Ошибка. При выборке кафедры с id = %v в БД: %v", idF, err.Error()))
		}

		if !exist {
			return nil, "Указанной кафедры не существует"
		}
		IDD = &idD
	}

	return &t.User{Login: login,
			Password:     GenerateMD5hash(password),
			FullName:     &fullName,
			IDFaculty:    &idF,
			IDDepartment: IDD,
			Role:         t.Verificator,
			IsActivated:  true},
		""
}
