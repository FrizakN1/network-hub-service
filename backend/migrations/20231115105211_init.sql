-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "House_type" (
    id serial PRIMARY KEY,
    name character varying NOT NULL,
    short_name character varying NOT NULL
);

CREATE TABLE IF NOT EXISTS "Street_type" (
    id serial PRIMARY KEY,
    name character varying NOT NULL,
    short_name character varying NOT NULL
);

CREATE TABLE IF NOT EXISTS "Street" (
    id serial PRIMARY KEY,
    name character varying NOT NULL,
    type_id integer NOT NULL,
    fias_id character varying,
    FOREIGN KEY (type_id) REFERENCES "Street_type"(id)
);

CREATE TABLE IF NOT EXISTS "House" (
    id serial PRIMARY KEY,
    name character varying NOT NULL,
    type_id integer NOT NULL,
    fias_id character varying,
    street_id integer NOT NULL,
    FOREIGN KEY (street_id) REFERENCES "Street"(id),
    FOREIGN KEY (type_id) REFERENCES "House_type"(id)
);

CREATE TABLE IF NOT EXISTS "House_files" (
    id serial PRIMARY KEY,
    house_id integer NOT NULL,
    file_path character varying NOT NULL UNIQUE,
    file_name character varying NOT NULL,
    upload_at bigint NOT NULL,
    in_archive boolean NOT NULL,
    FOREIGN KEY (house_id) REFERENCES "House"(id)
);

CREATE TABLE IF NOT EXISTS "Role" (
    id serial PRIMARY KEY,
    value character varying NOT NULL UNIQUE,
    translate_value character varying NOT NULL
);

CREATE TABLE IF NOT EXISTS "User" (
    id serial PRIMARY KEY ,
    role_id integer NOT NULL,
    login character varying NOT NULL UNIQUE,
    name character varying NOT NULL,
    password character varying NOT NULL,
    baned boolean DEFAULT false NOT NULL,
    created_at bigint NOT NULL,
    updated_at bigint,
    FOREIGN KEY (role_id) REFERENCES "Role"(id)
);

CREATE TABLE IF NOT EXISTS "House_history"(
    id bigserial PRIMARY KEY,
    description character varying NOT NULL,
    house_id integer NOT NULL,
    user_id integer NOT NULL,
    FOREIGN KEY (house_id) REFERENCES "House"(id),
    FOREIGN KEY (user_id) REFERENCES "User"(id)
);

CREATE TABLE IF NOT EXISTS "Session" (
    hash character varying PRIMARY KEY,
    user_id integer NOT NULL,
    created_at bigint,
    FOREIGN KEY (user_id) REFERENCES "User"(id)
);

CREATE TABLE IF NOT EXISTS "Node_owner" (
    id serial PRIMARY KEY,
    name character varying NOT NULL UNIQUE,
    created_at bigint
);

CREATE TABLE IF NOT EXISTS "Node_type" (
    id serial PRIMARY KEY,
    name character varying NOT NULL UNIQUE,
    created_at bigint
);

CREATE TABLE IF NOT EXISTS "Node" (
    id serial PRIMARY KEY,
    parent_id integer,
    house_id integer NOT NULL,
    type_id integer NOT NULL,
    owner_id integer NOT NULL,
    name character varying NOT NULL,
    zone character varying,
    placement text,
    supply text,
    access text,
    description text,
    created_at bigint NOT NULL,
    updated_at bigint,
    FOREIGN KEY (parent_id) REFERENCES "Node"(id),
    FOREIGN KEY (house_id) REFERENCES "House"(id),
    FOREIGN KEY (type_id) REFERENCES "Node_type"(id),
    FOREIGN KEY (owner_id) REFERENCES "Node_owner"(id)
);

CREATE TABLE IF NOT EXISTS "Node_files" (
    id serial PRIMARY KEY,
    node_id integer NOT NULL,
    file_path character varying NOT NULL UNIQUE,
    file_name character varying NOT NULL,
    upload_at bigint NOT NULL,
    in_archive boolean NOT NULL,
    is_preview_image boolean NOT NULL,
    FOREIGN KEY (node_id) REFERENCES "Node"(id)
);

INSERT INTO "Role" (value, translate_value)
VALUES
    ('admin', 'Админ'),
    ('user', 'Пользователь');

INSERT INTO "Street_type" (name, short_name)
VALUES
    ('Километр', 'км'),
    ('Переулок', 'пер'),
    ('Бульвар', 'б-р'),
    ('Магистраль', 'мгстр.'),
    ('Площадь', 'пл.'),
    ('Сквер', 'сквер'),
    ('Парк', 'парк'),
    ('Улица', 'ул'),
    ('Территория', 'тер.'),
    ('Территория', 'тер'),
    ('Шоссе', 'ш'),
    ('Шоссе', 'ш.'),
    ('Набережная', 'наб'),
    ('Тракт', 'тракт'),
    ('Разъезд', 'рзд.'),
    ('Микрорайон', 'мкр.'),
    ('Площадь', 'пл'),
    ('Проезд', 'пр-д'),
    ('Проспект', 'пр-кт');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "Node_files";
DROP TABLE IF EXISTS "Node";
DROP TABLE IF EXISTS "Node_owner";
DROP TABLE IF EXISTS "Node_type";
DROP TABLE IF EXISTS "Session";
DROP TABLE IF EXISTS "House_history";
DROP TABLE IF EXISTS "User";
DROP TABLE IF EXISTS "Role";
DROP TABLE IF EXISTS "House";
DROP TABLE IF EXISTS "Street";
DROP TABLE IF EXISTS "Street_type";
DROP TABLE IF EXISTS "House_type";

-- +goose StatementEnd
