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

	d.query["EDIT_REPORT_DATA"], err = d.db.Prepare(`
		UPDATE "Report_data" SET value = $2, description = $3 WHERE key = $1
	`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_REPORT_DATA"], err = d.db.Prepare(`
		SELECT * FROM "Report_data" ORDER BY id
	`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_ADDRESS_PARAMS"], err = d.db.Prepare(`
		SELECT hp.roof_type_id, hp.wiring_type_id, rt.value, wt.value
		FROM "House_param" AS hp
		LEFT JOIN "Roof_type" AS rt ON hp.roof_type_id = rt.id
		LEFT JOIN "Wiring_type" AS wt ON hp.wiring_type_id = wt.id
		WHERE house_id = $1
	`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_ADDRESSES_AMOUNTS"], err = d.db.Prepare(`
		WITH houses AS (
		    SELECT house_id
		    FROM (
		        SELECT house_id FROM "Node"
		        UNION ALL 
		        SELECT house_id FROM "House_files"
		    ) AS houses
		    GROUP BY house_id
		    ORDER BY house_id
		    OFFSET $1
		    LIMIT 20
		)
		SELECT 
            h.house_id,
            COUNT(DISTINCT f.id) AS files_count,
            COUNT(DISTINCT n.id) AS nodes_count,
            COUNT(DISTINCT hd.id) AS hardware_count
        FROM houses AS h 
        LEFT JOIN "House_files" AS f ON h.house_id = f.house_id
        LEFT JOIN "Node" AS n ON n.house_id = h.house_id AND n.is_delete = false
        LEFT JOIN "Hardware" AS hd ON hd.node_id = n.id AND hd.is_delete = false
        GROUP BY h.house_id
		ORDER BY h.house_id
	`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_ADDRESSES_AMOUNTS_BY_HOUSE_IDS"], err = d.db.Prepare(`
		SELECT 
            h.id AS house_id,
            COUNT(DISTINCT f.id) AS files_count,
            COUNT(DISTINCT n.id) AS nodes_count,
            COUNT(DISTINCT hd.id) AS hardware_count
        FROM (SELECT unnest($1::integer[]) AS id) AS h
        LEFT JOIN "House_files" AS f ON h.id = f.house_id 
        LEFT JOIN "Node" AS n ON n.house_id = h.id AND n.is_delete = false
        LEFT JOIN "Hardware" AS hd ON hd.node_id = n.id AND hd.is_delete = false
        GROUP BY h.id
		ORDER BY h.id
	`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["SET_HOUSE_PARAMS"], err = d.db.Prepare(`
		INSERT INTO "House_param"(house_id, roof_type_id, wiring_type_id) 
		VALUES ($1, $2, $3)
		ON CONFLICT(house_id) DO UPDATE SET roof_type_id = $2, wiring_type_id = $3
	`)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS"], err = d.db.Prepare(`
		SELECT e.*, n.name, hwt.value, COUNT(*) OVER()
		FROM "Event" AS e
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		ORDER BY e.created_at DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HOUSE_ALL"], err = d.db.Prepare(`
		SELECT e.*, n.name, hwt.value, COUNT(*) OVER()
		FROM "Event" AS e
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.house_id = $2
		ORDER BY e.created_at DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HOUSE_ONLY"], err = d.db.Prepare(`
		SELECT e.*, n.name, hwt.value, COUNT(*) OVER()
		FROM "Event" AS e
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.house_id = $2 AND e.node_id IS NULL AND e.hardware_id IS NULL
		ORDER BY e.created_at DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_NODE_ALL"], err = d.db.Prepare(`
		SELECT e.*, n.name, hwt.value, COUNT(*) OVER()
		FROM "Event" AS e
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.node_id = $2
		ORDER BY e.created_at DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_NODE_ONLY"], err = d.db.Prepare(`
		SELECT e.*, n.name, hwt.value, COUNT(*) OVER()
		FROM "Event" AS e
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.node_id = $2 AND e.hardware_id IS NULL
		ORDER BY e.created_at DESC
		OFFSET $1
		LIMIT 20;
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_EVENTS_HARDWARE"], err = d.db.Prepare(`
		SELECT e.*, n.name, hwt.value, COUNT(*) OVER()
		FROM "Event" AS e
		LEFT JOIN "Node" AS n ON e.node_id = n.id
		LEFT JOIN "Hardware" AS hw ON e.hardware_id = hw.id
		LEFT JOIN "Hardware_type" AS hwt ON hw.type_id = hwt.id
		WHERE e.hardware_id = $2
		ORDER BY e.created_at DESC
		OFFSET $1
		LIMIT 20
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

	d.query["GET_HARDWARE_FOR_INDEX"], err = d.db.Prepare(`
		SELECT hd.id, hdt.value, n.name, sw.name, hd.ip_address, n.house_id, hd.is_delete
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE"], err = d.db.Prepare(`
		SELECT hd.id, hd.node_id, hd.type_id, hd.switch_id, hd.ip_address, n.house_id, hdt.key, hdt.value, sw.name, n.name, COUNT(*) OVER()
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.is_delete = false
			AND ($2 = 0 OR n.house_id = $2)
			AND ($3 = 0 OR hd.node_id = $3)
		ORDER BY hd.id DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_HARDWARE_BY_IDS"], err = d.db.Prepare(`
		SELECT hd.id, hd.node_id, hd.type_id, hd.switch_id, hd.ip_address, n.house_id, hdt.key, hdt.value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
		JOIN "Hardware_type" AS hdt ON hd.type_id = hdt.id
		LEFT JOIN "Switch" AS sw ON hd.switch_id = sw.id
		WHERE hd.id = ANY($1)
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
		SELECT hd.*, n.house_id, hdt.key, hdt.value, sw.name, n.name
		FROM "Hardware" AS hd
		JOIN "Node" AS n ON hd.node_id = n.id
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

	d.query["GET_NODES_FOR_INDEX"], err = d.db.Prepare(`
		SELECT n.id, n.name, n.zone, no.value, n.house_id, n.is_delete, n.is_passive
		FROM "Node" AS n 
		JOIN "Node_owner" AS no ON n.owner_id = no.id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODES"], err = d.db.Prepare(`
		SELECT n.id, n.house_id, n.owner_id, n.name, n.zone, n.is_passive, n.placement, n.supply, no.value, nt.key, COUNT(*) OVER ()
		FROM "Node" AS n 
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		LEFT JOIN "Node_type" AS nt ON n.type_id = nt.id
		WHERE n.is_delete = false
			AND ($2 = false OR n.is_passive = false)
			AND ($3 = 0 OR house_id = $3)
		ORDER BY n.id DESC
		OFFSET $1
		LIMIT 20
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE"], err = d.db.Prepare(`
		SELECT n.*, nt.value, no.value, p.name
		FROM "Node" AS n 
		LEFT JOIN "Node_type" AS nt ON n.type_id = nt.id
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		LEFT JOIN "Node" AS p ON n.parent_id = p.id
		WHERE n.id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_NODE"], err = d.db.Prepare(`
		INSERT INTO "Node"(parent_id, house_id, type_id, owner_id, name, zone, placement, supply, access, description, created_at, updated_at, is_passive) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_NODE"], err = d.db.Prepare(`
		UPDATE "Node" SET parent_id = $2, type_id = $3, owner_id = $4, name = $5, zone = $6, placement = $7, supply = $8,
		                  access = $9, description = $10, updated_at = $11, house_id = $12, is_passive = $13
		WHERE id = $1	
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODES_BY_IDS"], err = d.db.Prepare(`
		SELECT n.id, n.house_id, n.owner_id, n.name, n.zone, n.is_passive, no.value
		FROM "Node" AS n 
		JOIN "Node_owner" AS no ON n.owner_id = no.id
		WHERE n.id = ANY($1)
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

	d.query["GET_ROOF_TYPES"], err = d.db.Prepare(`
		SELECT * FROM "Roof_type"
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_WIRING_TYPES"], err = d.db.Prepare(`
		SELECT * FROM "Wiring_type"
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["GET_NODE_TYPES"], err = d.db.Prepare(`
		SELECT id, value, created_at FROM "Node_type"
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
		INSERT INTO "Node_owner"(value, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_OWNERS"], err = d.db.Prepare(`
		UPDATE "Node_owner" SET value = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_ROOF_TYPES"], err = d.db.Prepare(`
		INSERT INTO "Roof_type"(value, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_ROOF_TYPES"], err = d.db.Prepare(`
		UPDATE "Roof_type" SET value = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_WIRING_TYPES"], err = d.db.Prepare(`
		INSERT INTO "Wiring_type"(value, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_WIRING_TYPES"], err = d.db.Prepare(`
		UPDATE "Wiring_type" SET value = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_NODE_TYPES"], err = d.db.Prepare(`
		INSERT INTO "Node_type"(value, created_at) VALUES ($1, $2)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_NODE_TYPES"], err = d.db.Prepare(`
		UPDATE "Node_type" SET value = $2 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_HARDWARE_TYPES"], err = d.db.Prepare(`
		INSERT INTO "Hardware_type"(key, value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_HARDWARE_TYPES"], err = d.db.Prepare(`
		UPDATE "Hardware_type" SET key = $2, value = $3 WHERE id = $1
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["CREATE_OPERATION_MODES"], err = d.db.Prepare(`
		INSERT INTO "Operation_mode"(key, value, created_at) VALUES ($1, $2, $3)
		RETURNING id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	d.query["EDIT_OPERATION_MODES"], err = d.db.Prepare(`
		UPDATE "Operation_mode" SET key = $2, value = $3 WHERE id = $1
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
		SELECT s.*, om.key, om.value 
		FROM "Switch" AS s
		LEFT JOIN "Operation_mode" AS om ON s.operation_mode_id = om.id
		ORDER BY s.id
    `)
	if err != nil {
		errorsList = append(errorsList, err)
	}

	return errorsList
}
