package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("cannot read migrations dir: %w", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, name := range sqlFiles {
		path := filepath.Join(migrationsDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read %s: %w", name, err)
		}

		fmt.Printf("Running migration: %s\n", name)
		if _, err := db.Exec(string(data)); err != nil {
			// Ignore "already exists" errors, fail on real errors
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("migration %s failed: %w", name, err)
			}
		}
	}
	return nil
}