package main

import (
	"fmt"
	"log/slog"
	"math/bits"
	"os"

	"github.com/dgraph-io/badger/v4"
)

func main() {
	dbPath := "../badger-data"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbPath = "badger-data" // try current dir if not found in parent
	}
	
	opts := badger.DefaultOptions(dbPath).WithReadOnly(true).WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		slog.Error("failed to open badger db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Scanning for valid promo codes (at least 2 bits set)...")

	count := 0
	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid() && count < 10; it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				if len(val) > 0 {
					mask := val[0]
					if bits.OnesCount8(mask) >= 2 {
						slog.Info("Valid promo code found", "code", string(item.Key()[6:]), "mask", fmt.Sprintf("%b", mask))
						count++
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		slog.Error("failed to scan badger db", "error", err)
		os.Exit(1)
	}

	if count == 0 {
		slog.Info("No valid promo codes found yet.")
	}
}
