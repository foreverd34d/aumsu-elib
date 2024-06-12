CREATE TABLE departments (
    department_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE specialties (
    specialty_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    department_id INTEGER NOT NULL REFERENCES departments
);

CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    name CHAR(4) NOT NULL,
    specialty_id INTEGER NOT NULL REFERENCES specialties
);

CREATE TABLE roles (
    role_id serial PRIMARY KEY,
    name varchar(20) NOT NULL
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    surname VARCHAR(30) NOT NULL,
    patronymic VARCHAR(30),
    login VARCHAR(30) NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    role_id integer NOT NULL REFERENCES roles,
    group_id INTEGER REFERENCES groups
);

CREATE TABLE sessions (
    session_id serial PRIMARY KEY,
    refresh_token char(64) NOT NULL,
    expires_at bigint NOT NULL,
    user_id integer REFERENCES users
);

CREATE TABLE disciplines (
    discipine_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    specialty_id INTEGER NOT NULL REFERENCES specialties
);

CREATE TABLE chapters (
    chapter_id SERIAL PRIMARY KEY,
    discipline_id INTEGER NOT NULL REFERENCES disciplines,
    module INTEGER NOT NULL CHECK(module > 0)
);

CREATE TABLE lessons (
    lesson_id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    chapter_id INTEGER NOT NULL REFERENCES chapters
);

CREATE TABLE material_types (
    type_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE materials (
    material_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    filepath VARCHAR(150) NOT NULL,
    type_id INTEGER NOT NULL REFERENCES material_types
);

CREATE TABLE lesson_materials (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL REFERENCES materials,
    lesson_id INTEGER NOT NULL REFERENCES lessons
);

CREATE TABLE authors (
    author_id serial PRIMARY KEY,
    surname VARCHAR(30) NOT NULL,
    name VARCHAR(30) NOT NULL,
    patronymic VARCHAR(30)
);

CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    author_id varchar(50) NOT NULL,
    filepath VARCHAR(150) NOT NULL
);

CREATE TABLE author_books (
    author_book_id serial PRIMARY KEY,
    author_id integer NOT NULL REFERENCES authors,
    book_id integer NOT NULL REFERENCES books
);

CREATE TABLE material_books (
    id SERIAL PRIMARY KEY,
    material_id INTEGER NOT NULL REFERENCES materials,
    book_id INTEGER NOT NULL REFERENCES books
);
