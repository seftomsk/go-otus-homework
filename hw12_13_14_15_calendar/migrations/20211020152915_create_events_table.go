package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upCreateEventsTable, downCreateEventsTable)
}

var eventsTable = `
CREATE TABLE events (
	id VARCHAR(128) PRIMARY KEY UNIQUE,
	user_id VARCHAR(128),
	title VARCHAR(255),
	description TEXT NULL,
	datetime TIMESTAMP NOT NULL UNIQUE,
	duration INT NOT NULL,
	remind_before INT NOT NULL
)
`

func upCreateEventsTable(tx *sql.Tx) error {
	if _, err := tx.Exec(eventsTable); err != nil {
		return fmt.Errorf("upCreateEventsTable: %w", err)
	}
	return nil
}

func downCreateEventsTable(tx *sql.Tx) error {
	if _, err := tx.Exec("DROP TABLE events"); err != nil {
		return fmt.Errorf("downCreateEventsTable: %w", err)
	}
	return nil
}
