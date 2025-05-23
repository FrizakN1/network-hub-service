-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "House_type" (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL,
    short_name character varying(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "Street_type" (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL,
    short_name character varying(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "Street" (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL,
    type_id integer NOT NULL,
    fias_id character varying(255),
    FOREIGN KEY (type_id) REFERENCES "Street_type"(id)
);

CREATE TABLE IF NOT EXISTS "House" (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL,
    type_id integer NOT NULL,
    fias_id character varying(255),
    street_id integer NOT NULL,
    FOREIGN KEY (street_id) REFERENCES "Street"(id),
    FOREIGN KEY (type_id) REFERENCES "House_type"(id)
);

CREATE TABLE IF NOT EXISTS "House_files" (
    id serial PRIMARY KEY,
    house_id integer NOT NULL,
    file_path character varying(255) NOT NULL UNIQUE,
    file_name character varying(255) NOT NULL,
    upload_at bigint NOT NULL,
    in_archive boolean NOT NULL,
    FOREIGN KEY (house_id) REFERENCES "House"(id)
);

CREATE TABLE IF NOT EXISTS "Node_owner" (
    id serial PRIMARY KEY,
    value character varying(255) NOT NULL UNIQUE,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Node_type" (
    id serial PRIMARY KEY,
    value character varying(255) NOT NULL UNIQUE,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Node" (
    id serial PRIMARY KEY,
    parent_id integer,
    house_id integer NOT NULL,
    type_id integer,
    owner_id integer NOT NULL,
    name character varying(255) NOT NULL,
    zone character varying(255),
    placement text,
    supply text,
    access text,
    description text,
    created_at bigint NOT NULL,
    updated_at bigint,
    is_delete boolean NOT NULL DEFAULT false,
    is_passive boolean NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES "Node"(id),
    FOREIGN KEY (house_id) REFERENCES "House"(id),
    FOREIGN KEY (type_id) REFERENCES "Node_type"(id),
    FOREIGN KEY (owner_id) REFERENCES "Node_owner"(id)
);

CREATE TABLE IF NOT EXISTS "Node_files" (
    id serial PRIMARY KEY,
    node_id integer NOT NULL,
    file_path character varying(255) NOT NULL UNIQUE,
    file_name character varying(255) NOT NULL,
    upload_at bigint NOT NULL,
    in_archive boolean NOT NULL,
    is_preview_image boolean NOT NULL,
    FOREIGN KEY (node_id) REFERENCES "Node"(id)
);

CREATE TABLE IF NOT EXISTS "Hardware_type" (
    id serial PRIMARY KEY,
    key character varying(255) NOT NULL UNIQUE,
    value character varying(255) NOT NULL,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Operation_mode" (
    id serial PRIMARY KEY,
    key character varying(255) NOT NULL UNIQUE,
    value character varying(255) NOT NULL,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Switch" (
    id serial PRIMARY KEY,
    name character varying(255) NOT NULL,
    operation_mode_id integer,
    community_read character varying(255),
    community_write character varying(255),
    port_amount integer NOT NULL,
    firmware_oid character varying(255),
    system_name_oid character varying(255),
    sn_oid character varying(255),
    save_config_oid character varying(255),
    port_desc_oid character varying(255),
    vlan_oid character varying(255),
    port_untagged_oid character varying(255),
    speed_oid character varying(255),
    battery_status_oid character varying(255),
    battery_charge_oid character varying(255),
    port_mode_oid character varying(255),
    uptime_oid character varying(255),
    created_at bigint NOT NULL,
    mac_oid character varying(255),
    FOREIGN KEY (operation_mode_id) REFERENCES "Operation_mode"(id)
);

CREATE TABLE IF NOT EXISTS "Hardware" (
    id serial PRIMARY KEY,
    node_id integer NOT NULL,
    type_id integer NOT NULL,
    switch_id integer,
    ip_address character varying(255),
    mgmt_vlan character varying(255),
    description character varying(255),
    created_at bigint NOT NULL,
    updated_at bigint,
    is_delete boolean NOT NULL DEFAULT false,
    FOREIGN KEY (node_id) REFERENCES "Node"(id),
    FOREIGN KEY (type_id) REFERENCES "Hardware_type"(id),
    FOREIGN KEY (switch_id) REFERENCES "Switch"(id)
);

CREATE TABLE IF NOT EXISTS "Hardware_files" (
    id serial PRIMARY KEY,
    hardware_id integer NOT NULL,
    file_path character varying(255) NOT NULL UNIQUE,
    file_name character varying(255) NOT NULL,
    upload_at bigint NOT NULL,
    in_archive boolean NOT NULL,
    FOREIGN KEY (hardware_id) REFERENCES "Hardware"(id)
);

CREATE TABLE IF NOT EXISTS "Event" (
    id bigserial PRIMARY KEY,
    house_id integer NOT NULL,
    node_id integer,
    hardware_id integer,
    user_id integer NOT NULL,
    description character varying(255) NOT NULL,
    created_at bigint NOT NULL,
    FOREIGN KEY (house_id) REFERENCES "House"(id),
    FOREIGN KEY (node_id) REFERENCES "Node"(id),
    FOREIGN KEY (hardware_id) REFERENCES "Hardware"(id)
);

CREATE TABLE IF NOT EXISTS "Roof_type" (
    id serial PRIMARY KEY,
    value character varying(255) NOT NULL UNIQUE,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Wiring_type" (
    id serial PRIMARY KEY,
    value character varying(255) NOT NULL UNIQUE,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "House_param" (
    id serial PRIMARY KEY,
    house_id integer NOT NULL UNIQUE,
    roof_type_id integer,
    wiring_type_id integer,
    FOREIGN KEY (house_id) REFERENCES "House"(id),
    FOREIGN KEY (roof_type_id) REFERENCES "Roof_type"(id),
    FOREIGN KEY (wiring_type_id) REFERENCES "Wiring_type"(id)
);

INSERT INTO "Operation_mode"(key, value, created_at)
VALUES
    ('eltex', 'как Eltex', floor(extract(epoch from now()))),
    ('dlink', 'как D-link', floor(extract(epoch from now())));

INSERT INTO "Hardware_type"(key, value, created_at)
VALUES
    ('switch', 'Коммутатор', floor(extract(epoch from now())));

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
DROP TABLE IF EXISTS "House_param";
DROP TABLE IF EXISTS "Wiring_type";
DROP TABLE IF EXISTS "Roof_type";
DROP TABLE IF EXISTS "Event";
DROP TABLE IF EXISTS "Hardware_files";
DROP TABLE IF EXISTS "Hardware";
DROP TABLE IF EXISTS "Switch";
DROP TABLE IF EXISTS "Operation_mode";
DROP TABLE IF EXISTS "Hardware_type";
DROP TABLE IF EXISTS "Node_files";
DROP TABLE IF EXISTS "Node";
DROP TABLE IF EXISTS "Node_owner";
DROP TABLE IF EXISTS "Node_type";
DROP TABLE IF EXISTS "Session";
DROP TABLE IF EXISTS "House_history";
DROP TABLE IF EXISTS "User";
DROP TABLE IF EXISTS "Role";
DROP TABLE IF EXISTS "House_files";
DROP TABLE IF EXISTS "House";
DROP TABLE IF EXISTS "Street";
DROP TABLE IF EXISTS "Street_type";
DROP TABLE IF EXISTS "House_type";

-- +goose StatementEnd
