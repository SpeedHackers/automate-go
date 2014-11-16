package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type DB struct {
	path string
	dbs  map[string]*leveldb.DB
}

func OpenDB(path string) (*DB, error) {
	safePath := filepath.FromSlash(path)
	database := &DB{path: safePath}
	database.dbs = make(map[string]*leveldb.DB)
	safePathInfo, err := os.Stat(safePath)
	if err != nil {
		err := os.MkdirAll(safePath, os.ModeDir|0777)
		if err != nil {
			return nil, err
		}
	} else {
		if !safePathInfo.IsDir() {
			return nil, errors.New("Not a directory")
		}
	}
	return database, nil
}

func (db *DB) GetDB(database string) (*leveldb.DB, error) {
	_, exists := db.dbs[database]
	if !exists {
		newdb, err := leveldb.OpenFile(filepath.FromSlash(db.path+"/"+database), nil)
		if err != nil {
			return nil, err
		}
		db.dbs[database] = newdb
	}
	return db.dbs[database], nil
}

func (db *DB) Set(database, key string, value interface{}) error {
	d, err := db.GetDB(database)
	if err != nil {
		return err
	}
	jval, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = d.Put([]byte(key), jval, nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(database, key string, dest interface{}) error {
	d, err := db.GetDB(database)
	if err != nil {
		return err
	}
	jval, err := d.Get([]byte(key), nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jval, dest)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Del(database, key string) error {
	d, err := db.GetDB(database)
	if err != nil {
		return err
	}
	return d.Delete([]byte(key), nil)
}

func (db *DB) GetAll(database string) ([][]byte, error) {
	num, err := db.Count(database)
	if err != nil {
		return nil, err
	}
	iter, err := db.NewIterator(database)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, num)
	i := 0
	for iter.Next() {
		out[i] = make([]byte, len(iter.Value()))
		copy(out[i], iter.Value())
		i++
	}
	iter.Release()
	return out, iter.Error()
}

func (db *DB) Exists(database, key string) (bool, error) {
	d, err := db.GetDB(database)
	if err != nil {
		return false, err
	}
	_, err = d.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

func (db *DB) NewIterator(database string) (iterator.Iterator, error) {
	d, err := db.GetDB(database)
	if err != nil {
		return nil, err
	}
	return d.NewIterator(nil, nil), nil
}

func (db *DB) Count(database string) (int, error) {
	iter, err := db.NewIterator(database)
	if err != nil {
		return 0, err
	}
	count := 0
	for iter.Next() {
		count++
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		return 0, err
	}
	return count, nil
}
