package cli

import (
	"bytescheme/common/log"
	"bytescheme/controller/generated/client/store"
	"bytescheme/controller/generated/models"
	"bytescheme/controller/shared"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

func setupStoreCommand() {
	storeCmd := &cobra.Command{
		Use:   "store get/set/delete ...",
		Short: "Runs the store operations",
	}
	rootCmd.AddCommand(storeCmd)

	// Get command
	storeCmdGet := &cobra.Command{
		Use:   "get",
		Short: "Get the value of a stored key or keys",
		Run:   storeCommandGet,
	}
	storeCmdGet.Flags().StringVarP(&storeCmdParams.host, "host", "s", storeCmdParams.host, "Host to be connected")
	storeCmdGet.Flags().IntVarP(&storeCmdParams.port, "port", "n", storeCmdParams.port, "Port to be connected")
	storeCmdGet.Flags().StringVarP(&storeCmdParams.apiKey, "apikey", "a", storeCmdParams.apiKey, "API key for the service access")
	storeCmdGet.Flags().StringVarP(&storeCmdParams.key, "key", "k", storeCmdParams.key, "Key of the data")
	storeCmdGet.Flags().BoolVarP(&storeCmdParams.isLocal, "local", "l", storeCmdParams.isLocal, "Connect to the DB locally")
	storeCmdGet.Flags().BoolVarP(&storeCmdParams.isPrefix, "prefix", "p", storeCmdParams.isPrefix, "If the key is prefix or not")
	storeCmdGet.MarkFlagRequired("key")
	storeCmd.AddCommand(storeCmdGet)

	// Set command
	storeCmdSet := &cobra.Command{
		Use:   "set",
		Short: "Set the value of a key",
		Run:   storeCommandSet,
	}
	storeCmdSet.Flags().StringVarP(&storeCmdParams.host, "host", "s", storeCmdParams.host, "Host to be connected")
	storeCmdSet.Flags().IntVarP(&storeCmdParams.port, "port", "n", storeCmdParams.port, "Port to be connected")
	storeCmdSet.Flags().StringVarP(&storeCmdParams.apiKey, "apikey", "a", storeCmdParams.apiKey, "API key for the service access")
	storeCmdSet.Flags().StringVarP(&storeCmdParams.key, "key", "k", storeCmdParams.key, "Key of the data")
	storeCmdSet.Flags().StringVarP(&storeCmdParams.value, "value", "v", storeCmdParams.value, "Value of the data")
	storeCmdSet.Flags().BoolVarP(&storeCmdParams.isLocal, "local", "l", storeCmdParams.isLocal, "Connect to the DB locally")
	storeCmdSet.MarkFlagRequired("key")
	storeCmdSet.MarkFlagRequired("value")
	storeCmd.AddCommand(storeCmdSet)

	// Set file command
	storeCmdSetf := &cobra.Command{
		Use:   "setf",
		Short: "Set the value of a key from a file",
		Run:   storeCommandSetFile,
	}
	storeCmdSetf.Flags().StringVarP(&storeCmdParams.host, "host", "s", storeCmdParams.host, "Host to be connected")
	storeCmdSetf.Flags().IntVarP(&storeCmdParams.port, "port", "n", storeCmdParams.port, "Port to be connected")
	storeCmdSetf.Flags().StringVarP(&storeCmdParams.apiKey, "apikey", "a", storeCmdParams.apiKey, "API key for the service access")
	storeCmdSetf.Flags().StringVarP(&storeCmdParams.key, "key", "k", storeCmdParams.key, "Key of the data")
	storeCmdSetf.Flags().StringVarP(&storeCmdParams.file, "file", "f", storeCmdParams.file, "Value file")
	storeCmdSetf.Flags().BoolVarP(&storeCmdParams.isLocal, "local", "l", storeCmdParams.isLocal, "Connect to the DB locally")
	storeCmdSetf.MarkFlagRequired("key")
	storeCmdSetf.MarkFlagRequired("file")
	storeCmd.AddCommand(storeCmdSetf)

	// Delete command
	storeCmdDelete := &cobra.Command{
		Use:   "delete",
		Short: "Delete a stored key or keys",
		Run:   storeCommandDelete,
	}
	storeCmdDelete.Flags().StringVarP(&storeCmdParams.host, "host", "s", storeCmdParams.host, "Host to be connected")
	storeCmdDelete.Flags().IntVarP(&storeCmdParams.port, "port", "n", storeCmdParams.port, "Port to be connected")
	storeCmdDelete.Flags().StringVarP(&storeCmdParams.apiKey, "apikey", "a", storeCmdParams.apiKey, "API key for the service access")
	storeCmdDelete.Flags().StringVarP(&storeCmdParams.key, "key", "k", storeCmdParams.key, "Key of the data")
	storeCmdDelete.Flags().BoolVarP(&storeCmdParams.isLocal, "local", "l", storeCmdParams.isLocal, "Connect to the DB locally")
	storeCmdDelete.Flags().BoolVarP(&storeCmdParams.isPrefix, "prefix", "p", storeCmdParams.isPrefix, "If the key is prefix or not")
	storeCmdDelete.MarkFlagRequired("key")
	storeCmd.AddCommand(storeCmdDelete)
}

func getStoreClient() store.ClientService {
	server := fmt.Sprintf("%s:%d", storeCmdParams.host, storeCmdParams.port)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	transport := httptransport.NewWithClient(server, "", []string{storeCmdParams.scheme}, client)
	return store.New(transport, strfmt.Default)
}

func getStoreAuth() runtime.ClientAuthInfoWriter {
	return httptransport.APIKeyAuth("Authorization", "header", storeCmdParams.apiKey)
}

func storeCommandGet(cmd *cobra.Command, args []string) {
	var keyValues models.KeyValues
	if storeCmdParams.isLocal {
		store, err := shared.CreateStore(false)
		if err != nil {
			log.Errorf("Unable to open store. Error: %s\n", err.Error())
			return
		}
		defer store.Close()
		keyValues, err = shared.ListStoreKeys(context.TODO(), store, storeCmdParams.key, storeCmdParams.isPrefix)
		if err != nil {
			log.Errorf("Unable to list store keys. Error: %s\n", err.Error())
			return
		}

	} else {

		client := getStoreClient()
		authParam := getStoreAuth()

		params := store.NewListStoreKeysParams()
		params.Key = storeCmdParams.key
		params.Prefix = &storeCmdParams.isPrefix
		ok, err := client.ListStoreKeys(params, authParam)
		if err != nil {
			log.Errorf("Unable to list store keys. Error: %s\n", err.Error())
			return
		}
		keyValues = ok.Payload
	}
	for _, keyValue := range keyValues {
		log.Infof("Key: %s, Value: %s\n", keyValue.Key, keyValue.Value)
	}
}

func storeCommandSet(cmd *cobra.Command, args []string) {
	keyValues := models.KeyValues{
		&models.KeyValue{
			Key:   storeCmdParams.key,
			Value: storeCmdParams.value,
		},
	}
	if storeCmdParams.isLocal {
		store, err := shared.CreateStore(false)
		if err != nil {
			log.Errorf("Unable to open store. Error: %s\n", err.Error())
			return
		}
		defer store.Close()
		keyValues, err = shared.UpdateStoreKeys(context.TODO(), store, keyValues)
		if err != nil {
			log.Errorf("Error occurred for key: %s, value: %s. Error: %s\n", storeCmdParams.key, storeCmdParams.value, err.Error())
			return
		}
	} else {
		client := getStoreClient()
		authParam := getStoreAuth()

		params := store.NewUpdateStoreKeysParams()
		params.Payload = keyValues
		ok, err := client.UpdateStoreKeys(params, authParam)
		if err != nil {
			log.Errorf("Unable to list store keys. Error: %s\n", err.Error())
			return
		}
		keyValues = ok.Payload
	}
	for _, keyValue := range keyValues {
		log.Infof("Key: %s, Value: %s\n", keyValue.Key, keyValue.Value)
	}
}

func storeCommandSetFile(cmd *cobra.Command, args []string) {
	if !strings.HasPrefix(storeCmdParams.file, "/") {
		path, err := os.Getwd()
		if err != nil {
			log.Errorf("Cannot open file %s. Error: %s\n", storeCmdParams.file, err.Error())
		}
		storeCmdParams.file = filepath.Join(path, storeCmdParams.file)
	}
	r, err := os.Open(storeCmdParams.file)
	if err != nil {
		log.Errorf("Cannot open file %s. Error: %s\n", storeCmdParams.file, err.Error())
		return
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("Cannot read file %s. Error: %s\n", storeCmdParams.file, err.Error())
	}
	storeCmdParams.value = string(b)
	keyValues := models.KeyValues{
		&models.KeyValue{
			Key:   storeCmdParams.key,
			Value: storeCmdParams.value,
		},
	}
	if storeCmdParams.isLocal {
		store, err := shared.CreateStore(false)
		if err != nil {
			log.Errorf("Unable to open store. Error: %s\n", err.Error())
			return
		}
		defer store.Close()

		keyValues, err = shared.UpdateStoreKeys(context.TODO(), store, keyValues)
		if err != nil {
			log.Errorf("Error occurred for key: %s, value: %s. Error: %s\n", storeCmdParams.key, storeCmdParams.value, err.Error())
			return
		}
	} else {
		client := getStoreClient()
		authParam := getStoreAuth()

		params := store.NewUpdateStoreKeysParams()
		params.Payload = keyValues
		ok, err := client.UpdateStoreKeys(params, authParam)
		if err != nil {
			log.Errorf("Unable to list store keys. Error: %s\n", err.Error())
			return
		}
		keyValues = ok.Payload
	}
	for _, keyValue := range keyValues {
		log.Infof("Key: %s, Value: %s\n", keyValue.Key, keyValue.Value)
	}
}

func storeCommandDelete(cmd *cobra.Command, args []string) {
	var keys models.Keys
	if storeCmdParams.isLocal {
		store, err := shared.CreateStore(false)
		if err != nil {
			log.Errorf("Unable to open store. Error: %s\n", err.Error())
			return
		}
		defer store.Close()
		keys, err = shared.DeleteStoreKeys(context.TODO(), store, storeCmdParams.key, storeCmdParams.isPrefix)
		if err != nil {
			log.Errorf("Error occurred for key: %s, value: %s. Error: %s\n", storeCmdParams.key, storeCmdParams.value, err.Error())
			return
		}
	} else {
		client := getStoreClient()
		authParam := getStoreAuth()

		params := store.NewDeleteStoreKeysParams()
		params.Key = storeCmdParams.key
		params.Prefix = &storeCmdParams.isPrefix

		ok, err := client.DeleteStoreKeys(params, authParam)
		if err != nil {
			log.Errorf("Unable to delete store key. Error: %s\n", err.Error())
			return
		}
		keys = ok.Payload
	}
	for _, key := range keys {
		log.Infof("Key %s deleted\n", key)
	}
}
