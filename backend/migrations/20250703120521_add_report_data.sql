-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "Report_data" (
    id serial PRIMARY KEY,
    key character varying(255) NOT NULL UNIQUE,
    value character varying(255) NOT NULL,
    description character varying(255)
);

INSERT INTO "Report_data" (key, value, description)
VALUES
    ('HN_SWITCH', 'Коммутатор D-Link DGS-1100-10/ME', 'Коммутатор для домового узла'),
    ('HN_SWITCH_POWER', '0.004', 'Мощность в кВт коммутатора домового узла'),
    ('BN_SWITCH', 'Коммутатор Eltex MES2324 FB', 'Коммутатор для районного и магистрального узлов'),
    ('BN_SWITCH_POWER', '0.045', 'Мощность в кВт коммутатора районного и магистрального узлов'),
    ('OPTICAL_RECEIVER', 'Оптич. приемник TVBS OR-826H', 'Оптический приемник'),
    ('OPTICAL_RECEIVER_POWER', '0.006', 'Мощность в кВт оптического приемника'),
    ('VOLTAGE_LEVEL', '0.22 кВ', 'Уровень напряжения'),
    ('CATEGORY_RELIABILITY_POWER_SUPPLY', 'III', 'Категория надежности электроснабжения');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "Report_data";
-- +goose StatementEnd
