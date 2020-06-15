package db

import (
	"bytescheme/common/util"
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"
)

// Store is the persistent key-value store
type Store struct {
	config *StoreConfig
	db     *badger.DB
}

// StoreConfig is the config for the Store
type StoreConfig struct {
	Filepath string
	ReadOnly bool
}

// KeyValue is a placeholder for key value pair
type KeyValue struct {
	Key   string        `json:"key"`
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

// NewStore creates an instance of the Store with the given config
func NewStore(config *StoreConfig) (*Store, error) {
	sCfg := &StoreConfig{}
	// Deep copy
	err := util.Convert(config, sCfg)
	if err != nil {
		return nil, err
	}
	opts := badger.DefaultOptions(sCfg.Filepath)
	opts.ReadOnly = sCfg.ReadOnly
	opts.ValueLogLoadingMode = options.FileIO
	opts.TableLoadingMode = options.FileIO
	fmt.Printf("Store config options %+v\n", opts)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	store := &Store{
		config: sCfg,
		db:     db,
	}
	util.ShutdownHandler.RegisterCloseable(store)
	return store, nil
}

// Close closes the store it graceully
func (store *Store) Close() error {
	return store.db.Close()
}

// Set stores a key and value
func (store *Store) Set(kv *KeyValue) error {
	if kv == nil {
		return fmt.Errorf("Invalid key-value")
	}
	return store.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(kv.Key), []byte(kv.Value))
		if kv.TTL >= time.Second {
			entry = entry.WithTTL(kv.TTL)
		}
		return txn.SetEntry(entry)
	})
}

// Get gets the value of a key
func (store *Store) Get(key string) (*string, error) {
	var strVal *string
	err := store.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}
		if err != nil {
			return err
		}
		return item.Value(func(value []byte) error {
			// This func with val would only be called if item.Value encounters no error.
			valueCopy := string(append([]byte{}, value...))
			strVal = &valueCopy
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return strVal, nil
}

// Delete deletes a key
func (store *Store) Delete(key string) error {
	return store.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// Sets sets multiple key value pairs
func (store *Store) Sets(keyValues []*KeyValue) ([]*KeyValue, error) {
	wb := store.db.NewWriteBatch()
	defer wb.Cancel()
	savedKeys := make([]*KeyValue, 0, len(keyValues))
	for idx := range keyValues {
		keyValue := keyValues[idx]
		entry := badger.NewEntry([]byte(keyValue.Key), []byte(keyValue.Value))
		if keyValue.TTL >= time.Second {
			entry = entry.WithTTL(keyValue.TTL)
		}
		err := wb.SetEntry(entry) // Will create txns as needed.
		if err != nil {
			return []*KeyValue{}, err
		}
		savedKeys = append(savedKeys, keyValue)
	}
	err := wb.Flush()
	if err != nil {
		return []*KeyValue{}, err
	}
	return savedKeys, nil
}

// Gets gets all key value pairs with the common prefix
func (store *Store) Gets(prefix string) ([]*KeyValue, error) {
	keyValues := []*KeyValue{}
	err := store.db.View(func(txn *badger.Txn) error {
		iterOptions := badger.IteratorOptions{
			Prefix: []byte(prefix),
		}
		iter := txn.NewIterator(iterOptions)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			err := item.Value(func(value []byte) error {
				valueCopy := string(append([]byte{}, value...))
				keyValues = append(keyValues, &KeyValue{
					Key:   string(item.Key()),
					Value: string(valueCopy),
				})
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return []*KeyValue{}, err
	}
	return keyValues, nil
}

// GetKeys returns keys for the prefix
func (store *Store) GetKeys(prefix string) ([]string, error) {
	keys := []string{}
	err := store.db.View(func(txn *badger.Txn) error {
		iterOptions := badger.IteratorOptions{
			Prefix: []byte(prefix),
		}
		iter := txn.NewIterator(iterOptions)
		defer iter.Close()
		for iter.Rewind(); iter.Valid(); iter.Next() {
			item := iter.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return keys, nil
}

// Deletes deletes all key value pairs with the common prefix
func (store *Store) Deletes(prefix string) ([]string, error) {
	deletedKeys := []string{}
	keyValues, err := store.Gets(prefix)
	if err != nil {
		return deletedKeys, err
	}
	wb := store.db.NewWriteBatch()
	defer wb.Cancel()
	for idx := range keyValues {
		keyValue := keyValues[idx]
		err := wb.Delete([]byte(keyValue.Key))
		if err != nil {
			return []string{}, err
		}
		deletedKeys = append(deletedKeys, keyValue.Key)
	}
	err = wb.Flush()
	if err != nil {
		return []string{}, err
	}
	return deletedKeys, nil
}

// DoInTxn invokes the function in transaction
func (store *Store) DoInTxn(fn func(txn *badger.Txn) error) error {
	// Start a writable transaction.
	txn := store.db.NewTransaction(true)
	defer txn.Discard()
	err := fn(txn)
	if err != nil {
		return err
	}
	return txn.Commit()
}
