package main

import (
	badger "github.com/dgraph-io/badger/v4"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}
