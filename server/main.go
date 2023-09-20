package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.formulabun.club/extractor/env"
	"go.formulabun.club/metadatadb"
)

func validateInput(data *[]metadatadb.File) error {
	contains := make(map[metadatadb.File]struct{})
	for i, d := range *data {
		if d.Filename == "" {
			return fmt.Errorf("Empty filename for entry %d: %v", i, d)
		}

		if d.Checksum == "" || len(d.Checksum) != 32 {
			return fmt.Errorf("Bad checksum for entry %d: %v", i, d)
		}

		_, ok := contains[d]
		if !ok {
			contains[d] = struct{}{}
		} else {
			return fmt.Errorf("Duplicate entry for entry %d: %v", i, d)
		}
	}
	return nil
}

func extractHandle(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var data []metadatadb.File
	err = json.Unmarshal(bs, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = validateInput(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		log.Println(err)
		return
	}

	go extract(&data)

	w.WriteHeader(http.StatusAccepted)
}

func main() {
	http.HandleFunc("/", extractHandle)
  log.Printf("hosting on %s\n", env.Port)
	log.Fatal(http.ListenAndServe(env.Port, nil))
}
