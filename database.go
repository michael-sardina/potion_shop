package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func setSelectedColor(db *sql.DB, color string) error {
	_, err := db.Exec(`UPDATE settings SET selected_color=? WHERE id=1`, color)
	return err
}

func purchaseTxn(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()

	var gold int
	var selected string
	if err := tx.QueryRow(`SELECT gold, selected_color FROM settings WHERE id=1`).Scan(&gold, &selected); err != nil {
		return err
	}
	var price int
	if err := tx.QueryRow(`SELECT price FROM prices WHERE color=?`, selected).Scan(&price); err != nil {
		return err
	}
	if gold < price {
		return nil
	}

	if _, err := tx.Exec(`UPDATE settings SET gold = gold - ? WHERE id=1`, price); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE inventory SET qty = qty + 1 WHERE color=?`, selected); err != nil {
		return err
	}

	return tx.Commit()
}

func loadStore(db *sql.DB) (Store, error) {
	var s Store
	s.Inventory = map[string]int{}
	s.Prices = map[string]int{}

	row := db.QueryRow(`SELECT gold, selected_color FROM settings WHERE id=1`)
	if err := row.Scan(&s.Gold, &s.SelectedColor); err != nil {
		return s, err
	}

	rows, err := db.Query(`SELECT color, qty FROM inventory`)
	if err != nil {
		return s, err
	}
	defer rows.Close()
	for rows.Next() {
		var c string
		var q int
		if err := rows.Scan(&c, &q); err != nil {
			return s, err
		}
		s.Inventory[c] = q
	}

	rows2, err := db.Query(`SELECT color, price FROM prices`)
	if err != nil {
		return s, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var c string
		var p int
		if err := rows2.Scan(&c, &p); err != nil {
			return s, err
		}
		s.Prices[c] = p
	}
	return s, rows2.Err()
}

func migrateAndSeed(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY CHECK (id=1),
			gold INTEGER NOT NULL,
			selected_color TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS inventory (
			color TEXT PRIMARY KEY,
			qty INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS prices (
			color TEXT PRIMARY KEY,
			price INTEGER NOT NULL
		);`,
	}

	for _, s := range stmts {
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}

	_, _ = db.Exec(`INSERT OR IGNORE INTO settings(id, gold, selected_color) VALUES(1, 200, 'red')`)
	for _, c := range []string{"red", "blue", "green"} {
		_, _ = db.Exec(`INSERT OR IGNORE INTO inventory(color, qty) VALUES(?, 0)`, c)
	}
	_, _ = db.Exec(`INSERT OR IGNORE INTO prices(color, price) VALUES
		('red', 50), ('blue', 40), ('green', 30)`)
	return nil
}
