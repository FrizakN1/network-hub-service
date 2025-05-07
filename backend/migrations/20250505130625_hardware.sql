-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Hardware_type" (
    id serial PRIMARY KEY,
    value character varying NOT NULL UNIQUE,
    translate_value character varying NOT NULL,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Operation_mode" (
    id serial PRIMARY KEY,
    value character varying NOT NULL UNIQUE,
    translate_value character varying NOT NULL,
    created_at bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS "Switch" (
    id serial PRIMARY KEY,
    name character varying NOT NULL,
    operation_mode_id integer,
    community_read character varying,
    community_write character varying,
    port_amount integer NOT NULL,
    firmware_oid character varying,
    system_name_oid character varying,
    sn_oid character varying,
    save_config_oid character varying,
    port_desc_oid character varying,
    vlan_oid character varying,
    port_untagged_oid character varying,
    speed_oid character varying,
    battery_status_oid character varying,
    battery_charge_oid character varying,
    port_mode_oid character varying,
    uptime_oid character varying,
    created_at bigint NOT NULL,
    FOREIGN KEY (operation_mode_id) REFERENCES "Operation_mode"(id)
);

CREATE TABLE IF NOT EXISTS "Hardware" (
    id serial PRIMARY KEY,
    node_id integer NOT NULL,
    type_id integer NOT NULL,
    switch_id integer,
    ip_address character varying,
    mgmt_vlan character varying,
    description character varying,
    created_at bigint NOT NULL,
    updated_at bigint,
    FOREIGN KEY (node_id) REFERENCES "Node"(id),
    FOREIGN KEY (type_id) REFERENCES "Hardware_type"(id),
    FOREIGN KEY (switch_id) REFERENCES "Switch"(id)
);

CREATE TABLE IF NOT EXISTS "Hardware_files" (
    id serial PRIMARY KEY,
    hardware_id integer NOT NULL,
    file_path character varying NOT NULL UNIQUE,
    file_name character varying NOT NULL,
    upload_at bigint NOT NULL,
    in_archive boolean NOT NULL,
    FOREIGN KEY (hardware_id) REFERENCES "Hardware"(id)
);

INSERT INTO "Operation_mode"(value, translate_value, created_at)
VALUES
    ('eltex', 'как Eltex', floor(extract(epoch from now()))),
    ('dlink', 'как D-link', floor(extract(epoch from now())));

INSERT INTO "Hardware_type"(value, translate_value, created_at)
VALUES
    ('switch', 'Коммутатор', floor(extract(epoch from now())));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "Hardware_files";
DROP TABLE IF EXISTS "Hardware";
DROP TABLE IF EXISTS "Switch";
DROP TABLE IF EXISTS "Operation_mode";
DROP TABLE IF EXISTS "Hardware_type";
-- +goose StatementEnd
