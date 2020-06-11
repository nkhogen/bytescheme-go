package tool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

var (
	importRegex           = regexp.MustCompile("import \\(")
	globalMiddlewareRegex = regexp.MustCompile("func setupGlobalMiddleware\\(handler http\\.Handler\\) http\\.Handler \\{[\n\\s\ra-zA-Z]+\\}")

	importRep = []byte(`import(
		"fmt"
		"strings"
		_ "bytescheme/controller/generated/statik"
		"github.com/rakyll/statik/fs"
		`)
	middlewareRep = []byte(`func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request received for %s\n", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/v1") {
			handler.ServeHTTP(w, r)
		} else {
			statikFS, err := fs.New()
			if err != nil {
				fmt.Printf("Cannot create statik FS. Error: %s\n", err.Error())
				statikFS = http.Dir(".")
			}
			file, err := statikFS.Open(r.URL.Path)
			if err == nil {
				file.Close()
			} else {
				r.URL.Path="index.html"
			}
			http.FileServer(statikFS).ServeHTTP(w, r)
		}
	})
}`)
)

func ReplaceGlobalMiddlewareFunc(fpath string) error {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	data = importRegex.ReplaceAll(data, importRep)
	data = globalMiddlewareRegex.ReplaceAll(data, middlewareRep)

	fmt.Println(string(data))

	return ioutil.WriteFile(fpath, data, 0)
}

func ProcessSwagger(fpath string) error {
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
	fmt.Println(data)
	iface, ok := data["paths"]
	if !ok {
		fmt.Printf("Key 'paths' not found\n")
		return nil
	}
	paths, ok := iface.(map[string]interface{})
	if !ok {
		return nil
	}
	for path, iface := range paths {
		methods, ok := iface.(map[string]interface{})
		if !ok {
			fmt.Printf("Map expected for path value\n")
			continue
		}
		for method, iface := range methods {
			fmt.Printf("Method %s %s\n", method, path)
			body, ok := iface.(map[string]interface{})
			if !ok {
				fmt.Printf("Map expected for method value")
				continue
			}
			iface, ok := body["operationId"]
			if !ok {
				fmt.Printf("Key 'operationId' not found\n")
				continue
			}
			fmt.Println(path, " ==>  ", iface)
		}
	}
	return nil
}
