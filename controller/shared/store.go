package shared

import (
	db "bytescheme/common/db"
	"bytescheme/controller/generated/models"
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	// StoreFilepath is the db file path
	StoreFilepath = "/tmp/service-test"

	// Store is the singleton store instance
	Store *db.Store
)

// InitStore instantiates the store
func InitStore() {
	var err error
	Store, err = CreateStore(false)
	if err != nil {
		panic(err)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		Store.Close()
	}()
}

// CreateStore instantiates the store
func CreateStore(readOnly bool) (*db.Store, error) {
	path := os.Getenv("CONTROLLER_STORE_PATH")
	if path == "" {
		path = StoreFilepath
	}
	storeConfig := &db.StoreConfig{
		Filepath: path,
		ReadOnly: readOnly,
	}
	store, err := db.NewStore(storeConfig)
	if err != nil {
		return nil, err
	}
	return store, nil
}

// ListStoreKeys returns the key value pairs for the given parameters
func ListStoreKeys(ctx context.Context, store *db.Store, key string, isPrefix bool) (models.KeyValues, error) {
	keyValues := models.KeyValues{}
	if isPrefix {
		kvs, err := store.Gets(key)
		if err != nil {
			return keyValues, err
		}
		for _, kv := range kvs {
			keyValues = append(keyValues, &models.KeyValue{
				Key:   kv.Key,
				Value: kv.Value,
			})
		}
	} else {
		value, err := store.Get(key)
		if err != nil {
			return keyValues, err
		}
		if value == nil {
			return keyValues, nil
		}
		keyValues = append(keyValues, &models.KeyValue{
			Key:   key,
			Value: *value,
		})
	}
	return keyValues, nil
}

// UpdateStoreKeys updates key value pairs
func UpdateStoreKeys(ctx context.Context, store *db.Store, keyValues models.KeyValues) (models.KeyValues, error) {
	rKeyValues := models.KeyValues{}
	dbKeyValues := make([]*db.KeyValue, 0, len(keyValues))
	for _, keyValue := range keyValues {
		dbKeyValues = append(dbKeyValues, &db.KeyValue{
			Key:   keyValue.Key,
			Value: keyValue.Value,
		})
	}
	dbKeyValues, err := store.Sets(dbKeyValues)
	if err != nil {
		return rKeyValues, err
	}
	for _, keyValue := range dbKeyValues {
		rKeyValues = append(rKeyValues, &models.KeyValue{
			Key:   keyValue.Key,
			Value: keyValue.Value,
		})
	}
	return rKeyValues, nil
}

// DeleteStoreKeys deletes keys or key depending on the parameter
func DeleteStoreKeys(ctx context.Context, store *db.Store, key string, isPrefix bool) (models.Keys, error) {
	deletedKeys := models.Keys{}
	if isPrefix {
		keys, err := store.Deletes(key)
		if err != nil {
			return deletedKeys, err
		}
		for _, key := range keys {
			deletedKeys = append(deletedKeys, key)
		}
	} else {
		err := store.Delete(key)
		if err != nil {
			return deletedKeys, err
		}
		deletedKeys = append(deletedKeys, key)
	}
	return deletedKeys, nil
}
