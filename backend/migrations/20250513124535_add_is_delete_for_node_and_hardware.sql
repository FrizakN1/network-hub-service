-- +goose Up
-- +goose StatementBegin
ALTER TABLE "Node"
ADD COLUMN "is_delete" boolean NOT NULL DEFAULT false;

ALTER TABLE "Hardware"
    ADD COLUMN "is_delete" boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "Hardware"
    DROP COLUMN "is_delete";

ALTER TABLE "Node"
    DROP COLUMN "is_delete";
-- +goose StatementEnd
