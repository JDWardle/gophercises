package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/JDWardle/gophercises/urlshort"
	"github.com/boltdb/bolt"
)

var (
	yamlFile string
	jsonFile string
	boltDB   string
)

func main() {
	flag.StringVar(&yamlFile, "yaml", "paths.yml", "path to a YAML paths file")
	flag.StringVar(&jsonFile, "json", "paths.json", "path to a JSON paths file")
	flag.StringVar(&boltDB, "db", "paths.db", "sets the path to a Bolt database file")
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/testing":        "https://twitter.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the yamlHandler as the fallback.
	json, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	// Build the BoltHandler using the JSON handler as the fallback.
	boltDB, err := bolt.Open(boltDB, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer boltDB.Close()

	if err = initializeBoltDB(boltDB); err != nil {
		panic(err)
	}

	boltHandler, err := urlshort.BoltHandler(boltDB, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func initializeBoltDB(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("PathToURL"))
		if err != nil {
			return err
		}

		if err = b.Put([]byte("/bolt"), []byte("https://github.com/boltdb/bolt")); err != nil {
			return err
		}
		if err = b.Put([]byte("/gophercises"), []byte("https://gophercises.com")); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
