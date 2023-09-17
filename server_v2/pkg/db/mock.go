package db

import "github.com/dgraph-io/badger/v4"

type MockBadgerDatabase struct {
	OnView   func(fn func(txn *badger.Txn) error) error
	OnUpdate func(fn func(txn *badger.Txn) error) error
}

func (db *MockBadgerDatabase) View(fn func(txn *badger.Txn) error) error {
	return db.OnView(fn)
}

func (db *MockBadgerDatabase) Update(fn func(txn *badger.Txn) error) error {
	return db.OnUpdate(fn)
}

type NoopBadgerDatabase struct{}

func (db *NoopBadgerDatabase) View(_ func(txn *badger.Txn) error) error {
	return nil
}

func (db *NoopBadgerDatabase) Update(_ func(txn *badger.Txn) error) error {
	return nil
}
