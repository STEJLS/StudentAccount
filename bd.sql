CREATE TABLE faculties
(
    id SERIAL PRIMARY KEY,
	name varchar NOT NULL UNIQUE,
    shortName varchar NOT NULL UNIQUE
);

CREATE TABLE departments
(
    id SERIAL PRIMARY KEY,
    id_faculty INTEGER NOT NULL UNIQUE,
	name varchar NOT NULL UNIQUE,
    shortName varchar NOT NULL UNIQUE,
    FOREIGN KEY(id_faculty) REFERENCES faculties(id) ON DELETE SET NULL
);

CREATE TABLE fieldsOfStudy
(
    id SERIAL PRIMARY KEY,
    id_department INTEGER NOT NULL,
    name varchar NOT NULL,
    alias varchar NOT NULL,
	code varchar NOT NULL,
    level smallint NOT NULL,
    FOREIGN KEY(id_department) REFERENCES departments(id) ON DELETE SET NULL
);
alter table fieldsofstudy add unique(id_department,code)

level 
0-бакалавр
1-специалист
2-магистр
3-аспирант

CREATE TABLE subjects
(
    id SERIAL PRIMARY KEY,
    id_department INTEGER NOT NULL,
    id_field INTEGER NOT NULL,
	name varchar NOT NULL,
    passType smallint NOT NULL,
    FOREIGN KEY(id_department) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY(id_field) REFERENCES fieldsOfStudy(id) ON DELETE SET NULL
);

ALTER TABLE subjects ADD UNIQUE(id_field,name,passType,id_department)
passType
0-экз
1-д.зачет
2-зачет

CREATE TABLE Users
(
    id SERIAL PRIMARY KEY,
    role smallint NOT NULL,
    login varchar NOT NULL UNIQUE,
    password varchar  NOT NULL,
    fullName varchar,
    isActivated BOOLEAN NOT NULL DEFAULT false,
    id_faculty INTEGER,
    id_department INTEGER,
    id_student INTEGER,

    FOREIGN KEY(id_faculty) REFERENCES faculties(id) ON DELETE SET NULL,
    FOREIGN KEY(id_department) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY(id_student) REFERENCES students(id) ON DELETE SET NULL
);

CREATE TABLE students
(
    id SERIAL PRIMARY KEY,
    id_field INTEGER NOT NULL,
    number varchar NOT NULL UNIQUE,
    team varchar NOT NULL,
    durationOfStudy smallint NOT NULL, 
    teamNumber smallint NOT NULL,
    FOREIGN KEY(id_field) REFERENCES fieldsOfStudy(id) ON DELETE SET NULL
);

CREATE TABLE marks
(
    id_student BIGINT NOT NULL,
    id_subject BIGINT NOT NULL,
    rating smallint NOT NULL,
    semester smallint NOT NULL,
    repass BOOLEAN NOT NULL DEFAULT false,
    FOREIGN KEY(id_student) REFERENCES students(id),
    FOREIGN KEY(id_subject) REFERENCES subjects(id)
);

CREATE TABLE articles
(
    id SERIAL PRIMARY KEY,
    id_student BIGINT NOT NULL,
    name varchar,
    journal varchar,
    biblioRecord varchar,
    type varchar,
    fileName varchar,   
    realFileName varchar,
    confirmed bool DEFAULT false,
    FOREIGN KEY(id_student) REFERENCES students(id)
);
ALTER TABLE subjects ADD UNIQUE(name,id_student)
Типы: конференция, конференция РИНЦ, статья РИНЦ, статья ВАК, статья Scopus, статья Web of Science




CREATE TABLE practicis
(
    id SERIAL PRIMARY KEY,
    id_student BIGINT NOT NULL,
    semester smallint NOT NULL,
    name varchar, 
    head varchar,
    company varchar, 
    begin_date date,
    end_date  date,
    rating smallint,
    FOREIGN KEY(id_student) REFERENCES students(id)
);
ALTER TABLE practicis ADD UNIQUE(id_student,semester)

CREATE TABLE courseworks
(
    id SERIAL PRIMARY KEY,
    id_student BIGINT NOT NULL,
    id_subject INTEGER NOT NULL,
    semester smallint NOT NULL,
    theme varchar, 
    head varchar,
    rating smallint,
    confirmed bool DEFAULT false,
    FOREIGN KEY(id_student) REFERENCES students(id),
    FOREIGN KEY(id_subject) REFERENCES subjects(id)
);

ALTER TABLE courseworks ADD UNIQUE(id_subject,semester)


-- команда создания скрипта БД
-- pg_dump.exe --host=localhost --port=5432 --username=postgres --schema-only --file=D:\schChema.sql bmstu