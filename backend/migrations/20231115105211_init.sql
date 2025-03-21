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

CREATE TABLE IF NOT EXISTS "Files" (
    id serial PRIMARY KEY,
    house_id integer NOT NULL,
    file_path character varying NOT NULL UNIQUE,
    file_name character varying NOT NULL,
    upload_at bigint NOT NULL,
    description character varying,
    in_archive boolean NOT NULL,
    FOREIGN KEY (house_id) REFERENCES "House"(id)
);

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

DROP TABLE IF EXISTS "House";
DROP TABLE IF EXISTS "Street";
DROP TABLE IF EXISTS "Street_type";
DROP TABLE IF EXISTS "House_type";

-- +goose StatementEnd
