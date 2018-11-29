package XMLconfig

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Config - это основная структура для парсинга xml файла
type Config struct {
	HTTP  Http     `xml:"http"`
	Admin Admin    `xml:"admin"`
	DB    DataBase `xml:"DataBase"`
}

// Http - это структура для парсинга
// информации об http из xml файла
type Http struct {
	XMLName xml.Name `xml:"http"`
	Port    int      `xml:"port,attr"`
	Host    string   `xml:"host,attr"`
}

// Http - это структура для парсинга
// информации об учетной записи админа из xml файла
type Admin struct {
	XMLName  xml.Name `xml:"admin"`
	Login    string   `xml:"login,attr"`
	Password string   `xml:"password,attr"`
}

// DataBase - это структура для парсинга
// информации об базеданных из xml файла
type DataBase struct {
	XMLName  xml.Name `xml:"DataBase"`
	Driver   string   `xml:"driver"`
	User     string   `xml:"user"`
	Password string   `xml:"password"`
	DBname   string   `xml:"dbname"`
	Host     string   `xml:"host"`
	Port     int      `xml:"port"`
	SSLmode  string   `xml:"sslmode"`
}

// Get - это функция парсит xml конфиг, находящийся в файле "source"
// а также проверяет его на правильность
func Get(source string) Config {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Фатал. При открытии xml файла(%v) для парсинга: ", source) + err.Error())
	}

	var config Config
	err = xml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Фатал. При анмаршалинге xml файла(%q): ", source) + err.Error())
	}

	log.Printf("Инфо. Файл %q успешно расперсен.", source)

	err = validating(config)
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

// Validating - это функция которая проверяет введенную информацию из конфига
func validating(config Config) error {
	if config.HTTP.Port < 1024 || config.HTTP.Port >= 65535 {
		return fmt.Errorf("Фатал. Не валидный номер http порта(от 1024 до 65535), а вы ввели %v", config.HTTP.Port)
	}

	if strings.ContainsAny(config.DB.DBname, "/\\.\"*<>:|?$,'") {
		return fmt.Errorf("Фатал. Не валидное имя базы данных(не должно быть символов /, \\, ., \", *, <, >, :, |, ?, $), введено: %q", config.DB.DBname)
	}

	if config.DB.Port < 1024 || config.DB.Port >= 65535 {
		return fmt.Errorf("Фатал. Не валидный номер http порта(от 1024 до 65535), а вы ввели %v", config.HTTP.Port)
	}

	if config.DB.User == "" {
		return fmt.Errorf("Фатал. Не валидное имя пользователя базы данных")
	}

	if config.DB.Password == "" {
		return fmt.Errorf("Фатал. Не валидный пароль от базы данных")
	}

	if config.DB.DBname == "" {
		return fmt.Errorf("Фатал. Пустое имя базы")
	}

	log.Printf("Инфо. Конфиг успешно прошел проверку.")
	return nil
}
