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
    role_id integer NOT NULL REFERENCES roles,
    group_id INTEGER REFERENCES groups
);

CREATE TABLE users_credentials (
    user_credential_id serial PRIMARY KEY,
    login varchar(30) UNIQUE NOT NULL,
    password_hash char(64) NOT NULL,
    user_id integer UNIQUE NOT NULL REFERENCES users
);

CREATE INDEX users_credentials_login_idx ON users_credentials (login);

CREATE TABLE sessions (
    session_id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users,
    logged_in_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    logged_out_at timestamp
);

CREATE TABLE tokens (
    token_id serial PRIMARY KEY,
    refresh_token char(64) UNIQUE NOT NULL,
    expires_at bigint NOT NULL,
    session_id integer REFERENCES sessions
);

CREATE INDEX tokens_refresh_token_idx ON tokens (refresh_token);

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

CREATE TABLE lesson_types (
    lesson_type_id serial PRIMARY KEY,
    name varchar(30) NOT NULL
);

CREATE TABLE lessons (
    lesson_id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    plan_filepath text NOT NULL,
    compendium_filepath text NOT NULL,
    presentation_filepath text NOT NULL,
    chapter_id INTEGER NOT NULL REFERENCES chapters,
    type_id integer NOT NULL REFERENCES lesson_types
);

CREATE TABLE material_types (
    type_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE materials (
    material_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    filepath text NOT NULL,
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

CREATE TABLE publishers (
    publisher_id serial PRIMARY KEY,
    name varchar(50) NOT NULL
);

CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    filepath VARCHAR(150) NOT NULL,
    release_year integer NOT NULL,
    publisher_id integer NOT NULL REFERENCES publishers
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
