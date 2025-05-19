package database

import (
	"database/sql"
	"fmt"
	"os"
)

type Database interface {
	Connect() error
	PrepareQuery() []error
	GetQuery(key string) (*sql.Stmt, bool)
}

type DefaultDatabase struct {
	db    *sql.DB
	query map[string]*sql.Stmt
}

func (d *DefaultDatabase) Connect() error {
	var err error

	d.db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME")))
	if err != nil {
		return err
	}

	return d.db.Ping()
}

func (d *DefaultDatabase) GetQuery(key string) (*sql.Stmt, bool) {
	stmt, ok := d.query[key]

	return stmt, ok
}

func (d *DefaultDatabase) PrepareQuery() []error {
	var err error
	errorsList := make([]error, 0)
	d.query = make(map[string]*sql.Stmt)

	d.query["STREET_TYPE"], err = d.db.Prepare(`SELECT * FROM "Street_type" ORDER BY id`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["HOUSE_TYPE"], err = d.db.Prepare(`SELECT * FROM "House_type" ORDER BY id`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SUGGESTIONS"], err = d.db.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.id, h.name, h.type_id, ht.short_name, 
		       (SELECT COUNT(*) FROM "House_files" AS f
			                        WHERE f.house_id = h.id ),
				(SELECT COUNT(*) FROM "Node" AS n
										WHERE n.house_id = h.id ),
		    (SELECT COUNT(*) FROM "Hardware" AS hd
		                     JOIN "Node" AS n ON hd.node_id = n.id
										WHERE n.house_id = h.id )
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE s.name ILIKE '%' || $1 || '%'
          AND (h.name ILIKE '%' || $2 || '%' OR $2 = '')
        ORDER BY 
            CASE 
				WHEN h.name = $2 THEN 0               -- точное полное совпадение (например, "3" == "3")
				WHEN h.name ~ ('^' || $2 || '[^0-9]') THEN 1  -- точное числовое совпадение с префиксом (например, "3а" при "3")
				WHEN h.name ILIKE $2 || '%' THEN 2     -- начинается с номера (например, "3" соответствует "3а")
				ELSE 3                                 -- частичное совпадение
    		END,
			LENGTH(h.name),                           -- сортировка по длине
			h.name
        OFFSET $3
		LIMIT $4
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SUGGESTIONS_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(h.name)
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        WHERE s.name ILIKE '%' || $1 || '%'
          AND (h.name ILIKE '%' || $2 || '%' OR $2 = '')
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSE"], err = d.db.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.name, h.type_id, ht.short_name
        FROM "House" AS h
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE h.id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSES"], err = d.db.Prepare(`
		SELECT s.name, s.type_id, st.short_name, h.id, h.name, h.type_id, ht.short_name, 
		        (SELECT COUNT(*) FROM "House_files" AS f
			                        	WHERE f.house_id = h.id ),
				(SELECT COUNT(*) FROM "Node" AS n
										WHERE n.house_id = h.id ),
				(SELECT COUNT(*) FROM "Hardware" AS hd
								 JOIN "Node" AS n ON hd.node_id = n.id
								 WHERE n.house_id = h.id )
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
        WHERE EXISTS (
			SELECT 1 
			FROM "House_files" AS f
			WHERE f.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Node" AS n
			WHERE n.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Hardware" AS hd
			JOIN "Node" AS n ON hd.node_id = n.id
			WHERE n.house_id = h.id
		)
        OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSES_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*)
        FROM "Street" AS s
        JOIN "House" AS h ON s.id = h.street_id
        WHERE EXISTS (
			SELECT 1 
			FROM "House_files" AS f
			WHERE f.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Node" AS n
			WHERE n.house_id = h.id
		) OR EXISTS (
			SELECT 1 
			FROM "Hardware" AS hd
			JOIN "Node" AS n ON hd.node_id = n.id
			WHERE n.house_id = h.id
		)
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS"], err = d.db.Prepare(`
		SELECT e.*, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		ORDER BY e.created_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Event"
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HOUSE_ALL"], err = d.db.Prepare(`
		SELECT e.*, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.house_id = $1
		ORDER BY e.created_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HOUSE_ALL_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Event" WHERE house_id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HOUSE_ONLY"], err = d.db.Prepare(`
		SELECT e.*, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.house_id = $1 AND e.node_id IS NULL AND e.hardware_id IS NULL
		ORDER BY e.created_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HOUSE_ONLY_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Event"
		WHERE house_id = $1 AND node_id IS NULL AND hardware_id IS NULL
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_NODE_ALL"], err = d.db.Prepare(`
		SELECT e.*, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.node_id = $1
		ORDER BY e.created_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_NODE_ALL_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*)
		FROM "Event"
		WHERE node_id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_NODE_ONLY"], err = d.db.Prepare(`
		SELECT e.*, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.node_id = $1 AND e.hardware_id IS NULL
		ORDER BY e.created_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_NODE_ONLY_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Event"
		WHERE node_id = $1 AND hardware_id IS NULL
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HARDWARE"], err = d.db.Prepare(`
		SELECT e.*, s.name, st.short_name, h.name, ht.short_name, n.name, hwt.translate_value
		FROM "Event" AS e
		JOIN "House" AS h ON e.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.hardware_id = $1
		ORDER BY e.created_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HARDWARE_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Event" WHERE hardware_id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_EVENT"], err = d.db.Prepare(`
		INSERT INTO "Event"(house_id, node_id, hardware_id, user_id, description, created_at) VALUES ($1, $2, $3, $4, $5, $6)
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_FILE_HOUSES"], err = d.db.Prepare(`
		INSERT INTO "House_files"(house_id, file_path, file_name, upload_at, in_archive) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_FILE_HARDWARE"], err = d.db.Prepare(`
		INSERT INTO "Hardware_files"(hardware_id, file_path, file_name, upload_at, in_archive) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_FILE_NODES"], err = d.db.Prepare(`
		INSERT INTO "Node_files"(node_id, file_path, file_name, upload_at, in_archive, is_preview_image) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSE_FILES"], err = d.db.Prepare(`
		SELECT * FROM "House_files" WHERE house_id = $1
		ORDER BY upload_at DESC 
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["ARCHIVE_FILE_HOUSES"], err = d.db.Prepare(`
		UPDATE "House_files" SET in_archive = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["ARCHIVE_FILE_NODES"], err = d.db.Prepare(`
		UPDATE "Node_files" SET in_archive = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["ARCHIVE_FILE_HARDWARE"], err = d.db.Prepare(`
		UPDATE "Hardware_files" SET in_archive = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["DELETE_FILE_HOUSES"], err = d.db.Prepare(`
		DELETE FROM "House_files" WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["DELETE_FILE_NODES"], err = d.db.Prepare(`
		DELETE FROM "Node_files" WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE_FILES"], err = d.db.Prepare(`
		SELECT * FROM "Node_files" WHERE node_id = $1 AND is_preview_image = $2
		ORDER BY upload_at DESC 
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE_FILES"], err = d.db.Prepare(`
		SELECT * FROM "Hardware_files" WHERE hardware_id = $1
		ORDER BY upload_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["DELETE_FILE_HARDWARE"], err = d.db.Prepare(`
		DELETE FROM "Hardware_files" WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE"], err = d.db.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.is_delete = false
		ORDER BY hd.id DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Hardware"
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSE_HARDWARE"], err = d.db.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE n.house_id = $1 AND hd.is_delete = false
		ORDER BY hd.id DESC
		OFFSET $2
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSE_HARDWARE_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(hd.*) 
		FROM "Hardware" AS hd 
		JOIN "Node" AS n ON hd.node_id = n.id
		WHERE n.house_id = $1 AND hd.is_delete = false
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE_HARDWARE"], err = d.db.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.node_id = $1 AND hd.is_delete = false
		ORDER BY hd.id DESC
		OFFSET $2
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE_HARDWARE_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) 
		FROM "Hardware"
		WHERE node_id = $1 AND is_delete = false
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SEARCH_HARDWARE"], err = d.db.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE (n.name ILIKE '%' || $1 || '%'
			OR hdt.translate_value ILIKE '%' || $1 || '%'
			OR sw.name ILIKE '%' || $1 || '%'
			OR hd.ip_address ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%')
			AND hd.is_delete = false
		ORDER BY hd.id DESC
		OFFSET $2
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SEARCH_HARDWARE_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(hd.*)
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE (n.name ILIKE '%' || $1 || '%'
			OR hdt.translate_value ILIKE '%' || $1 || '%'
			OR sw.name ILIKE '%' || $1 || '%'
			OR hd.ip_address ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%')
			AND hd.is_delete = false
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_HARDWARE"], err = d.db.Prepare(`
		INSERT INTO "Hardware" (node_id, type_id, switch_id, ip_address, mgmt_vlan, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_HARDWARE"], err = d.db.Prepare(`
		UPDATE "Hardware" SET node_id = $2, type_id = $3, switch_id = $4, ip_address = $5, mgmt_vlan = $6, description = $7,
		                      updated_at = $8
		WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE_BY_ID"], err = d.db.Prepare(`
		SELECT hd.*, s.name, st.short_name, h.id, h.name, ht.short_name, hdt.value, hdt.translate_value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE_FILES"], err = d.db.Prepare(`
		SELECT * FROM "Hardware_files" WHERE hardware_id = $1
		ORDER BY upload_at DESC
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODES"], err = d.db.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.is_delete = false
		ORDER BY id DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODES_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*) FROM "Node" WHERE is_delete = false
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSE_NODES"], err = d.db.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.house_id = $1 AND n.is_delete = false
		ORDER BY id DESC
		OFFSET $2
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HOUSE_NODES_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(*)
		FROM "Node"
		WHERE house_id = $1 AND is_delete = false
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE"], err = d.db.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_NODE"], err = d.db.Prepare(`
		INSERT INTO "Node"(parent_id, house_id, type_id, owner_id, name, zone, placement, supply, access, description, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_NODE"], err = d.db.Prepare(`
		UPDATE "Node" SET parent_id = $2, type_id = $3, owner_id = $4, name = $5, zone = $6, placement = $7, supply = $8,
		                  access = $9, description = $10, updated_at = $11, house_id = $12
		WHERE id = $1	
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SEARCH_NODES"], err = d.db.Prepare(`
		SELECT n.*, s.name, st.short_name, h.name, ht.short_name, nt.name, no.name, (SELECT p.name FROM "Node" AS p WHERE p.id = n.parent_id) AS parent_name
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE (n.name ILIKE '%' || $1 || '%'
			OR nt.name ILIKE '%' || $1 || '%'
			OR no.name ILIKE '%' || $1 || '%'
			OR n.zone ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%')
			AND n.is_delete = false
		ORDER BY id
		OFFSET $2
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SEARCH_NODES_COUNT"], err = d.db.Prepare(`
		SELECT COUNT(n.*)
		FROM "Node" AS n 
		JOIN "House" AS h ON n.house_id = h.id
        JOIN "Street" AS s ON s.id = h.street_id
        JOIN "Street_type" AS st ON s.type_id = st.id
        JOIN "House_type" AS ht ON h.type_id = ht.id
		JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE (n.name ILIKE '%' || $1 || '%'
			OR nt.name ILIKE '%' || $1 || '%'
			OR no.name ILIKE '%' || $1 || '%'
			OR n.zone ILIKE '%' || $1 || '%'
			OR s.name ILIKE '%' || $1 || '%'
			OR st.short_name ILIKE '%' || $1 || '%'
			OR h.name ILIKE '%' || $1 || '%'
			OR ht.short_name ILIKE '%' || $1 || '%'
			OR (st.short_name || ' ' || s.name || ', ' || ht.short_name || ' ' || h.name) ILIKE '%' || $1 || '%')
			AND n.is_delete = false 
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["DELETE_NODE"], err = d.db.Prepare(`
		UPDATE "Node" SET is_delete = true WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_OWNERS"], err = d.db.Prepare(`
		SELECT * FROM "Node_owner"
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE_TYPES"], err = d.db.Prepare(`
		SELECT * FROM "Node_type"
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE_TYPES"], err = d.db.Prepare(`
		SELECT * FROM "Hardware_type" ORDER BY id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_OPERATION_MODES"], err = d.db.Prepare(`
		SELECT * FROM "Operation_mode" ORDER BY id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_OWNERS"], err = d.db.Prepare(`
		INSERT INTO "Node_owner"(name, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_OWNERS"], err = d.db.Prepare(`
		UPDATE "Node_owner" SET name = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_NODE_TYPES"], err = d.db.Prepare(`
		INSERT INTO "Node_type"(name, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_NODE_TYPES"], err = d.db.Prepare(`
		UPDATE "Node_type" SET name = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_HARDWARE_TYPES"], err = d.db.Prepare(`
		INSERT INTO "Hardware_type"(value, translate_value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_HARDWARE_TYPES"], err = d.db.Prepare(`
		UPDATE "Hardware_type" SET value = $2, translate_value = $3 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_OPERATION_MODES"], err = d.db.Prepare(`
		INSERT INTO "Operation_mode"(value, translate_value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_OPERATION_MODES"], err = d.db.Prepare(`
		UPDATE "Operation_mode" SET value = $2, translate_value = $3 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_SWITCH"], err = d.db.Prepare(`
		INSERT INTO "Switch"(name, operation_mode_id, community_read, community_write, port_amount, firmware_oid, 
		                     system_name_oid, sn_oid, save_config_oid, port_desc_oid, vlan_oid, port_untagged_oid, 
		                     speed_oid, battery_status_oid, battery_charge_oid, port_mode_oid, uptime_oid, created_at, mac_oid) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_SWITCH"], err = d.db.Prepare(`
		UPDATE "Switch" SET name = $2, operation_mode_id = $3, community_read = $4, community_write = $5, port_amount = $6,
		                    firmware_oid = $7, system_name_oid = $8, sn_oid = $9, save_config_oid = $10, port_desc_oid = $11,
		                    vlan_oid = $12, port_untagged_oid = $13, speed_oid = $14, battery_status_oid = $15, battery_charge_oid = $16,
		                    port_mode_oid = $17, uptime_oid = $18, mac_oid = $19
		WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_SWITCHES"], err = d.db.Prepare(`
		SELECT s.*, om.value, om.translate_value 
		FROM "Switch" AS s
		LEFT JOIN "Operation_mode" AS om ON s.operation_mode_id = om.id
		ORDER BY s.id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	return errorsList
}
