package database

import (
	"os"
	"path/filepath"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/dgraph-io/badger/v4"
)

// OpenBadger opens a BadgerDB instance with production settings (read-only by default as requested).
func OpenBadger() (*badger.DB, error) {
	// Find BadgerDB path - should be in root/badger-data
	cwd, _ := os.Getwd()
	badgerPath := filepath.Join(cwd, "..", "badger-data")
	if _, err := os.Stat(badgerPath); os.IsNotExist(err) {
		badgerPath = filepath.Join(cwd, "badger-data")
	}
	
	opts := badger.DefaultOptions(badgerPath).WithReadOnly(true).WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, apperrors.NewInternal("database: open badger", err)
	}
	return db, nil
}
