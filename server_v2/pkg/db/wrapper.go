package db

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	KEY_DELIMITER = "."
	ENTRY_PREFIX  = ""
)

type (
	Kind string
	KV   struct {
		K []byte
		V []byte
	}
	SupportedType string
)

type AssembleDisassemble[T any] struct {
	A func(t *T, b []byte)
	D func(t *T) []byte
}

type Mapping[T any] map[string]AssembleDisassemble[T]

type Database[T any] interface {
	Create(id string, object *T) error
	Get(id string) (*T, error)
	Delete(id string) error
}

type BadgerDatabase interface {
	View(fn func(txn *badger.Txn) error) error
	Update(fn func(txn *badger.Txn) error) error
}

type BadgerDatabaseWrapper[T any] struct {
	db       BadgerDatabase
	kind     Kind
	mappings Mapping[T]
}

func New[T any](db BadgerDatabase, kind Kind, mappings Mapping[T]) *BadgerDatabaseWrapper[T] {
	return &BadgerDatabaseWrapper[T]{
		db:       db,
		kind:     kind,
		mappings: mappings,
	}
}

func concat(first, second string) string {
	return first + KEY_DELIMITER + second
}

func (db *BadgerDatabaseWrapper[T]) entry(id string) []byte {
	return []byte(concat(ENTRY_PREFIX, concat(string(db.kind), id)))
}

func (db *BadgerDatabaseWrapper[T]) Create(id string, object *T) error {
	key := db.entry(id)
	err := db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)

		// Key not found is the desired error here.
		if err == badger.ErrKeyNotFound {
			return nil
		}

		if err != nil {
			return err
		}

		return &AlreadyExistsError{
			id: id,
		}
	})
	if err != nil {
		return err
	}

	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, []byte{})
		if err != nil {
			return err
		}

		for field, mapping := range db.mappings {
			err := txn.Set([]byte(concat(string(db.kind), concat(id, field))), mapping.D(object))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *BadgerDatabaseWrapper[T]) Get(id string) (*T, error) {
	t := new(T)
	err := db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(db.entry(id))
		if err == badger.ErrKeyNotFound {
			return &NotFoundError{
				id: id,
			}
		}

		if err != nil {
			return err
		}

		for name, mapping := range db.mappings {
			item, err := txn.Get([]byte(concat(concat(string(db.kind), id), name)))
			if err != nil {
				return err
			}

			err = item.Value(func(val []byte) error {
				mapping.A(t, val)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (db *BadgerDatabaseWrapper[T]) Delete(id string) error {
	key := db.entry(id)
	err := db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	if err != nil {
		return err
	}

	return db.db.Update(func(txn *badger.Txn) error {

		err := txn.Delete(key)
		if err != nil {
			return err
		}

		for field := range db.mappings {
			err := txn.Delete([]byte(concat(string(db.kind), concat(id, field))))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *BadgerDatabaseWrapper[T]) disassemble(id string, object *T) []KV {
	var index int
	kvs := make([]KV, len(db.mappings))
	for field, mapping := range db.mappings {
		kvs[index] = KV{
			K: []byte(concat(string(db.kind), concat(id, field))),
			V: mapping.D(object),
		}
		index++
	}
	return kvs
}

func (db *BadgerDatabaseWrapper[T]) assemble(id string, kvs []KV) (*T, error) {
	t := new(T)
	for _, kv := range kvs {
		var index int
		for i := len(kv.K) - 1; kv.K[i] != 0x2e; i-- {
			index = i
		}

		mapping, ok := db.mappings[string(kv.K[index:])]
		if !ok {
			return t, fmt.Errorf("invalid field")
		}

		mapping.A(t, kv.V)
	}

	return t, nil
}
