package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
)

func (d *Database) Migrate() error {
	fmt.Println("Migrating database...")

	client := stdlib.OpenDBFromPool(d.Pool)
	defer client.Close()

	driver, err := pgx.WithInstance(client, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgress driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres",
		driver,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run up migrations: %w", err)
		}
		fmt.Println("No new migrations to run")
	}

	fmt.Println("Database migrated successfully")
	return nil
}
