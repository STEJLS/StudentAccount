--
-- PostgreSQL database dump
--

-- Dumped from database version 10.5
-- Dumped by pg_dump version 10.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id integer NOT NULL,
    id_student bigint NOT NULL,
    name character varying,
    journal character varying,
    bibliorecord character varying,
    type character varying,
    filename character varying,
    realfilename character varying,
    confirmed boolean DEFAULT false
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- Name: articles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.articles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.articles_id_seq OWNER TO postgres;

--
-- Name: articles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.articles_id_seq OWNED BY public.articles.id;


--
-- Name: courseworks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.courseworks (
    id integer NOT NULL,
    id_student bigint NOT NULL,
    id_subject integer NOT NULL,
    semester smallint NOT NULL,
    theme character varying,
    head character varying,
    rating smallint,
    confirmed boolean DEFAULT false
);


ALTER TABLE public.courseworks OWNER TO postgres;

--
-- Name: courseworks_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.courseworks_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.courseworks_id_seq OWNER TO postgres;

--
-- Name: courseworks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.courseworks_id_seq OWNED BY public.courseworks.id;


--
-- Name: departments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.departments (
    id integer NOT NULL,
    id_faculty integer NOT NULL,
    name character varying NOT NULL,
    shortname character varying NOT NULL
);


ALTER TABLE public.departments OWNER TO postgres;

--
-- Name: departments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.departments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.departments_id_seq OWNER TO postgres;

--
-- Name: departments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.departments_id_seq OWNED BY public.departments.id;


--
-- Name: faculties; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.faculties (
    id integer NOT NULL,
    name character varying NOT NULL,
    shortname character varying NOT NULL
);


ALTER TABLE public.faculties OWNER TO postgres;

--
-- Name: faculties_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.faculties_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.faculties_id_seq OWNER TO postgres;

--
-- Name: faculties_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.faculties_id_seq OWNED BY public.faculties.id;


--
-- Name: fieldsofstudy; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.fieldsofstudy (
    id integer NOT NULL,
    id_department integer NOT NULL,
    code character varying NOT NULL,
    name character varying NOT NULL,
    alias character varying NOT NULL,
    level smallint NOT NULL
);


ALTER TABLE public.fieldsofstudy OWNER TO postgres;

--
-- Name: fieldsofstudy_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.fieldsofstudy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.fieldsofstudy_id_seq OWNER TO postgres;

--
-- Name: fieldsofstudy_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.fieldsofstudy_id_seq OWNED BY public.fieldsofstudy.id;


--
-- Name: marks; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.marks (
    id_student bigint NOT NULL,
    id_subject bigint NOT NULL,
    rating smallint NOT NULL,
    semester smallint NOT NULL,
    repass boolean DEFAULT false NOT NULL
);


ALTER TABLE public.marks OWNER TO postgres;

--
-- Name: practicis; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.practicis (
    id integer NOT NULL,
    id_student bigint NOT NULL,
    semester smallint NOT NULL,
    name character varying,
    head character varying,
    company character varying,
    begin_date date,
    end_date date,
    rating smallint
);


ALTER TABLE public.practicis OWNER TO postgres;

--
-- Name: practicis_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.practicis_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.practicis_id_seq OWNER TO postgres;

--
-- Name: practicis_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.practicis_id_seq OWNED BY public.practicis.id;


--
-- Name: students; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.students (
    id integer NOT NULL,
    id_field integer NOT NULL,
    number character varying NOT NULL,
    durationofstudy smallint NOT NULL,
    team character varying NOT NULL,
    teamnumber smallint NOT NULL
);


ALTER TABLE public.students OWNER TO postgres;

--
-- Name: students_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.students_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.students_id_seq OWNER TO postgres;

--
-- Name: students_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.students_id_seq OWNED BY public.students.id;


--
-- Name: subjects; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.subjects (
    id integer NOT NULL,
    id_department integer NOT NULL,
    id_field integer NOT NULL,
    name character varying NOT NULL,
    passtype smallint NOT NULL
);


ALTER TABLE public.subjects OWNER TO postgres;

--
-- Name: subjects_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.subjects_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.subjects_id_seq OWNER TO postgres;

--
-- Name: subjects_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.subjects_id_seq OWNED BY public.subjects.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    role smallint NOT NULL,
    login character varying NOT NULL,
    password character varying NOT NULL,
    fullname character varying,
    isactivated boolean DEFAULT false NOT NULL,
    id_faculty integer,
    id_department integer,
    id_student integer
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: articles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles ALTER COLUMN id SET DEFAULT nextval('public.articles_id_seq'::regclass);


--
-- Name: courseworks id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courseworks ALTER COLUMN id SET DEFAULT nextval('public.courseworks_id_seq'::regclass);


--
-- Name: departments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments ALTER COLUMN id SET DEFAULT nextval('public.departments_id_seq'::regclass);


--
-- Name: faculties id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.faculties ALTER COLUMN id SET DEFAULT nextval('public.faculties_id_seq'::regclass);


--
-- Name: fieldsofstudy id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fieldsofstudy ALTER COLUMN id SET DEFAULT nextval('public.fieldsofstudy_id_seq'::regclass);


--
-- Name: practicis id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.practicis ALTER COLUMN id SET DEFAULT nextval('public.practicis_id_seq'::regclass);


--
-- Name: students id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students ALTER COLUMN id SET DEFAULT nextval('public.students_id_seq'::regclass);


--
-- Name: subjects id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subjects ALTER COLUMN id SET DEFAULT nextval('public.subjects_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: articles articles_name_id_student_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_name_id_student_key UNIQUE (name, id_student);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);


--
-- Name: courseworks courseworks_id_subject_semester_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courseworks
    ADD CONSTRAINT courseworks_id_subject_semester_key UNIQUE (id_subject, semester);


--
-- Name: courseworks courseworks_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courseworks
    ADD CONSTRAINT courseworks_pkey PRIMARY KEY (id);


--
-- Name: departments departments_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_name_key UNIQUE (name);


--
-- Name: departments departments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_pkey PRIMARY KEY (id);


--
-- Name: departments departments_shortname_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_shortname_key UNIQUE (shortname);


--
-- Name: faculties faculties_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.faculties
    ADD CONSTRAINT faculties_name_key UNIQUE (name);


--
-- Name: faculties faculties_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.faculties
    ADD CONSTRAINT faculties_pkey PRIMARY KEY (id);


--
-- Name: faculties faculties_shortname_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.faculties
    ADD CONSTRAINT faculties_shortname_key UNIQUE (shortname);


--
-- Name: fieldsofstudy fieldsofstudy_id_department_code_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fieldsofstudy
    ADD CONSTRAINT fieldsofstudy_id_department_code_key UNIQUE (id_department, code);


--
-- Name: fieldsofstudy fieldsofstudy_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fieldsofstudy
    ADD CONSTRAINT fieldsofstudy_pkey PRIMARY KEY (id);


--
-- Name: practicis practicis_id_student_semester_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.practicis
    ADD CONSTRAINT practicis_id_student_semester_key UNIQUE (id_student, semester);


--
-- Name: practicis practicis_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.practicis
    ADD CONSTRAINT practicis_pkey PRIMARY KEY (id);


--
-- Name: students students_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_number_key UNIQUE (number);


--
-- Name: students students_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_pkey PRIMARY KEY (id);


--
-- Name: subjects subjects_id_field_name_passtype_id_department_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_id_field_name_passtype_id_department_key UNIQUE (id_field, name, passtype, id_department);


--
-- Name: subjects subjects_id_field_name_passtype_id_department_key1; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_id_field_name_passtype_id_department_key1 UNIQUE (id_field, name, passtype, id_department);


--
-- Name: subjects subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);


--
-- Name: users users_login_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_login_key UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: articles articles_id_student_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_id_student_fkey FOREIGN KEY (id_student) REFERENCES public.students(id);


--
-- Name: courseworks courseworks_id_student_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courseworks
    ADD CONSTRAINT courseworks_id_student_fkey FOREIGN KEY (id_student) REFERENCES public.students(id);


--
-- Name: courseworks courseworks_id_subject_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.courseworks
    ADD CONSTRAINT courseworks_id_subject_fkey FOREIGN KEY (id_subject) REFERENCES public.subjects(id);


--
-- Name: departments departments_id_faculty_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_id_faculty_fkey FOREIGN KEY (id_faculty) REFERENCES public.faculties(id) ON DELETE SET NULL;


--
-- Name: fieldsofstudy fieldsofstudy_id_department_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.fieldsofstudy
    ADD CONSTRAINT fieldsofstudy_id_department_fkey FOREIGN KEY (id_department) REFERENCES public.departments(id) ON DELETE SET NULL;


--
-- Name: marks marks_id_student_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.marks
    ADD CONSTRAINT marks_id_student_fkey FOREIGN KEY (id_student) REFERENCES public.students(id);


--
-- Name: marks marks_id_subject_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.marks
    ADD CONSTRAINT marks_id_subject_fkey FOREIGN KEY (id_subject) REFERENCES public.subjects(id);


--
-- Name: practicis practicis_id_student_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.practicis
    ADD CONSTRAINT practicis_id_student_fkey FOREIGN KEY (id_student) REFERENCES public.students(id);


--
-- Name: students students_id_field_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.students
    ADD CONSTRAINT students_id_field_fkey FOREIGN KEY (id_field) REFERENCES public.fieldsofstudy(id) ON DELETE SET NULL;


--
-- Name: subjects subjects_id_department_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_id_department_fkey FOREIGN KEY (id_department) REFERENCES public.departments(id) ON DELETE SET NULL;


--
-- Name: subjects subjects_id_field_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subjects
    ADD CONSTRAINT subjects_id_field_fkey FOREIGN KEY (id_field) REFERENCES public.fieldsofstudy(id) ON DELETE SET NULL;


--
-- Name: users users_id_department_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_id_department_fkey FOREIGN KEY (id_department) REFERENCES public.departments(id) ON DELETE SET NULL;


--
-- Name: users users_id_faculty_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_id_faculty_fkey FOREIGN KEY (id_faculty) REFERENCES public.faculties(id) ON DELETE SET NULL;


--
-- Name: users users_id_student_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_id_student_fkey FOREIGN KEY (id_student) REFERENCES public.students(id) ON DELETE SET NULL;


--
-- PostgreSQL database dump complete
--

