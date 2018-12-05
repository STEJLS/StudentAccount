package types

// User - пользователь
type User struct {
	ID           int
	Role         int
	Login        string
	Password     string
	FullName     *string
	IsActivated  bool
	IDFaculty    *int
	IDDepartment *int
	IDStudent    *int
}

// Faculty - факультет
type Faculty struct {
	ID        int
	Name      string
	ShortName string
}

// Department - кафедра
type Department struct {
	ID        int
	IDFaculty int
	Name      string
	ShortName string
}

// FieldOfStudy - направление подготовки
type FieldOfStudy struct {
	ID           int
	IDDepartment int
	Code         string
	Name         string
	Alias        string
	Level        int
}

// Subject - предмет
type Subject struct {
	ID           int
	IDDepartment int
	IDField      int
	Name         string
	PassType     int
}

// Student - студент
type Student struct {
	ID              int
	IDfield         int
	Number          string
	Team            string
	DurationOfStudy int
	GroupNumber     int
}

// Marks - оценки
type Marks struct {
	IDuser    int
	IDSubject int
	Rating    int
	Repass    bool
}

// ResponseMarks - оценки для фронта
type ResponseMarks struct {
	Subject  string
	Rating   int
	PassType int
	Semester int
	Repass   bool
}

// Article - статья
type Article struct {
	id           int
	IDStudent    int
	Name         string
	Journal      string
	BiblioRecord string
	ArticleType  string
	FileName     string
	RealFileName string
	Confirmed    bool
}

const (
	// Admin - администратор
	Admin int = iota
	// Verificator - верификатор документов
	Verificator
	// StudentRole - студент
	StudentRole
)

// ResponseMessage - структура, которая представляет собой форму ответа каждого обработчика
type ResponseMessage struct {
	Сompleted bool
	Message   string
	Body      interface{}
}
