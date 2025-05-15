-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Event" (
    id bigserial PRIMARY KEY,
    house_id integer NOT NULL,
    node_id integer,
    hardware_id integer,
    user_id integer NOT NULL,
    description character varying NOT NULL,
    created_at bigint NOT NULL,
    FOREIGN KEY (house_id) REFERENCES "House"(id),
    FOREIGN KEY (node_id) REFERENCES "Node"(id),
    FOREIGN KEY (hardware_id) REFERENCES "Hardware"(id)
);

ALTER TABLE "Switch"
    ADD COLUMN "mac_oid" character varying;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "Switch"
    DROP COLUMN "mac_oid";

DROP TABLE IF EXISTS "Event";
-- +goose StatementEnd
