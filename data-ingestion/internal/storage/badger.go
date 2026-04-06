package storage

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
)

type BadgerClient struct {
	DB *badger.DB
}

func NewBadgerDB(dir string) (*BadgerClient, error) {
	opts := badger.DefaultOptions(dir).
		WithLoggingLevel(badger.WARNING) // Keep logs clean
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerClient{DB: db}, nil
}

func (b *BadgerClient) Close() error {
	return b.DB.Close()
}

// MarkBatch processes a slice of codes in a single transaction for extreme speed
func (b *BadgerClient) MarkBatch(codes []string, bitOffset int) error {
	mask := byte(1 << bitOffset)

	// We use a retry loop because Badger throws ErrConflict if there are concurrent modifications.
	// (Our hash-partitioned workers make this rare, but it's best practice).
	for retries := 0; retries < 3; retries++ {
		err := b.DB.Update(func(txn *badger.Txn) error {
			for _, code := range codes {
				key := []byte("promo:" + code)
				var val byte = 0

				item, err := txn.Get(key)
				if err == nil {
					// Key exists, grab the current byte
					err = item.Value(func(v []byte) error {
						val = v[0]
						return nil
					})
					if err != nil {
						return err
					}
				} else if !errors.Is(err, badger.ErrKeyNotFound) {
					return err // Some other DB error
				}

				// Apply the bitmask for this file (e.g., File 0 = 001, File 1 = 010)
				val = val | mask

				if err := txn.Set(key, []byte{val}); err != nil {
					return err
				}
			}
			return nil
		})

		if err == nil {
			return nil // Success
		}
		if !errors.Is(err, badger.ErrConflict) {
			return err // Fail on non-conflict errors
		}
	}
	return errors.New("max retries reached for badger transaction")
}