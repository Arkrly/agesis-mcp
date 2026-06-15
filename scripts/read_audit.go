package main

import (
	"fmt"
	"log"
	"os"

	"go.etcd.io/bbolt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/read_audit.go <path-to-audit.db>")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("audit"))
		if b == nil {
			return fmt.Errorf("audit bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			fmt.Printf("%s: %s\n", k, v)
			return nil
		})
	})

	if err != nil {
		log.Fatal(err)
	}
}
