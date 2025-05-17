-- +goose Up
CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TABLE person (
    id serial primary key,
    name varchar(50) not null,
    surname varchar(50) not null,
    patronymic varchar(50),
    age SMALLINT,
    gender gender,
    nationality varchar(2),
    UNIQUE (name, surname)
);

-- +goose Down
DROP TABLE person;

DROP TYPE gender;
