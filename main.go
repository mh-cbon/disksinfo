// Package disksinfo provides the list of partitions
package main

import (
	"encoding/json"
	"os"

	"github.com/mh-cbon/disksinfo/diskinfo"
)

func main() {
	loader := diskinfo.NewMultiOsLoader()
	p, err := loader.Load()
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(p); err != nil {
		panic(err)
	}
}
