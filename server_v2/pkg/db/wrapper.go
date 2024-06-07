package db

import (
	"regexp"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	KEY_DELIMITER    = string(rune(0x0))
	ENTRY_PREFIX     = string(rune(0x1))
	TAG_PREFIX       = string(rune(0x2))
	TAG_INDEX_PREFIX = string(rune(0x3))
)

var verifyID = regexp.MustCompile(`^[A-Za-z0-9_@!%+$.-]+$`)

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

type Peek interface {
	Peek(id string) error
}

type Database[T any] interface {
	Create(id string, object *T, tags ...string) error
	Update(id string, object *T) error
	Tags(id string, tags ...string) error
	ListByTag(tag string) ([]string, error)
	Get(id string) (*T, error)
	Delete(id string) error
	Peek
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

func (db *BadgerDatabaseWrapper[T]) prefixTag(id string) []byte {
	return []byte(concat(concat(TAG_PREFIX, concat(string(db.kind), id)), ""))
}

func (db *BadgerDatabaseWrapper[T]) prefixIndex(id string) []byte {
	return []byte(concat(concat(TAG_INDEX_PREFIX, concat(string(db.kind), id)), ""))
}

func (db *BadgerDatabaseWrapper[T]) tag(tag string, id string) []byte {
	return []byte(concat(TAG_PREFIX, concat(string(db.kind), concat(tag, id))))
}

func (db *BadgerDatabaseWrapper[T]) index(tag string, id string) []byte {
	return []byte(concat(TAG_INDEX_PREFIX, concat(string(db.kind), concat(id, tag))))
}

func (db *BadgerDatabaseWrapper[T]) validate(key string, field string) error {
	ok := verifyID.MatchString(key)
	if !ok {
		return &InvalidKeyError{
			field: field,
		}
	}

	return nil
}

func (db *BadgerDatabaseWrapper[T]) Create(id string, object *T, tags ...string) error {
	err := db.validate(id, "identifier")
	if err != nil {
		return err
	}

	for _, tag := range tags {
		err := db.validate(tag, "tag")
		if err != nil {
			return err
		}
	}

	key := db.entry(id)
	err = db.db.View(func(txn *badger.Txn) error {
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

		err = db.addTags(txn, id, tags...)
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

func (db *BadgerDatabaseWrapper[T]) Update(id string, object *T) error {
	err := db.validate(id, "identifier")
	if err != nil {
		return err
	}

	key := db.entry(id)
	err = db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)

		if err == badger.ErrKeyNotFound {
			return &NotFoundError{
				id: id,
			}
		}

		return nil
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

func (db *BadgerDatabaseWrapper[T]) addTags(txn *badger.Txn, id string, tags ...string) error {
	for _, tag := range tags {
		err := txn.Set(db.tag(tag, id), []byte(id))
		if err != nil {
			return err
		}

		err = txn.Set(db.index(tag, id), []byte(tag))
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BadgerDatabaseWrapper[T]) removeTags(txn *badger.Txn, id string) error {
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := db.prefixIndex(id)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		err := item.Value(func(val []byte) error {
			return txn.Delete(db.tag(string(val), id))
		})
		if err != nil {
			return err
		}

		err = txn.Delete(item.Key())
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *BadgerDatabaseWrapper[T]) Tags(id string, tags ...string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return db.addTags(txn, id, tags...)
	})
}

func (db *BadgerDatabaseWrapper[T]) ListByTag(tag string) ([]string, error) {
	var tags []string
	db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := db.prefixTag(tag)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				tags = append(tags, string(val))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return tags, nil
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

func (db *BadgerDatabaseWrapper[T]) Peek(id string) error {
	return db.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(db.entry(id))
		if err == badger.ErrKeyNotFound {
			return &NotFoundError{
				id: id,
			}
		}

		if err != nil {
			return err
		}

		return nil
	})
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

		return db.removeTags(txn, id)
	})
}
