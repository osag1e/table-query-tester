package query

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/osag1e/table-query-tester/internal/model"
	migrate "github.com/rubenv/sql-migrate"
)

func TestBookStore(t *testing.T) {
	txdb.Register("txdb", "postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")

	sqlDB, err := sql.Open("txdb", "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	err = runMigrations(sqlDB)
	if err != nil {
		t.Fatal(err)
	}

	bookStore := &BookStore{
		DB: sqlDB,
	}

	book := model.Books{
		Title:  "TableQueryTester",
		Author: "Osagie",
		Price:  20.29,
	}

	insertedBook, err := bookStore.InsertBook(&book)
	if err != nil {
		t.Fatalf("InsertBook returned an unexpected error: %v", err)
	}

	if insertedBook.ID == uuid.Nil {
		t.Errorf("Expected book ID to be set, but it was empty")
	}
}

func runMigrations(db *sql.DB) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	migrationsDir := filepath.Join(currentDir, "db_test_script")

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	if err := dropTables(db); err != nil {
		return fmt.Errorf("error dropping existing tables: %v", err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("error applying migrations: %v", err)
	}

	log.Printf("Applied %d migrations!\n", n)

	return nil
}

func dropTables(db *sql.DB) error {
	statements := []string{
		"DROP TABLE IF EXISTS store.books;",
	}

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}
